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

func TestSolver(t *testing.T) {

	tests := map[string]struct {
		rules    SolverRules
		expected SolverRules
	}{
		"01-Sequential-SingleLinear": {
			rules: SolverRules{
				NewLinearSequentialRule("Hourly", 10*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 10 * time.Hour}, StartAmount: 0, EndAmount: 10.0},
			},
		},
		"02-Sequential-MultipleLinear": {
			rules: SolverRules{
				NewLinearSequentialRule("A", 4*time.Hour, 1.0, MetaData{}),
				NewLinearSequentialRule("B", 5*time.Hour, 2.0, MetaData{}),
				NewLinearSequentialRule("C", 7*time.Hour, 3.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 4 * time.Hour}, StartAmount: 0, EndAmount: 4.0},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 4 * time.Hour, To: 9 * time.Hour}, StartAmount: 0, EndAmount: 10.0},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 9 * time.Hour, To: 16 * time.Hour}, StartAmount: 0, EndAmount: 21.0},
			},
		},
		"03-Sequential-SingleFixedRate": {
			rules: SolverRules{
				NewFixedRateSequentialRule("A", 4*time.Hour, 3.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 4 * time.Hour}, StartAmount: 3.0, EndAmount: 3.0},
			},
		},
		"04-Sequential-MultipleFixedRate": {
			rules: SolverRules{
				NewFixedRateSequentialRule("A", 4*time.Hour, 3.0, MetaData{}),
				NewFixedRateSequentialRule("B", 5*time.Hour, 4.0, MetaData{}),
				NewFixedRateSequentialRule("C", 7*time.Hour, 5.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 4 * time.Hour}, StartAmount: 3.0, EndAmount: 3.0},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 4 * time.Hour, To: 9 * time.Hour}, StartAmount: 4.0, EndAmount: 4.0},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 9 * time.Hour, To: 16 * time.Hour}, StartAmount: 5.0, EndAmount: 5.0},
			},
		},
		"05-Sequential-MixedLinearAndFixedRate": {
			rules: SolverRules{
				NewLinearSequentialRule("A", 4*time.Hour, 1.0, MetaData{}),
				NewFixedRateSequentialRule("B", 5*time.Hour, 2.0, MetaData{}),
				NewLinearSequentialRule("C", 7*time.Hour, 3.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 4 * time.Hour}, StartAmount: 0, EndAmount: 4.0},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 4 * time.Hour, To: 9 * time.Hour}, StartAmount: 2.0, EndAmount: 2.0},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 9 * time.Hour, To: 16 * time.Hour}, StartAmount: 0, EndAmount: 21.0},
			},
		},
		"06-Fixed-MultipleLinearBasic": {
			rules: SolverRules{
				NewLinearFixedRule("A", RelativeTimeSpan{0 * time.Hour, 2 * time.Hour}, 1.0, MetaData{}),
				NewLinearFixedRule("B", RelativeTimeSpan{2 * time.Hour, 5 * time.Hour}, 2.0, MetaData{}),
				NewLinearFixedRule("C", RelativeTimeSpan{5 * time.Hour, 9 * time.Hour}, 3.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 2 * time.Hour}, StartAmount: 0, EndAmount: 2.0},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 2 * time.Hour, To: 5 * time.Hour}, StartAmount: 0, EndAmount: 6.0},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Hour, To: 9 * time.Hour}, StartAmount: 0, EndAmount: 12.0},
			},
		},
		"07-Fixed-MultipleLinearOverlapping": {
			rules: SolverRules{
				NewLinearFixedRule("A", RelativeTimeSpan{2 * time.Hour, 6 * time.Hour}, 1.0, MetaData{}),
				NewLinearFixedRule("B", RelativeTimeSpan{1 * time.Hour, 5 * time.Hour}, 2.0, MetaData{}),
				NewLinearFixedRule("X", RelativeTimeSpan{18 * time.Hour, 20 * time.Hour}, 2.0, MetaData{}),
				NewLinearFixedRule("C", RelativeTimeSpan{3 * time.Hour, 9 * time.Hour}, 3.0, MetaData{}),
				NewLinearFixedRule("D", RelativeTimeSpan{7 * time.Hour, 15 * time.Hour}, 4.0, MetaData{}),
				NewLinearFixedRule("Y", RelativeTimeSpan{19 * time.Hour, 25 * time.Hour}, 2.0, MetaData{}),
				NewLinearFixedRule("E", RelativeTimeSpan{0 * time.Hour, 7 * time.Hour}, 5.0, MetaData{}),
				NewLinearFixedRule("F", RelativeTimeSpan{10 * time.Hour, 17 * time.Hour}, 6.0, MetaData{}),
				NewLinearFixedRule("Z", RelativeTimeSpan{18 * time.Hour, 53 * time.Hour}, 2.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "E", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 1 * time.Hour}, StartAmount: 0.0, EndAmount: 5.0},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 1 * time.Hour, To: 2 * time.Hour}, StartAmount: 0.0, EndAmount: 2.0},
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 2 * time.Hour, To: 6 * time.Hour}, StartAmount: 0.0, EndAmount: 4.0},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 6 * time.Hour, To: 9 * time.Hour}, StartAmount: 0.0, EndAmount: 9.0},
				{RuleName: "D", RelativeTimeSpan: RelativeTimeSpan{From: 9 * time.Hour, To: 15 * time.Hour}, StartAmount: 0.0, EndAmount: 24.0},
				{RuleName: "F", RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Hour, To: 17 * time.Hour}, StartAmount: 0.0, EndAmount: 12.0},
			},
		},
		"08-Mixed-MultipleLinearOverlappingAndSequential": {
			rules: SolverRules{
				NewLinearFixedRule("A", RelativeTimeSpan{2 * time.Hour, 6 * time.Hour}, 1.0, MetaData{}),
				NewLinearFixedRule("B", RelativeTimeSpan{1 * time.Hour, 5 * time.Hour}, 2.0, MetaData{}),
				NewLinearFixedRule("X", RelativeTimeSpan{18 * time.Hour, 20 * time.Hour}, 2.0, MetaData{}),
				NewLinearFixedRule("C", RelativeTimeSpan{3 * time.Hour, 9 * time.Hour}, 3.0, MetaData{}),
				NewLinearFixedRule("D", RelativeTimeSpan{7 * time.Hour, 15 * time.Hour}, 4.0, MetaData{}),
				NewLinearFixedRule("Y", RelativeTimeSpan{19 * time.Hour, 25 * time.Hour}, 2.0, MetaData{}),
				NewLinearFixedRule("E", RelativeTimeSpan{0 * time.Hour, 7 * time.Hour}, 5.0, MetaData{}),
				NewLinearFixedRule("F", RelativeTimeSpan{10 * time.Hour, 17 * time.Hour}, 6.0, MetaData{}),
				NewLinearFixedRule("Z", RelativeTimeSpan{18 * time.Hour, 53 * time.Hour}, 2.0, MetaData{}),
				NewLinearSequentialRule("Q", 10*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "E", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 1 * time.Hour}, StartAmount: 0.0, EndAmount: 5.0},
				{RuleName: "B", RelativeTimeSpan: RelativeTimeSpan{From: 1 * time.Hour, To: 2 * time.Hour}, StartAmount: 0.0, EndAmount: 2.0},
				{RuleName: "A", RelativeTimeSpan: RelativeTimeSpan{From: 2 * time.Hour, To: 6 * time.Hour}, StartAmount: 0.0, EndAmount: 4.0},
				{RuleName: "C", RelativeTimeSpan: RelativeTimeSpan{From: 6 * time.Hour, To: 9 * time.Hour}, StartAmount: 0.0, EndAmount: 9.0},
				{RuleName: "D", RelativeTimeSpan: RelativeTimeSpan{From: 9 * time.Hour, To: 15 * time.Hour}, StartAmount: 0.0, EndAmount: 24.0},
				{RuleName: "F", RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Hour, To: 17 * time.Hour}, StartAmount: 0.0, EndAmount: 12.0},
				{RuleName: "Q", RelativeTimeSpan: RelativeTimeSpan{From: 17 * time.Hour, To: 18 * time.Hour}, StartAmount: 0.0, EndAmount: 1.0},
				{RuleName: "X", RelativeTimeSpan: RelativeTimeSpan{From: 18 * time.Hour, To: 20 * time.Hour}, StartAmount: 0.0, EndAmount: 4.0},
				{RuleName: "Y", RelativeTimeSpan: RelativeTimeSpan{From: 20 * time.Hour, To: 25 * time.Hour}, StartAmount: 0.0, EndAmount: 10.0},
				{RuleName: "Z", RelativeTimeSpan: RelativeTimeSpan{From: 25 * time.Hour, To: 53 * time.Hour}, StartAmount: 0.0, EndAmount: 56.0},
				{RuleName: "Q", RelativeTimeSpan: RelativeTimeSpan{From: 53 * time.Hour, To: 62 * time.Hour}, StartAmount: 0.0, EndAmount: 9.0},
			},
		},

		"10-LinearWithNonPaying": {
			rules: SolverRules{
				NewNonPayingFixedRule("Morning", RelativeTimeSpan{2 * time.Hour, 6 * time.Hour}, MetaData{}),
				NewLinearSequentialRule("Hourly", 10*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 2 * time.Hour}, StartAmount: 0, EndAmount: 2.0},
				{RuleName: "Morning", RelativeTimeSpan: RelativeTimeSpan{From: 2 * time.Hour, To: 6 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 6 * time.Hour, To: 14 * time.Hour}, StartAmount: 0, EndAmount: 8.0},
			},
		},
		"11-LinearWithFlatrate": {
			rules: SolverRules{
				NewFlatRateFixedRule("Morning", RelativeTimeSpan{2 * time.Hour, 6 * time.Hour}, 3.0, MetaData{}),
				NewLinearSequentialRule("Hourly", 10*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 5 * time.Hour}, StartAmount: 0, EndAmount: 5.0},
				{RuleName: "Morning", RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Hour, To: 6 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 6 * time.Hour, To: 11 * time.Hour}, StartAmount: 0, EndAmount: 5.0},
			},
		},
		"12-LinearWithNonActivatedFlatrate": {
			rules: SolverRules{
				NewFlatRateFixedRule("Morning", RelativeTimeSpan{2 * time.Hour, 6 * time.Hour}, 4.0, MetaData{}),
				NewLinearSequentialRule("Hourly", 10*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 10 * time.Hour}, StartAmount: 0, EndAmount: 10.0},
			},
		},

		// 12 - Overlapping multiple calendar flatrate with linear in the middle and multiple flatrate activation
		"20-MultipleCalendarFlatrate": {
			rules: SolverRules{
				NewFlatRateFixedRule("Morning", RelativeTimeSpan{2 * time.Hour, 6 * time.Hour}, 3.0, MetaData{}),
				NewFlatRateFixedRule("Evening", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
				NewFlatRateFixedRule("Daily", RelativeTimeSpan{2 * time.Hour, 14 * time.Hour}, 7.0, MetaData{}),
				NewLinearSequentialRule("Hourly", 12*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 5 * time.Hour}, StartAmount: 0, EndAmount: 5.0},
				{RuleName: "Morning", RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Hour, To: 6 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 6 * time.Hour, To: 10 * time.Hour}, StartAmount: 0, EndAmount: 4.0},
				{RuleName: "Daily", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Hour, To: 14 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 14 * time.Hour, To: 17 * time.Hour}, StartAmount: 0, EndAmount: 3.0},
			},
		},
		// 13 - Almost the same as above but with different time spans to make evening flatrate more advantageous than daily
		"21-MultipleCalendarFlatrate": {
			rules: SolverRules{
				NewFlatRateFixedRule("Morning", RelativeTimeSpan{2 * time.Hour, 6 * time.Hour}, 3.0, MetaData{}),
				NewFlatRateFixedRule("Evening", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
				NewFlatRateFixedRule("Daily", RelativeTimeSpan{3 * time.Hour, 14 * time.Hour}, 7.0, MetaData{}),
				NewLinearSequentialRule("Hourly", 12*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 5 * time.Hour}, StartAmount: 0, EndAmount: 5.0},
				{RuleName: "Morning", RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Hour, To: 6 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 6 * time.Hour, To: 10 * time.Hour}, StartAmount: 0, EndAmount: 4.0},
				{RuleName: "Evening", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Hour, To: 11 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 11 * time.Hour, To: 12 * time.Hour}, StartAmount: 0, EndAmount: 1.0},
				{RuleName: "Daily", RelativeTimeSpan: RelativeTimeSpan{From: 12 * time.Hour, To: 14 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 14 * time.Hour, To: 16 * time.Hour}, StartAmount: 0, EndAmount: 2.0},
			},
		},

		// 30 - Figma test case with single flat rate rule (activating the flatrate)
		"30-Figma-SingleFlatRule-1": {
			rules: SolverRules{
				NewFlatRateFixedRule("Morning", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
				NewLinearSequentialRule("Hourly", 20*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 10 * time.Hour}, StartAmount: 0, EndAmount: 10.0},
				{RuleName: "Morning", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Hour, To: 11 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 11 * time.Hour, To: 21 * time.Hour}, StartAmount: 0, EndAmount: 10.0},
			},
		},
		// 31 - Figma test case with single flat rate rule (activating the flatrate)
		"31-Figma-SingleFlatRule-2": {
			rules: SolverRules{
				NewFlatRateFixedRule("Morning", RelativeTimeSpan{3 * time.Hour, 7 * time.Hour}, 3.0, MetaData{}),
				NewLinearSequentialRule("Hourly", 20*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 6 * time.Hour}, StartAmount: 0, EndAmount: 6.0},
				{RuleName: "Morning", RelativeTimeSpan: RelativeTimeSpan{From: 6 * time.Hour, To: 7 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 7 * time.Hour, To: 21 * time.Hour}, StartAmount: 0, EndAmount: 14.0},
			},
		},
		// 32 - Figma test case with single flat rate rule (activating the flatrate)
		"32-Figma-SingleFlatRule-3": {
			rules: SolverRules{
				NewFlatRateFixedRule("Morning", RelativeTimeSpan{2 * time.Hour, 6 * time.Hour}, 3.0, MetaData{}),
				NewLinearSequentialRule("Hourly", 20*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 5 * time.Hour}, StartAmount: 0, EndAmount: 5.0},
				{RuleName: "Morning", RelativeTimeSpan: RelativeTimeSpan{From: 5 * time.Hour, To: 6 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 6 * time.Hour, To: 21 * time.Hour}, StartAmount: 0, EndAmount: 15.0},
			},
		},
		// 33 - Figma test case with single flat rate rule (flatrate is not activated)
		"33-Figma-SingleFlatRule-4": {
			rules: SolverRules{
				NewFlatRateFixedRule("Morning", RelativeTimeSpan{-1 * time.Hour, 3 * time.Hour}, 3.0, MetaData{}),
				NewLinearSequentialRule("Hourly", 20*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 20 * time.Hour}, StartAmount: 0, EndAmount: 20.0},
			},
		},
		// 34 - Figma test case with multiple flat rate rules (activating both flatrate)
		"34-Figma-MultiFlatRate-1": {
			rules: SolverRules{
				NewFlatRateFixedRule("Morning", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
				NewFlatRateFixedRule("Evening", RelativeTimeSpan{6 * time.Hour, 20 * time.Hour}, 8.0, MetaData{}),
				NewLinearSequentialRule("Hourly", 20*time.Hour, 1.0, MetaData{}),
			},
			expected: SolverRules{
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 0 * time.Hour, To: 10 * time.Hour}, StartAmount: 0, EndAmount: 10.0},
				{RuleName: "Morning", RelativeTimeSpan: RelativeTimeSpan{From: 10 * time.Hour, To: 11 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 11 * time.Hour, To: 15 * time.Hour}, StartAmount: 0, EndAmount: 4.0},
				{RuleName: "Evening", RelativeTimeSpan: RelativeTimeSpan{From: 15 * time.Hour, To: 20 * time.Hour}, StartAmount: 0, EndAmount: 0},
				{RuleName: "Hourly", RelativeTimeSpan: RelativeTimeSpan{From: 20 * time.Hour, To: 26 * time.Hour}, StartAmount: 0, EndAmount: 6.0},
			},
		},
	}

	for name, testcase := range tests {
		t.Run(name, func(t *testing.T) {
			solver := NewSolver()
			solver.SetWindow(time.Now(), time.Duration(48*time.Hour))
			solver.AppendMany(testcase.rules...)

			solver.Solve()

			if solver.solvedRules.Len() != len(testcase.expected) {
				t.Errorf("SolveAndAppend expected %v rules, got %v", len(testcase.expected), solver.solvedRules.Len())
			}
			i := 0
			solver.solvedRules.Ascend(func(rule *SolverRule) bool {
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
				return i <= len(testcase.expected)
			})
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
			for i := range testcase.rules {
				solver.solvedRules.ReplaceOrInsert(&testcase.rules[i])
			}

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
			flatRateRule:         NewFlatRateFixedRule("FlatRate", RelativeTimeSpan{8 * time.Hour, 12 * time.Hour}, 4.0, MetaData{}),
			extraRule:            SolverRule{},
			solvedRules:          SolverRules{},
			expectedActivation:   false,
			expectedActivationAt: 0,
		},
		// 1 - Activation with one extra rule fully in range
		"1-ActivationWithOneExtraRuleFullyInRange": {
			flatRateRule:         NewFlatRateFixedRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			extraRule:            NewLinearSequentialRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}),
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 2 - Activation with multiple extra rules
		"2-ActivationWithMultipleExtraRules": {
			flatRateRule:         NewFlatRateFixedRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			extraRule:            NewLinearSequentialRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(4 * time.Hour),
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 3 - No activation, extra rules do not meet the required amount
		"3-NoActivationExtraRulesInsufficient": {
			flatRateRule:         NewFlatRateFixedRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			extraRule:            NewLinearSequentialRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(5 * time.Hour),
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 4 - Activation with overlapping extra rules
		"4-ActivationWithOverlappingExtraRules": {
			flatRateRule:         NewFlatRateFixedRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			extraRule:            NewLinearSequentialRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(8 * time.Hour),
			expectedActivation:   false,
			expectedActivationAt: 0,
		},
		// 5 - Activation with one extra rule fully in range
		"5-ActivationWithOneExtraRuleFullyInRangeDuplicate": {
			flatRateRule:         NewFlatRateFixedRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			solvedRules:          SolverRules{NewLinearSequentialRule("ExtraRule", 16*time.Hour, 1.0, MetaData{})},
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 6 - Activation with multiple extra rules
		"6-ActivationWithMultipleExtraRulesDuplicate": {
			flatRateRule:         NewFlatRateFixedRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			solvedRules:          SolverRules{NewLinearSequentialRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(4 * time.Hour)},
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 7 - No activation, extra rules do not meet the required amount
		"7-NoActivationExtraRulesInsufficientDuplicate": {
			flatRateRule:         NewFlatRateFixedRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			solvedRules:          SolverRules{NewLinearSequentialRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(5 * time.Hour)},
			expectedActivation:   true,
			expectedActivationAt: 10 * time.Hour,
		},
		// 8 - Activation with overlapping extra rules
		"8-ActivationWithOverlappingExtraRulesDuplicate": {
			flatRateRule:         NewFlatRateFixedRule("FlatRate", RelativeTimeSpan{7 * time.Hour, 11 * time.Hour}, 3.0, MetaData{}),
			solvedRules:          SolverRules{NewLinearSequentialRule("ExtraRule", 16*time.Hour, 1.0, MetaData{}).Shift(8 * time.Hour)},
			expectedActivation:   false,
			expectedActivationAt: 0,
		},
		// 9 - Activation with one extra rule fully in range
		"9-ActivationWithOneExtraRuleFullyInRange": {
			flatRateRule: NewFlatRateFixedRule("FlatRate", RelativeTimeSpan{6 * time.Hour, 20 * time.Hour}, 8.0, MetaData{}),
			solvedRules: SolverRules{
				NewLinearSequentialRule("A", 10*time.Hour, 1.0, MetaData{}),
				NewLinearSequentialRule("B", 1*time.Hour, 0, MetaData{}).Shift(10 * time.Hour),
				NewLinearSequentialRule("C", 10*time.Hour, 1.0, MetaData{}).Shift(11 * time.Hour),
			},
			expectedActivation:   true,
			expectedActivationAt: 15 * time.Hour,
		},
		// 10 - Activation with one extra rule fully in range
		"10-ActivationWithOneExtraRuleFullyInRange": {
			flatRateRule: NewFlatRateFixedRule("FlatRate", RelativeTimeSpan{6 * time.Hour, 20 * time.Hour}, 8.0, MetaData{}),
			solvedRules: SolverRules{
				NewLinearSequentialRule("A", 10*time.Hour, 1.0, MetaData{}),
				NewLinearSequentialRule("B", 1*time.Hour, 0, MetaData{}).Shift(10 * time.Hour),
			},
			extraRule:            NewLinearSequentialRule("C", 10*time.Hour, 1.0, MetaData{}).Shift(11 * time.Hour),
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
