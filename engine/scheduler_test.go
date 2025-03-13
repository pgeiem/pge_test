package engine

import (
	"testing"
)

func TestSchedulerAppend(t *testing.T) {
	tests := []struct {
		name     string
		entries  []SchedulerEntry
		expected []SchedulerEntry
	}{
		{
			name: "Non-overlapping entries",
			entries: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq1"}},
				{RelativeTimeSpan{From: 20, To: 30}, &TariffSequence{Name: "Seq2"}},
				{RelativeTimeSpan{From: 40, To: 50}, &TariffSequence{Name: "Seq3"}},
			},
			expected: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq1"}},
				{RelativeTimeSpan{From: 20, To: 30}, &TariffSequence{Name: "Seq2"}},
				{RelativeTimeSpan{From: 40, To: 50}, &TariffSequence{Name: "Seq3"}},
			},
		},
		{
			name: "Overlapping entries",
			entries: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq1"}},
				{RelativeTimeSpan{From: 20, To: 30}, &TariffSequence{Name: "Seq2"}},
				{RelativeTimeSpan{From: 5, To: 25}, &TariffSequence{Name: "Seq3"}},
			},
			expected: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq1"}},
				{RelativeTimeSpan{From: 10, To: 20}, &TariffSequence{Name: "Seq3"}},
				{RelativeTimeSpan{From: 20, To: 30}, &TariffSequence{Name: "Seq2"}},
			},
		},
		{
			name: "Completely overlapping entry",
			entries: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq1"}},
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq2"}},
			},
			expected: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq1"}},
			},
		},
		{
			name: "Entry within another entry",
			entries: []SchedulerEntry{
				{RelativeTimeSpan{From: 10, To: 20}, &TariffSequence{Name: "Seq2"}},
				{RelativeTimeSpan{From: 0, To: 30}, &TariffSequence{Name: "Seq1"}},
			},
			expected: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq1"}},
				{RelativeTimeSpan{From: 10, To: 20}, &TariffSequence{Name: "Seq2"}},
				{RelativeTimeSpan{From: 20, To: 30}, &TariffSequence{Name: "Seq1"}},
			},
		},
		{
			name: "Multiple overlapping entries",
			entries: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq1"}},
				{RelativeTimeSpan{From: 5, To: 15}, &TariffSequence{Name: "Seq2"}},
				{RelativeTimeSpan{From: 10, To: 20}, &TariffSequence{Name: "Seq3"}},
			},
			expected: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq1"}},
				{RelativeTimeSpan{From: 10, To: 15}, &TariffSequence{Name: "Seq2"}},
				{RelativeTimeSpan{From: 15, To: 20}, &TariffSequence{Name: "Seq3"}},
			},
		},
		{
			name: "Entry completely before another entry",
			entries: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq1"}},
				{RelativeTimeSpan{From: 10, To: 20}, &TariffSequence{Name: "Seq2"}},
			},
			expected: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq1"}},
				{RelativeTimeSpan{From: 10, To: 20}, &TariffSequence{Name: "Seq2"}},
			},
		},
		{
			name: "Entry completely after another entry",
			entries: []SchedulerEntry{
				{RelativeTimeSpan{From: 10, To: 20}, &TariffSequence{Name: "Seq1"}},
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq2"}},
			},
			expected: []SchedulerEntry{
				{RelativeTimeSpan{From: 0, To: 10}, &TariffSequence{Name: "Seq2"}},
				{RelativeTimeSpan{From: 10, To: 20}, &TariffSequence{Name: "Seq1"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scheduler := NewScheduler()
			for _, entry := range tt.entries {
				scheduler.Append(entry)
			}

			var got []SchedulerEntry
			scheduler.entries.Ascend(func(entry SchedulerEntry) bool {
				got = append(got, entry)
				return true
			})

			if len(got) != len(tt.expected) {
				t.Errorf("expected %d entries, got %d", len(tt.expected), len(got))
			}

			for i, entry := range got {
				if entry.From != tt.expected[i].From || entry.To != tt.expected[i].To || entry.Sequence.Name != tt.expected[i].Sequence.Name {
					t.Errorf("expected entry %v, got %v", tt.expected[i], entry)
				}
			}
		})
	}
}
