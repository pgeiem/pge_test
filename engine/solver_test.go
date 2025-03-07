package engine

import (
	"testing"
	"time"
)

type TestSolverRule struct {
	From                 time.Duration
	To                   time.Duration
	StartAmount          Amount
	EndAmount            Amount
	RuleName             string
	StartTimePolicy      StartTimePolicy
	RuleResolutionPolicy RuleResolutionPolicy
}

func (r TestSolverRule) Name() string {
	return r.RuleName
}

func (r TestSolverRule) RelativeTo(now time.Time) (time.Duration, time.Duration) {
	return r.From, r.To
}

func (r TestSolverRule) Policies() (StartTimePolicy, RuleResolutionPolicy) {
	return r.StartTimePolicy, r.RuleResolutionPolicy
}

func TestSolveVsSingle(t *testing.T) {

	tests := map[string]struct {
		lpRule               SolverRule
		hpRule               SolverRule
		ruleResolutionPolicy RuleResolutionPolicy
		expected             SolverRules
	}{
		// 0 - No conflict, Resolve policy
		"0-NoConflictAfter-Resolve": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute},
			hpRule:               SolverRule{From: 20 * time.Minute, To: 30 * time.Minute},
			ruleResolutionPolicy: ResolvePolicy,
			expected:             SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
		},
		// 1 - No conflict, Resolve policy
		"1-NoConflictBefore-Resolve": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute},
			hpRule:               SolverRule{From: 5 * time.Minute, To: 10 * time.Minute},
			ruleResolutionPolicy: ResolvePolicy,
			expected:             SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
		},
		// 2 - HPRule overlapping end of rule, Resolve policy
		"2-OverlapEnd-Resolve": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: ResolvePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute}},
		},
		// 3 - HPRule overlapping end of rule, Resolve policy
		"3-OverlapEndStep-Resolve": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 100},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: ResolvePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 100},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, StartAmount: 0, EndAmount: 0}},
		},
		// 4 - HPRule overlapping end of rule, Resolve policy
		"4-OverlapLinear-Resolve": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 200},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: ResolvePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 150},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, StartAmount: 0, EndAmount: 50}},
		},
		// 5 - HPRule overlapping beginning of rule, Resolve policy
		"5-OverlapBegining-Resolve": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute},
			hpRule:               SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			ruleResolutionPolicy: ResolvePolicy,
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 25 * time.Minute}},
		},
		// 6 - HPRule overlapping beginning of rule, Resolve policy
		"6-OverlapBeginingStep-Resolve": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 100},
			hpRule:               SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			ruleResolutionPolicy: ResolvePolicy,
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 25 * time.Minute, StartAmount: 100, EndAmount: 100}},
		},
		// 7 - HPRule overlapping beginning of rule, Resolve policy
		"7-OverlapBeginingLinear-Resolve": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 200},
			hpRule:               SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			ruleResolutionPolicy: ResolvePolicy,
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 25 * time.Minute, StartAmount: 100, EndAmount: 200}},
		},
		// 8 - HPRule overlapping middle of rule, Resolve policy
		"8-OverlapMiddle-Resolve": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 30 * time.Minute},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: ResolvePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute},
				SolverRule{From: 25 * time.Minute, To: 40 * time.Minute}},
		},
		// 9 - HPRule overlapping middle of rule, Resolve policy
		"9-OverlapMiddleStep-Resolve": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, StartAmount: 100, EndAmount: 100},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: ResolvePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 100},
				SolverRule{From: 25 * time.Minute, To: 40 * time.Minute, StartAmount: 0, EndAmount: 0}},
		},
		// 10 - HPRule overlapping middle of rule, Resolve policy
		"10-OverlapMiddleLinear-Resolve": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, StartAmount: 100, EndAmount: 200},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: ResolvePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 125},
				SolverRule{From: 25 * time.Minute, To: 40 * time.Minute, StartAmount: 0, EndAmount: 75}},
		},
		// 11 - No conflict Truncate policy
		"11-NoConflictAfter-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute},
			hpRule:               SolverRule{From: 20 * time.Minute, To: 30 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected:             SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
		},
		// 12 - No conflict Truncate policy
		"12-NoConflictBefore-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute},
			hpRule:               SolverRule{From: 5 * time.Minute, To: 10 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected:             SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
		},
		// 13 - HPRule overlapping end of rule, Truncate policy
		"13-OverlapEnd-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute}},
		},
		// 14 - HPRule overlapping end of rule, Truncate policy
		"14-OverlapEndStep-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 100},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 100}},
		},
		// 15 - HPRule overlapping end of rule, Truncate policy
		"15-OverlapLinear-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 200},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 150}},
		},
		// 16 - HPRule overlapping beginning of rule, Truncate policy
		"16-OverlapBegining-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute},
			hpRule:               SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 20 * time.Minute}},
		},
		// 17 - HPRule overlapping beginning of rule, Truncate policy
		"17-OverlapBeginingStep-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 100},
			hpRule:               SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 20 * time.Minute, StartAmount: 0, EndAmount: 0}},
		},
		// 18 - HPRule overlapping beginning of rule, Truncate policy
		"18-OverlapBeginingLinear-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 200},
			hpRule:               SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 20 * time.Minute, StartAmount: 0, EndAmount: 50}},
		},
		// 19 - HPRule overlapping middle of rule, Truncate policy
		"19-OverlapMiddle-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 30 * time.Minute},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute}},
		},
		// 20 - HPRule overlapping middle of rule, Truncate policy
		"20-OverlapMiddleStep-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, StartAmount: 100, EndAmount: 100},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 100},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, StartAmount: 0, EndAmount: 0}},
		},
		// 21 - HPRule overlapping middle of rule, Truncate policy
		"21-OverlapMiddleLinear-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, StartAmount: 100, EndAmount: 200},
			hpRule:               SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 125},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, StartAmount: 0, EndAmount: 25}},
		},
		// 22 - HPRule fully overlap rule, Truncate policy
		"22-FullOverlap-Truncate": {
			lpRule:               SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, StartAmount: 100, EndAmount: 200},
			hpRule:               SolverRule{From: 5 * time.Minute, To: 35 * time.Minute},
			ruleResolutionPolicy: TruncatePolicy,
			expected:             SolverRules{},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver(time.Now())
			solver.currentStartOffset = 45 * time.Minute
			out := solver.solveVsSingle(testcase.lpRule, testcase.hpRule, testcase.ruleResolutionPolicy)
			if len(out) != len(testcase.expected) {
				t.Errorf("solveVsSingle expected %v rules, got %v", len(testcase.expected), len(out))
			} else {
				for i := range out {

					if out[i].From != testcase.expected[i].From || out[i].To != testcase.expected[i].To {
						t.Errorf("solveVsSingle expected rule %v, got %v", testcase.expected, out)
					}
					if out[i].StartAmount != testcase.expected[i].StartAmount || out[i].EndAmount != testcase.expected[i].EndAmount {
						t.Errorf("solveVsSingle expected rule %v, got %v", testcase.expected, out)
					}
				}
			}
		})
	}
}

type ExpectedSolverRule struct {
	SolverRule
	Name string
}

func TestSolveAndAppend(t *testing.T) {

	DummyRule := TestSolverRule{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy}

	tests := map[string]struct {
		rules    []SolvableRule
		expected []ExpectedSolverRule
	}{
		// 0 - No conflict, empty rules§
		"0-NoConflictEmptyRules": {
			rules: []SolvableRule{
				&DummyRule,
			},
			expected: []ExpectedSolverRule{
				{Name: "A", SolverRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
			},
		},
		// 1 - No conflict, multiple rules
		"1-NoConflictMultipleRules": {
			rules: []SolvableRule{
				TestSolverRule{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "B", From: 20 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "C", From: 30 * time.Minute, To: 40 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: []ExpectedSolverRule{
				{Name: "A", SolverRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
				{Name: "B", SolverRule: SolverRule{From: 20 * time.Minute, To: 30 * time.Minute}},
				{Name: "C", SolverRule: SolverRule{From: 30 * time.Minute, To: 40 * time.Minute}},
			},
		},
		// 2 - Overlapping rules, Resolve policy
		"2-OverlapResolvePolicy": {
			rules: []SolvableRule{
				TestSolverRule{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: []ExpectedSolverRule{
				{Name: "A", SolverRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
				{Name: "B", SolverRule: SolverRule{From: 20 * time.Minute, To: 30 * time.Minute}},
				{Name: "C", SolverRule: SolverRule{From: 30 * time.Minute, To: 40 * time.Minute}},
			},
		},
		// 3 - Overlapping rules, Truncate policy
		"3-OverlapTruncatePolicy": {
			rules: []SolvableRule{
				TestSolverRule{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: []ExpectedSolverRule{
				{Name: "A", SolverRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
				{Name: "B", SolverRule: SolverRule{From: 20 * time.Minute, To: 25 * time.Minute}},
				{Name: "C", SolverRule: SolverRule{From: 25 * time.Minute, To: 35 * time.Minute}},
			},
		},
		// 4 - Overlapping rules, Delete policy
		"4-OverlapDeletePolicy": {
			rules: []SolvableRule{
				TestSolverRule{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: []ExpectedSolverRule{
				{Name: "A", SolverRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
				{Name: "C", SolverRule: SolverRule{From: 25 * time.Minute, To: 35 * time.Minute}},
			},
		},
		// 5 - Overlapping rules, mixed policies
		"5-OverlapMixedPolicies": {
			rules: []SolvableRule{
				TestSolverRule{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: []ExpectedSolverRule{
				{Name: "A", SolverRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
				{Name: "B", SolverRule: SolverRule{From: 20 * time.Minute, To: 30 * time.Minute}},
				{Name: "C", SolverRule: SolverRule{From: 30 * time.Minute, To: 35 * time.Minute}},
			},
		},
		// 6 - Overlapping rules, Truncate policy, Shiftable policy
		"6-OverlapTruncateShiftablePolicy": {
			rules: []SolvableRule{
				TestSolverRule{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "B", From: 30 * time.Minute, To: 40 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "C", From: 0 * time.Minute, To: 10 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				TestSolverRule{RuleName: "D", From: 0 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				TestSolverRule{RuleName: "Z", From: 60 * time.Minute, To: 90 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "E", From: 5 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				TestSolverRule{RuleName: "F", From: 0 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
			},
			expected: []ExpectedSolverRule{
				{Name: "C", SolverRule: SolverRule{From: 0 * time.Minute, To: 10 * time.Minute}},
				{Name: "A", SolverRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
				{Name: "D", SolverRule: SolverRule{From: 20 * time.Minute, To: 25 * time.Minute}},
				{Name: "E", SolverRule: SolverRule{From: 25 * time.Minute, To: 30 * time.Minute}},
				{Name: "B", SolverRule: SolverRule{From: 30 * time.Minute, To: 40 * time.Minute}},
				{Name: "F", SolverRule: SolverRule{From: 40 * time.Minute, To: 55 * time.Minute}},
				{Name: "Z", SolverRule: SolverRule{From: 60 * time.Minute, To: 90 * time.Minute}},
			},
		},
		// 7 - Overlapping rules, Truncate policy, Shiftable policy
		"7-OverlapTruncateShiftablePolicy": {
			rules: []SolvableRule{
				TestSolverRule{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "B", From: 30 * time.Minute, To: 40 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "C", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				TestSolverRule{RuleName: "D", From: 0 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				TestSolverRule{RuleName: "Z", From: 60 * time.Minute, To: 90 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				TestSolverRule{RuleName: "E", From: 5 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
			},
			expected: []ExpectedSolverRule{
				{Name: "C", SolverRule: SolverRule{From: 0 * time.Minute, To: 10 * time.Minute}},
				{Name: "A", SolverRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
				{Name: "D", SolverRule: SolverRule{From: 20 * time.Minute, To: 30 * time.Minute}},
				{Name: "B", SolverRule: SolverRule{From: 30 * time.Minute, To: 40 * time.Minute}},
				{Name: "D", SolverRule: SolverRule{From: 40 * time.Minute, To: 45 * time.Minute}},
				{Name: "E", SolverRule: SolverRule{From: 45 * time.Minute, To: 60 * time.Minute}},
				{Name: "Z", SolverRule: SolverRule{From: 60 * time.Minute, To: 90 * time.Minute}},
			},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver(time.Now())
			solver.Append(testcase.rules...)

			if solver.rules.Len() != len(testcase.expected) {
				t.Errorf("SolveAndAppend expected %v rules, got %v", len(testcase.expected), solver.rules.Len())
			} else {
				i := 0
				solver.rules.Ascend(func(rule SolverRule) bool {
					expected := testcase.expected[i]
					if rule.From != expected.From || rule.To != expected.To {
						t.Errorf("SolveAndAppend time error, expected from %s to %s, got  from %s to %s", expected.From, expected.To, rule.From, rule.To)
					}
					if rule.StartAmount != expected.StartAmount || rule.EndAmount != expected.EndAmount {
						t.Errorf("SolveAndAppend amount error, expected rule %v, got %v", expected, rule)
					}
					// Test name
					if rule.Name() != expected.Name {
						t.Errorf("SolveAndAppend mismatch names, expected name %s, got %s", expected.Name, rule.Name())
					}
					i++
					return true
				})
			}
		})
	}

	//Additional test to check that link to original rule is kept while running the solver
	t.Run("OriginalRuleLinkCheck", func(t *testing.T) {
		solver := NewSolver(time.Now())
		testcase := tests["0-NoConflictEmptyRules"]
		solver.Append(testcase.rules...)

		DummyRule.RuleName = "P" //Modify the original rule after the solver has been run
		if solver.rules.Len() != 1 {
			t.Errorf("SolveAndAppend expected %v rules, got %v", len(testcase.expected), solver.rules.Len())
		}
		// Check that the original rule change is reflected in the solver output
		rule, _ := solver.rules.Min()
		if rule.Name() != "P" {
			t.Errorf("SolveAndAppend mismatch names, expected name P, got %s", rule.Name())
		}
	})
}
