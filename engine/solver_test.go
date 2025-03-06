package engine

import (
	"testing"
	"time"
)

func TestSolveVsSingle(t *testing.T) {

	tests := map[string]struct {
		rule     SolverRule
		hpRule   SolverRule
		expected SolverRules
	}{
		// 0 - No conflict, Resolve policy
		"0-NoConflictAfter-Resolve": {
			rule:     SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
			hpRule:   SolverRule{From: 20 * time.Minute, To: 30 * time.Minute},
			expected: SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy}},
		},
		// 1 - No conflict, Resolve policy
		"1-NoConflictBefore-Resolve": {
			rule:     SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
			hpRule:   SolverRule{From: 5 * time.Minute, To: 10 * time.Minute},
			expected: SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy}},
		},
		// 2 - HPRule overlapping end of rule, Resolve policy
		"2-OverlapEnd-Resolve": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: ResolvePolicy}},
		},
		// 3 - HPRule overlapping end of rule, Resolve policy
		"3-OverlapEndStep-Resolve": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 100},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 100},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 0, EndAmount: 0}},
		},
		// 4 - HPRule overlapping end of rule, Resolve policy
		"4-OverlapLinear-Resolve": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 200},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 150},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 0, EndAmount: 50}},
		},
		// 5 - HPRule overlapping beginning of rule, Resolve policy
		"5-OverlapBegining-Resolve": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: ResolvePolicy}},
		},
		// 6 - HPRule overlapping beginning of rule, Resolve policy
		"6-OverlapBeginingStep-Resolve": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 100},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 100}},
		},
		// 7 - HPRule overlapping beginning of rule, Resolve policy
		"7-OverlapBeginingLinear-Resolve": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 200},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 200}},
		},
		// 8 - HPRule overlapping middle of rule, Resolve policy
		"8-OverlapMiddle-Resolve": {
			rule:   SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
				SolverRule{From: 25 * time.Minute, To: 40 * time.Minute, RuleResolutionPolicy: ResolvePolicy}},
		},
		// 9 - HPRule overlapping middle of rule, Resolve policy
		"9-OverlapMiddleStep-Resolve": {
			rule:   SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 100},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 100},
				SolverRule{From: 25 * time.Minute, To: 40 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 0, EndAmount: 0}},
		},
		// 10 - HPRule overlapping middle of rule, Resolve policy
		"10-OverlapMiddleLinear-Resolve": {
			rule:   SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 200},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 100, EndAmount: 125},
				SolverRule{From: 25 * time.Minute, To: 40 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartAmount: 0, EndAmount: 75}},
		},
		// 11 - No conflict Truncate policy
		"11-NoConflictAfter-Truncate": {
			rule:     SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy},
			hpRule:   SolverRule{From: 20 * time.Minute, To: 30 * time.Minute},
			expected: SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy}},
		},
		// 12 - No conflict Truncate policy
		"12-NoConflictBefore-Truncate": {
			rule:     SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy},
			hpRule:   SolverRule{From: 5 * time.Minute, To: 10 * time.Minute},
			expected: SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy}},
		},
		// 13 - HPRule overlapping end of rule, Truncate policy
		"13-OverlapEnd-Truncate": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy}},
		},
		// 14 - HPRule overlapping end of rule, Truncate policy
		"14-OverlapEndStep-Truncate": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 100, EndAmount: 100},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 100, EndAmount: 100}},
		},
		// 15 - HPRule overlapping end of rule, Truncate policy
		"15-OverlapLinear-Truncate": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 100, EndAmount: 200},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 100, EndAmount: 150}},
		},
		// 16 - HPRule overlapping beginning of rule, Truncate policy
		"16-OverlapBegining-Truncate": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy}},
		},
		// 17 - HPRule overlapping beginning of rule, Truncate policy
		"17-OverlapBeginingStep-Truncate": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 100, EndAmount: 100},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 0, EndAmount: 0}},
		},
		// 18 - HPRule overlapping beginning of rule, Truncate policy
		"18-OverlapBeginingLinear-Truncate": {
			rule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 100, EndAmount: 200},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 0, EndAmount: 50}},
		},
		// 19 - HPRule overlapping middle of rule, Truncate policy
		"19-OverlapMiddle-Truncate": {
			rule:   SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: TruncatePolicy}},
		},
		// 20 - HPRule overlapping middle of rule, Truncate policy
		"20-OverlapMiddleStep-Truncate": {
			rule:   SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 100, EndAmount: 100},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 100, EndAmount: 100},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 0, EndAmount: 0}},
		},
		// 21 - HPRule overlapping middle of rule, Truncate policy
		"21-OverlapMiddleLinear-Truncate": {
			rule:   SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 100, EndAmount: 200},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 100, EndAmount: 125},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 0, EndAmount: 25}},
		},
		// 22 - HPRule fully overlap rule, Truncate policy
		"22-FullOverlap-Truncate": {
			rule:     SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartAmount: 100, EndAmount: 200},
			hpRule:   SolverRule{From: 5 * time.Minute, To: 35 * time.Minute},
			expected: SolverRules{},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver(time.Now())
			solver.currentStartOffset = 45 * time.Minute
			out := solver.solveVsSingle(testcase.rule, testcase.hpRule)
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
					if out[i].StartTimePolicy != testcase.expected[i].StartTimePolicy || out[i].RuleResolutionPolicy != testcase.expected[i].RuleResolutionPolicy {
						t.Errorf("solveVsSingle expected rule %v, got %v", testcase.expected, out)
					}
				}
			}
		})
	}
}
func TestSolveAndAppend(t *testing.T) {

	tests := map[string]struct {
		rules    []SolverRule
		expected []SolverRule
	}{
		// 0 - No conflict, empty rules§
		"0-NoConflictEmptyRules": {
			rules: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
			},
		},
		// 1 - No conflict, multiple rules
		"1-NoConflictMultipleRules": {
			rules: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{Name: "B", From: 20 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{Name: "C", From: 30 * time.Minute, To: 40 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{Name: "B", From: 20 * time.Minute, To: 30 * time.Minute},
				{Name: "C", From: 30 * time.Minute, To: 40 * time.Minute},
			},
		},
		// 2 - Overlapping rules, Resolve policy
		"2-OverlapResolvePolicy": {
			rules: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{Name: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{Name: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{Name: "B", From: 20 * time.Minute, To: 30 * time.Minute},
				{Name: "C", From: 30 * time.Minute, To: 40 * time.Minute},
			},
		},
		// 3 - Overlapping rules, Truncate policy
		"3-OverlapTruncatePolicy": {
			rules: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{Name: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{Name: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{Name: "B", From: 20 * time.Minute, To: 25 * time.Minute},
				{Name: "C", From: 25 * time.Minute, To: 35 * time.Minute},
			},
		},
		// 4 - Overlapping rules, Delete policy
		"4-OverlapDeletePolicy": {
			rules: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				{Name: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				{Name: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{Name: "C", From: 25 * time.Minute, To: 35 * time.Minute},
			},
		},
		// 5 - Overlapping rules, mixed policies
		"5-OverlapMixedPolicies": {
			rules: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				{Name: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{Name: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{Name: "B", From: 20 * time.Minute, To: 30 * time.Minute},
				{Name: "C", From: 30 * time.Minute, To: 35 * time.Minute},
			},
		},
		// 6 - Overlapping rules, Truncate policy, Shiftable policy
		"6-OverlapTruncateShiftablePolicy": {
			rules: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{Name: "B", From: 30 * time.Minute, To: 40 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{Name: "C", From: 0 * time.Minute, To: 10 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				{Name: "D", From: 0 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				{Name: "Z", From: 60 * time.Minute, To: 90 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{Name: "E", From: 5 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				{Name: "F", From: 0 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
			},
			expected: []SolverRule{
				{Name: "C", From: 0 * time.Minute, To: 10 * time.Minute},
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{Name: "D", From: 20 * time.Minute, To: 25 * time.Minute},
				{Name: "E", From: 25 * time.Minute, To: 30 * time.Minute},
				{Name: "B", From: 30 * time.Minute, To: 40 * time.Minute},
				{Name: "F", From: 40 * time.Minute, To: 55 * time.Minute},
				{Name: "Z", From: 60 * time.Minute, To: 90 * time.Minute},
			},
		},
		// 7 - Overlapping rules, Truncate policy, Shiftable policy
		"7-OverlapTruncateShiftablePolicy": {
			rules: []SolverRule{
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{Name: "B", From: 30 * time.Minute, To: 40 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{Name: "C", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				{Name: "D", From: 0 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				{Name: "Z", From: 60 * time.Minute, To: 90 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{Name: "E", From: 5 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
			},
			expected: []SolverRule{
				{Name: "C", From: 0 * time.Minute, To: 10 * time.Minute},
				{Name: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{Name: "D", From: 20 * time.Minute, To: 30 * time.Minute},
				{Name: "B", From: 30 * time.Minute, To: 40 * time.Minute},
				{Name: "D", From: 40 * time.Minute, To: 45 * time.Minute},
				{Name: "E", From: 45 * time.Minute, To: 60 * time.Minute},
				{Name: "Z", From: 60 * time.Minute, To: 90 * time.Minute},
			},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver(time.Now())
			for _, rule := range testcase.rules {
				solver.SolveAndAppend(rule)
			}
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
					if rule.Name != expected.Name {
						t.Errorf("SolveAndAppend mismatch names, expected rule %s, got %s", expected.Name, rule.Name)
					}
					i++
					return true
				})
			}
		})
	}
}
