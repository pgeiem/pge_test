package engine

import (
	"fmt"
	"time"
)

type TariffDefinition struct {
	Quotas    QuotaInventory
	NonPaying NonPayingInventory
	Sequences TariffSequenceInventory
	Config    TariffConfig
}

type TariffConfig struct {
	// Window of time to consider for the tarif computation
	Window time.Duration `yaml:"window"`
}

func DefaultConfig() TariffConfig {
	return TariffConfig{
		Window: time.Duration(48) * time.Hour,
	}
}

func (td TariffDefinition) Compute(now time.Time, history []AssignedRight) {

	now = now.Truncate(time.Second)

	// Update the quotas depending on the history
	td.Quotas.Update(now, history)

	// Solve all sequences
	td.Sequences.Solve(now, td.Config.Window, td.NonPaying)

	// Merge all sequences together
	rules, _ := td.Sequences.Merge(now, td.Config.Window) //TODO handle error if needed

	out := rules.GenerateOutput(now, true)
	json, _ := out.ToJson()
	fmt.Println(string(json))

}
