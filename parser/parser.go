package parser

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/iem-rd/quoteengine/engine"
)

type TariffDescription struct {
	Version   string   `yaml:"version"`
	NonPaying ast.Node `yaml:"nonpaying"`
	Quotas    ast.Node `yaml:"quotas"`
	Sequences ast.Node `yaml:"sequences"`
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

func ParseTariffDefinition(r io.Reader) (engine.TariffDefinition, error) {
	var tariff engine.TariffDefinition

	//validate := validator.New()

	dec := yaml.NewDecoder(r,
		yaml.Strict(),
		//yaml.Validator(validate),
	)

	// Parse YAML into a temporary root level only struct
	var desc TariffDescription
	err := dec.Decode(&desc)
	if err != nil {
		return tariff, err
	}

	// Check the version
	if desc.Version != "0.1" {
		return tariff, fmt.Errorf("invalid tariff version: %s", desc.Version)
	}

	// Decode the nonpaying section
	err = yaml.NodeToValue(desc.NonPaying, &tariff.NonPaying, decoderOptions()...)
	if err != nil {
		return tariff, fmt.Errorf("failed to parse nonpaying section: %w", err)
	}

	// Decode the quotas section
	err = yaml.NodeToValue(desc.Quotas, &tariff.Quotas, decoderOptions()...)
	if err != nil {
		return tariff, fmt.Errorf("failed to parse quotas section: %w", err)
	}

	// Decode the sequences section
	tariff.Sequences, err = parseSequences(desc.Sequences, tariff.Quotas)
	if err != nil {
		return tariff, fmt.Errorf("failed to parse sequences section: %w", err)
	}

	//yaml.NodeToValue(desc.Sequences, &tariff.Sequences, decoderOptions()...)

	return tariff, nil
}

func ParseTariffDefinitionString(s string) (engine.TariffDefinition, error) {
	return ParseTariffDefinition(strings.NewReader(s))
}

func ParseTariffDefinitionFile(filename string) (engine.TariffDefinition, error) {
	f, err := os.Open(filename)
	if err != nil {
		return engine.TariffDefinition{}, err
	}
	defer f.Close()
	return ParseTariffDefinition(f)
}
