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
			lpRule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
			hpRule:   SolverRule{From: 20 * time.Minute, To: 30 * time.Minute},
			expected: SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
		},
		// 1 - No conflict, Resolve policy
		"1-NoConflictBefore-Resolve": {
			lpRule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
			hpRule:   SolverRule{From: 5 * time.Minute, To: 10 * time.Minute},
			expected: SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
		},
		// 2 - HPRule overlapping end of rule, Resolve policy
		"2-OverlapEnd-Resolve": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute}},
		},
		// 3 - HPRule overlapping end of rule, Resolve policy
		"3-OverlapEndStep-Resolve": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 100},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, StartAmount: 0, EndAmount: 0}},
		},
		// 4 - HPRule overlapping end of rule, Resolve policy
		"4-OverlapLinear-Resolve": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 150},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, StartAmount: 0, EndAmount: 50}},
		},
		// 5 - HPRule overlapping beginning of rule, Resolve policy
		"5-OverlapBegining-Resolve": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 25 * time.Minute}},
		},
		// 6 - HPRule overlapping beginning of rule, Resolve policy
		"6-OverlapBeginingStep-Resolve": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 25 * time.Minute, StartAmount: 100, EndAmount: 100}},
		},
		// 7 - HPRule overlapping beginning of rule, Resolve policy
		"7-OverlapBeginingLinear-Resolve": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 25 * time.Minute, StartAmount: 100, EndAmount: 200}},
		},
		// 8 - HPRule overlapping middle of rule, Resolve policy
		"8-OverlapMiddle-Resolve": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute},
				SolverRule{From: 25 * time.Minute, To: 40 * time.Minute}},
		},
		// 9 - HPRule overlapping middle of rule, Resolve policy
		"9-OverlapMiddleStep-Resolve": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 100},
				SolverRule{From: 25 * time.Minute, To: 40 * time.Minute, StartAmount: 0, EndAmount: 0}},
		},
		// 10 - HPRule overlapping middle of rule, Resolve policy
		"10-OverlapMiddleLinear-Resolve": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: ResolvePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 125},
				SolverRule{From: 25 * time.Minute, To: 40 * time.Minute, StartAmount: 0, EndAmount: 75}},
		},
		// 11 - No conflict Truncate policy
		"11-NoConflictAfter-Truncate": {
			lpRule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy},
			hpRule:   SolverRule{From: 20 * time.Minute, To: 30 * time.Minute},
			expected: SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
		},
		// 12 - No conflict Truncate policy
		"12-NoConflictBefore-Truncate": {
			lpRule:   SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy},
			hpRule:   SolverRule{From: 5 * time.Minute, To: 10 * time.Minute},
			expected: SolverRules{SolverRule{From: 10 * time.Minute, To: 20 * time.Minute}},
		},
		// 13 - HPRule overlapping end of rule, Truncate policy
		"13-OverlapEnd-Truncate": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute}},
		},
		// 14 - HPRule overlapping end of rule, Truncate policy
		"14-OverlapEndStep-Truncate": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 100}},
		},
		// 15 - HPRule overlapping end of rule, Truncate policy
		"15-OverlapLinear-Truncate": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 150}},
		},
		// 16 - HPRule overlapping beginning of rule, Truncate policy
		"16-OverlapBegining-Truncate": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 20 * time.Minute}},
		},
		// 17 - HPRule overlapping beginning of rule, Truncate policy
		"17-OverlapBeginingStep-Truncate": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 20 * time.Minute, StartAmount: 0, EndAmount: 0}},
		},
		// 18 - HPRule overlapping beginning of rule, Truncate policy
		"18-OverlapBeginingLinear-Truncate": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 20 * time.Minute, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 5 * time.Minute, To: 15 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 15 * time.Minute, To: 20 * time.Minute, StartAmount: 0, EndAmount: 50}},
		},
		// 19 - HPRule overlapping middle of rule, Truncate policy
		"19-OverlapMiddle-Truncate": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute}},
		},
		// 20 - HPRule overlapping middle of rule, Truncate policy
		"20-OverlapMiddleStep-Truncate": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, StartAmount: 100, EndAmount: 100, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 100},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, StartAmount: 0, EndAmount: 0}},
		},
		// 21 - HPRule overlapping middle of rule, Truncate policy
		"21-OverlapMiddleLinear-Truncate": {
			lpRule: SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: TruncatePolicy},
			hpRule: SolverRule{From: 15 * time.Minute, To: 25 * time.Minute},
			expected: SolverRules{
				SolverRule{From: 10 * time.Minute, To: 15 * time.Minute, StartAmount: 100, EndAmount: 125},
				SolverRule{From: 25 * time.Minute, To: 30 * time.Minute, StartAmount: 0, EndAmount: 25}},
		},
		// 22 - HPRule fully overlap rule, Truncate policy
		"22-FullOverlap-Truncate": {
			lpRule:   SolverRule{From: 10 * time.Minute, To: 30 * time.Minute, StartAmount: 100, EndAmount: 200, RuleResolutionPolicy: TruncatePolicy},
			hpRule:   SolverRule{From: 5 * time.Minute, To: 35 * time.Minute},
			expected: SolverRules{},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver()
			solver.SetWindow(time.Now(), time.Duration(48*time.Hour))
			solver.currentRelativeStartOffset = 45 * time.Minute
			out := solver.solveVsSingle(testcase.lpRule, testcase.hpRule)
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

	DummyRule := SolverRule{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy}

	tests := map[string]struct {
		rules    SolverRules
		expected SolverRules
	}{
		// 0 - No conflict, empty rules§
		"0-NoConflictEmptyRules": {
			rules: SolverRules{
				DummyRule,
			},
			expected: SolverRules{
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute},
			},
		},
		// 1 - No conflict, multiple rules
		"1-NoConflictMultipleRules": {
			rules: SolverRules{
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "B", From: 20 * time.Minute, To: 30 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "C", From: 30 * time.Minute, To: 40 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: SolverRules{
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{RuleName: "B", From: 20 * time.Minute, To: 30 * time.Minute},
				{RuleName: "C", From: 30 * time.Minute, To: 40 * time.Minute},
			},
		},
		// 2 - Overlapping rules, Resolve policy
		"2-OverlapResolvePolicy": {
			rules: SolverRules{
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: SolverRules{
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{RuleName: "B", From: 20 * time.Minute, To: 30 * time.Minute},
				{RuleName: "C", From: 30 * time.Minute, To: 40 * time.Minute},
			},
		},
		// 3 - Overlapping rules, Truncate policy
		"3-OverlapTruncatePolicy": {
			rules: SolverRules{
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: SolverRules{
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{RuleName: "B", From: 20 * time.Minute, To: 25 * time.Minute},
				{RuleName: "C", From: 25 * time.Minute, To: 35 * time.Minute},
			},
		},
		// 4 - Overlapping rules, Delete policy
		"4-OverlapDeletePolicy": {
			rules: SolverRules{
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: SolverRules{
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{RuleName: "C", From: 25 * time.Minute, To: 35 * time.Minute},
			},
		},
		// 5 - Overlapping rules, mixed policies
		"5-OverlapMixedPolicies": {
			rules: SolverRules{
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: DeletePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "B", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: ResolvePolicy, StartTimePolicy: FixedPolicy},
				{RuleName: "C", From: 25 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: FixedPolicy},
			},
			expected: SolverRules{
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{RuleName: "B", From: 20 * time.Minute, To: 30 * time.Minute},
				{RuleName: "C", From: 30 * time.Minute, To: 35 * time.Minute},
			},
		},
		// 6 - Overlapping rules, Truncate policy, Shiftable policy
		"6-OverlapTruncateShiftablePolicy": {
			rules: SolverRules{
				NewAbsoluteFlatRateRule("A", 10*time.Minute, 20*time.Minute, MustParseAmount("0")),
				NewAbsoluteFlatRateRule("B", 30*time.Minute, 40*time.Minute, MustParseAmount("0")),
				{RuleName: "C", From: 0 * time.Minute, To: 10 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				{RuleName: "D", From: 0 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				NewAbsoluteFlatRateRule("Z", 60*time.Minute, 90*time.Minute, MustParseAmount("0")),
				{RuleName: "E", From: 5 * time.Minute, To: 15 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				{RuleName: "F", From: 0 * time.Minute, To: 20 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
			},
			expected: SolverRules{
				{RuleName: "C", From: 0 * time.Minute, To: 10 * time.Minute},
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{RuleName: "D", From: 20 * time.Minute, To: 25 * time.Minute},
				{RuleName: "E", From: 25 * time.Minute, To: 30 * time.Minute},
				{RuleName: "B", From: 30 * time.Minute, To: 40 * time.Minute},
				{RuleName: "F", From: 40 * time.Minute, To: 55 * time.Minute},
				{RuleName: "Z", From: 60 * time.Minute, To: 90 * time.Minute},
			},
		},
		// 7 - Overlapping rules, Truncate policy, Shiftable policy
		"7-OverlapTruncateShiftablePolicy": {
			rules: SolverRules{
				NewAbsoluteFlatRateRule("A", 10*time.Minute, 20*time.Minute, MustParseAmount("0")),
				NewAbsoluteFlatRateRule("B", 30*time.Minute, 40*time.Minute, MustParseAmount("0")),
				{RuleName: "C", From: 15 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				{RuleName: "D", From: 0 * time.Minute, To: 35 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
				NewAbsoluteFlatRateRule("Z", 60*time.Minute, 90*time.Minute, MustParseAmount("0")),
				{RuleName: "E", From: 5 * time.Minute, To: 25 * time.Minute, RuleResolutionPolicy: TruncatePolicy, StartTimePolicy: ShiftablePolicy},
			},
			expected: SolverRules{
				{RuleName: "C", From: 0 * time.Minute, To: 10 * time.Minute},
				{RuleName: "A", From: 10 * time.Minute, To: 20 * time.Minute},
				{RuleName: "D", From: 20 * time.Minute, To: 30 * time.Minute},
				{RuleName: "B", From: 30 * time.Minute, To: 40 * time.Minute},
				{RuleName: "D", From: 40 * time.Minute, To: 45 * time.Minute},
				{RuleName: "E", From: 45 * time.Minute, To: 60 * time.Minute},
				{RuleName: "Z", From: 60 * time.Minute, To: 90 * time.Minute},
			},
		},

		// 8 - Overlapping multiple calendar flatrate with linear in the middle and multiple flatrate activation
		"10-MultipleCalendarFlatrate": {
			rules: SolverRules{
				NewAbsoluteFlatRateRule("Day", 2*time.Hour, 14*time.Hour, MustParseAmount("7.0")),
				NewAbsoluteFlatRateRule("Morning", 2*time.Hour+1, 6*time.Hour, MustParseAmount("3.0")),
				NewAbsoluteFlatRateRule("Evening", 7*time.Hour, 10*time.Hour, MustParseAmount("4.0")),
				NewAbsoluteFlatRateRule("NotIntersecting", 24*time.Hour, 25*time.Hour, MustParseAmount("4.0")),
				NewRelativeLinearRule("Hourly", 10*time.Hour, MustParseAmount("1.0")),
			},
			expected: SolverRules{
				{RuleName: "Hourly", From: 0 * time.Hour, To: 3 * time.Hour, StartAmount: AmountZero, EndAmount: MustParseAmount("3.0")},
				{RuleName: "Morning", From: 3 * time.Hour, To: 6 * time.Hour, StartAmount: AmountZero, EndAmount: AmountZero},
				{RuleName: "Hourly", From: 6 * time.Hour, To: 9 * time.Hour, StartAmount: AmountZero, EndAmount: MustParseAmount("3.0")},
				{RuleName: "Evening", From: 9 * time.Hour, To: 10 * time.Hour, StartAmount: AmountZero, EndAmount: AmountZero},
				{RuleName: "Hourly", From: 10 * time.Hour, To: 11 * time.Hour, StartAmount: AmountZero, EndAmount: MustParseAmount("1.0")},
				{RuleName: "Day", From: 11 * time.Hour, To: 14 * time.Hour, StartAmount: AmountZero, EndAmount: AmountZero},
				{RuleName: "Hourly", From: 14 * time.Hour, To: 17 * time.Hour, StartAmount: AmountZero, EndAmount: MustParseAmount("3.0")},
			},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver()
			solver.SetWindow(time.Now(), time.Duration(48*time.Hour))
			solver.Append(testcase.rules...)

			if solver.rules.Len() != len(testcase.expected) {
				t.Errorf("SolveAndAppend expected %v rules, got %v", len(testcase.expected), solver.rules.Len())
			} /*else  {*/
			i := 0
			solver.rules.Ascend(func(rule *SolverRule) bool {
				expected := testcase.expected[i]
				if rule.From != expected.From || rule.To != expected.To {
					t.Errorf("SolveAndAppend time error, expected from %s to %s, got  from %s to %s", expected.From, expected.To, rule.From, rule.To)
				}
				if rule.StartAmount != expected.StartAmount || rule.EndAmount != expected.EndAmount {
					t.Errorf("SolveAndAppend amount error, expected rule %v, got %v", expected, rule)
				}
				// Test name
				if rule.Name() != expected.Name() {
					t.Errorf("SolveAndAppend mismatch names, expected name %s, got %s", expected.Name(), rule.Name())
				}
				i++
				return true
			})
			/*}*/
		})
	}
	/*
		//Additional test to check that link to original rule is kept while running the solver
		t.Run("OriginalRuleLinkCheck", func(t *testing.T) {
			solver := NewSolver()
			solver.SetWindow(time.Now(), time.Duration(48*time.Hour))
			testcase := tests["0-NoConflictEmptyRules"]
			solver.AppendSolverRules(testcase.rules...)

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
	*/
}
