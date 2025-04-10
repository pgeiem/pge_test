package engine

import (
	"fmt"
	"time"
)

type TariffDefinition struct {
	Quotas    QuotaInventory
	NonPaying AbsoluteNonPayingRules
	Sequences TariffSequenceInventory
	Config    TariffConfig
}

type TariffConfig struct {
	// Window of time to consider for the tarif computation
	Window time.Duration `yaml:"window"`
	Limits TariffLimits  `yaml:",inline"`
}

func DefaultConfig() TariffConfig {
	return TariffConfig{
		Window: time.Duration(48) * time.Hour,
	}
}

func (td TariffDefinition) Compute(now time.Time, history []AssignedRight) Output {

	now = now.Local().Truncate(time.Second)
	fmt.Println("Now is", now)

	// Update the quotas depending on the history
	td.Quotas.Update(now, history)

	// Solve all sequences
	td.Sequences.Solve(now, td.Config.Window, td.NonPaying)

	// Merge all sequences together
	rules, _ := td.Sequences.Merge(now, td.Config.Window) //TODO handle error if needed

	rules = rules.ApplyLimits(td.Config.Limits)

	return rules.GenerateOutput(now, true)
}
