package parser

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	yaml "github.com/goccy/go-yaml"
)

type TariffDefinition struct {
	Quotas QuotaInventory `yaml:"quotas"`
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

func ParseTariffDefinition(r io.Reader) (TariffDefinition, error) {
	validate := validator.New()

	dec := yaml.NewDecoder(r,
		yaml.Strict(),
		yaml.CustomUnmarshaler(unmarshalDuration),
		yaml.CustomUnmarshaler(unmarshalTimeDuration),
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

/*func ParseTariffDefinitionAST(r io.Reader) (TariffDefinition, error) {
yaml

}*/
