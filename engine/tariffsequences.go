package engine

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type TariffSequence struct {
	Name           string
	ValidityPeriod RecurrentTimeSpan
	Quota          Quota
	NonPaying      NonPayingInventory
	Rules          SolvableRules
	Solver         Solver
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

func (ts TariffSequence) Solve(now time.Time, window time.Duration) {
	ts.Solver.SetWindow(now, window)
	//TODO append NonPaying rules
	/*for i := range ts.NonPaying {
		ts.NonPaying[i].ToSolverRules(now, now.Add(window), ts.Solver.Append)
		ts.Solver.Append(ts.NonPaying[i])
	}*/
	for i := range ts.Rules {
		ts.Rules[i].ToSolverRules(now, now.Add(window), ts.Solver.Append)
	}

	//FIXME PGE
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

		sb.WriteString("    NonPaying:\n")
		for _, r := range s.NonPaying {
			sb.WriteString("      - ")
			sb.WriteString(r.String())
			sb.WriteString("\n")
		}

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
func (inventory TariffSequenceInventory) merge(now time.Time, window time.Duration) (SolverRules, error) { //TODO remove error
	var out SolverRules

	if len(inventory) == 0 {
		return out, nil
	}

	// If there is only one sequence, return its rules directly skipping merging
	if len(inventory) == 1 {
		fmt.Println("Single sequence, skipping merging")
		return inventory[0].Solver.ExtractRulesInRange(RelativeTimeSpan{0, window}), nil
	}

	// Create a scheduler and solve all sequences excepted the last one
	scheduler := NewScheduler()
	scheduler.SetWindow(now, window)
	for i := range (inventory)[:len(inventory)-1] {
		scheduler.AddSequence(&inventory[i])
	}
	// Add latest sequences. Lowest priority sequence must always match the window as it is the default one
	scheduler.Append(SchedulerEntry{
		RelativeTimeSpan: RelativeTimeSpan{0, window},
		Sequence:         &inventory[len(inventory)-1],
	})
	fmt.Println("Scheduler:", scheduler.String())

	// Merge all sequences
	scheduler.entries.Ascend(func(entry SchedulerEntry) bool {
		out = append(out, entry.Sequence.Solver.ExtractRulesInRange(entry.RelativeTimeSpan)...)
		return true
	})
	return out, nil
}

func (inventory TariffSequenceInventory) Solve(now time.Time, window time.Duration) {
	for i := range inventory {
		inventory[i].Solve(now, window)
	}

	//TODO handle error
	rules, _ := inventory.merge(now, window)
	out := rules.GenerateOutput(now, true)
	json, _ := out.ToJson()
	fmt.Println(string(json))
}

func (out *TariffSequenceInventory) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {

	// Temporarily unmarshal the sequences section in a temporary struct
	temp := []struct {
		Name           string             `yaml:"name"`
		ValidityPeriod RecurrentTimeSpan  `yaml:",inline"`
		Quota          string             `yaml:"quota,"`
		NonPayingRules NonPayingInventory `yaml:"nonpaying"`
		Rules          SolvableRules      `yaml:"rules"`
	}{}
	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	// Convert the temporary struct into the final struct and link the referred quotas
	*out = make(TariffSequenceInventory, 0, len(temp))
	for _, n := range temp {

		seq := NewTariffSequence()
		seq.Name = n.Name
		seq.ValidityPeriod = n.ValidityPeriod
		seq.NonPaying = n.NonPayingRules
		seq.Rules = n.Rules

		// Search the coresponding quota
		if n.Quota != "" {
			quotas, ok := ctx.Value("quotas").(QuotaInventory)
			if !ok {
				return fmt.Errorf("quotas not found in context")
			}

			quota, exists := quotas[n.Quota]
			if !exists {
				return fmt.Errorf("unknown quota: %s", n.Quota)
			}
			seq.Quota = quota
		}
		*out = append(*out, seq)
	}

	return nil
}
