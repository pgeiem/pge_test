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

type LinearSequentialRule struct {
	BaseRule   `yaml:",inline"`
	Duration   time.Duration `yaml:"duration"`
	HourlyRate Amount        `yaml:"hourlyrate"`
}

func (r LinearSequentialRule) ToSolverRules(from, to time.Time, appender func(SolverRule)) {
	appender(NewLinearSequentialRule(r.RuleName, r.Duration, r.HourlyRate, r.Meta))
}

func (r LinearSequentialRule) String() string {
	return fmt.Sprintf("RelativeLinearRule %s", r.RuleName)
}

type FixedRateSequentialRule struct {
	BaseRule `yaml:",inline"`
	Duration time.Duration `yaml:"duration"`
	Amount   Amount        `yaml:"amount"`
	Repeat   int           `yaml:"repeat"`
}

func (r FixedRateSequentialRule) ToSolverRules(from, to time.Time, appender func(SolverRule)) {
	if r.Repeat < 1 {
		r.Repeat = 1
	}
	for i := 0; i < r.Repeat; i++ {
		solverRule := NewFixedRateSequentialRule(r.RuleName, r.Duration, r.Amount, r.Meta)
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Repetition no%d", i))
		appender(solverRule)
	}
}

func (r FixedRateSequentialRule) String() string {
	return fmt.Sprintf("RelativeFlatRateRule %s", r.RuleName)
}

type LinearFixedRule struct {
	BaseRule          `yaml:",inline"`
	RecurrentTimeSpan `yaml:",inline"`
	HourlyRate        Amount `yaml:"hourlyrate"`
}

// Unrolling the recurrent segment into a list of solver rules
func (r LinearFixedRule) ToSolverRules(from, to time.Time, appender func(SolverRule)) {
	cnt := 0
	r.RecurrentTimeSpan.BetweenIterator(from, to, func(timespan AbsTimeSpan) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		solverRule := NewLinearFixedRule(r.RuleName, ts, r.HourlyRate, r.Meta)
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Occurence no%d", cnt))
		appender(solverRule)
		cnt++
		return true
	})
}

func (r LinearFixedRule) String() string {
	return fmt.Sprintf("AbsoluteLinearRule %s", r.RuleName)
}

type FixedRateFixedRule struct {
	BaseRule          `yaml:",inline"`
	RecurrentTimeSpan `yaml:",inline"`
	Amount            Amount `yaml:"amount"`
}

func (r FixedRateFixedRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	cnt := 0
	r.RecurrentTimeSpan.BetweenIterator(from, to, func(timespan AbsTimeSpan) bool {
		fmt.Println("####### FixedRateFixedRule", timespan, cnt)
		ts := timespan.ToRelativeTimeSpan(from)
		solverRule := NewFixedRateFixedRule(r.RuleName, ts, r.Amount, r.Meta)
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Occurence no%d", cnt))
		iterator(solverRule)
		cnt++
		return true
	})
}

func (r FixedRateFixedRule) String() string {
	return fmt.Sprintf("AbsoluteFlatRateRule %s", r.RuleName)
}

type FlatRateFixedRule struct {
	BaseRule          `yaml:",inline"`
	RecurrentTimeSpan `yaml:",inline"`
	Amount            Amount `yaml:"amount"`
}

func (r FlatRateFixedRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	cnt := 0
	r.RecurrentTimeSpan.BetweenIterator(from, to, func(timespan AbsTimeSpan) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		solverRule := NewFlatRateFixedRule(r.RuleName, ts, r.Amount, r.Meta)
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Occurence no%d", cnt))
		iterator(solverRule)
		cnt++
		return true
	})
}

func (r FlatRateFixedRule) String() string {
	return fmt.Sprintf("FlatRateFixedRule %s", r.RuleName)
}

type NonPayingFixedRule struct {
	BaseRule          `yaml:",inline"`
	RecurrentTimeSpan `yaml:",inline"`
}

type AbsoluteNonPayingRules []NonPayingFixedRule

func (r NonPayingFixedRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	cnt := 0
	r.RecurrentTimeSpan.BetweenIterator(from, to, func(timespan AbsTimeSpan) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		solverRule := NewNonPayingFixedRule(r.RuleName, ts, r.Meta)
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Occurence no%d", cnt))
		iterator(solverRule)
		cnt++
		return true
	})
}

func (r NonPayingFixedRule) String() string {
	return fmt.Sprintf("AbsoluteNonPayingRule %s", r.RuleName)

}

func (rules *SolvableRules) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {

	temp := []struct {
		LinearSequentialRate    *LinearSequentialRule    `yaml:"linear"`
		FixedRateSequentialRule *FixedRateSequentialRule `yaml:"fixedrate"`
		LinearFixedRule         *LinearFixedRule         `yaml:"abslinear"`
		FlatRateFixedRule       *FlatRateFixedRule       `yaml:"absflatrate"`
		FixedRateFixedRule      *FixedRateFixedRule      `yaml:"absfixedrate"`
		NonPayingFixedRule      *NonPayingFixedRule      `yaml:"nonpaying"`
	}{}

	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	//fmt.Println("SolvableRules UnmarshalYAML", temp)

	*rules = make(SolvableRules, 0, len(temp))
	for _, t := range temp {
		quota := SolvableRule(nil)
		// TODO return an error if both DurationQuota and CounterQuota are set
		if t.LinearSequentialRate != nil {
			quota = t.LinearSequentialRate
		} else if t.FixedRateSequentialRule != nil {
			quota = t.FixedRateSequentialRule
		} else if t.LinearFixedRule != nil {
			quota = t.LinearFixedRule
		} else if t.FlatRateFixedRule != nil {
			quota = t.FlatRateFixedRule
		} else if t.FixedRateFixedRule != nil {
			quota = t.FixedRateFixedRule
		} else if t.NonPayingFixedRule != nil {
			quota = t.NonPayingFixedRule
		}
		if quota != nil {
			*rules = append(*rules, quota)
		}
	}

	//fmt.Println("SolvableRules UnmarshalYAML", rules)
	return nil
}
