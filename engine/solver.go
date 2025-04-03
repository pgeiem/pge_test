package engine

import (
	"fmt"
	"time"

	"github.com/google/btree"
)

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

func (rule *SolverRule) And(timespan RelativeTimeSpan) (*SolverRule, bool) {
	// rule is fully inside timespan, then rule is not touched
	if rule.From >= timespan.From && rule.To <= timespan.To {
		return rule, true
		// rule is longer than timespan (timespan fully inside rule), then rule beginning and end are truncated
	} else if rule.From <= timespan.From && rule.To >= timespan.To {
		r := rule.TruncateAfter(timespan.To).TruncateBefore(timespan.From)
		r.Trace = append(r.Trace, "truncated for sequence merging")
		return &r, true
		// rule is partially at the end of timespan, then rule end is truncated
	} else if rule.From < timespan.To && rule.To >= timespan.To {
		r := rule.TruncateAfter(timespan.To)
		r.Trace = append(r.Trace, "truncated for sequence merging")
		return &r, true
		// rule is partially at the end beginning of timespan, then rule beginning is truncated
	} else if rule.From <= timespan.From && rule.To > timespan.From {
		r := rule.TruncateBefore(timespan.From)
		r.Trace = append(r.Trace, "truncated for sequence merging")
		return &r, true
	} else {
		return nil, false
	}
}

func (rule *SolverRule) DurationForAmount(amount Amount) time.Duration {
	fmt.Println(" >> DurationForAmount", rule.Name(), amount, rule.StartAmount, rule.EndAmount, rule.Duration())
	return rule.From + time.Duration(float64(rule.Duration())*float64(amount)/float64(rule.EndAmount-rule.StartAmount))
}

type Solver struct {
	now                        time.Time
	window                     time.Duration
	flatrateRules              *btree.BTreeG[*SolverRule]
	solvedRules                *btree.BTreeG[*SolverRule]
	fixedRules                 *btree.BTreeG[*SolverRule]
	shiftableRules             []*SolverRule
	currentRelativeStartOffset time.Duration
}

func NewSolver() Solver {

	// Sorting function for B-Tree storing all solved rules segments
	// sorted by start time, then by start amount, then by end time
	RulesLess := func(i, j *SolverRule) bool {
		if i.From == j.From {
			if i.StartAmount == j.StartAmount {
				return i.To > j.To
			}
			return i.StartAmount < j.StartAmount
		}
		return i.From < j.From
	}

	return Solver{
		//rules:       btree.NewG(2, RulesLess),
		flatrateRules: btree.NewG(2, RulesLess),
		solvedRules:   btree.NewG(2, RulesLess),
		fixedRules:    btree.NewG(2, RulesLess),
	}
}

func (s *Solver) SetWindow(now time.Time, window time.Duration) {
	s.now = now
	s.window = window
}

func (s *Solver) AppendMany(rules ...SolverRule) {
	for i := range rules {
		s.Append(&rules[i])
	}
}

func (s *Solver) Append(rule *SolverRule) {
	if rule.ActivationAmount > 0 {
		// flatrate rules are stored in a b-tree
		s.flatrateRules.ReplaceOrInsert(rule)
	} else if rule.StartTimePolicy == FixedPolicy {
		// fixed rules are stored in a sorted b-tree
		s.fixedRules.ReplaceOrInsert(rule)
	} else {
		// shiftable rules are stored in a list in the appended order
		s.shiftableRules = append(s.shiftableRules, rule)
	}
}

func (s *Solver) Solve() {

	fmt.Println("Solving rules...")
	fmt.Println("  >> flatrates")
	s.flatrateRules.Ascend(func(rule *SolverRule) bool {
		fmt.Println("    >>", rule.Name(), rule.From, rule.To, rule.ActivationAmount)
		return true
	})
	fmt.Println("  >> fixed rules")
	s.fixedRules.Ascend(func(rule *SolverRule) bool {
		fmt.Println("    >>", rule.Name(), rule.From, rule.To, rule.StartAmount, rule.EndAmount)
		return true
	})
	fmt.Println("  >> shiftable rules")
	for i := range s.shiftableRules {
		fmt.Println("    >>", s.shiftableRules[i].Name(), s.shiftableRules[i].From, s.shiftableRules[i].To, s.shiftableRules[i].StartAmount, s.shiftableRules[i].EndAmount)
	}

	// TODO solve fixed rules first

	// Solve all shiftable rules against all fixed rules
	for i := range s.shiftableRules {
		s.solveVsFixedRules(s.shiftableRules[i])
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

	fmt.Println(" >> solveVsSingle", lpRule.Name(), "vs", hpRule.Name())
	lpRule.Trace = append(lpRule.Trace, fmt.Sprintf("solve against %s", hpRule.Name()))

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

func (s *Solver) buildFixedRulesList(lpRule *SolverRule) *btree.BTreeG[*SolverRule] {

	fixedRules := s.fixedRules.Clone()

	s.flatrateRules.Ascend(func(flatRateRule *SolverRule) bool {
		// Check if the higer priority rule is activated, if not skip the rule
		activatedAfter, activated := s.findFlatRateActivationTime(flatRateRule, lpRule)
		fmt.Println(" >> findFlatRateActivationTime", flatRateRule.Name(), activatedAfter, activated)
		//Skip rules if not activated
		if !activated {
			return true
		}

		// If the higer priority is activated after a certains time (flatrate), then truncate the rule before this activation time
		if activatedAfter > 0 {
			tmp := flatRateRule.TruncateBefore(activatedAfter)
			flatRateRule = &tmp
			fmt.Println(" >> truncate flatrate rule before flatrate activation", activatedAfter)
		}

		fixedRules.ReplaceOrInsert(flatRateRule)
		fmt.Println(" >> append flatrate rule", flatRateRule)
		return true
	})

	fmt.Println(" >> fixed rules list", fixedRules.Len(), "rules")
	fixedRules.Ascend(func(rule *SolverRule) bool {
		fmt.Println("    >>", rule.Name(), rule.From, rule.To, rule.StartAmount, rule.EndAmount)
		return true
	})

	return fixedRules
}

// Solve the rule against a collection of fixed rules resolving the conflict according to rules policy
// a collection of new rules is returned and current rule is not changed
func (s *Solver) solveVsFixedRules(lpRule *SolverRule) {

	/*rulesfifo := arrayqueue.New[*SolverRule]()
	rulesfifo.Enqueue(lpRule)

	for !rulesfifo.Empty() {*/

	//rule, _ := rulesfifo.Dequeue()

	fmt.Println("\n------\nSolving rule", lpRule.Name(), "from", lpRule.From, "to", lpRule.To)

	// Shift the rule if needed using the current start offset
	_, startOffset := s.sumAllSolvedRules()
	tmp := lpRule.Shift(startOffset)
	lpRule = &tmp
	fmt.Println(" >> shift rule", lpRule.Name(), "at", s.currentRelativeStartOffset)

	solved := false
	for !solved {
		fmt.Println("---")

		// Build a list of fixed rules including all fixed rules and activated flatrates
		fixedRules := s.buildFixedRulesList(lpRule)

		// Iterate over the previously built fixedrules list and solve the current rule against each of them
		fixedRules.Ascend(func(hpRule *SolverRule) bool {
			// Solve the lower priority rule against the higher priority rule
			ret, changed := s.solveVsSingle(*lpRule, hpRule)
			fmt.Println(" >> solveVsSingle Result", ret, changed)
			solved = !changed
			if changed {
				switch len(ret) {
				case 0: // Rule deleted, exit the loop as there is nothing to solve anymore
					lpRule = nil
					solved = true
				case 1: // Rule Shifted or truncated, continue the solving process with the new rule
					lpRule = &ret[0]
				case 2: // Rule splitted
					lpRule = &ret[1]                         // right part is the new rule to solve
					s.solvedRules.ReplaceOrInsert(&(ret[0])) // Left part may be inserted in the new rules collection
					s.solvedRules.ReplaceOrInsert(hpRule)    //
					fmt.Println(" >> append splitted fixed rule", ret[0])
				}
			}
			return !changed
		})

		if fixedRules.Len() == 0 {
			solved = true
		}
	}

	if lpRule != nil && lpRule.Duration() > time.Duration(0) {
		// Insert the last rule part in the new rules collection
		s.solvedRules.ReplaceOrInsert(lpRule)
	}

	//}
}

func (s *Solver) sumAllSolvedRules() (Amount, time.Duration) {
	var amountSum Amount
	var durationSum time.Duration

	// Loop over all already solved entries and determine the end amount and end duration
	s.solvedRules.Ascend(func(rule *SolverRule) bool {
		amountSum += rule.EndAmount
		durationSum = rule.To // last rule duration is the total duration
		return true
	})

	return amountSum, durationSum
}

/*
func (s *Solver) addAndSolveFixedRule(lpRule *SolverRule) {

	s.fixedRules.Ascend(func(hpRule *SolverRule) bool {
		ret, changed := s.solveVsSingle(*lpRule, hpRule)
		if changed {
			switch len(ret) {
			case 0: // Rule deleted
				lpRule = nil
				return false
			case 1: // Rule Shifted, truncated or untouched
				lpRule = &ret[0]
			case 2: // Rule splitted
				lpRule = &ret[1]                        // right part is the new rule to solve
				s.fixedRules.ReplaceOrInsert(&(ret[0])) // Left part may be inserted in the new rules collection
				s.currentRelativeAmountOffset += ret[0].EndAmount
				fmt.Println(" >> append splitted fixed rule", ret[0])
			}
		}
		return true
	})

	// Insert the last rule part in the new rules collection
	if lpRule != nil && lpRule.Duration() > time.Duration(0) {
		s.fixedRules.ReplaceOrInsert(lpRule)
		fmt.Println(" >> append final fixed rule", lpRule)
	}
}
*/

/*
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
			newRule.StartAmount, newRule.EndAmount = 0, 0
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
			//fmt.Println(" >> solveVsSingle Result", newRules, changed)
			return !changed // Stop the iterator loop if the rule has been changed
		})
	}

	return newRules, changed
}
*/

/*
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
*/
/*
func (s *Solver) solvedRulesSumAmountInRange(timespan RelativeTimeSpan, extraRules ...*SolverRule) Amount {
	var sumAmount Amount
	// Sum amount of all solved rules included in the timespan range
	s.solvedRules.Ascend(func(rule *SolverRule) bool {
		r, _ := rule.And(timespan)
		if r != nil {
			sumAmount += r.EndAmount
		}
		return true
	})
	// Sum amount of all extra rules included in the timespan rnage
	for _, rule := range extraRules {
		r, _ := rule.And(timespan)
		if r != nil {
			sumAmount += r.EndAmount
		}
	}
	return sumAmount
}*/

func (s *Solver) findFlatRateActivationTime(flatRateRule *SolverRule, extraRules ...*SolverRule) (time.Duration, bool) {
	var sumAmount Amount
	var activated bool
	var activatedAfter time.Duration

	// Trivial case, if activation amount is 0, then always activated from the beginning
	if flatRateRule.ActivationAmount == 0 {
		return 0, true
	}

	//fmt.Println(" >> findFlatRateActivationTime for", flatRateRule.Name(), "with amount", flatRateRule.ActivationAmount)

	processRule := func(rule *SolverRule) (time.Duration, bool) {
		r, _ := rule.And(flatRateRule.RelativeTimeSpan)
		if r != nil {
			//fmt.Println("    >> processRule", rule.Name(), sumAmount, "+", r.EndAmount, "vs", flatRateRule.ActivationAmount)
			if sumAmount+r.EndAmount > flatRateRule.ActivationAmount {
				return r.DurationForAmount(flatRateRule.ActivationAmount - sumAmount), true
			}
			sumAmount += r.EndAmount
		}
		return 0, false
	}

	// Sum amount of all solved rules included in the timespan range
	s.solvedRules.Ascend(func(rule *SolverRule) bool {
		activatedAfter, activated = processRule(rule)
		return !activated
	})

	// Sum amount of all extra rules included in the timespan rnage
	if !activated {
		for _, rule := range extraRules {
			activatedAfter, activated = processRule(rule)

		}
	}
	return activatedAfter, activated
}

/*
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

	//fmt.Println(" >> IsIntersectingFlatRate", relativeRule.Name(), "vs", flatRateRule.Name(), intersect, s.currentRelativeAmountOffset, s.activatedFlatRatesSum, s.currentRelativeStartOffset)

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

*/

func (s *Solver) ExtractRulesInRange(timespan RelativeTimeSpan) SolverRules {
	var out SolverRules
	s.solvedRules.Ascend(func(rule *SolverRule) bool {
		r, _ := rule.And(timespan)
		if r != nil {
			out = append(out, *r)
		}
		return true
	})
	return out
}
