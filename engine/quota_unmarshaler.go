package engine

import (
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

// TODO rewrite this based on ast.Node
func unmarshalQuota(quota *Quota, data []byte) error {
	temp := struct {
		Type string `yaml:"type"`
	}{}

	err := yaml.Unmarshal(data, &temp)
	if err != nil {
		return fmt.Errorf("failed to parse quota type: %w", err)
	}

	switch strings.ToLower(temp.Type) {
	case "duration":
		q := struct {
			DurationQuota `yaml:",inline"`
			Type          string `yaml:"type"`
		}{}
		if err := yaml.UnmarshalWithOptions(data, &q, decoderOptions()...); err != nil {
			return fmt.Errorf("failed to parse duration quota: %w", err)
		}
		*quota = &q.DurationQuota
	case "counter":
		q := struct {
			CounterQuota `yaml:",inline"`
			Type         string `yaml:"type"`
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
