package engine

import "strings"

type NonPayingRule struct {
	Name             string `yaml:"name"`
	RecurrentSegment `yaml:",inline"`
}

// Stringer for NonPayingRule display rule name and segment
func (npr NonPayingRule) String() string {
	return npr.Name + ": " + npr.RecurrentSegment.String()
}

type NonPayingInventory []NonPayingRule

// Stringer for NonPayingInventory display all rules as a dashed list
func (npi NonPayingInventory) String() string {
	var sb strings.Builder
	sb.WriteString("NonPaying:\n")
	for _, r := range npi {
		sb.WriteString(" - ")
		sb.WriteString(r.String())
		sb.WriteString("\n")
	}
	return sb.String()
}
