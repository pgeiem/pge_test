package engine

import (
	"strings"
	"time"
)

// Non paying rules
type NonPayingRule struct {
	RuleName         string `yaml:"name"`
	RecurrentSegment `yaml:",inline"`
}

// Stringer for NonPayingRule display rule name and segment
func (npr NonPayingRule) String() string {
	return npr.RuleName + ": " + npr.RecurrentSegment.String()
}

func (npr NonPayingRule) Name() string {
	return npr.RuleName
}

func (npr NonPayingRule) RelativeToWindow(from, to time.Time, iterator func(RelativeTimeSpan) bool) {
	npr.RecurrentSegment.BetweenIterator(from, to, func(s Segment) bool {
		return iterator(s.ToRelativeTimeSpan(from))
	})
}

func (npr NonPayingRule) Policies() (StartTimePolicy, RuleResolutionPolicy) {
	return FixedPolicy, TruncatePolicy
}

type NonPayingInventory []NonPayingRule

// Stringer for NonPayingInventory display all rules as a dashed list
func (npi NonPayingInventory) String() string {
	var sb strings.Builder
	sb.WriteString("NonPaying:\n")
	for _, r := range npi {
		sb.WriteString(" - ")
		sb.WriteString(r.String())
		sb.WriteString("\n")
	}
	return sb.String()
}
