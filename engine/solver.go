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
	//fmt.Println(" >> DurationForAmount", rule.Name(), amount, rule.StartAmount, rule.EndAmount, rule.Duration())
	return rule.From + time.Duration(float64(rule.Duration())*float64(amount)/float64(rule.EndAmount-rule.StartAmount))
}

func (rule SolverRule) TruncateAfterAmount(amount Amount) SolverRule {
	if amount >= rule.EndAmount {
		return rule
	}
	rule.To = rule.DurationForAmount(amount)
	rule.EndAmount = amount
	if rule.StartAmount > rule.EndAmount {
		rule.StartAmount = rule.EndAmount
	}
	fmt.Println(" >> TruncateAfterAmount", amount, "->", rule.To)
	rule.Trace = append(rule.Trace, fmt.Sprintf("truncate after amount %s", amount.String()))
	return rule
}

type Solver struct {
	now            time.Time
	window         time.Duration
	flatrateRules  *btree.BTreeG[*SolverRule]
	fixedRules     *btree.BTreeG[*SolverRule]
	shiftableRules []*SolverRule
	solvedRules    *btree.BTreeG[*SolverRule]
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

// TODO remove this
func (s *Solver) AppendByValue(rule SolverRule) {
	s.Append(&rule)
}

func (s *Solver) Append(rule *SolverRule) {
	if rule.ActivationAmount > 0 {
		// flatrate rules are stored in a b-tree
		s.flatrateRules.ReplaceOrInsert(rule)
	} else if rule.StartTimePolicy == FixedPolicy {
		// fixed rules are stored in a sorted b-tree
		s.solveAndAppend(rule, s.fixedRules)
		//s.fixedRules.ReplaceOrInsert(rule)
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

	// Solve potential continuous fixed rules starting from t=0
	s.SolveContinousFixedRules(time.Duration(0))

	// Solve all shiftable rules against all fixed rules
	for i := range s.shiftableRules {
		s.solveShiftableVsFixedRules(s.shiftableRules[i])
	}

	// Solve potential continuous fixed rules at the end of the last rule
	_, start := s.sumAllSolvedRules()
	s.SolveContinousFixedRules(start)
}

// SolveContinousFixedRules appends all potentially continous fixed rules in the solved rules collection
func (s *Solver) SolveContinousFixedRules(start time.Duration) {
	time := start
	s.fixedRules.Ascend(func(rule *SolverRule) bool {
		// Skip rules which are completely before the current time
		if rule.To <= time {
			return true
		}
		if rule.From <= time {
			s.solvedRules.ReplaceOrInsert(rule)
			time = rule.To
			return true
		}
		return false
	})
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

	default:
		panic(fmt.Errorf("unhandled solving policy %v between %v and %v", lpRule.RuleResolutionPolicy, lpRule, hpRule))
	}

	return SolverRules{}, true
}

// buildFixedRulesList builds a list of fixed rules including all fixed rules and activated flatrates
// This list potnetialy contains some activated flatrates which are valid only for the current rule
// The list must be regenerated for each rule to be solved
func (s *Solver) buildFixedRulesList(lpRule *SolverRule) *btree.BTreeG[*SolverRule] {

	fixedRules := s.fixedRules.Clone()

	s.flatrateRules.Ascend(func(flatRateRule *SolverRule) bool {
		// Check if the higer priority rule is activated, if not skip the rule
		activatedAfter, activated := s.findFlatRateActivationTime(flatRateRule, lpRule)
		//fmt.Println(" >> findFlatRateActivationTime", flatRateRule.Name(), activatedAfter, activated)
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

		s.solveAndAppend(flatRateRule, fixedRules)

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
func (s *Solver) solveShiftableVsFixedRules(lpRule *SolverRule) {

	fmt.Println("\n------\nSolving rule", lpRule.Name(), "from", lpRule.From, "to", lpRule.To)

	// Shift the rule using the current start offset
	_, startOffset := s.sumAllSolvedRules()
	tmp := lpRule.Shift(startOffset)
	lpRule = &tmp
	fmt.Println(" >> shift rule", lpRule.Name(), "at", startOffset)

	solved := false
	for !solved {
		fmt.Println("---")

		// Build a list of fixed rules including all fixed rules and activated flatrates
		fixedRules := s.buildFixedRulesList(lpRule)

		// Iterate over the previously built fixedrules list and solve the current rule against each of them
		fixedRules.Ascend(func(hpRule *SolverRule) bool {
			// Solve the lower priority rule against the higher priority rule
			ret, changed := s.solveVsSingle(*lpRule, hpRule)
			//fmt.Println(" >> solveVsSingle Result", ret, changed)
			solved = !changed
			if changed {
				switch len(ret) {
				case 0: // Rule deleted, exit the loop as there is nothing to solve anymore
					lpRule = nil
					solved = true
				case 1: // Rule Shifted or truncated, continue the solving process with the new rule
					lpRule = &ret[0]
					s.solvedRules.ReplaceOrInsert(hpRule) // when the rule is shifted we also insert the higher priority rule
				case 2: // Rule splitted
					lpRule = &ret[1]                         // right part is the new rule to solve
					s.solvedRules.ReplaceOrInsert(&(ret[0])) // Left part may be inserted in the new rules collection
					s.solvedRules.ReplaceOrInsert(hpRule)    // when the rule is splitted we also insert the higher priority rule
					fmt.Println(" >> append splitted fixed rule", ret[0], "and higher priority rule", hpRule)
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
}

// TODO fix comment
func (s *Solver) solveAndAppend(lpRule *SolverRule, collection *btree.BTreeG[*SolverRule]) {

	var newRules []*SolverRule
	//fmt.Println("------ Solving rule", lpRule.Name, "from", lpRule.From, "to", lpRule.To)

	// Loop over all rules in the collection and solve the current rule against each of them
	collection.Ascend(func(hpRule *SolverRule) bool {
		fmt.Println(" >> Solving fixed rule", lpRule.Name(), "vs", hpRule.Name())
		ret, _ := s.solveVsSingle(*lpRule, hpRule)
		switch len(ret) {
		case 0: // Rule deleted
			lpRule = nil
			return false
		case 1: // Rule Shifted or untouched
			lpRule = &ret[0]
		case 2: // Rule splitted
			newRules = append(newRules, &ret[0]) // Left part may be inserted in the new rules
			lpRule = &ret[1]                     // right part is the new rule to solve
		}
		return true
	})

	// Insert the last rule in the new rules collection
	if lpRule != nil && lpRule.Duration() > time.Duration(0) { //TODO check why we can get nil value here (bitonto test)
		newRules = append(newRules, lpRule)
	}

	// Effectively insert all parts of the resolved rules in the rules collection
	for i := range newRules {
		if newRules[i].Duration() > time.Duration(0) {
			collection.ReplaceOrInsert(newRules[i])
		}
	}
}

// sumAllSolvedRules returns the sum of all currently solved rules amounts and the total duration
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

// findFlatRateActivationTime returns the time when the flat rate rule is activated
// and the boolean indicating if the rule is activated
// The flat rate rule is activated when the sum of all rules from solver solvedRules in the timespan corresponding to the
// flat rate rule is greater than the activation amount. This point may be somewhere in the middle of the flat rate rule.
// The function returns the time when the flat rate rule is activated
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
