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
)

//TOOD: merge RuleResolutionPolicy with StartTimePolicy as shiftable is usefull only with truncate ?

// Define the solver rule
type SolverRule struct {
	// Starting point in time
	From time.Duration
	// End point in time
	To time.Duration
	// Amount in cents at the beginning of the rule segment (non 0 values are step)
	StartAmount Amount
	// Amount in cents at the end of the rule segment
	EndAmount Amount
	// Trace buffer for debugging all rule changes
	Trace string
	// Rule type reported to output used for tariff details
	Type string
	// OriginalRule is the original rule, from which the SolverRule was derived
	OriginalRule *SolvableRule
}

func (rule SolverRule) Duration() time.Duration {
	return rule.To - rule.From
}

func (rule SolverRule) Name() string {
	if rule.OriginalRule == nil {
		return "nil"
	}

	return (*rule.OriginalRule).Name()
}

func (rule SolverRule) String() string {
	return fmt.Sprintf("%s(%s) From %s to %s,\tAmount %d to %d,\tType %s",
		rule.Name(), rule.Trace, rule.From.String(), rule.To.String(), rule.StartAmount, rule.EndAmount, rule.Type)
}

// Shift the rule to the new start time, the new rule is returned and current rule is not changed
func (rule SolverRule) Shift(from time.Duration) SolverRule {
	rule.To = from + rule.Duration()
	rule.From = from
	rule.Trace += "_S"
	return rule
}

func (rule SolverRule) TruncateAfter(after time.Duration) SolverRule {
	ruleA := rule
	ruleA.To = after
	ruleA.Trace += "_TA"
	if rule.Duration() != time.Duration(0) {
		ruleA.EndAmount = Amount(int64(rule.EndAmount-rule.StartAmount)*int64(ruleA.Duration())/int64(rule.Duration())) + rule.StartAmount
	}
	return ruleA
}

func (rule SolverRule) TruncateBefore(before time.Duration) SolverRule {
	ruleA := rule
	ruleA.From = before
	ruleA.Trace += "_TB"
	ruleA.StartAmount = 0

	if rule.Duration() != time.Duration(0) {
		ruleA.EndAmount = Amount(int64(rule.EndAmount-rule.StartAmount) * int64(ruleA.Duration()) / int64(rule.Duration()))
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
	ruleB.Trace += "_B"
	ruleB.StartAmount = 0
	ruleB.EndAmount = rule.EndAmount - ruleA.EndAmount

	return SolverRules{ruleA, ruleB}
}

// Define a collection of solver rule
type SolverRules []SolverRule

type SolvableRule interface {
	// Name returns the name of the rule
	Name() string
	// RelativeTo returns the rule relative start/end to a given time
	RelativeTo(now time.Time) (time.Duration, time.Duration)
	// Policies returns the policies of the rule
	Policies() (StartTimePolicy, RuleResolutionPolicy)
}

type Solver struct {
	now                time.Time
	rules              *btree.BTreeG[SolverRule]
	currentStartOffset time.Duration
}

func NewSolver(now time.Time) *Solver {

	// Sorting function for B-Tree storing all solved rules segments
	Less := func(i, j SolverRule) bool {
		return i.From < j.From
	}

	return &Solver{
		now:   now,
		rules: btree.NewG(2, Less),
	}
}

// Add some SolvableRule to the solver
func (s *Solver) Append(rules ...SolvableRule) {
	for i := range rules { //Don't use iterator so we can get pointer on original rule
		var r SolverRule
		r.From, r.To = rules[i].RelativeTo(s.now)
		r.OriginalRule = &rules[i]
		startTimePolicy, ruleResolutionPolicy := rules[i].Policies()
		s.solveAndAppend(r, startTimePolicy, ruleResolutionPolicy)
	}
}

// Solve the rule against an Higer Priority Rule resolving the conflict according to rule policy
// a collection of new rules containing 0, 1, or 2 rules is returned and current rule is not changed
func (s *Solver) solveVsSingle(lprule, hpRule SolverRule, ruleResolutionPolicy RuleResolutionPolicy) SolverRules {

	// trivial case, both rules don't overlap
	if (hpRule.To <= lprule.From) ||
		(hpRule.From >= lprule.To) {
		return SolverRules{lprule}
	}

	switch ruleResolutionPolicy {

	// both rules overlap at least slightly, if policy is 'remove' then remove the low priority rule
	case DeletePolicy:
		fmt.Println("## DeletePolicy", lprule.Name(), "vs", hpRule.Name())
		return SolverRules{}

	// both rules overlap at least slightly, if policy is 'resolve' then try to split rule to fill the holes
	case ResolvePolicy:
		fmt.Println("## ResolvePolicy", lprule.Name(), "vs", hpRule.Name())

		// high priority rule is before low priority rule, then low priority rule is simply shifted
		if hpRule.From <= lprule.From {
			return SolverRules{lprule.Shift(hpRule.To)}
		} else {
			// high priority rule is after or in middle of low priority rule, then low priority rule is split and shifted
			return lprule.Split(hpRule.From, hpRule.To)
		}

	// both rules overlap at least slightly, if policy is 'truncate' then truncate overlapping rule
	case TruncatePolicy:
		fmt.Println("## TruncatePolicy", lprule.Name(), "vs", hpRule.Name())

		// high priority rule is partially after low priority rule, then low priority rule end is truncated
		if hpRule.From >= lprule.From && hpRule.To >= lprule.To {
			return SolverRules{lprule.TruncateAfter(hpRule.From)}
		}

		// high priority rule is partially before low priority rule, then low priority rule end is truncated
		if hpRule.From <= lprule.From && hpRule.To <= lprule.To {
			return SolverRules{lprule.TruncateBefore(hpRule.To)}
		}

		// high priority rule completely overlap low priority rule, then remove the low priority rule
		if hpRule.From <= lprule.From && hpRule.To >= lprule.To {
			return SolverRules{}
		}

		// high priority rule is in middle of low priority rule, then low priority rule middle is truncated
		if hpRule.From >= lprule.From && hpRule.To <= lprule.To {
			return lprule.TruncateBetween(hpRule.From, hpRule.To)
		}

	default:
		panic(fmt.Errorf("unhandled solving policy %v between %v and %v", ruleResolutionPolicy, lprule, hpRule))
	}

	return SolverRules{}
}

// Solve the rule against a collection of Higer Priority Rule resolving the conflict according to rules policy
// a collection of new rules is returned and current rule is not changed
func (s *Solver) solveAndAppend(lpRule SolverRule, startTimePolicy StartTimePolicy, ruleResolutionPolicy RuleResolutionPolicy) {

	var newRules SolverRules
	//fmt.Println("------ Solving rule", lpRule.Name, "from", lpRule.From, "to", lpRule.To)

	// Shift the rule if needed using the current start offset
	if startTimePolicy == ShiftablePolicy {
		//fmt.Println("### Shifting rule", lpRule.Name, "from", lpRule.From, "to", lpRule.From+s.currentStartOffset)
		lpRule = lpRule.Shift(s.currentStartOffset)
		s.currentStartOffset += lpRule.Duration()
	}

	// Loop over all rules in the collection and solve the current rule against each of them
	s.rules.Ascend(func(hpRule SolverRule) bool {
		fmt.Println("### Solving rule", lpRule.Name(), "vs", hpRule.Name())
		ret := s.solveVsSingle(lpRule, hpRule, ruleResolutionPolicy)
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
		return true
	})

	// Insert the last rule part in the new rules collection
	newRules = append(newRules, lpRule)

	// Effectively insert all parts of the resolved rules in the rules collection
	for _, rule := range newRules {
		if rule.Duration() > time.Duration(0) {
			s.rules.ReplaceOrInsert(rule)
		}
	}
}

/*
// Fully solve the rules collection
func (rules *SolverRules) Solve() {

	var solvedRules SolverRules
	for _, lpRule := range *rules {
		solvedRules = append(solvedRules, lpRule.SolveVsMany(solvedRules)...)
	}

	solvedRules.Sort()
	*rules = solvedRules
	rules.RemoveZeroDuration()
	rules.Filter()
}
*/
