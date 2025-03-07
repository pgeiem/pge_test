package engine

import (
	"context"
	"fmt"
	"time"
)

type RealtiveLinearPayingRule struct {
	RuleName   string        `yaml:"name"`
	Duration   time.Duration `yaml:"duration"`
	HourlyRate Amount        `yaml:"hourlyrate"`
}

func (r RealtiveLinearPayingRule) Name() string {
	return r.RuleName
}

/*func (r RealtiveLinearPayingRule) RelativeTo(now time.Time) (time.Duration, time.Duration) {
	return time.Duration(0), r.Duration
}*/

func (r RealtiveLinearPayingRule) RelativeToWindow(from, to time.Time, iterator func(RelativeTimeSpan) bool) {
	iterator(RelativeTimeSpan{From: 0, To: r.Duration})
}

func (r RealtiveLinearPayingRule) Policies() (StartTimePolicy, RuleResolutionPolicy) {
	return ShiftablePolicy, ResolvePolicy
}

func (r RealtiveLinearPayingRule) String() string {
	return fmt.Sprintf("Linear %s: hourlyrate %s, duration %s", r.Name(), r.HourlyRate.ToString(), r.Duration)
}

type RelativeFlatRatePayingRule struct {
	RuleName string        `yaml:"name"`
	Duration time.Duration `yaml:"duration"`
	Amount   Amount        `yaml:"amount"`
}

func (r RelativeFlatRatePayingRule) Name() string {
	return r.RuleName
}

/*func (r RelativeFlatRatePayingRule) RelativeTo(now time.Time) (time.Duration, time.Duration) {
	return time.Duration(0), r.Duration
}*/

func (r RelativeFlatRatePayingRule) RelativeToWindow(from, to time.Time, iterator func(RelativeTimeSpan) bool) {
	iterator(RelativeTimeSpan{From: 0, To: r.Duration})
}

func (r RelativeFlatRatePayingRule) Policies() (StartTimePolicy, RuleResolutionPolicy) {
	return ShiftablePolicy, TruncatePolicy
}

func (r RelativeFlatRatePayingRule) String() string {
	return fmt.Sprintf("FlatRate %s: amount %s, duration %s", r.Name(), r.Amount.ToString(), r.Duration)
}

type RelativeTariffRulesSequence []SolvableRule

func (sequ *RelativeTariffRulesSequence) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {

	temp := []struct {
		LinearRule   *RealtiveLinearPayingRule   `yaml:"linear"`
		FlatRateRule *RelativeFlatRatePayingRule `yaml:"flatrate"`
	}{}

	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	*sequ = make(RelativeTariffRulesSequence, 0, len(temp))
	for _, t := range temp {
		rule := SolvableRule(nil)
		// TODO return an error if both LinearRule and FlatRateRule or none are set
		if t.LinearRule != nil {
			rule = t.LinearRule
		} else if t.FlatRateRule != nil {
			rule = t.FlatRateRule
		}
		if rule != nil {
			*sequ = append(*sequ, rule)
		}
	}
	return nil
}
