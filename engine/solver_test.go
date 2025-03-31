package engine

import (
	"testing"
	"time"
)

func TestSolveVsSingle(t *testing.T) {

	tests := map[string]struct {
		lpRule   SolverRule
		hpRule   SolverRule
		expected SolverRules
	}{
		// 0 - No conflict, Resolve policy
		"0-NoConflictAfter-Resolve": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 20 * time.Minute, To: 30 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
			},
		},
		// 1 - No conflict, Resolve policy
		"1-NoConflictBefore-Resolve": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Minute, To: 10 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
			},
		},
		// 2 - HPRule overlapping end of rule, Resolve policy
		"2-OverlapEnd-Resolve": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}},
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 30 * time.Minute}},
			},
		},
		// 3 - HPRule overlapping end of rule, Resolve policy
		"3-OverlapEndStep-Resolve": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}, StartAmount: 100, EndAmount: 100},
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 30 * time.Minute}, StartAmount: 0, EndAmount: 0},
			},
		},
		// 4 - HPRule overlapping end of rule, Resolve policy
		"4-OverlapLinear-Resolve": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}, StartAmount: 100, EndAmount: 150},
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 30 * time.Minute}, StartAmount: 0, EndAmount: 50},
			},
		},
		// 5 - HPRule overlapping beginning of rule, Resolve policy
		"5-OverlapBegining-Resolve": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Minute, To: 15 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			},
		},
		// 6 - HPRule overlapping beginning of rule, Resolve policy
		"6-OverlapBeginingStep-Resolve": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Minute, To: 15 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}, StartAmount: 100, EndAmount: 100}},
		},

		// 7 - HPRule overlapping beginning of rule, Resolve policy
		"7-OverlapBeginingLinear-Resolve": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Minute, To: 15 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}, StartAmount: 100, EndAmount: 200}},
		},

		// 8 - HPRule overlapping middle of rule, Resolve policy
		"8-OverlapMiddle-Resolve": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 30 * time.Minute}, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}},
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 40 * time.Minute}},
			},
		},
		// 9 - HPRule overlapping middle of rule, Resolve policy
		"9-OverlapMiddleStep-Resolve": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 30 * time.Minute}, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}, StartAmount: 100, EndAmount: 100},
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 40 * time.Minute}, StartAmount: 0, EndAmount: 0}},
		},
		// 10 - HPRule overlapping middle of rule, Resolve policy
		"10-OverlapMiddleLinear-Resolve": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 30 * time.Minute}, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}, StartAmount: 100, EndAmount: 125},
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 40 * time.Minute}, StartAmount: 0, EndAmount: 75},
			},
		},
		// 11 - No conflict Truncate policy
		"11-NoConflictAfter-Truncate": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 20 * time.Minute, To: 30 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
			},
		},
		// 12 - No conflict Truncate policy
		"12-NoConflictBefore-Truncate": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Minute, To: 10 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
			},
		},
		// 13 - HPRule overlapping end of rule, Truncate policy
		"13-OverlapEnd-Truncate": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}},
			},
		},
		// 14 - HPRule overlapping end of rule, Truncate policy
		"14-OverlapEndStep-Truncate": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: TruncatePolicy}, hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}, StartAmount: 100, EndAmount: 100},
			},
		},
		// 15 - HPRule overlapping end of rule, Truncate policy
		"15-OverlapLinear-Truncate": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}, StartAmount: 100, EndAmount: 150},
			},
		},
		// 16 - HPRule overlapping beginning of rule, Truncate policy
		"16-OverlapBegining-Truncate": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Minute, To: 15 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 20 * time.Minute}},
			},
		},
		// 17 - HPRule overlapping beginning of rule, Truncate policy
		"17-OverlapBeginingStep-Truncate": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Minute, To: 15 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 20 * time.Minute}, StartAmount: 0, EndAmount: 0},
			},
		},
		// 18 - HPRule overlapping beginning of rule, Truncate policy
		"18-OverlapBeginingLinear-Truncate": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Minute, To: 15 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 20 * time.Minute}, StartAmount: 0, EndAmount: 50},
			},
		},
		// 19 - HPRule overlapping middle of rule, Truncate policy
		"19-OverlapMiddle-Truncate": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 30 * time.Minute}, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}},
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 30 * time.Minute}},
			},
		},
		// 20 - HPRule overlapping middle of rule, Truncate policy
		"20-OverlapMiddleStep-Truncate": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 30 * time.Minute}, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}, StartAmount: 100, EndAmount: 100},
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 30 * time.Minute}, StartAmount: 0, EndAmount: 0}},
		},
		// 21 - HPRule overlapping middle of rule, Truncate policy
		"21-OverlapMiddleLinear-Truncate": {
			lpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 30 * time.Minute}, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			expected: SolverRules{
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}, StartAmount: 100, EndAmount: 125},
				SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 30 * time.Minute}, StartAmount: 0, EndAmount: 25}},
		},
		// 22 - HPRule fully overlap rule, Truncate policy
		"22-FullOverlap-Truncate": {
			lpRule:   SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 30 * time.Minute}, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: TruncatePolicy},
			hpRule:   SolverRule{RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Minute, To: 35 * time.Minute}},
			expected: SolverRules{},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver()
			solver.SetWindow(time.Now(), time.Duration(48*time.Hour))
			solver.currentRelativeStartOffset = 45 * time.Minute
			out, _ := solver.solveVsSingle(testcase.lpRule, &testcase.hpRule)
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

func TestAppend(t *testing.T) {

	tests := map[string]struct {
		rules    SolverRules
		expected SolverRules
	}{
		// 0 - No conflict, empty rules§
		"0-NoConflictEmptyRules": {
			rules: SolverRules{
				SolverRule{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
			},
		},
		// 1 - No conflict, multiple rules
		"1-NoConflictMultipleRules": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 20 * time.Minute, To: 30 * time.Minute}, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 30 * time.Minute, To: 40 * time.Minute}, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 20 * time.Minute, To: 30 * time.Minute}},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 30 * time.Minute, To: 40 * time.Minute}},
			},
		},
		// 2 - Overlapping rules, Resolve policy
		"2-OverlapResolvePolicy": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 35 * time.Minute}, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 20 * time.Minute, To: 30 * time.Minute}},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 30 * time.Minute, To: 40 * time.Minute}},
			},
		},
		// 3 - Overlapping rules, Truncate policy
		"3-OverlapTruncatePolicy": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 35 * time.Minute}, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 20 * time.Minute, To: 25 * time.Minute}},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 35 * time.Minute}},
			},
		},
		// 4 - Overlapping rules, Delete policy
		"4-OverlapDeletePolicy": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 35 * time.Minute}, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 35 * time.Minute}},
			},
		},
		// 5 - Overlapping rules, mixed policies
		"5-OverlapMixedPolicies": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 35 * time.Minute}, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 20 * time.Minute, To: 30 * time.Minute}},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 30 * time.Minute, To: 35 * time.Minute}},
			},
		},
		// 6 - Overlapping rules, Truncate policy, Shiftable policy
		"6-OverlapTruncateShiftablePolicy": {
			rules: SolverRules{
				NewAbsoluteNonPaying("A", RelativeTimeSpan{10 * time.Minute, 20 * time.Minute}, MetaData{}),
				NewAbsoluteNonPaying("B", RelativeTimeSpan{30 * time.Minute, 40 * time.Minute}, MetaData{}),
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Minute, To: 10 * time.Minute}, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				{RuleName: "D", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Minute, To: 15 * time.Minute}, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				NewAbsoluteNonPaying("Z", RelativeTimeSpan{60 * time.Minute, 90 * time.Minute}, MetaData{}),
				{RuleName: "E", RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Minute, To: 15 * time.Minute}, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				{RuleName: "F", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Minute, To: 20 * time.Minute}, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
			},
			expected: SolverRules{
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Minute, To: 10 * time.Minute}},
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
				{RuleName: "D", RelativeTimeSpan: RelativeTimeSpan{From: 20 * time.Minute, To: 25 * time.Minute}},
				{RuleName: "E", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 30 * time.Minute}},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 30 * time.Minute, To: 40 * time.Minute}},
				{RuleName: "F", RelativeTimeSpan: RelativeTimeSpan{From: 40 * time.Minute, To: 55 * time.Minute}},
				{RuleName: "Z", RelativeTimeSpan: RelativeTimeSpan{From: 60 * time.Minute, To: 90 * time.Minute}},
			},
		},
		// 7 - Overlapping rules, Truncate policy, Shiftable policy
		"7-OverlapTruncateShiftablePolicy": {
			rules: SolverRules{
				NewAbsoluteNonPaying("A", RelativeTimeSpan{10 * time.Minute, 20 * time.Minute}, MetaData{}),
				NewAbsoluteNonPaying("B", RelativeTimeSpan{30 * time.Minute, 40 * time.Minute}, MetaData{}),
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				{RuleName: "D", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Minute, To: 35 * time.Minute}, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				NewAbsoluteNonPaying("Z", RelativeTimeSpan{60 * time.Minute, 90 * time.Minute}, MetaData{}),
				{RuleName: "E", RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Minute, To: 25 * time.Minute}, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
			},
			expected: SolverRules{
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Minute, To: 10 * time.Minute}},
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
				{RuleName: "D", RelativeTimeSpan: RelativeTimeSpan{From: 20 * time.Minute, To: 30 * time.Minute}},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 30 * time.Minute, To: 40 * time.Minute}},
				{RuleName: "D", RelativeTimeSpan: RelativeTimeSpan{From: 40 * time.Minute, To: 45 * time.Minute}},
				{RuleName: "E", RelativeTimeSpan: RelativeTimeSpan{From: 45 * time.Minute, To: 60 * time.Minute}},
				{RuleName: "Z", RelativeTimeSpan: RelativeTimeSpan{From: 60 * time.Minute, To: 90 * time.Minute}},
			},
		},
		// 8 - Overlapping rules, Truncate policy, Shiftable policy
		"8-OverlapTruncateShiftablePolicy": {
			rules: SolverRules{
				NewAbsoluteLinearRule("A", RelativeTimeSpan{2 * time.Hour, 4 * time.Hour}, 2.0, MetaData{}),
				NewRelativeLinearRule("Hourly", 10*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 2 * time.Hour}, StartAmount: 0, EndAmount: 2.0},
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 2 * time.Hour, To: 4 * time.Hour}, StartAmount: 0, EndAmount: 4.0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 4 * time.Hour, To: 12 * time.Hour}, StartAmount: 0, EndAmount: 8.0},
			},
		},
		// 10 - Overlapping multiple calendar flatrate with linear in the middle and multiple flatrate activation
		"10-MultipleCalendarFlatrate": {
			rules: SolverRules{
				NewAbsoluteFlatRateRule("Morning", RelativeTimeSpan{2 * time.Hour, 6 * time.Hour}, 3.0, MetaData{}),
				NewAbsoluteFlatRateRule("Evening", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
				NewRelativeLinearRule("Hourly", 10*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 3 * time.Hour}, StartAmount: 0, EndAmount: 3.0},
				{RuleName: "Morning", RelativeTimeSpan: RelativeTimeSpan{From: 3 * time.Hour, To: 6 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 6 * time.Hour, To: 9 * time.Hour}, StartAmount: 0, EndAmount: 3.0},
				{RuleName: "Evening", RelativeTimeSpan: RelativeTimeSpan{From: 9 * time.Hour, To: 11 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 11 * time.Hour, To: 15 * time.Hour}, StartAmount: 0, EndAmount: 4.0},
			},
		},

		// 12 - Overlapping multiple calendar flatrate with linear in the middle and multiple flatrate activation
		// This one is not working because the daily flat rate amount should not be summed with others flatrates
		// Overlaped flatrates are not supported
		/*"12-MultipleCalendarFlatrate": {
			rules: SolverRules{
				NewAbsoluteFlatRateRule("Morning", 2*time.Hour, 6*time.Hour, 3.0),
				NewAbsoluteFlatRateRule("Evening", 7*time.Hour, 11*time.Hour, 3.0),
				NewAbsoluteFlatRateRule("Daily", 2*time.Hour, 14*time.Hour, 7.0),
				NewRelativeLinearRule("Hourly", 10*time.Hour, 1.0),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 3 * time.Hour, StartAmount: 0, EndAmount: 3.0},
				{RuleName: "Morning", RelativeTimeSpan: RelativeTimeSpan{From: 3 * time.Hour, To: 6 * time.Hour, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 6 * time.Hour, To: 9 * time.Hour, StartAmount: 0, EndAmount: 3.0},
				{RuleName: "Evening", RelativeTimeSpan: RelativeTimeSpan{From: 9 * time.Hour, To: 11 * time.Hour, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 11 * time.Hour, To: 12 * time.Hour, StartAmount: 0, EndAmount: 4.0},
				{RuleName: "Daily", RelativeTimeSpan: RelativeTimeSpan{From: 12 * time.Hour, To: 14 * time.Hour, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 14 * time.Hour, To: 17 * time.Hour, StartAmount: 0, EndAmount: 4.0},
			},
		},
		*/
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver()
			solver.SetWindow(time.Now(), time.Duration(48*time.Hour))
			solver.AppendMany(testcase.rules...)

			if solver.rules.Len() != len(testcase.expected) {
				t.Errorf("SolveAndAppend expected %v rules, got %v", len(testcase.expected), solver.rules.Len())
			} else {
				i := 0
				solver.rules.Ascend(func(rule *SolverRule) bool {
					expected := testcase.expected[i]
					if rule.From != expected.From || rule.To != expected.To {
						t.Errorf("SolveAndAppend(%d) time error, expected rule %v, got %v", i, expected, rule)
					}
					if rule.StartAmount != expected.StartAmount || rule.EndAmount != expected.EndAmount {
						t.Errorf("SolveAndAppend(%d) amount error, expected rule %v, got %v", i, expected, rule)
					}
					// Test name
					if rule.Name() != expected.Name() {
						t.Errorf("SolveAndAppend(%d) mismatch names, expected name %s, got %s", i, expected.Name(), rule.Name())
					}
					i++
					return true
				})
				//Testing output
				/*out := solver.GenerateOutput(true)
				fmt.Println(out)
				data, err := out.ToJson()
				if err != nil {
					t.Errorf("SolveAndAppend(%d) error converting to json %v", i, err)
				} else {
					fmt.Println(string(data))
				}*/
			}
		})
	}
}

func TestFindIntersectPositionFlatRate(t *testing.T) {
	tests := map[string]struct {
		relativeRule          SolverRule
		flatRateRule          SolverRule
		realtiveStartOffset   time.Duration
		relativeAmountOffset  Amount
		activatedFlatRatesSum Amount
		expectedIntersect     bool
		expected              time.Duration
	}{
		// 0 - No intersection below
		"0-NoIntersectionBelow": {
			flatRateRule:          NewAbsoluteFlatRateRule("A", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 4.0, MetaData{}),
			relativeRule:          NewRelativeLinearRule("B", 3*time.Hour, 1.0, MetaData{}),
			realtiveStartOffset:   8 * time.Hour,
			relativeAmountOffset:  0,
			activatedFlatRatesSum: 0,
			expectedIntersect:     false,
			expected:              0,
		},
		// 1 - No intersection below
		"1-NoIntersectionAbove": {
			flatRateRule:          NewAbsoluteFlatRateRule("A", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 4.0, MetaData{}),
			relativeRule:          NewRelativeLinearRule("B", 3*time.Hour, 1.0, MetaData{}),
			realtiveStartOffset:   8 * time.Hour,
			relativeAmountOffset:  4.0,
			activatedFlatRatesSum: 0,
			expectedIntersect:     false,
			expected:              0,
		},
		// 2 - No intersection before
		"2-NoIntersectionBefore": {
			flatRateRule:          NewAbsoluteFlatRateRule("A", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 4.0, MetaData{}),
			relativeRule:          NewRelativeLinearRule("B", 3*time.Hour, 1.0, MetaData{}),
			realtiveStartOffset:   2 * time.Hour,
			relativeAmountOffset:  2.0,
			activatedFlatRatesSum: 0,
			expectedIntersect:     false,
			expected:              0,
		},
		// 3 - No intersection after
		"3-NoIntersectionAfter": {
			flatRateRule:          NewAbsoluteFlatRateRule("A", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 4.0, MetaData{}),
			relativeRule:          NewRelativeLinearRule("B", 3*time.Hour, 1.0, MetaData{}),
			realtiveStartOffset:   12 * time.Hour,
			relativeAmountOffset:  4.0,
			activatedFlatRatesSum: 0,
			expectedIntersect:     false,
			expected:              0,
		},
		// 4 - Intersection in the middle
		"4-IntersectionInMiddle": {
			flatRateRule:          NewAbsoluteFlatRateRule("A", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 4.0, MetaData{}),
			relativeRule:          NewRelativeLinearRule("B", 4*time.Hour, 1.0, MetaData{}),
			realtiveStartOffset:   8 * time.Hour,
			relativeAmountOffset:  2.0,
			activatedFlatRatesSum: 0,
			expectedIntersect:     true,
			expected:              10 * time.Hour,
		},
		// 5 - Intersection in the middle
		"5-IntersectionInMiddle": {
			flatRateRule:          NewAbsoluteFlatRateRule("A", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 2.0, MetaData{}),
			relativeRule:          NewRelativeLinearRule("B", 4*time.Hour, 1.0, MetaData{}),
			realtiveStartOffset:   8 * time.Hour,
			relativeAmountOffset:  2.0,
			activatedFlatRatesSum: 2.0,
			expectedIntersect:     true,
			expected:              10 * time.Hour,
		},
		// 6 - Intersection in the middle
		"6-IntersectionInMiddle": {
			flatRateRule:          NewAbsoluteFlatRateRule("A", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 2.0, MetaData{}),
			relativeRule:          NewRelativeLinearRule("B", 4*time.Hour, 1.0, MetaData{}),
			realtiveStartOffset:   6 * time.Hour,
			relativeAmountOffset:  2.0,
			activatedFlatRatesSum: 2.0,
			expectedIntersect:     true,
			expected:              8 * time.Hour,
		},
		// 7 - Intersection in the beginning
		"7-IntersectionInBeginning": {
			flatRateRule:          NewAbsoluteFlatRateRule("A", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 2.0, MetaData{}),
			relativeRule:          NewRelativeLinearRule("B", 3*time.Hour, 1.0, MetaData{}),
			realtiveStartOffset:   9 * time.Hour,
			relativeAmountOffset:  3.0,
			activatedFlatRatesSum: 2.0,
			expectedIntersect:     true,
			expected:              10 * time.Hour,
		},
		// 8 - Intersection in the beginning
		"8-IntersectionInBeginning": {
			flatRateRule:          NewAbsoluteFlatRateRule("A", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 2.0, MetaData{}),
			relativeRule:          NewRelativeLinearRule("B", 4*time.Hour, 1.0, MetaData{}),
			realtiveStartOffset:   4 * time.Hour,
			relativeAmountOffset:  0,
			activatedFlatRatesSum: 2.0,
			expectedIntersect:     true,
			expected:              8 * time.Hour,
		},
		// 9 - Intersection in the end must be considered as not intersection
		"9-IntersectionInEnd": {
			flatRateRule:          NewAbsoluteFlatRateRule("A", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 2.0, MetaData{}),
			relativeRule:          NewRelativeLinearRule("B", 4*time.Hour, 1.0, MetaData{}),
			realtiveStartOffset:   12 * time.Hour,
			relativeAmountOffset:  0,
			activatedFlatRatesSum: 4.0,
			expectedIntersect:     false,
			expected:              0,
		},
		// 10 - Intersection
		//>> FindIntersectPositionFlatRate Hourly vs Evening => 10h0m0s | Hourly(6h0m0s -> 13h0m0s; 0.000 -> 7.000) Evening(7h0m0s -> 11h0m0s; 0.000 -> 0.000) | 3.000 3.000 0s
		//>> FindIntersectPositionFlatRate Hourly vs Evening => 9h0m0s | Hourly(6h0m0s -> 13h0m0s; 0.000 -> 7.000) Evening(7h0m0s -> 11h0m0s; 0.000 -> 0.000) | 3.000 3.000 0s

		// >> FindIntersectPositionFlatRate Hourly vs Evening => 10h0m0s | Hourly(6h0m0s -> 13h0m0s; 0.000 -> 7.000) Evening(7h0m0s -> 11h0m0s; 0.000 -> 0.000) | 3.000 3.000 0s | 7.000 3.000 10.000 6h0m0s
		// >> FindIntersectPositionFlatRate Hourly vs Evening => 9h0m0s | Hourly(6h0m0s -> 13h0m0s; 0.000 -> 7.000) Evening(7h0m0s -> 11h0m0s; 0.000 -> 0.000) | 3.000 3.000 0s | 6.000 3.000 10.000 6h0m0s
		"10-IntersectionInMiddle": {
			flatRateRule:          NewAbsoluteFlatRateRule("Evening", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			relativeRule:          SolverRule{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 6 * time.Hour, To: 13 * time.Hour}, StartAmount: 0, EndAmount: 7.0, StartTimePolicy: ShiftablePolicy, RuleResolutionPolicy: ResolvePolicy},
			realtiveStartOffset:   0, //3 * time.Hour,
			relativeAmountOffset:  3.0,
			activatedFlatRatesSum: 3.0,
			expectedIntersect:     true,
			expected:              9 * time.Hour,
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver()
			solver.SetWindow(time.Now(), time.Duration(48*time.Hour))
			solver.currentRelativeStartOffset = testcase.realtiveStartOffset
			solver.currentRelativeAmountOffset = testcase.relativeAmountOffset
			solver.activatedFlatRatesSum = testcase.activatedFlatRatesSum

			intersect := solver.IsIntersectingFlatRate(&testcase.relativeRule, &testcase.flatRateRule)
			if intersect != testcase.expectedIntersect {
				t.Errorf("IsIntersectingFlatRate expected %v, got %v", testcase.expectedIntersect, intersect)
			}
			out := solver.FindIntersectPositionFlatRate(&testcase.relativeRule, &testcase.flatRateRule)
			if out != testcase.expected {
				t.Errorf("FindIntersectPositionFlatRate expected %v, got %v", testcase.expected, out)
			}
		})
	}
}

func TestExtractRange(t *testing.T) {
	tests := map[string]struct {
		rules    SolverRules
		timespan RelativeTimeSpan
		expected SolverRules
	}{
		// 0 - No rules in range
		"0-NoRulesInRange": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
			},
			timespan: RelativeTimeSpan{From: 30 * time.Minute, To: 40 * time.Minute},
			expected: SolverRules{},
		},
		// 1 - One rule fully in range
		"1-OneRuleFullyInRange": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
			},
			timespan: RelativeTimeSpan{From: 5 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
			},
		},
		// 2 - One rule partially in range at the beginning
		"2-OneRulePartiallyInRangeAtBeginning": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
			},
			timespan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 20 * time.Minute}},
			},
		},
		// 3 - One rule partially in range at the end
		"3-OneRulePartiallyInRangeAtEnd": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
			},
			timespan: RelativeTimeSpan{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 15 * time.Minute}},
			},
		},
		// 4 - One rule longer than range
		"4-OneRuleLongerThanRange": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 30 * time.Minute}},
			},
			timespan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 25 * time.Minute}},
			},
		},
		// 5 - Multiple rules in range
		"5-MultipleRulesInRange": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 35 * time.Minute}},
			},
			timespan: RelativeTimeSpan{From: 5 * time.Minute, To: 30 * time.Minute},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 30 * time.Minute}},
			},
		},
		// 6 - Multiple rules with one fully in range and one partially in range
		"6-MultipleRulesWithOneFullyInRangeAndOnePartiallyInRange": {
			rules: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Minute, To: 20 * time.Minute}},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 35 * time.Minute}},
			},
			timespan: RelativeTimeSpan{From: 15 * time.Minute, To: 30 * time.Minute},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Minute, To: 20 * time.Minute}},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Minute, To: 30 * time.Minute}},
			},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver()
			solver.SetWindow(time.Now(), time.Duration(48*time.Hour))
			solver.AppendMany(testcase.rules...)

			out := solver.ExtractRulesInRange(testcase.timespan)
			if len(out) != len(testcase.expected) {
				t.Errorf("ExtractRange expected %v rules, got %v", len(testcase.expected), len(out))
			} else {
				for i := range out {
					if out[i].From != testcase.expected[i].From || out[i].To != testcase.expected[i].To {
						t.Errorf("ExtractRange expected rule %v, got %v", testcase.expected, out)
					}
					if out[i].StartAmount != testcase.expected[i].StartAmount || out[i].EndAmount != testcase.expected[i].EndAmount {
						t.Errorf("ExtractRange expected rule %v, got %v", testcase.expected, out)
					}
				}
			}
		})
	}
}

func TestFindFlatRateActivationTime(t *testing.T) {
	tests := map[string]struct {
		flatRateRule         SolverRule
		solvedRules          SolverRules
		extraRule            SolverRule
		expectedActivation   bool
		expectedActivationAt time.Duration
	}{
		// 0 - No activation, no extra rules
		"0-NoActivationNoExtraRules": {
			flatRateRule:         NewAbsoluteFlatRateRule("FlatRate", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 4.0, MetaData{}),
			extraRule:            SolverRule{},
			solvedRules:          SolverRules{},
			expectedActivation:   false,
			expectedActivationAt: 0,
		},
		// 1 - Activation with one extra rule fully in range
		"1-ActivationWithOneExtraRuleFullyInRange": {
			flatRateRule:         NewAbsoluteFlatRateRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			extraRule:            NewRelativeLinearRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}),
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 2 - Activation with multiple extra rules
		"2-ActivationWithMultipleExtraRules": {
			flatRateRule:         NewAbsoluteFlatRateRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			extraRule:            NewRelativeLinearRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(4 * time.Hour),
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 3 - No activation, extra rules do not meet the required amount
		"3-NoActivationExtraRulesInsufficient": {
			flatRateRule:         NewAbsoluteFlatRateRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			extraRule:            NewRelativeLinearRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(5 * time.Hour),
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 4 - Activation with overlapping extra rules
		"4-ActivationWithOverlappingExtraRules": {
			flatRateRule:         NewAbsoluteFlatRateRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			extraRule:            NewRelativeLinearRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(8 * time.Hour),
			expectedActivation:   false,
			expectedActivationAt: 0,
		},
		// 5 - Activation with one extra rule fully in range
		"5-ActivationWithOneExtraRuleFullyInRangeDuplicate": {
			flatRateRule:         NewAbsoluteFlatRateRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			solvedRules:          SolverRules{NewRelativeLinearRule("ExtraRule", 16*time.Hour, 1.0, MetaData{})},
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 6 - Activation with multiple extra rules
		"6-ActivationWithMultipleExtraRulesDuplicate": {
			flatRateRule:         NewAbsoluteFlatRateRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			solvedRules:          SolverRules{NewRelativeLinearRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(4 * time.Hour)},
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 7 - No activation, extra rules do not meet the required amount
		"7-NoActivationExtraRulesInsufficientDuplicate": {
			flatRateRule:         NewAbsoluteFlatRateRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			solvedRules:          SolverRules{NewRelativeLinearRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(5 * time.Hour)},
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 8 - Activation with overlapping extra rules
		"8-ActivationWithOverlappingExtraRulesDuplicate": {
			flatRateRule:         NewAbsoluteFlatRateRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			solvedRules:          SolverRules{NewRelativeLinearRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(8 * time.Hour)},
			expectedActivation:   false,
			expectedActivationAt: 0,
		},
		// 9 - Activation with one extra rule fully in range
		"9-ActivationWithOneExtraRuleFullyInRange": {
			flatRateRule: NewAbsoluteFlatRateRule("FlatRate", RelativeTimeSpan{6 * time.Hour, 20 * time.Hour}, 8.0, MetaData{}),
			solvedRules: SolverRules{
				NewRelativeLinearRule("A", 10*time.Hour, 1.0, MetaData{}),
				NewRelativeLinearRule("B", 1*time.Hour, 0, MetaData{}).Shift(10 * time.Hour),
				NewRelativeLinearRule("C", 10*time.Hour, 1.0, MetaData{}).Shift(11 * time.Hour),
			},
			expectedActivation:   true,
			expectedActivationAt: 15 * time.Hour,
		},
		// 10 - Activation with one extra rule fully in range
		"10-ActivationWithOneExtraRuleFullyInRange": {
			flatRateRule: NewAbsoluteFlatRateRule("FlatRate", RelativeTimeSpan{6 * time.Hour, 20 * time.Hour}, 8.0, MetaData{}),
			solvedRules: SolverRules{
				NewRelativeLinearRule("A", 10*time.Hour, 1.0, MetaData{}),
				NewRelativeLinearRule("B", 1*time.Hour, 0, MetaData{}).Shift(10 * time.Hour),
			},
			extraRule:            NewRelativeLinearRule("C", 10*time.Hour, 1.0, MetaData{}).Shift(11 * time.Hour),
			expectedActivation:   true,
			expectedActivationAt: 15 * time.Hour,
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver()
			solver.SetWindow(time.Now(), time.Duration(48*time.Hour))
			for i := range testcase.solvedRules {
				solver.solvedRules.ReplaceOrInsert(&testcase.solvedRules[i])
			}

			activationAt, activated := solver.findFlatRateActivationTime(&testcase.flatRateRule, &testcase.extraRule)
			if activated != testcase.expectedActivation {
				t.Errorf("findFlatRateActivationTime expected activation %v, got %v", testcase.expectedActivation, activated)
			}
			if activationAt != testcase.expectedActivationAt {
				t.Errorf("findFlatRateActivationTime expected activation at %v, got %v", testcase.expectedActivationAt, activationAt)
			}
		})
	}
}
