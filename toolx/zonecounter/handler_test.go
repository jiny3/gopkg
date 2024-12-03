package zonecounter

import (
	"reflect"
	"testing"
)

func TestZones_Discretization(t *testing.T) {
	tests := []struct {
		name     string
		input    Zones
		expected Zones
	}{
		{
			name: "Simple case",
			input: Zones{
				{Begin: 1, End: 4, Value: 1},
				{Begin: 2, End: 8, Value: 1},
				{Begin: 3, End: 6, Value: 1},
			},
			expected: Zones{
				{Begin: 1, End: 2, Value: 1},
				{Begin: 2, End: 3, Value: 2},
				{Begin: 3, End: 4, Value: 3},
				{Begin: 4, End: 6, Value: 2},
				{Begin: 6, End: 8, Value: 1},
			},
		},
		{
			name: "No overlap",
			input: Zones{
				{Begin: 1, End: 3, Value: 1},
				{Begin: 3, End: 6, Value: 2},
			},
			expected: Zones{
				{Begin: 1, End: 3, Value: 1},
				{Begin: 3, End: 6, Value: 2},
			},
		},
		{
			name: "No overlap 2",
			input: Zones{
				{Begin: 1, End: 3, Value: 1},
				{Begin: 3, End: 6, Value: 1},
			},
			expected: Zones{
				{Begin: 1, End: 6, Value: 1},
			},
		},
		{
			name: "Full overlap",
			input: Zones{
				{Begin: 1, End: 5, Value: 2},
				{Begin: 1, End: 5, Value: 3},
			},
			expected: Zones{
				{Begin: 1, End: 5, Value: 5},
			},
		},
		{
			name:     "Edge case with no zones",
			input:    Zones{},
			expected: Zones{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.DiscretizationWithMerge()
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Discretization() = %v, want %v", got, tt.expected)
			}
		})
	}
}
