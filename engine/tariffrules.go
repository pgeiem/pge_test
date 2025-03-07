package engine

// DurationType represents the different type of parking duration
type DurationType string

const (
	FreeDuration      DurationType = "free"
	NonPayingDuration DurationType = "nonpaying"
	PayingDuration    DurationType = "paying"
	BannedDuration    DurationType = "banned"
)

/*
type TariffRuleInventory []SolvableRule

type RealtiveLinearPayingRule struct {
	RuleName 	 string `yaml:"name"`
	Duration time.Duration `yaml:"duration"`
	HourlyRate Amount `yaml:"hourlyrate"`
}

func (r RealtiveLinearPayingRule) Name() string {
	return r.RuleName
}

func (r RealtiveLinearPayingRule) RelativeTo(now time.Time) (time.Duration, time.Duration) {
	return now, r.Duration + now
}

func (r RealtiveLinearPayingRule) Policies() (StartTimePolicy, RuleResolutionPolicy) {
	return StartTimePolicy, ResolvePolicy
}

type RelativeFlatRatePayingRule struct {
	RecurrentSegment
	Name   string `yaml:"name"`
	Duration time.Duration `yaml:"duration"`
	Amount Amount `yaml:"amount"`
}

func (r RealtiveLinearPayingRule) Name() string {
	return r.RuleName
}

func (r RealtiveLinearPayingRule) RelativeTo(now time.Time) (time.Duration, time.Duration) {
	return now, r.Duration + now
}

func (r RealtiveLinearPayingRule) Policies() (StartTimePolicy, RuleResolutionPolicy) {
	return StartTimePolicy, TruncatePolicy
}
/*
type FreeRule struct {
	RecurrentSegment
	Name string `yaml:"name"`
}

type BannedRule struct {
	RecurrentSegment
	Name string `yaml:"name"`
}
*/

/*
func (r *TariffRule) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {

	temp := []struct {
		DurationQuota *DurationQuota `yaml:"duration"`
		CounterQuota  *CounterQuota  `yaml:"counter"`
	}{}

	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	*qi = make(QuotaInventory)
	for _, t := range temp {
		quota := Quota(nil)
		// TODO return an error if both DurationQuota and CounterQuota are set
		if t.DurationQuota != nil {
			quota = t.DurationQuota
		} else if t.CounterQuota != nil {
			quota = t.CounterQuota
		}
		if quota != nil {
			if quota.GetName() == "" {
				return fmt.Errorf("missing quota name")
			}
			(*qi)[quota.GetName()] = quota
		}
	}
	return nil
}
*/
