package engine

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/btree"
	"github.com/iem-rd/quote-engine/timeutils"
)

type SchedulerEntry struct {
	timeutils.RelativeTimeSpan
	Sequence *TariffSequence
}

func (entry SchedulerEntry) TruncateAfter(after time.Duration) SchedulerEntry {
	entry.To = after
	return entry
}

func (entry SchedulerEntry) TruncateBefore(before time.Duration) SchedulerEntry {
	entry.From = before
	return entry
}

func (entry SchedulerEntry) TruncateBetween(truncateStart, truncateEnd time.Duration) SchedulerEntries {
	return SchedulerEntries{entry.TruncateAfter(truncateStart), entry.TruncateBefore(truncateEnd)}
}

func (entry SchedulerEntry) String() string {
	return entry.RelativeTimeSpan.String() + " " + entry.Sequence.Name
}

type SchedulerEntries []SchedulerEntry

func (entries SchedulerEntries) String() string {
	var sb strings.Builder
	sb.WriteString("Scheduler:\n")
	for _, e := range entries {
		sb.WriteString(" - ")
		sb.WriteString(e.Sequence.Name)
		sb.WriteString(" ")
		sb.WriteString(e.String())
		sb.WriteString("\n")
	}
	return sb.String()
}

type Scheduler struct {
	now     time.Time
	window  time.Duration
	entries *btree.BTreeG[SchedulerEntry]
}

func NewScheduler() Scheduler {

	// Sorting function for B-Tree storing all solver segments
	RulesLess := func(i, j SchedulerEntry) bool {
		return i.From < j.From
	}

	return Scheduler{
		entries: btree.NewG(2, RulesLess),
	}
}

func (s *Scheduler) String() string {
	var sb strings.Builder
	sb.WriteString("Scheduler Entries:\n")
	s.entries.Ascend(func(entry SchedulerEntry) bool {
		sb.WriteString("  - ")
		sb.WriteString(entry.String())
		sb.WriteString("\n")
		return true
	})
	return sb.String()
}

func (s *Scheduler) SetWindow(now time.Time, window time.Duration) {
	s.now = now
	s.window = window
}

// Solve the rule against an Higer Priority Rule resolving the conflict according to rule policy
// a collection of new rules containing 0, 1, or 2 rules is returned and current rule is not changed
// the second return value is true if the rule has intersected and has been changed, false if untouched
func (s *Scheduler) solveVsSingle(lpSpan SchedulerEntry, hpSpan SchedulerEntry) (SchedulerEntries, bool) {

	// trivial case, both rules don't overlap
	if (hpSpan.To <= lpSpan.From) ||
		(hpSpan.From >= lpSpan.To) {
		return SchedulerEntries{lpSpan}, false
	}

	// high priority rule is partially after low priority rule, then low priority rule end is truncated
	if hpSpan.From >= lpSpan.From && hpSpan.To >= lpSpan.To {
		return SchedulerEntries{lpSpan.TruncateAfter(hpSpan.From)}, true
	}

	// high priority rule is partially before low priority rule, then low priority rule end is truncated
	if hpSpan.From <= lpSpan.From && hpSpan.To <= lpSpan.To {
		return SchedulerEntries{lpSpan.TruncateBefore(hpSpan.To)}, true
	}

	// high priority rule completely overlap low priority rule, then remove the low priority rule
	if hpSpan.From <= lpSpan.From && hpSpan.To >= lpSpan.To {
		return SchedulerEntries{}, true
	}

	// high priority rule is in middle of low priority rule, then low priority rule middle is truncated
	if hpSpan.From >= lpSpan.From && hpSpan.To <= lpSpan.To {
		return lpSpan.TruncateBetween(hpSpan.From, hpSpan.To), true
	}

	return SchedulerEntries{}, true
}

// Solve the rule against a collection of Higer Priority Rule resolving the conflict according to rules policy
// a collection of new rules is returned and current rule is not changed
func (s *Scheduler) Append(lpEntry SchedulerEntry) {
	var newEntries SchedulerEntries

	fmt.Println("Solving scheduler entry", lpEntry, "from", lpEntry.From, "to", lpEntry.To)

	// Loop over all entries and solve the current entry against each of them
	s.entries.Ascend(func(hpEntry SchedulerEntry) bool {
		ret, _ := s.solveVsSingle(lpEntry, hpEntry)
		switch len(ret) {
		case 0: // Entry deleted
			lpEntry = SchedulerEntry{}
			return false
		case 1: // Entry truncated or untouched
			lpEntry = ret[0]
		case 2: // Entry splitted
			newEntries = append(newEntries, ret[0]) // Left part may be inserted in the new rules
			lpEntry = ret[1]                        // right part is the new rule to solve
		}
		return true
	})

	//Insert the last rule in the new entries list
	newEntries = append(newEntries, lpEntry)

	// Effectively insert all parts of the resolved rules in the rules collection
	for _, entry := range newEntries {
		if entry.Duration() > time.Duration(0) {
			s.entries.ReplaceOrInsert(entry)
		}
	}
}

func (s *Scheduler) AddSequence(seq *TariffSequence) {
	for _, seg := range seq.ValidityPeriod.Between(s.now, s.now.Add(s.window)) {
		timespan := seg.ToRelativeTimeSpan(s.now)
		if timespan.From < 0 {
			timespan.From = 0
		}
		s.Append(SchedulerEntry{timespan, seq})
	}
}
