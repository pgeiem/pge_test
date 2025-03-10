package engine

import (
	"context"
	"fmt"
	"time"
)

type SolvableRule interface {
	ToSolverRules(from, to time.Time, iterator func(SolverRule))
	String() string
}

type SolvableRules []SolvableRule

type RelativeLinearRule struct {
	RuleName   string        `yaml:"name"`
	Duration   time.Duration `yaml:"duration"`
	HourlyRate Amount        `yaml:"hourlyrate"`
	MetaData   MetaData      `yaml:"meta"`
}

func (r RelativeLinearRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	iterator(NewRelativeLinearRule(r.RuleName, r.Duration, r.HourlyRate))
}

func (r RelativeLinearRule) String() string {
	return fmt.Sprintf("RelativeLinearRule %s", r.RuleName)
}

type RelativeFlatRateRule struct {
	RuleName string        `yaml:"name"`
	Duration time.Duration `yaml:"duration"`
	Amount   Amount        `yaml:"amount"`
	MetaData MetaData      `yaml:"meta"`
}

func (r RelativeFlatRateRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	iterator(NewRelativeFlatRateRule(r.RuleName, r.Duration, r.Amount))
}

func (r RelativeFlatRateRule) String() string {
	return fmt.Sprintf("RelativeFlatRateRule %s", r.RuleName)
}

type AbsoluteLinearRule struct {
	RuleName         string `yaml:"name"`
	RecurrentSegment `yaml:",inline"`
	HourlyRate       Amount   `yaml:"hourlyrate"`
	MetaData         MetaData `yaml:"meta"`
}

// Unrolling the recurrent segment into a list of solver rules
func (r AbsoluteLinearRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	r.RecurrentSegment.BetweenIterator(from, to, func(timespan Segment) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		iterator(NewAbsoluteLinearRule(r.RuleName, ts, r.HourlyRate))
		return true
	})
}

func (r AbsoluteLinearRule) String() string {
	return fmt.Sprintf("AbsoluteLinearRule %s", r.RuleName)
}

type AbsoluteFlatRateRule struct {
	RuleName         string `yaml:"name"`
	RecurrentSegment `yaml:",inline"`
	Amount           Amount   `yaml:"amount"`
	MetaData         MetaData `yaml:"meta"`
}

func (r AbsoluteFlatRateRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	r.RecurrentSegment.BetweenIterator(from, to, func(timespan Segment) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		iterator(NewAbsoluteFlatRateRule(r.RuleName, ts, r.Amount))
		return true
	})
}

func (r AbsoluteFlatRateRule) String() string {
	return fmt.Sprintf("AbsoluteFlatRateRule %s", r.RuleName)
}

type AbsoluteNonPayingRule struct {
	RuleName         string `yaml:"name"`
	RecurrentSegment `yaml:",inline"`
	MetaData         MetaData `yaml:"meta"`
}

func (r AbsoluteNonPayingRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	r.RecurrentSegment.BetweenIterator(from, to, func(timespan Segment) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		iterator(NewAbsoluteNonPaying(r.RuleName, ts))
		return true
	})
}

func (r AbsoluteNonPayingRule) String() string {
	return fmt.Sprintf("AbsoluteNonPayingRule %s", r.RuleName)

}

func (rules *SolvableRules) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {

	temp := []struct {
		RelativeLinearRule *RelativeLinearRule    `yaml:"linear"`
		RelativeFlatRate   *RelativeFlatRateRule  `yaml:"flatrate"`
		AbsoluteLinearRule *AbsoluteLinearRule    `yaml:"abslinear"`
		AbsoluteFlatRate   *AbsoluteFlatRateRule  `yaml:"absflatrate"`
		AbsoluteNonPaying  *AbsoluteNonPayingRule `yaml:"nonpaying"`
	}{}

	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	fmt.Println("SolvableRules UnmarshalYAML", temp)

	*rules = make(SolvableRules, 0, len(temp))
	for _, t := range temp {
		quota := SolvableRule(nil)
		// TODO return an error if both DurationQuota and CounterQuota are set
		if t.RelativeLinearRule != nil {
			quota = t.RelativeLinearRule
		} else if t.RelativeFlatRate != nil {
			quota = t.RelativeFlatRate
		} else if t.AbsoluteLinearRule != nil {
			quota = t.AbsoluteLinearRule
		} else if t.AbsoluteFlatRate != nil {
			quota = t.AbsoluteFlatRate
		} else if t.AbsoluteNonPaying != nil {
			quota = t.AbsoluteNonPaying
		}
		if quota != nil {
			*rules = append(*rules, quota)
		}
	}

	fmt.Println("SolvableRules UnmarshalYAML", rules)
	return nil
}
