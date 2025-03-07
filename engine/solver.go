package engine

import (
	"fmt"
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
	KeepPolicy     RuleResolutionPolicy = "keep"
)

//TOOD: merge RuleResolutionPolicy with StartTimePolicy as shiftable is usefull only with truncate ?

// Define the solver rule
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
	// Rule type reported to output used for tariff details
	//Type string

	StartTimePolicy StartTimePolicy

	RuleResolutionPolicy RuleResolutionPolicy

	Meta interface{}
}

// Define a collection of solver rule
type SolverRules []SolverRule

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
		RuleResolutionPolicy: KeepPolicy,
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
		rule.StartAmount != 0 // is not non-paying
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
	rule.Trace = append(rule.Trace, "shift to", from.String())
	return rule
}

func (rule SolverRule) TruncateAfter(after time.Duration) SolverRule {
	ruleA := rule
	ruleA.To = after
	ruleA.Trace = append(rule.Trace, "truncate after", after.String())
	if rule.Duration() != time.Duration(0) {
		ruleA.EndAmount = InterpolAmount(rule, ruleA.Duration())
	}
	return ruleA
}

func (rule SolverRule) TruncateBefore(before time.Duration) SolverRule {
	ruleA := rule
	ruleA.From = before
	ruleA.Trace = append(rule.Trace, "truncate before", before.String())
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
	ruleB.Trace = append(rule.Trace, "truncate split between", splitStart.String(), "and", splitEnd.String())

	ruleB.StartAmount = 0
	ruleB.EndAmount = rule.EndAmount - ruleA.EndAmount

	return SolverRules{ruleA, ruleB}
}

type Solver struct {
	now                         time.Time //TODO rework name for now - window inconsistent with between from/to
	window                      time.Duration
	rules                       *btree.BTreeG[*SolverRule]
	absFlatRates                *btree.BTreeG[*SolverRule]
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

	// Sorting function for B-Tree storing absolute flatrate rules segments
	AmountLess := func(i, j *SolverRule) bool {
		return i.EndAmount < j.EndAmount
	}

	return Solver{
		rules:        btree.NewG(2, RulesLess),
		absFlatRates: btree.NewG(2, AmountLess),
	}
}

func (s *Solver) SetWindow(now time.Time, window time.Duration) {
	s.now = now
	s.window = window
}

func (s *Solver) AppendSolverRules(rules ...SolverRule) {
	for i := range rules {
		s.solveAndAppend(rules[i])
	}
}

// Solve the rule against an Higer Priority Rule resolving the conflict according to rule policy
// a collection of new rules containing 0, 1, or 2 rules is returned and current rule is not changed
func (s *Solver) solveVsSingle(lpRule, hpRule SolverRule) SolverRules {

	// trivial case, both rules don't overlap
	if (hpRule.To <= lpRule.From) ||
		(hpRule.From >= lpRule.To) {
		return SolverRules{lpRule}
	}

	switch lpRule.RuleResolutionPolicy {

	// Keep policy, do nothing
	case KeepPolicy:
		return SolverRules{lpRule}

	// both rules overlap at least slightly, if policy is 'remove' then remove the low priority rule
	case DeletePolicy:
		fmt.Println(" DeletePolicy", lpRule.Name(), "vs", hpRule.Name())
		return SolverRules{}

	// both rules overlap at least slightly, if policy is 'resolve' then try to split rule to fill the holes
	case ResolvePolicy:
		fmt.Println(" ResolvePolicy", lpRule.Name(), "vs", hpRule.Name(), lpRule, hpRule)

		// high priority rule is before low priority rule, then low priority rule is simply shifted
		if hpRule.From <= lpRule.From {
			return SolverRules{lpRule.Shift(hpRule.To)}
		} else {
			// high priority rule is after or in middle of low priority rule, then low priority rule is split and shifted
			return lpRule.Split(hpRule.From, hpRule.To)
		}

	// both rules overlap at least slightly, if policy is 'truncate' then truncate overlapping rule
	case TruncatePolicy:
		fmt.Println(" TruncatePolicy", lpRule.Name(), "vs", hpRule.Name(), lpRule, hpRule)

		// high priority rule is partially after low priority rule, then low priority rule end is truncated
		if hpRule.From >= lpRule.From && hpRule.To >= lpRule.To {
			return SolverRules{lpRule.TruncateAfter(hpRule.From)}
		}

		// high priority rule is partially before low priority rule, then low priority rule end is truncated
		if hpRule.From <= lpRule.From && hpRule.To <= lpRule.To {
			return SolverRules{lpRule.TruncateBefore(hpRule.To)}
		}

		// high priority rule completely overlap low priority rule, then remove the low priority rule
		if hpRule.From <= lpRule.From && hpRule.To >= lpRule.To {
			return SolverRules{}
		}

		// high priority rule is in middle of low priority rule, then low priority rule middle is truncated
		if hpRule.From >= lpRule.From && hpRule.To <= lpRule.To {
			return lpRule.TruncateBetween(hpRule.From, hpRule.To)
		}

		//TODO add policiy for IntersectSolve

	default:
		panic(fmt.Errorf("unhandled solving policy %v between %v and %v", lpRule.RuleResolutionPolicy, lpRule, hpRule))
	}

	return SolverRules{}
}

/*
func (s *Solver) SelectFlatRate(lpRule SolverRule) {
	s.rules.Ascend(func(rule SolverRule) bool {
		if rule.IsAbsoluteFlatRate() {

			return false
		}
		return true
	})
	return flatRate
}*/

// Solve the rule against a collection of Higer Priority Rule resolving the conflict according to rules policy
// a collection of new rules is returned and current rule is not changed
func (s *Solver) solveAndAppend(lpRule SolverRule) {

	var newRules SolverRules
	//var sumAmount Amount
	var incRelativeStartOffset time.Duration
	var incRelativeAmountOffset Amount

	fmt.Println("------ ", lpRule.Name())
	//fmt.Println("Solving rule", lpRule.Name(), "from", lpRule.From, "to", lpRule.To)

	// Shift the rule if needed using the current start offset
	if lpRule.StartTimePolicy == ShiftablePolicy {
		//fmt.Println("### Shifting rule", lpRule.Name, "from", lpRule.From, "to", lpRule.From+s.currentStartOffset)
		lpRule = lpRule.Shift(s.currentRelativeStartOffset)
		incRelativeStartOffset = lpRule.Duration()
		incRelativeAmountOffset = lpRule.EndAmount
	}

	// Loop over all rules in the collection and solve the current rule against each of them
	s.rules.Ascend(func(hpRule *SolverRule) bool {
		tmpHpRule := *hpRule
		skip := false
		fmt.Println("Solving rule", lpRule.Name(), "vs", hpRule.Name() /*, s.currentRelativeAmountOffset, s.currentRelativeStartOffset*/)

		// If lpRule may intersect with hpRule flatrate
		if lpRule.StartTimePolicy == ShiftablePolicy && hpRule.IsAbsoluteFlatRate() {
			isBest, intersectAfter := s.IsBestFlatRateAvailable(&lpRule, hpRule)
			skip = !isBest
			if isBest {
				s.activatedFlatRatesSum += hpRule.StartAmount
				tmpHpRule = hpRule.TruncateBefore(intersectAfter)
			}
		}
		fmt.Println(" skip", skip, "activatedFlatRatesSum", s.activatedFlatRatesSum)

		if !skip {
			ret := s.solveVsSingle(lpRule, tmpHpRule)
			fmt.Println(" >>", len(ret), "rules returned:", ret)
			switch len(ret) {
			case 0: // Rule deleted
				lpRule = SolverRule{}
				return false
			case 1: // Rule Shifted or untouched
				lpRule = ret[0]
			case 2: // Rule splitted
				newRules = append(newRules, ret[0]) // Left part may be inserted in the new rules collection
				lpRule = ret[1]                     // right part is the new rule to solve
			}
		}
		return true
	})

	// Insert the last rule part in the new rules collection
	newRules = append(newRules, lpRule)

	// Effectively insert all parts of the resolved rules in the rules collection
	for _, rule := range newRules {
		if rule.Duration() > time.Duration(0) {
			s.rules.ReplaceOrInsert(&rule)
			if rule.IsAbsoluteFlatRate() {
				s.absFlatRates.ReplaceOrInsert(&rule)
			}
		}
	}

	// Update the current start/amount offset used for relative rules
	s.currentRelativeStartOffset += incRelativeStartOffset
	s.currentRelativeAmountOffset += incRelativeAmountOffset
}

func (s *Solver) IsIntersectingFlatRate(relativeRule, flatRateRule *SolverRule) (bool, time.Duration) {
	var intersectAfter time.Duration
	intersect := flatRateRule.IsAbsoluteFlatRate() && relativeRule.IsRelative()
	flatrateAmount := flatRateRule.StartAmount + s.activatedFlatRatesSum
	// StartAmount bigger than flatrate amount
	intersect = intersect && !(s.currentRelativeAmountOffset+relativeRule.StartAmount > flatrateAmount)
	// EndAmount lower than flatrate amount
	intersect = intersect && !(s.currentRelativeAmountOffset+relativeRule.EndAmount < flatrateAmount)
	// Start after the flatrate end
	intersect = intersect && !(s.currentRelativeStartOffset+relativeRule.From > flatRateRule.To)
	// End before the flatrate start
	intersect = intersect && !(s.currentRelativeStartOffset+relativeRule.To < flatRateRule.From)

	if intersect {
		intersectAfter = time.Duration(float64(relativeRule.Duration()) * float64(flatrateAmount-relativeRule.StartAmount) / float64(relativeRule.EndAmount-relativeRule.StartAmount))
	}
	fmt.Println(" >> IsIntersectingFlatRate", relativeRule.Name(), "vs", flatRateRule.Name(), intersect, intersectAfter /*, relativeRule, flatRateRule*/)

	return intersect, intersectAfter
}

func (s *Solver) IsBestFlatRateAvailable(lpRule, hpRule *SolverRule) (bool, time.Duration) {
	cpt := 0
	minAmount := AmountMax
	bestRule := hpRule
	bestIntersectAfter := time.Duration(0)
	s.rules.Ascend(func(rule *SolverRule) bool {
		isIntersect, intersectAfter := s.IsIntersectingFlatRate(lpRule, rule)
		if isIntersect {
			cpt++
			if rule.StartAmount < minAmount {
				minAmount = rule.StartAmount
				bestRule = rule
				bestIntersectAfter = intersectAfter
			}
		}
		return true
	})
	fmt.Println(" >>>> cpt", cpt, "minAmount", minAmount, "bestName", bestRule.Name(), "isBest", bestRule == hpRule)
	return bestRule == hpRule, bestIntersectAfter //TODO call IntersectAfter function from here
}

func (s *Solver) SumAmount(until time.Duration) Amount {
	var sum Amount
	s.rules.Ascend(func(rule *SolverRule) bool {
		if rule.From < until && rule.To > until {
			sum += InterpolAmount(*rule, until-rule.From)
			return false
		}
		if rule.To > until {
			return false
		}
		sum += rule.EndAmount
		return true
	})
	return sum
}
