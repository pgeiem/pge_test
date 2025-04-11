package engine

import (
	"fmt"
	"math"
	"time"

	"github.com/iem-rd/quote-engine/table"
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
	// Amount for which this flatrate rule is active
	ActivationAmount Amount
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

func NewLinearSequentialRule(name string, duration time.Duration, hourlyRate Amount, meta MetaData) SolverRule {
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

func NewFixedRateSequentialRule(name string, duration time.Duration, amount Amount, meta MetaData) SolverRule {
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

func NewLinearFixedRule(name string, timespan RelativeTimeSpan, hourlyRate Amount, meta MetaData) SolverRule {
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
		RuleResolutionPolicy: TruncatePolicy,
		DurationType:         DurationTypeFromAmount(hourlyRate),
	}
}
func NewFixedRateFixedRule(name string, timespan RelativeTimeSpan, amount Amount, meta MetaData) SolverRule {
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

func NewFlatRateFixedRule(name string, timespan RelativeTimeSpan, amount Amount, meta MetaData) SolverRule {
	if !timespan.IsValid() {
		panic(fmt.Errorf("invalid rule timespan %v", timespan))
	}
	return SolverRule{
		RuleName:             name,
		Meta:                 meta,
		RelativeTimeSpan:     timespan,
		StartAmount:          0,
		EndAmount:            0,
		ActivationAmount:     amount,
		StartTimePolicy:      FixedPolicy,
		RuleResolutionPolicy: TruncatePolicy,
		DurationType:         DurationTypeFromAmount(amount),
	}
}

func NewNonPayingFixedRule(name string, timespan RelativeTimeSpan, meta MetaData) SolverRule {
	r := NewFlatRateFixedRule(name, timespan, 0, meta)
	r.DurationType = NonPayingDuration
	return r
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

type TariffLimits struct {
	// MaxAmount is the maximum amount allowed for the rules
	MaxAmount Amount `yaml:"maxamount"`
	// MaxDuration is the maximum duration allowed for the rules
	MaxDuration time.Duration `yaml:"maxduration"`
}

func (limits TariffLimits) String() string {
	return fmt.Sprintf("MaxAmount %f, MaxDuration %s", limits.MaxAmount, limits.MaxDuration)
}

func (limits *TariffLimits) AddOffset(offsetAmout Amount, offsetDuration time.Duration) {
	if limits.MaxAmount > 0 {
		limits.MaxAmount += offsetAmout
	}
	if limits.MaxDuration > 0 {
		limits.MaxDuration += offsetDuration
	}
}

func (rules SolverRules) ApplyLimits(limits TariffLimits) SolverRules {
	if limits.MaxAmount == 0 && limits.MaxDuration == 0 {
		return rules
	}

	fmt.Println(" >> Applying limits to", len(rules), "rules", limits)
	sumAmount := Amount(0)
	overflow := false
	out := SolverRules{}
	for _, rule := range rules {

		fmt.Println("   >> Rule", rule, sumAmount)

		// check max duration limit
		if rule.DurationType != NonPayingDuration && limits.MaxDuration > 0 {
			if rule.From > limits.MaxDuration {
				fmt.Println("   >> maxDuration reached, rule skipped", rule)
				overflow = true
			} else if rule.To > limits.MaxDuration {
				rule = rule.TruncateAfter(limits.MaxDuration)
				fmt.Println("   >> maxDuration reached, rule truncated", rule)
				overflow = true
			}
		}

		// check max amount limit
		if limits.MaxAmount > 0 {
			if sumAmount+rule.EndAmount > limits.MaxAmount {
				rule = rule.TruncateAfterAmount(limits.MaxAmount - sumAmount)
				fmt.Println("   >> maxAmount reached, rule truncated", rule)
				overflow = true
			}
		}

		out = append(out, rule)
		sumAmount += rule.EndAmount

		if overflow {
			break
		}
	}
	return out
}

func (rules *SolverRules) GenerateOutput(now time.Time, detailed bool) Output {
	var out Output
	var previous SolverRule

	out.Now = now

	fmt.Println("Generating output for", len(*rules), "rules")
	for _, rule := range *rules {
		fmt.Println("   Rule", rule)
		// If there is a gap between the previous rule and the current one this is the end of the output
		if previous.To != rule.From {
			fmt.Println("   >> Gap detected, end of output", previous, rule)
			break
		}
		seg := OutputSegment{
			Duration:     int(math.Round(rule.To.Seconds() - previous.To.Seconds())),
			Amount:       rule.EndAmount.Simplify(),
			Islinear:     !rule.IsFlatRate(),
			DurationType: rule.DurationType,
			Meta:         rule.Meta,
		}
		if detailed {
			seg.SegName = rule.Name()
			seg.Trace = rule.Trace
		}
		out.Table = append(out.Table, seg)
		previous = rule
	}
	return out
}

// sumAll returns the sum of all rules amounts and the total duration
func (rules SolverRules) SumAll() (Amount, time.Duration) {
	var amountSum Amount
	var durationSum time.Duration

	// Loop over all rules and determine the end amount and end duration
	for i := range rules {
		rule := rules[i]
		amountSum += rule.EndAmount
		durationSum = rule.To // last rule duration is the total duration
	}

	return amountSum, durationSum
}

// PrintAsTable prints the rules as a table using table view
func (rules SolverRules) PrintAsTable(title string, now time.Time) {

	tbl := StartRulesTable(title, now)
	for _, rule := range rules {
		tbl.AddRule(&rule)
	}
	tbl.Print()
}

type RulesTable struct {
	now   time.Time
	tbl   table.Table
	title string
	empty bool
}

func StartRulesTable(title string, now time.Time) *RulesTable {
	tbl := table.New("Name", "From", "To", "Duration", "From (abs)", "To (abs)", "StartAmount", "EndAmount", "IsLinear", "ActivAm", "Type")
	t := RulesTable{
		now:   now,
		title: title,
		tbl:   tbl,
		empty: true,
	}
	table.SetDefaultTheme(&t.tbl)
	return &t
}

func (t *RulesTable) AddRule(rule *SolverRule) {

	dateToString := func(now time.Time, delta time.Duration) string {
		if now.IsZero() {
			return ""
		}
		return now.Add(delta).Format("2006-01-02 15:04:05")
	}

	t.tbl.AddRow(rule.Name(),
		rule.From.String(),
		rule.To.String(),
		rule.Duration().String(),
		dateToString(t.now, rule.From),
		dateToString(t.now, rule.To),
		rule.StartAmount.String(),
		rule.EndAmount.String(),
		fmt.Sprintf("%t", rule.IsFlatRate()),
		rule.ActivationAmount.String(),
		string(rule.DurationType),
	)
	t.empty = false
}

func (t *RulesTable) Print() {
	if !t.empty {
		fmt.Println()
		table.TitleTheme().Println(t.title)
		t.tbl.Print()
	}
}
