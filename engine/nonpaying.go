package engine

import (
	"strings"
)

type NonPayingInventory []AbsoluteNonPayingRule

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
