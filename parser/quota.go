package parser

import (
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/iem-rd/quoteengine/engine"
)

// TODO rewrite this based on ast.Node
func unmarshalQuota(quota *engine.Quota, data []byte) error {
	temp := struct {
		Type string `yaml:"type"`
	}{}

	err := yaml.Unmarshal(data, &temp)
	if err != nil {
		return fmt.Errorf("failed to parse quota type: %w", err)
	}

	switch temp.Type {
	case "duration":
		q := struct {
			engine.DurationQuota `yaml:",inline"`
			Type                 string `yaml:"type"`
		}{}
		if err := yaml.UnmarshalWithOptions(data, &q, decoderOptions()...); err != nil {
			return fmt.Errorf("failed to parse duration quota: %w", err)
		}
		*quota = &q.DurationQuota
	case "counter":
		q := struct {
			engine.CounterQuota `yaml:",inline"`
			Type                string `yaml:"type"`
		}{}
		if err := yaml.UnmarshalWithOptions(data, &q, decoderOptions()...); err != nil {
			return fmt.Errorf("failed to parse counter quota: %w", err)
		}
		*quota = &q.CounterQuota
	default:
		return fmt.Errorf("unknown quota type: %s", temp.Type)
	}

	return nil
}
