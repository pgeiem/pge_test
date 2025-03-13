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

/*
// Loop over all sequences from start until window end and define at each time the applicable sequence

	func (tsi TariffSequenceInventory) ResolveSequenceApplicability(now time.Time, window time.Duration) (PrioritizedSequences, error) {
		var out PrioritizedSequences
		var t time.Duration
		for t < window {

			// Loop over all sequences by priority order
			for _, s := range tsi {
				var relspan RelativeTimeSpan
				valid := false
				fmt.Println("Sequence", s.Name)
				// Check if the sequence has a validity period
				if s.ValidityPeriod.IsValid() {
					// Check if the sequence is applicable at this instant
					within, timespan, err := s.ValidityPeriod.IsWithin(now.Add(t))
					fmt.Println("Within", within, timespan, err)
					if err != nil {
						return nil, err
					}
					if within {
						relspan = timespan.ToRelativeTimeSpan(now)
						if relspan.From < 0 {
							relspan.From = 0
						}
					}
					valid = within
				} else {
					// If the sequence has no validity period, it is always applicable (default sequence)
					relspan = RelativeTimeSpan{
						From: 0,
						To:   window,
					}
					valid = true
				}
				// TODO check if sequence quota or condition is also met

				if valid {
					out = append(out, PrioritizedSequence{
						RelativeTimeSpan: relspan,
						Sequence:         &s,
					})
					t = relspan.To
				}
			}
		}
		return out, nil
	}

	func (tsi TariffSequenceInventory) Merge(now time.Time, window time.Duration) (SolverRules, error) {
		var out SolverRules
		prio, err := tsi.ResolveSequenceApplicability(now, window)
		fmt.Println(prio)
		if err != nil {
			return out, err
		}
		for _, s := range prio {
			out = append(out, s.Sequence.Solver.ExtractRulesInRange(s.RelativeTimeSpan)...)
		}
		return out, nil
	}
*/

/*
func (inventory *TariffSequenceInventory) Solve(now time.Time, window time.Duration) {
	now = now.Truncate(time.Second)
	for i := range *inventory {
		(*inventory)[i].Solve(now, window)
	}

	//TODO handle error
	rules, _ := inventory.Merge(now, window)
	out := rules.GenerateOutput(now, true)
	json, _ := out.ToJson()
	fmt.Println(string(json))
}*/

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
