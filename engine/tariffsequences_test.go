package engine

/*
func RecurrentSegmentMustParse(start, end string) RecurrentSegment {
	seg, err := NewRecurrentSegmentFromPatterns(start, end)
	if err != nil {
		panic(err)
	}
	return seg
}

func TestResolveSequenceApplicability(t *testing.T) {

	inventory := TariffSequenceInventory{
		{
			Name:           "Evening",
			ValidityPeriod: RecurrentSegmentMustParse("pattern(2023/10/* 12:00:00)", "pattern(2023/10/* 18:00:00)"),
		},
		{
			Name:           "Morning",
			ValidityPeriod: RecurrentSegmentMustParse("pattern(2023/10/* 08:00:00)", "pattern(2023/10/* 14:00:00)"),
		},
		{
			Name:           "Default",
			ValidityPeriod: RecurrentSegment{},
		},
	}

	tests := []struct {
		name   string
		now    time.Time
		window time.Duration
		want   PrioritizedSequences
	}{
		{
			name:   "Test case 1",
			now:    time.Now(),
			window: 2 * time.Hour,
			want:   PrioritizedSequences{}, // Define the expected result
		},
		{
			name:   "Test case 2",
			now:    time.Now().Add(1 * time.Hour),
			window: 3 * time.Hour,
			want:   PrioritizedSequences{}, // Define the expected result
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := inventory.ResolveSequenceApplicability(tt.now, tt.window)
			if err != nil {
				t.Errorf("ResolveSequenceApplicability() error = %v", err)
				return
			}
			if !comparePrioritizedSequences(got, tt.want) {
				t.Errorf("ResolveSequenceApplicability() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare two PrioritizedSequences
func comparePrioritizedSequences(a, b PrioritizedSequences) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Sequence.Name != b[i].Sequence.Name {
			return false
		}
	}
	return true
}
*/
