package engine

import "time"

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

	// Update the quotas depending on the history
	td.Quotas.Update(now, history)

}
