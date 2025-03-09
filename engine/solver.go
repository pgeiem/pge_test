package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/btree"
)

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
	//TODO replace From/To by a RelativeTimeSpan
	// Starting/End point in time
	From time.Duration
	// End point in time
	To time.Duration
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
}

// Define a collection of solver rule
type SolverRules []SolverRule

type MetaData map[string]interface{}

func NewRelativeLinearRule(name string, duration time.Duration, hourlyRate Amount) SolverRule {
	return SolverRule{
		RuleName:             name,
		From:                 time.Duration(0),
		To:                   duration,
		StartAmount:          0,
		EndAmount:            Amount(float64(hourlyRate) * duration.Hours()),
		StartTimePolicy:      ShiftablePolicy,
		RuleResolutionPolicy: ResolvePolicy,
	}
}

func NewRelativeFlatRateRule(name string, duration time.Duration, amount Amount) SolverRule {
	return SolverRule{
		RuleName:             name,
		From:                 time.Duration(0),
		To:                   duration,
		StartAmount:          amount,
		EndAmount:            amount,
		StartTimePolicy:      ShiftablePolicy,
		RuleResolutionPolicy: TruncatePolicy,
	}
}

func NewAbsoluteFlatRateRule(name string, from, to time.Duration, amount Amount) SolverRule {
	if from > to {
		panic(fmt.Errorf("invalid rule duration %v to %v", from, to))
	}
	return SolverRule{
		RuleName:             name,
		From:                 from,
		To:                   to,
		StartAmount:          amount,
		EndAmount:            amount,
		StartTimePolicy:      FixedPolicy,
		RuleResolutionPolicy: TruncatePolicy,
	}
}

func NewAbsoluteNonPaying(name string, from, to time.Duration) SolverRule {
	if from > to {
		panic(fmt.Errorf("invalid rule duration %v to %v", from, to))
	}
	return SolverRule{
		RuleName:             name,
		From:                 from,
		To:                   to,
		StartAmount:          0,
		EndAmount:            0,
		StartTimePolicy:      FixedPolicy,
		RuleResolutionPolicy: TruncatePolicy,
	}
}

func (rule SolverRule) Duration() time.Duration {
	return rule.To - rule.From
}

func (rule SolverRule) IsFlatRate() bool {
	return rule.StartAmount == rule.EndAmount
}

func (rule SolverRule) IsAbsoluteFlatRate() bool {
	return rule.IsFlatRate() && //is FLatRate
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

func InterpolAmountNoOffset(rule SolverRule, at time.Duration) Amount {
	return Amount(float64(rule.EndAmount-rule.StartAmount) * float64(at) / float64(rule.Duration()))
}

func InterpolAmount(rule SolverRule, at time.Duration) Amount {
	return InterpolAmountNoOffset(rule, at) + rule.StartAmount
}

// Shift the rule to the new start time, the new rule is returned and current rule is not changed
func (rule SolverRule) Shift(from time.Duration) SolverRule {
	rule.To = from + rule.Duration()
	rule.From = from
	rule.Trace = append(rule.Trace, fmt.Sprintf("shift to %s", from.String()))
	return rule
}

func (rule SolverRule) TruncateAfter(after time.Duration) SolverRule {
	ruleA := rule
	ruleA.To = after
	ruleA.Trace = append(rule.Trace, fmt.Sprintf("truncate after %s", after.String()))
	if rule.Duration() != time.Duration(0) {
		ruleA.EndAmount = InterpolAmount(rule, ruleA.Duration())
	}
	return ruleA
}

func (rule SolverRule) TruncateBefore(before time.Duration) SolverRule {
	ruleA := rule
	ruleA.From = before
	ruleA.Trace = append(rule.Trace, fmt.Sprintf("truncate before %s", before.String()))
	ruleA.StartAmount = 0

	if rule.Duration() != time.Duration(0) {
		ruleA.EndAmount = InterpolAmountNoOffset(rule, ruleA.Duration())
		// TODO: review all amount type with float or decimal library
		// ex: go run . processor -f samples/tariff_PortValais_t3.json -n 2024-11-28T16:13:00
		//     go run . player -f output/table.json -a 150
	}
	return ruleA
}

func (rule SolverRule) TruncateBetween(truncateStart, truncateEnd time.Duration) SolverRules {
	return SolverRules{rule.TruncateAfter(truncateStart), rule.TruncateBefore(truncateEnd)}
}

// Split the rule in two parts inserting a hole in middle from splitStart to splitEnd, the new
// rules are returned and current rule is not changed
func (rule SolverRule) Split(splitStart, splitEnd time.Duration) SolverRules {

	// Shorten the rule part before the inserted hole
	ruleA := rule.TruncateAfter(splitStart)

	// Shorten and shift the rule part after the inserted hole
	ruleB := rule
	ruleB.From = splitEnd
	ruleB.To = rule.To + splitEnd - splitStart
	ruleB.Trace = append(rule.Trace, fmt.Sprintf("truncate split between %s and %s", splitStart.String(), splitEnd.String()))

	ruleB.StartAmount = 0
	ruleB.EndAmount = rule.EndAmount - ruleA.EndAmount

	return SolverRules{ruleA, ruleB}
}

type Solver struct {
	now                         time.Time //TODO rework name for now - window inconsistent with between from/to
	window                      time.Duration
	rules                       *btree.BTreeG[*SolverRule]
	flatrates                   *btree.BTreeG[*SolverRule]
	currentRelativeStartOffset  time.Duration
	currentRelativeAmountOffset Amount
	activatedFlatRatesSum       Amount
}

func NewSolver() Solver {

	// Sorting function for B-Tree storing all solved rules segments
	// sorted by start time, then by start amount, then by end time
	RulesLess := func(i, j *SolverRule) bool {
		if i.From == j.From {
			if i.StartAmount == j.StartAmount {
				return i.To < j.To
			}
			return i.StartAmount < j.StartAmount
		}
		return i.From < j.From
	}

	return Solver{
		rules:     btree.NewG(2, RulesLess),
		flatrates: btree.NewG(2, RulesLess),
	}
}

func (s *Solver) SetWindow(now time.Time, window time.Duration) {
	s.now = now
	s.window = window
}

func (s *Solver) Append(rules ...SolverRule) {
	for i := range rules {
		if rules[i].IsAbsoluteFlatRate() {
			s.flatrates.ReplaceOrInsert(&rules[i])
		} else {
			s.solveAndAppend(rules[i])
		}
	}
}

// Solve the rule against an Higer Priority Rule resolving the conflict according to rule policy
// a collection of new rules containing 0, 1, or 2 rules is returned and current rule is not changed
// the second return value is true if the rule has intersected and has been changed, false if untouched
func (s *Solver) solveVsSingle(lpRule SolverRule, hpRule *SolverRule) (SolverRules, bool) {

	// trivial case, both rules don't overlap
	if (hpRule.To <= lpRule.From) ||
		(hpRule.From >= lpRule.To) {
		return SolverRules{lpRule}, false
	}

	fmt.Println("   ## solveVsSingle", lpRule.Name(), "vs", hpRule.Name())

	switch lpRule.RuleResolutionPolicy {

	// both rules overlap at least slightly, if policy is 'remove' then remove the low priority rule
	case DeletePolicy:
		fmt.Println("    DeletePolicy", lpRule.Name(), "vs", hpRule.Name())
		return SolverRules{}, true

	// both rules overlap at least slightly, if policy is 'resolve' then try to split rule to fill the holes
	case ResolvePolicy:
		fmt.Println("   ResolvePolicy", lpRule.Name(), "vs", hpRule.Name(), lpRule, hpRule)

		// high priority rule is before low priority rule, then low priority rule is simply shifted
		if hpRule.From <= lpRule.From {
			return SolverRules{lpRule.Shift(hpRule.To)}, true
		} else {
			// high priority rule is after or in middle of low priority rule, then low priority rule is split and shifted
			return lpRule.Split(hpRule.From, hpRule.To), true
		}

	// both rules overlap at least slightly, if policy is 'truncate' then truncate overlapping rule
	case TruncatePolicy:
		fmt.Println("   TruncatePolicy", lpRule.Name(), "vs", hpRule.Name(), lpRule, hpRule)

		// high priority rule is partially after low priority rule, then low priority rule end is truncated
		if hpRule.From >= lpRule.From && hpRule.To >= lpRule.To {
			return SolverRules{lpRule.TruncateAfter(hpRule.From)}, true
		}

		// high priority rule is partially before low priority rule, then low priority rule end is truncated
		if hpRule.From <= lpRule.From && hpRule.To <= lpRule.To {
			return SolverRules{lpRule.TruncateBefore(hpRule.To)}, true
		}

		// high priority rule completely overlap low priority rule, then remove the low priority rule
		if hpRule.From <= lpRule.From && hpRule.To >= lpRule.To {
			return SolverRules{}, true
		}

		// high priority rule is in middle of low priority rule, then low priority rule middle is truncated
		if hpRule.From >= lpRule.From && hpRule.To <= lpRule.To {
			return lpRule.TruncateBetween(hpRule.From, hpRule.To), true
		}

		//TODO add policiy for IntersectSolve

	default:
		panic(fmt.Errorf("unhandled solving policy %v between %v and %v", lpRule.RuleResolutionPolicy, lpRule, hpRule))
	}

	return SolverRules{}, true
}

func (s *Solver) SolveVsAll(lpRule SolverRule) (SolverRules, bool) {
	var newRules SolverRules
	var changed bool

	fmt.Println("\n ## SolveVsAll", lpRule)

	// Check if we have a flatrate to apply
	bestFlatRate := s.GetBestFlatRate(&lpRule)
	if bestFlatRate != nil {
		intersectAt := s.FindIntersectPositionFlatRate(&lpRule, bestFlatRate)
		newRule := bestFlatRate.TruncateBefore(intersectAt)
		fmt.Println(" >> adding flatrate based new rule", newRule)
		s.activatedFlatRatesSum += bestFlatRate.EndAmount
		if newRule.Duration() > time.Duration(0) {
			// Add a new rule based on the activated flatrate
			newRule.StartAmount, newRule.EndAmount = AmountZero, AmountZero
			newRule.Trace = append(newRule.Trace, fmt.Sprintf("derivated from flatrate %s, crossed by %s", bestFlatRate.Name(), lpRule.Name()))
			s.rules.ReplaceOrInsert(&newRule)
		}
	}

	// Loop over all rules in the collection and solve the current rule against each of them
	if s.rules.Len() == 0 {
		newRules = append(newRules, lpRule)
	} else {
		s.rules.Ascend(func(hpRule *SolverRule) bool {
			newRules, changed = s.solveVsSingle(lpRule, hpRule)
			fmt.Println(" >> solveVsSingle Result", newRules, changed)
			return !changed // Stop the iterator loop if the rule has been changed
		})
	}

	return newRules, changed
}

// Solve the rule against a collection of Higer Priority Rule resolving the conflict according to rules policy
// a collection of new rules is returned and current rule is not changed
func (s *Solver) solveAndAppend(lpRule SolverRule) {

	var incRelativeStartOffset time.Duration
	//var incRelativeAmountOffset Amount
	fmt.Println("\n------\nSolving rule", lpRule.Name(), "from", lpRule.From, "to", lpRule.To)

	// Shift the rule if needed using the current start offset
	if lpRule.StartTimePolicy == ShiftablePolicy {
		fmt.Println(" >> shift rule", lpRule.Name(), "at", s.currentRelativeStartOffset)
		lpRule = lpRule.Shift(s.currentRelativeStartOffset)
		incRelativeStartOffset = lpRule.Duration()
	}

	// Loop over all rules in the collection and solve the current rule against each of them
	run := true
	for run {
		//var ret SolverRules
		ret, changed := s.SolveVsAll(lpRule)
		run = changed
		fmt.Println(" SolveVsAll Result", ret)
		switch len(ret) {
		case 0: // Rule deleted
			lpRule = SolverRule{}
			run = false
		case 1: // Rule Shifted or untouched
			lpRule = ret[0]
		case 2: // Rule splitted
			lpRule = ret[1]                    // right part is the new rule to solve
			s.rules.ReplaceOrInsert(&(ret[0])) // Left part may be inserted in the new rules collection
			s.currentRelativeAmountOffset += ret[0].EndAmount
			fmt.Println(" >> append splitted rule", ret[0])
		}
	}

	// Insert the last rule part in the new rules collection
	if lpRule.Duration() > time.Duration(0) {
		s.rules.ReplaceOrInsert(&lpRule)
		s.currentRelativeAmountOffset += lpRule.EndAmount
		fmt.Println(" >> append final rule", lpRule)
	}

	// Update the current start offset used for relative rules
	s.currentRelativeStartOffset += incRelativeStartOffset
}

func (s *Solver) IsIntersectingFlatRate(relativeRule, flatRateRule *SolverRule) bool {
	flatrateAmount := flatRateRule.EndAmount + s.activatedFlatRatesSum
	intersect := flatRateRule.IsAbsoluteFlatRate() && relativeRule.IsRelative()
	// StartAmount bigger than flatrate amount
	intersect = intersect && (s.currentRelativeAmountOffset+relativeRule.StartAmount < flatrateAmount)
	// EndAmount lower than flatrate amount
	intersect = intersect && (s.currentRelativeAmountOffset+relativeRule.EndAmount >= flatrateAmount)
	// Start after the flatrate end
	intersect = intersect && (s.currentRelativeStartOffset+relativeRule.From < flatRateRule.To)
	// End before the flatrate start
	intersect = intersect && (s.currentRelativeStartOffset+relativeRule.To >= flatRateRule.From)

	//fmt.Println(" >> IsIntersectingFlatRate", relativeRule.Name(), "vs", flatRateRule.Name(), intersect, s.currentRelativeAmountOffset, s.activatedFlatRatesSum, s.currentRelativeStartOffset /*, relativeRule, flatRateRule*/)

	return intersect
}

func (s *Solver) FindIntersectPositionFlatRate(relativeRule, flatRateRule *SolverRule) time.Duration {
	var out time.Duration
	if s.IsIntersectingFlatRate(relativeRule, flatRateRule) {
		flatrateAmount := flatRateRule.EndAmount + s.activatedFlatRatesSum
		relativeStartAmount := s.currentRelativeAmountOffset + relativeRule.StartAmount
		relativeEndAmount := s.currentRelativeAmountOffset + relativeRule.EndAmount
		relativeFrom := s.currentRelativeStartOffset + relativeRule.From

		out = time.Duration(float64(flatrateAmount-relativeStartAmount)/float64(relativeEndAmount-relativeStartAmount)*float64(relativeRule.To-relativeRule.From) + float64(relativeFrom))
		fmt.Println(" >> FindIntersectPositionFlatRate", relativeRule.Name(), "vs", flatRateRule.Name(), "=>", out, "|", relativeRule, flatRateRule, "|", s.currentRelativeAmountOffset, s.activatedFlatRatesSum, s.currentRelativeStartOffset, "|", flatrateAmount, relativeStartAmount, relativeEndAmount, relativeFrom)
	}
	return out
}

func (s *Solver) GetBestFlatRate(lpRule *SolverRule) *SolverRule {
	var bestRule *SolverRule
	minAmount := AmountMax
	s.flatrates.Ascend(func(flatRateRule *SolverRule) bool {
		if s.IsIntersectingFlatRate(lpRule, flatRateRule) {
			flatrateAmount := flatRateRule.StartAmount + s.activatedFlatRatesSum
			if flatrateAmount < minAmount {
				minAmount = flatRateRule.StartAmount
				bestRule = flatRateRule
			}
		}
		return true
	})
	if bestRule != nil {
		fmt.Println(" >>>> GetBestFlatRate for", lpRule, "is", bestRule.Name())
	} else {
		fmt.Println(" >>>> GetBestFlatRate for", lpRule, "is nil")
	}
	return bestRule
}

type OutputSegment struct {
	SegName  string   `json:"n,omitempty"`
	Trace    []string `json:"dbg,omitempty"`
	At       int      `json:"t"`
	Amount   Amount   `json:"a"`
	Islinear bool     `json:"l"`
	Meta     MetaData `json:"m,omitempty"`
}

func (seg OutputSegment) String() string {
	at := time.Duration(seg.At) * time.Second
	return fmt.Sprintf(" - %s: %s (isLinear %t)", at, seg.Amount, seg.Islinear)
}

type OutputSegments []OutputSegment

func (segs OutputSegments) String() string {
	out := "OutputSegments:\n"
	for i := range segs {
		out += segs[i].String() + "\n"
	}
	return out
}

type Output struct {
	Now   time.Time      `json:"now"`
	Table OutputSegments `json:"table"`
}

func (segs Output) ToJson() ([]byte, error) {
	return json.Marshal(segs)
}

func (s *Solver) GenerateOutput(detailed bool) Output {
	var out Output
	var previous OutputSegment

	out.Now = s.now

	s.rules.Ascend(func(rule *SolverRule) bool {
		fmt.Println("Rule", rule)
		seg := OutputSegment{
			At:       int(math.Round(rule.To.Seconds())),
			Amount:   rule.EndAmount + previous.Amount,
			Islinear: !rule.IsFlatRate(),
			Meta:     rule.Meta,
		}
		if detailed {
			seg.SegName = rule.Name()
			seg.Trace = rule.Trace
		}
		out.Table = append(out.Table, seg)
		previous = seg
		return true
	})
	return out
}
