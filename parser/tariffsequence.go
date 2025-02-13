package parser

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
	//Rules           []TariffRule `yaml:"rules"`
}

func parseSequences(node ast.Node, quotas engine.QuotaInventory) (engine.TariffSequenceInventory, error) {
	var out engine.TariffSequenceInventory

	temp := []TariffSequence{}
	err := yaml.NodeToValue(node, &temp, decoderOptions()...)
	if err != nil {
		return out, fmt.Errorf("failed to parse sequences section: %w", err)
	}

	out = make(engine.TariffSequenceInventory, 0, len(temp))
	for i, n := range temp {
		fmt.Println(">>>", i, n)

		// Find the quota
		quota, exists := quotas[n.Quota]
		if !exists {
			return out, fmt.Errorf("unknown quota: %s", n.Quota)
		}

		seq := engine.TariffSequence{
			Name:           n.Name,
			ValidityPeriod: n.ValidityPeriod,
			Quota:          quota,
		}
		out = append(out, seq)
	}

	return out, nil
}
