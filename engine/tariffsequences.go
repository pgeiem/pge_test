package engine

import (
	"context"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml/ast"
)

type TariffSequence struct {
	Name           string             `yaml:"name"`
	ValidityPeriod RecurrentSegment   `yaml:",inline"`
	Quota          Quota              `yaml:"quota,"`
	NonPaying      NonPayingInventory `yaml:"nonpaying"`
	//Rules           []TariffRule `yaml:"rules"`
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

func (ts *TariffSequence) UnmarshalYAML(unmarshal func(interface{}) error) error {
	temp := struct {
		Name             string `yaml:"name"`
		RecurrentSegment `yaml:",inline"`
		Quota            string `yaml:"quota,"`
	}{}

	err := unmarshal(&temp)
	if err != nil {
		return err
	}
	*ts = TariffSequence{Name: temp.Name,
		ValidityPeriod: temp.RecurrentSegment,
		Quota:          nil, //FIXME: Quota(temp.Quota),
	}
	return nil
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
	}
	return sb.String()
}

func (out *TariffSequenceInventory) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {

	// Temporarily unmarshal the sequences section in a temporary struct
	temp := []struct {
		Name           string             `yaml:"name"`
		ValidityPeriod RecurrentSegment   `yaml:",inline"`
		Quota          string             `yaml:"quota,"`
		NonPaying      NonPayingInventory `yaml:"nonpaying"`
		Rules          []ast.Node         `yaml:"rules"`
	}{}
	err := unmarshal(&temp)
	if err != nil {
		return fmt.Errorf("failed to parse sequences section: %w", err)
	}

	// Convert the temporary struct into the final struct and link the referred quotas
	*out = make(TariffSequenceInventory, 0, len(temp))
	for _, n := range temp {

		seq := TariffSequence{
			Name:           n.Name,
			ValidityPeriod: n.ValidityPeriod,
			NonPaying:      n.NonPaying,
		}

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
