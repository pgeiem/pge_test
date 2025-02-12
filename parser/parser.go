package parser

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	yaml "github.com/goccy/go-yaml"
)

type TariffDefinition struct {
	Quotas QuotaInventory `yaml:"quotas"`
	//NonPaying NonPayingInventory `yaml:"nonpaying"`
}

func unmarshalTimeDuration(duration *time.Duration, data []byte) error {
	d, err := ParseDuration(string(data))
	if err != nil {
		return err
	}
	*duration = d.toDuration()
	return nil
}

func unmarshalDuration(duration *Duration, data []byte) error {
	d, err := ParseDuration(string(data))
	if err != nil {
		return err
	}
	*duration = d
	return nil
}

func unmarshalRecurrentDate(rec *RecurrentDate, data []byte) error {
	tmp, err := ParseRecurrentDate(string(data))
	if err != nil {
		return err
	}
	*rec = tmp
	return nil
}

func unmarshalQuota(quota *Quota, data []byte) error {
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

func decoderOptions() []yaml.DecodeOption {
	return []yaml.DecodeOption{
		yaml.Strict(),
		yaml.CustomUnmarshaler(unmarshalDuration),
		yaml.CustomUnmarshaler(unmarshalTimeDuration),
		yaml.CustomUnmarshaler(unmarshalRecurrentDate),
		yaml.CustomUnmarshaler(unmarshalQuota),
	}
}

func ParseTariffDefinition(r io.Reader) (TariffDefinition, error) {
	validate := validator.New()

	dec := yaml.NewDecoder(r,
		yaml.Strict(),
		yaml.CustomUnmarshaler(unmarshalDuration),
		yaml.CustomUnmarshaler(unmarshalTimeDuration),
		yaml.CustomUnmarshaler(unmarshalRecurrentDate),
		yaml.CustomUnmarshaler(unmarshalQuota),
		yaml.Validator(validate),
	)

	var tariff TariffDefinition
	if err := dec.Decode(&tariff); err != nil {
		return tariff, err
	}
	return tariff, nil
}

func ParseTariffDefinitionString(s string) (TariffDefinition, error) {
	return ParseTariffDefinition(strings.NewReader(s))
}

func ParseTariffDefinitionFile(filename string) (TariffDefinition, error) {
	f, err := os.Open(filename)
	if err != nil {
		return TariffDefinition{}, err
	}
	defer f.Close()
	return ParseTariffDefinition(f)
}
