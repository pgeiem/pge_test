package engine

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/iem-rd/quote-engine/table"
	"github.com/iem-rd/quote-engine/timeutils"
)

type TariffSequence struct {
	Name           string
	ValidityPeriod timeutils.RecurrentTimeSpan
	Quota          Quota
	Rules          SolvableRules
	Solver         Solver
	Limits         TariffLimits
}

// New TariffSequence from a name, a recurrent segment and a quota
func NewTariffSequence() TariffSequence {
	return TariffSequence{
		Solver: NewSolver(),
	}
}

// Stringer for TariffSequence display sequence name, segment potential attached quota and list the rules
func (ts TariffSequence) String() string {
	var sb strings.Builder
	sb.WriteString(ts.Name)
	sb.WriteString(": ")
	sb.WriteString(ts.ValidityPeriod.String())
	if ts.Quota != nil {
		sb.WriteString(" ")
		sb.WriteString(ts.Quota.String())
	}
	return sb.String()
}

func (ts TariffSequence) Solve(now time.Time, window time.Duration, globalNonpaying AbsoluteNonPayingRules) {
	fmt.Println()
	table.TitleTheme().Println("Solving sequence", ts.Name)

	ts.Solver.SetWindow(now, window)
	// Append first all global nonpaying rules...
	for i := range globalNonpaying {
		globalNonpaying[i].ToSolverRules(now, now.Add(window), ts.Solver.Append)
	}
	// ... then the sequence rules
	for i := range ts.Rules {
		ts.Rules[i].ToSolverRules(now, now.Add(window), ts.Solver.Append)
	}
	ts.Solver.Solve()
}

type TariffSequenceInventory []TariffSequence

// Stringer for TariffSequenceInventory display all sequences as a dashed list
func (tsi TariffSequenceInventory) String() string {
	var sb strings.Builder
	sb.WriteString("TariffSequences:\n")
	for _, s := range tsi {
		sb.WriteString(" - ")
		sb.WriteString(s.String())
		sb.WriteString("\n")

		sb.WriteString("    Rules:\n")
		for _, r := range s.Rules {
			sb.WriteString("     - ")
			sb.WriteString(r.String())
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// Merge all sequences into a single list of rules
func (inventory TariffSequenceInventory) Merge(now time.Time, window time.Duration) (SolverRules, error) { //TODO remove error
	var out SolverRules

	if len(inventory) == 0 {
		return out, nil
	}

	// If there is only one sequence, return its rules directly, skipping merging
	if len(inventory) == 1 {
		fmt.Println("Single sequence, skipping merging")
		return inventory[0].Solver.ExtractRulesInRange(timeutils.RelativeTimeSpan{From: 0, To: window}), nil
	}

	// Create a scheduler and solve all sequences excepted the last one
	scheduler := NewScheduler()
	scheduler.SetWindow(now, window)
	for i := range (inventory)[:len(inventory)-1] {
		scheduler.AddSequence(&inventory[i])
	}
	// Add latest sequences. Lowest priority sequence must always match the window as it's the default one
	scheduler.Append(SchedulerEntry{
		RelativeTimeSpan: timeutils.RelativeTimeSpan{From: 0, To: window},
		Sequence:         &inventory[len(inventory)-1],
	})
	fmt.Println("Scheduler:", scheduler.String())

	// Merge all sequences
	scheduler.entries.Ascend(func(entry SchedulerEntry) bool {
		rules := entry.Sequence.Solver.ExtractRulesInRange(entry.RelativeTimeSpan)
		fmt.Println("\nMerging", entry.Sequence.Name, len(rules), "rules between", entry.RelativeTimeSpan, "to", len(out), "rules already in output")

		rules.PrintAsTable(fmt.Sprintf("Rules from %s before applying limits (%d rules):", entry.Sequence.Name, len(rules)), now)

		// Calcul the position of the rules in the output and apply the sequence limits
		limits := entry.Sequence.Limits
		offsetAmout, offsetDuration := out.SumAll()
		limits.AddOffset(offsetAmout, offsetDuration)
		rules = rules.ApplyLimits(limits)
		// FIXME: the limits are applied for each sheduler entries but should be applied only once for all scheduler entries from the same sequence
		// For example if one sequence has 2 entries, the limits are applied twice instead of once globally

		rules.PrintAsTable(fmt.Sprintf("Rules from %s with limits applied (%d rules):", entry.Sequence.Name, len(rules)), now)

		out = append(out, rules...)
		return true
	})

	return out, nil
}

func (inventory TariffSequenceInventory) Solve(now time.Time, window time.Duration, globalNonpaying AbsoluteNonPayingRules) {
	//Solve all sequences individually
	for i := range inventory {
		inventory[i].Solve(now, window, globalNonpaying)
	}
}

func (out *TariffSequenceInventory) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {

	// Temporarily unmarshal the sequences section in a temporary struct
	temp := []struct {
		Name           string                      `yaml:"name"`
		ValidityPeriod timeutils.RecurrentTimeSpan `yaml:",inline"`
		Quota          string                      `yaml:"quota,"`
		Rules          SolvableRules               `yaml:"rules"`
		Limits         TariffLimits                `yaml:",inline"`
	}{}
	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	// Convert the temporary struct into the final struct and link the referred quotas
	*out = make(TariffSequenceInventory, 0, len(temp))
	for i, n := range temp {

		seq := NewTariffSequence()
		seq.Name = n.Name
		seq.ValidityPeriod = n.ValidityPeriod
		seq.Rules = n.Rules
		seq.Limits = n.Limits

		// Some validity check
		isValidityPeriodValid := n.ValidityPeriod.Start != nil && n.ValidityPeriod.End != nil
		isLastSequence := i == len(temp)-1
		if !isValidityPeriodValid && !isLastSequence {
			return fmt.Errorf("validity period is not valid for sequence %s", n.Name)
		}
		if isValidityPeriodValid && isLastSequence {
			// Last sequence must have an empty validity period
			return fmt.Errorf("last sequence must have an empty validity period")
		}

		// Search the coresponding quota
		if n.Quota != "" {
			// Get the quota from the context
			quota, exists := ContextGetQuotaByName(ctx, n.Quota)
			if !exists {
				return fmt.Errorf("unknown quota: %s", n.Quota)
			}
			seq.Quota = quota
		}
		*out = append(*out, seq)
	}

	return nil
}
