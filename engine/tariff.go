package engine

type TariffDefinition struct {
	Quotas    QuotaInventory          `yaml:"quotas"`
	NonPaying NonPayingInventory      `yaml:"nonpaying"`
	Sequences TariffSequenceInventory `yaml:"sequences"`
}
