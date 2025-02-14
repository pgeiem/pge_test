package parser

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/iem-rd/quoteengine/engine"
)

type ParserTariffRoot struct {
	Version   string   `yaml:"version"`
	NonPaying ast.Node `yaml:"nonpaying"`
	Quotas    ast.Node `yaml:"quotas"`
	Sequences ast.Node `yaml:"sequences"`
}

func decoderOptions() []yaml.DecodeOption {
	//validate := validator.New()
	return []yaml.DecodeOption{
		yaml.Strict(),
		yaml.CustomUnmarshaler(unmarshalDuration),
		yaml.CustomUnmarshaler(unmarshalTimeDuration),
		yaml.CustomUnmarshaler(unmarshalRecurrentDate),
		yaml.CustomUnmarshaler(unmarshalQuota),
		//yaml.Validator(validate),
	}
}

// nodeToValueNodeToValueContext converts node to the value pointed to by v with context.Context.
func nodeToValueContext(ctx context.Context, node ast.Node, v interface{}, opts ...yaml.DecodeOption) error {
	var buf bytes.Buffer
	if err := yaml.NewDecoder(&buf, opts...).DecodeFromNodeContext(ctx, node, v); err != nil {
		return err
	}
	return nil
}

func ParseTariffDefinition(data []byte) (engine.TariffDefinition, error) {
	/* This parse the YAML file root level manually, section after section because
	we need to parse them in a specific order. For example the quotas need to be
	parsed before any section that references them. */

	var tariff engine.TariffDefinition

	// Parse YAML into a temporary root level only struct
	var desc ParserTariffRoot
	err := yaml.UnmarshalWithOptions(data, &desc, yaml.Strict())
	if err != nil {
		return tariff, err
	}

	// Check the version
	if desc.Version != "0.1" {
		return tariff, fmt.Errorf("invalid tariff version: %s", desc.Version)
	}

	ctx := context.Background()

	// Decode the nonpaying section
	err = nodeToValueContext(ctx, desc.NonPaying, &tariff.NonPaying, decoderOptions()...)
	if err != nil {
		return tariff, fmt.Errorf("failed to parse nonpaying section: %w", err)
	}

	// Decode the quotas section
	err = nodeToValueContext(ctx, desc.Quotas, &tariff.Quotas, decoderOptions()...)
	if err != nil {
		return tariff, fmt.Errorf("failed to parse quotas section: %w", err)
	}
	ctx = context.WithValue(ctx, "quotas", tariff.Quotas)

	// Decode the sequences section
	//tariff.Sequences, err = parseSequences(desc.Sequences, tariff.Quotas)
	err = nodeToValueContext(ctx, desc.Sequences, &tariff.Sequences, decoderOptions()...)
	if err != nil {
		return tariff, fmt.Errorf("failed to parse sequences section: %w", err)
	}

	return tariff, nil
}

func ParseTariffDefinitionFile(filename string) (engine.TariffDefinition, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return engine.TariffDefinition{}, err
	}
	return ParseTariffDefinition(data)
}
