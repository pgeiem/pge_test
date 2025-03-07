package engine

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type TariffSequence struct {
	Name           string
	ValidityPeriod RecurrentSegment
	Quota          Quota
	NonPaying      NonPayingInventory
	RelativeRules  RelativeTariffRulesSequence
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
	//TODO append FlatRate rules
	for i := range ts.NonPaying {
		ts.Solver.Append(ts.NonPaying[i])
	}
	ts.Solver.Append(ts.RelativeRules...)
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

		sb.WriteString("    RelativeRules:\n")
		for _, r := range s.RelativeRules {
			sb.WriteString("     - ")
			sb.WriteString(r.String())
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

/*
type PrioritizedSequence struct {
	RelativeTimeSpan
	Sequence *TariffSequence
}

type PrioritizedSequences []PrioritizedSequence

// Loop over all sequences from start until window end and define at each time the applicable sequence
func (tsi TariffSequenceInventory) ResolveSequenceApplicability(now time.Time, window time.Duration) (PrioritizedSequences, error) {
	var out PrioritizedSequences
	t := now
	for t.Before(now.Add(window)) {

		// Loop over all sequences by priority order
		for _, s := range tsi {
			// Check if the sequence is applicable at this instant
			within, timespan, err := s.ValidityPeriod.IsWithinWithSegment(now)
			if err != nil {
				return nil, err
			}
			// TODO check if sequence quota or condition is also met
			if within {
				relspan := timespan.ToRelativeTimeSpan(now)
				out = append(out, PrioritizedSequence{
					RelativeTimeSpan: relspan,
					Sequence:         &s,
				})
				t = timespan.End
			}
		}
	}
	return out, nil
}
*/

func (inventory *TariffSequenceInventory) Solve(now time.Time, window time.Duration) {
	for i := range *inventory {
		(*inventory)[i].Solve(now, window)
	}
	//TODO resolve sequence applicability
}

func (out *TariffSequenceInventory) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {

	// Temporarily unmarshal the sequences section in a temporary struct
	temp := []struct {
		Name           string                      `yaml:"name"`
		ValidityPeriod RecurrentSegment            `yaml:",inline"`
		Quota          string                      `yaml:"quota,"`
		NonPayingRules NonPayingInventory          `yaml:"nonpaying"`
		RelativeRules  RelativeTariffRulesSequence `yaml:"rules"`
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
		seq.RelativeRules = n.RelativeRules

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
