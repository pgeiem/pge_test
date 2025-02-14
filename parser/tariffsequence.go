package parser

/*
import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/iem-rd/quoteengine/engine"
)

type TariffSequence struct {
	Name           string                  `yaml:"name"`
	ValidityPeriod engine.RecurrentSegment `yaml:",inline"`
	Quota          string                  `yaml:"quota,"`
	Rules          []ast.Node              `yaml:"rules"`
}

// Parser for TariffSequence section
func parseSequences(node ast.Node, quotas engine.QuotaInventory) (engine.TariffSequenceInventory, error) {
	var out engine.TariffSequenceInventory

	// Temporarily unmarshal the sequences section in a temporary struct
	temp := []TariffSequence{}
	err := yaml.NodeToValue(node, &temp, decoderOptions()...)
	if err != nil {
		return out, fmt.Errorf("failed to parse sequences section: %w", err)
	}

	// Convert the temporary struct into the final struct and link the referred quotas
	out = make(engine.TariffSequenceInventory, 0, len(temp))
	for _, n := range temp {

		seq := engine.TariffSequence{
			Name:           n.Name,
			ValidityPeriod: n.ValidityPeriod,
		}

		fmt.Println(">>>", n.Rules)

		// Search the coresponding quota
		if n.Quota != "" {
			quota, exists := quotas[n.Quota]
			if !exists {
				return out, fmt.Errorf("unknown quota: %s", n.Quota)
			}
			seq.Quota = quota
		}
		out = append(out, seq)
	}

	return out, nil
}
*/
