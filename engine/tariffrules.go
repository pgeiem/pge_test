package engine

import (
	"context"
	"fmt"
	"time"
)

type BaseRule struct {
	RuleName string `yaml:"name"`
	Meta     MetaData
}

type SolvableRule interface {
	ToSolverRules(from, to time.Time, iterator func(SolverRule))
	String() string
}

type SolvableRules []SolvableRule

type RelativeLinearRule struct {
	BaseRule   `yaml:",inline"`
	Duration   time.Duration `yaml:"duration"`
	HourlyRate Amount        `yaml:"hourlyrate"`
}

func (r RelativeLinearRule) ToSolverRules(from, to time.Time, appender func(SolverRule)) {
	appender(NewRelativeLinearRule(r.RuleName, r.Duration, r.HourlyRate, r.Meta))
}

func (r RelativeLinearRule) String() string {
	return fmt.Sprintf("RelativeLinearRule %s", r.RuleName)
}

type RelativeFlatRateRule struct {
	BaseRule `yaml:",inline"`
	Duration time.Duration `yaml:"duration"`
	Amount   Amount        `yaml:"amount"`
}

func (r RelativeFlatRateRule) ToSolverRules(from, to time.Time, appender func(SolverRule)) {
	appender(NewRelativeFlatRateRule(r.RuleName, r.Duration, r.Amount, r.Meta))
}

func (r RelativeFlatRateRule) String() string {
	return fmt.Sprintf("RelativeFlatRateRule %s", r.RuleName)
}

type AbsoluteLinearRule struct {
	BaseRule          `yaml:",inline"`
	RecurrentTimeSpan `yaml:",inline"`
	HourlyRate        Amount `yaml:"hourlyrate"`
}

// Unrolling the recurrent segment into a list of solver rules
func (r AbsoluteLinearRule) ToSolverRules(from, to time.Time, appender func(SolverRule)) {
	cnt := 0
	r.RecurrentTimeSpan.BetweenIterator(from, to, func(timespan AbsTimeSpan) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		solverRule := NewAbsoluteLinearRule(r.RuleName, ts, r.HourlyRate, r.Meta)
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Occurence no%d", cnt))
		appender(solverRule)
		cnt++
		return true
	})
}

func (r AbsoluteLinearRule) String() string {
	return fmt.Sprintf("AbsoluteLinearRule %s", r.RuleName)
}

type AbsoluteFlatRateRule struct {
	BaseRule          `yaml:",inline"`
	RecurrentTimeSpan `yaml:",inline"`
	Amount            Amount `yaml:"amount"`
}

func (r AbsoluteFlatRateRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	cnt := 0
	r.RecurrentTimeSpan.BetweenIterator(from, to, func(timespan AbsTimeSpan) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		solverRule := NewAbsoluteFlatRateRule(r.RuleName, ts, r.Amount, r.Meta)
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Occurence no%d", cnt))
		iterator(solverRule)
		cnt++
		return true
	})
}

func (r AbsoluteFlatRateRule) String() string {
	return fmt.Sprintf("AbsoluteFlatRateRule %s", r.RuleName)
}

type AbsoluteNonPayingRule struct {
	BaseRule          `yaml:",inline"`
	RecurrentTimeSpan `yaml:",inline"`
}

func (r AbsoluteNonPayingRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	cnt := 0
	r.RecurrentTimeSpan.BetweenIterator(from, to, func(timespan AbsTimeSpan) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		solverRule := NewAbsoluteNonPaying(r.RuleName, ts, r.Meta)
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Occurence no%d", cnt))
		iterator(solverRule)
		cnt++
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
