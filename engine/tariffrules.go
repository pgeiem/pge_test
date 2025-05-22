package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/iem-rd/quote-engine/timeutils"
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
	BaseRule
	Quota      Quota
	Duration   time.Duration
	HourlyRate Amount
}

func (r LinearSequentialRule) ToSolverRules(from, to time.Time, appender func(SolverRule)) {
	solverRule := NewLinearSequentialRule(r.RuleName, r.Duration, r.HourlyRate, r.Meta)
	solverRule.Quota = r.Quota
	appender(solverRule)
}

func (r LinearSequentialRule) String() string {
	return fmt.Sprintf("RelativeLinearRule %s", r.RuleName)
}

func (r *LinearSequentialRule) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {
	var ok bool
	temp := struct {
		BaseRule   `yaml:",inline"`
		QuotaName  string        `yaml:"quota"`
		Duration   time.Duration `yaml:"duration"`
		HourlyRate Amount        `yaml:"hourlyrate"`
	}{}

	// Unmarshal the base rule
	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	// Set the fields of the LinearSequentialRule
	r.BaseRule = temp.BaseRule
	r.Quota, ok = ContextGetQuotaByName(ctx, temp.QuotaName)
	if !ok {
		return fmt.Errorf("unknown quota: %s", temp.QuotaName)
	}
	r.Duration = temp.Duration
	r.HourlyRate = temp.HourlyRate
	return nil
}

type FixedRateSequentialRule struct {
	BaseRule
	Quota    Quota
	Duration time.Duration
	Amount   Amount
	Repeat   int
}

func (r FixedRateSequentialRule) ToSolverRules(from, to time.Time, appender func(SolverRule)) {
	if r.Repeat < 1 {
		r.Repeat = 1
	}
	for i := 0; i < r.Repeat; i++ {
		solverRule := NewFixedRateSequentialRule(r.RuleName, r.Duration, r.Amount, r.Meta)
		solverRule.Quota = r.Quota
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Repetition no%d", i))
		appender(solverRule)
	}
}

func (r FixedRateSequentialRule) String() string {
	return fmt.Sprintf("RelativeFlatRateRule %s", r.RuleName)
}

func (r *FixedRateSequentialRule) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {
	var ok bool
	temp := struct {
		BaseRule  `yaml:",inline"`
		QuotaName string        `yaml:"quota"`
		Duration  time.Duration `yaml:"duration"`
		Amount    Amount        `yaml:"amount"`
		Repeat    int           `yaml:"repeat"`
	}{}

	// Unmarshal the base rule
	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	// Set the fields of the FixedRateSequentialRule
	r.BaseRule = temp.BaseRule
	r.Quota, ok = ContextGetQuotaByName(ctx, temp.QuotaName)
	if !ok {
		return fmt.Errorf("unknown quota: %s", temp.QuotaName)
	}
	r.Duration = temp.Duration
	r.Amount = temp.Amount
	r.Repeat = temp.Repeat
	return nil
}

type LinearFixedRule struct {
	BaseRule
	timeutils.RecurrentTimeSpan
	Quota      Quota
	HourlyRate Amount
}

// Unrolling the recurrent segment into a list of solver rules
func (r LinearFixedRule) ToSolverRules(from, to time.Time, appender func(SolverRule)) {
	cnt := 0
	r.RecurrentTimeSpan.BetweenIterator(from, to, func(timespan timeutils.AbsTimeSpan) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		solverRule := NewLinearFixedRule(r.RuleName, ts, r.HourlyRate, r.Meta)
		solverRule.Quota = r.Quota
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Occurence no%d", cnt))
		appender(solverRule)
		cnt++
		return true
	})
}

func (r LinearFixedRule) String() string {
	return fmt.Sprintf("AbsoluteLinearRule %s", r.RuleName)
}

func (r *LinearFixedRule) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {
	var ok bool
	temp := struct {
		BaseRule                    `yaml:",inline"`
		timeutils.RecurrentTimeSpan `yaml:",inline"`
		QuotaName                   string `yaml:"quota"`
		HourlyRate                  Amount `yaml:"hourlyrate"`
	}{}

	// Unmarshal the base rule
	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	// Set the fields of the LinearFixedRule
	r.BaseRule = temp.BaseRule
	r.RecurrentTimeSpan = temp.RecurrentTimeSpan
	r.Quota, ok = ContextGetQuotaByName(ctx, temp.QuotaName)
	if !ok {
		return fmt.Errorf("unknown quota: %s", temp.QuotaName)
	}
	r.HourlyRate = temp.HourlyRate
	return nil
}

type FixedRateFixedRule struct {
	BaseRule
	timeutils.RecurrentTimeSpan
	Quota  Quota
	Amount Amount
}

func (r FixedRateFixedRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	cnt := 0
	r.RecurrentTimeSpan.BetweenIterator(from, to, func(timespan timeutils.AbsTimeSpan) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		solverRule := NewFixedRateFixedRule(r.RuleName, ts, r.Amount, r.Meta)
		solverRule.Quota = r.Quota
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Occurence no%d", cnt))
		iterator(solverRule)
		cnt++
		return true
	})
}

func (r FixedRateFixedRule) String() string {
	return fmt.Sprintf("AbsoluteFlatRateRule %s", r.RuleName)
}

func (r *FixedRateFixedRule) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {
	var ok bool
	temp := struct {
		BaseRule                    `yaml:",inline"`
		timeutils.RecurrentTimeSpan `yaml:",inline"`
		QuotaName                   string `yaml:"quota"`
		Amount                      Amount `yaml:"amount"`
	}{}

	// Unmarshal the base rule
	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	// Set the fields of the FixedRateFixedRule
	r.BaseRule = temp.BaseRule
	r.RecurrentTimeSpan = temp.RecurrentTimeSpan
	r.Quota, ok = ContextGetQuotaByName(ctx, temp.QuotaName)
	if !ok {
		return fmt.Errorf("unknown quota: %s", temp.QuotaName)
	}
	r.Amount = temp.Amount
	return nil
}

type FlatRateFixedRule struct {
	BaseRule
	timeutils.RecurrentTimeSpan
	Quota  Quota
	Amount Amount
}

func (r FlatRateFixedRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	cnt := 0
	r.RecurrentTimeSpan.BetweenIterator(from, to, func(timespan timeutils.AbsTimeSpan) bool {
		ts := timespan.ToRelativeTimeSpan(from)
		solverRule := NewFlatRateFixedRule(r.RuleName, ts, r.Amount, r.Meta)
		solverRule.Quota = r.Quota
		solverRule.Trace = append(solverRule.Trace, fmt.Sprintf("Occurence no%d", cnt))
		iterator(solverRule)
		cnt++
		return true
	})
}

func (r FlatRateFixedRule) String() string {
	return fmt.Sprintf("FlatRateFixedRule %s", r.RuleName)
}

func (r *FlatRateFixedRule) UnmarshalYAML(ctx context.Context, unmarshal func(interface{}) error) error {
	var ok bool
	temp := struct {
		BaseRule                    `yaml:",inline"`
		timeutils.RecurrentTimeSpan `yaml:",inline"`
		QuotaName                   string `yaml:"quota"`
		Amount                      Amount `yaml:"amount"`
	}{}

	// Unmarshal the base rule
	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	// Set the fields of the FlatRateFixedRule
	r.BaseRule = temp.BaseRule
	r.RecurrentTimeSpan = temp.RecurrentTimeSpan
	r.Quota, ok = ContextGetQuotaByName(ctx, temp.QuotaName)
	if !ok {
		return fmt.Errorf("unknown quota: %s", temp.QuotaName)
	}
	r.Amount = temp.Amount
	return nil
}

type NonPayingFixedRule struct {
	BaseRule                    `yaml:",inline"`
	timeutils.RecurrentTimeSpan `yaml:",inline"`
}

type AbsoluteNonPayingRules []NonPayingFixedRule

func (r NonPayingFixedRule) ToSolverRules(from, to time.Time, iterator func(SolverRule)) {
	cnt := 0
	r.RecurrentTimeSpan.BetweenIterator(from, to, func(timespan timeutils.AbsTimeSpan) bool {
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
		rule := SolvableRule(nil)
		// TODO return an error if more then one field is set
		if t.LinearSequentialRate != nil {
			rule = t.LinearSequentialRate
		} else if t.FixedRateSequentialRule != nil {
			rule = t.FixedRateSequentialRule
		} else if t.LinearFixedRule != nil {
			rule = t.LinearFixedRule
		} else if t.FlatRateFixedRule != nil {
			rule = t.FlatRateFixedRule
		} else if t.FixedRateFixedRule != nil {
			rule = t.FixedRateFixedRule
		} else if t.NonPayingFixedRule != nil {
			rule = t.NonPayingFixedRule
		}
		if rule != nil {
			*rules = append(*rules, rule)
		}
	}

	//fmt.Println("SolvableRules UnmarshalYAML", rules)
	return nil
}
