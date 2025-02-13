package engine

import "strings"

type TariffSequence struct {
	Name           string           `yaml:"name"`
	ValidityPeriod RecurrentSegment `yaml:",inline"`
	Quota          Quota            `yaml:"quota,"`
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
