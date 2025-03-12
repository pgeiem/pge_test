package engine

import (
	"fmt"
	"time"
)

// DurationType represents the different type of parking duration
type DurationType string

const (
	FreeDuration      DurationType = "free"
	NonPayingDuration DurationType = "nonpaying"
	PayingDuration    DurationType = "paying"
	BannedDuration    DurationType = "banned"
)

func (dt DurationType) MarshalText() ([]byte, error) {
	if dt == "" {
		return []byte{}, nil
	}
	conv := map[DurationType]string{
		FreeDuration:      "f",
		NonPayingDuration: "np",
		PayingDuration:    "p",
		BannedDuration:    "b",
	}
	if val, ok := conv[dt]; ok {
		return []byte(val), nil
	}
	return nil, fmt.Errorf("unknown duration type %s", dt)
}

// StartTimePolicy defines the policy used to move or not the beginning of the rule
type StartTimePolicy string // Todo replace by int32

const (
	ShiftablePolicy StartTimePolicy = "shiftable"
	FixedPolicy     StartTimePolicy = "fixed"
)

// RuleResolutionPolicy defines the policy used to solve the full rule duration
type RuleResolutionPolicy string // Todo replace by int32

const (
	TruncatePolicy RuleResolutionPolicy = "truncate"
	ResolvePolicy  RuleResolutionPolicy = "resolve"
	DeletePolicy   RuleResolutionPolicy = "delete"
)

//TOOD: merge RuleResolutionPolicy with StartTimePolicy as shiftable is usefull only with truncate ?

// Define the solver rule
// SolverRule represents a rule used in the solver engine.
type SolverRule struct {
	RuleName string
	// Starting/End point in time
	RelativeTimeSpan
	// Amount in cents at the beginning of the rule segment (non 0 values are step)
	StartAmount Amount
	// Amount in cents at the end of the rule segment
	EndAmount Amount
	// Trace buffer for debugging all rule changes
	Trace []string
	// StartTimePolicy defines the policy for determining the start time of the rule.
	StartTimePolicy StartTimePolicy
	// RuleResolutionPolicy defines the policy for resolving rule conflicts.
	RuleResolutionPolicy RuleResolutionPolicy
	// Meta holds additional metadata related to the rule.
	Meta MetaData
	// DurationType defines the type of duration for each rules, this is required to build duration details in the output
	DurationType DurationType
}

// Define a collection of solver rule
type SolverRules []SolverRule

type MetaData map[string]interface{}

func DurationTypeFromAmount(amount Amount) DurationType {
	if amount == 0 {
		return FreeDuration
	}
	return PayingDuration
}

func NewRelativeLinearRule(name string, duration time.Duration, hourlyRate Amount, meta MetaData) SolverRule {
	return SolverRule{
		RuleName:             name,
		Meta:                 meta,
		RelativeTimeSpan:     RelativeTimeSpan{From: time.Duration(0), To: duration},
		StartAmount:          0,
		EndAmount:            Amount(float64(hourlyRate) * duration.Hours()),
		StartTimePolicy:      ShiftablePolicy,
		RuleResolutionPolicy: ResolvePolicy,
		DurationType:         DurationTypeFromAmount(hourlyRate),
	}
}

func NewRelativeFlatRateRule(name string, duration time.Duration, amount Amount, meta MetaData) SolverRule {
	return SolverRule{
		RuleName:             name,
		Meta:                 meta,
		RelativeTimeSpan:     RelativeTimeSpan{From: time.Duration(0), To: duration},
		StartAmount:          amount,
		EndAmount:            amount,
		StartTimePolicy:      ShiftablePolicy,
		RuleResolutionPolicy: ResolvePolicy,
		DurationType:         DurationTypeFromAmount(amount),
	}
}

func NewAbsoluteLinearRule(name string, timespan RelativeTimeSpan, hourlyRate Amount, meta MetaData) SolverRule {
	if !timespan.IsValid() {
		panic(fmt.Errorf("invalid rule timespan %v", timespan))
	}
	return SolverRule{
		RuleName:             name,
		Meta:                 meta,
		RelativeTimeSpan:     timespan,
		StartAmount:          0,
		EndAmount:            Amount(float64(hourlyRate) * timespan.Duration().Hours()),
		StartTimePolicy:      FixedPolicy,
		RuleResolutionPolicy: ResolvePolicy,
		DurationType:         DurationTypeFromAmount(hourlyRate),
	}
}

func NewAbsoluteFlatRateRule(name string, timespan RelativeTimeSpan, amount Amount, meta MetaData) SolverRule {
	if !timespan.IsValid() {
		panic(fmt.Errorf("invalid rule timespan %v", timespan))
	}
	return SolverRule{
		RuleName:             name,
		Meta:                 meta,
		RelativeTimeSpan:     timespan,
		StartAmount:          amount,
		EndAmount:            amount,
		StartTimePolicy:      FixedPolicy,
		RuleResolutionPolicy: TruncatePolicy,
		DurationType:         DurationTypeFromAmount(amount),
	}
}

func NewAbsoluteNonPaying(name string, timespan RelativeTimeSpan, meta MetaData) SolverRule {
	if !timespan.IsValid() {
		panic(fmt.Errorf("invalid rule timespan %v", timespan))
	}
	return SolverRule{
		RuleName:             name,
		Meta:                 meta,
		RelativeTimeSpan:     timespan,
		StartAmount:          0,
		EndAmount:            0,
		StartTimePolicy:      FixedPolicy,
		RuleResolutionPolicy: TruncatePolicy,
		DurationType:         NonPayingDuration,
	}
}

func (rule SolverRule) Duration() time.Duration {
	return rule.To - rule.From
}

func (rule SolverRule) IsFlatRate() bool {
	return rule.StartAmount == rule.EndAmount
}

func (rule SolverRule) IsAbsoluteFlatRate() bool {
	return rule.IsFlatRate() && //is FlatRate
		rule.StartTimePolicy == FixedPolicy && // is Absolute
		rule.EndAmount != 0 // is not non-paying
}

func (rule SolverRule) IsRelative() bool {
	return rule.StartTimePolicy == ShiftablePolicy
}

func (rule SolverRule) Name() string {
	return rule.RuleName
}

func (rule SolverRule) String() string {
	return fmt.Sprintf("%s(%s -> %s; %s -> %s)",
		rule.Name(), rule.From.String(), rule.To.String(), rule.StartAmount, rule.EndAmount)
}
