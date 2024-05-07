package main

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAdd tests the Add function
func TestAdd(t *testing.T) {
	got := Add(1, 2)
	want := 3
	if got != want {
		t.Errorf("Add(1, 2) = %d; want %d", got, want)
	}
}

func TestSumArray(t *testing.T) {
	got := SumArray([]int{1, 2, 3, 4, 5})
	want := 15

	if got != want {
		t.Errorf("SumArray([1,2,3,4,5]) = %d; want %d", got, want)
	}
}
func TestCalculateDailyValues(t *testing.T) {
	type test struct {
		name     string
		input    RequestValues
		expected float64
	}

	tests := []test{
		{
			name: "All coins and bills present",
			input: RequestValues{
				Euro200: [5]int{1, 2, 3, 4, 5},
				Euro100: [5]int{1, 2, 3, 4, 5},
				Euro50:  [5]int{1, 2, 3, 4, 5},
				Euro20:  [5]int{1, 2, 3, 4, 5},
				Euro10:  [5]int{1, 2, 3, 4, 5},
				Euro5:   [5]int{1, 2, 3, 4, 5},
				Euro2:   [5]int{1, 2, 3, 4, 5},
				Euro1:   [5]int{1, 2, 3, 4, 5},
				Cent50:  [5]int{1, 2, 3, 4, 5},
				Cent20:  [5]int{1, 2, 3, 4, 5},
				Cent10:  [5]int{1, 2, 3, 4, 5},
				Cent5:   [5]int{1, 2, 3, 4, 5},
				Cent2:   [5]int{1, 2, 3, 4, 5},
				Cent1:   [5]int{1, 2, 3, 4, 5},
			},
			expected: 5833.2,
		},
		{
			name: "No coins or bills present",
			input: RequestValues{
				Euro200: [5]int{0, 0, 0, 0, 0},
				Euro100: [5]int{0, 0, 0, 0, 0},
				Euro50:  [5]int{0, 0, 0, 0, 0},
				Euro20:  [5]int{0, 0, 0, 0, 0},
				Euro10:  [5]int{0, 0, 0, 0, 0},
				Euro5:   [5]int{0, 0, 0, 0, 0},
				Euro2:   [5]int{0, 0, 0, 0, 0},
				Euro1:   [5]int{0, 0, 0, 0, 0},
				Cent50:  [5]int{0, 0, 0, 0, 0},
				Cent20:  [5]int{0, 0, 0, 0, 0},
				Cent10:  [5]int{0, 0, 0, 0, 0},
				Cent5:   [5]int{0, 0, 0, 0, 0},
				Cent2:   [5]int{0, 0, 0, 0, 0},
				Cent1:   [5]int{0, 0, 0, 0, 0},
			},
			expected: 0.0,
		},
		{
			name: "Only cents",
			input: RequestValues{
				Euro200: [5]int{0, 0, 0, 0, 0},
				Euro100: [5]int{0, 0, 0, 0, 0},
				Euro50:  [5]int{0, 0, 0, 0, 0},
				Euro20:  [5]int{0, 0, 0, 0, 0},
				Euro10:  [5]int{0, 0, 0, 0, 0},
				Euro5:   [5]int{0, 0, 0, 0, 0},
				Euro2:   [5]int{0, 0, 0, 0, 0},
				Euro1:   [5]int{0, 0, 0, 0, 0},
				Cent50:  [5]int{1, 2, 3, 4, 5},
				Cent20:  [5]int{1, 2, 3, 4, 5},
				Cent10:  [5]int{1, 2, 3, 4, 5},
				Cent5:   [5]int{1, 2, 3, 4, 5},
				Cent2:   [5]int{1, 2, 3, 4, 5},
				Cent1:   [5]int{1, 2, 3, 4, 5},
			},
			expected: 13.2,
		},
		{
			name: "Only 50 cents",
			input: RequestValues{
				Euro200: [5]int{0, 0, 0, 0, 0},
				Euro100: [5]int{0, 0, 0, 0, 0},
				Euro50:  [5]int{0, 0, 0, 0, 0},
				Euro20:  [5]int{0, 0, 0, 0, 0},
				Euro10:  [5]int{0, 0, 0, 0, 0},
				Euro5:   [5]int{0, 0, 0, 0, 0},
				Euro2:   [5]int{0, 0, 0, 0, 0},
				Euro1:   [5]int{0, 0, 0, 0, 0},
				Cent50:  [5]int{1, 2, 3, 4, 5},
				Cent20:  [5]int{0, 0, 0, 0, 0},
				Cent10:  [5]int{0, 0, 0, 0, 0},
				Cent5:   [5]int{0, 0, 0, 0, 0},
				Cent2:   [5]int{0, 0, 0, 0, 0},
				Cent1:   [5]int{0, 0, 0, 0, 0},
			},
			expected: 7.5,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if result := CalculateDailyValues(tc.input); !almostEqual(result, tc.expected) {
				t.Fatalf("CalculateDailyValues() returned %v, want %v", result, tc.expected)
			}
		})
	}
}

// Comparison of two floats is not straightforward so we're using a tolerance based comparison
func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= 1e-6
}

func TestCalculateRollValues(t *testing.T) {
	type test struct {
		name     string
		input    RollValues
		expected float64
	}

	tests := []test{
		{
			name: "All zero values",
			input: RollValues{
				Euro2:  [2]int{0, 0},
				Euro1:  [2]int{0, 0},
				Cent50: [2]int{0, 0},
				Cent20: [2]int{0, 0},
				Cent10: [2]int{0, 0},
				Cent5:  [2]int{0, 0},
				Cent2:  [2]int{0, 0},
				Cent1:  [2]int{0, 0},
			},
			expected: 0.0,
		},
		{
			name: "Only one roll per value",
			input: RollValues{
				Euro2:  [2]int{1, 0},
				Euro1:  [2]int{1, 0},
				Cent50: [2]int{1, 0},
				Cent20: [2]int{1, 0},
				Cent10: [2]int{1, 0},
				Cent5:  [2]int{1, 0},
				Cent2:  [2]int{1, 0},
				Cent1:  [2]int{1, 0},
			},
			expected: 111.0,
		},
		{
			name: "All one values",
			input: RollValues{
				Euro2:  [2]int{1, 1},
				Euro1:  [2]int{1, 1},
				Cent50: [2]int{1, 1},
				Cent20: [2]int{1, 1},
				Cent10: [2]int{1, 1},
				Cent5:  [2]int{1, 1},
				Cent2:  [2]int{1, 1},
				Cent1:  [2]int{1, 1},
			},
			expected: 222.0,
		},
		{
			name: "Random values",
			input: RollValues{
				Euro2:  [2]int{3, 2},
				Euro1:  [2]int{1, 4},
				Cent50: [2]int{5, 6},
				Cent20: [2]int{7, 8},
				Cent10: [2]int{9, 10},
				Cent5:  [2]int{10, 9},
				Cent2:  [2]int{8, 7},
				Cent1:  [2]int{6, 5},
			},
			expected: 859.00,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateRollValues(tt.input)

			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}
func TestCalculateBoxValues(t *testing.T) {
	cases := []struct {
		name     string
		box      BoxValues
		expected float64
	}{
		{"All boxes are empty", BoxValues{Euro2: [1]int{0}, Euro1: [1]int{0}, Cent50: [1]int{0}, Cent20: [1]int{0}, Cent10: [1]int{0}, Cent2: [1]int{0}, Cent1: [1]int{0}}, 0},
		{"All boxes are full", BoxValues{Euro2: [1]int{1}, Euro1: [1]int{1}, Cent50: [1]int{1}, Cent20: [1]int{1},
			Cent10: [1]int{1}, Cent5: [1]int{1}, Cent2: [1]int{1}, Cent1: [1]int{1}}, 336.0},
		{"Only Euro2 boxes are full", BoxValues{Euro2: [1]int{1}, Euro1: [1]int{0}, Cent50: [1]int{0},
			Cent20: [1]int{0}, Cent10: [1]int{0}, Cent2: [1]int{0}, Cent1: [1]int{0}}, 150},
		{"Only Euro1 boxes are full", BoxValues{Euro2: [1]int{0}, Euro1: [1]int{1}, Cent50: [1]int{0},
			Cent20: [1]int{0}, Cent10: [1]int{0}, Cent2: [1]int{0}, Cent1: [1]int{0}}, 75},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := CalculateBoxValues(c.box)
			if result != c.expected {
				t.Errorf("Expected %f but got %f", c.expected, result)
			}
		})
	}
}
func TestCalculateValuesForCashCounts(t *testing.T) {
	var tests = []struct {
		input          RequestPayload
		wantTotalValue float64
		wantBoxValue   float64
		wantRollValue  float64
		wantDiffValue  float64
	}{
		{RequestPayload{}, 0, 0, 0, 0},
		{RequestPayload{
			RequestValidation: RequestValidation{},
			RequestValues:     RequestValues{},
			BoxValues: BoxValues{
				Euro2: [1]int{1},
			},
			RollValues:  RollValues{},
			PayloadType: 1,
		}, 0, 150, 0, 150},
		{RequestPayload{
			RequestValidation: RequestValidation{
				TargetValue: "253.00",
			},
			RequestValues: RequestValues{},
			BoxValues: BoxValues{
				Euro2: [1]int{1},
			},
			RollValues: RollValues{
				Euro2: [2]int{1, 1},
			},
			PayloadType: 1,
		}, 0, 150, 100, -3},
	}
	for _, tt := range tests {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got1, got2, got3, got4 := CalculateValuesForCashCounts(tt.input)
			if got1 != tt.wantTotalValue {
				t.Errorf("CalculateValuesForCashCounts() = %v, want %v", got1, tt.wantTotalValue)
			}
			if got2 != tt.wantBoxValue {
				t.Errorf("CalculateValuesForCashCounts() = %v, want %v", got2, tt.wantBoxValue)
			}
			if got3 != tt.wantRollValue {
				t.Errorf("CalculateValuesForCashCounts() = %v, want %v", got3, tt.wantRollValue)
			}
			if got4 != tt.wantDiffValue {
				t.Errorf("CalculateValuesForCashCounts() = %v, want %v", got4, tt.wantDiffValue)
			}
		})
	}
}

func TestCalculateTotalValue(t *testing.T) {
	tests := []struct {
		name     string
		input    RequestPayload
		expected ResponsePayload
	}{
		{
			name: "Test with default values",
			input: RequestPayload{
				RequestValidation: RequestValidation{},
				RequestValues:     RequestValues{},
				BoxValues:         BoxValues{},
				RollValues:        RollValues{},
				PayloadType:       0,
			},
			expected: ResponsePayload{
				ResponseValues: ResponseValues{
					TotalValue:      "0.00",
					DifferenceValue: "0.00",
				},
				PayloadType: 2,
			},
		},
		{
			name: "Test with custom values",
			input: RequestPayload{
				RequestValidation: RequestValidation{
					TargetValue: "50.00",
				},
				RequestValues: RequestValues{
					Euro200: [5]int{},
					Euro100: [5]int{},
					Euro50:  [5]int{},
					Euro20:  [5]int{},
					Euro10:  [5]int{10, 0, 0, 0, 0},
					Euro5:   [5]int{},
					Euro2:   [5]int{},
					Euro1:   [5]int{},
					Cent50:  [5]int{},
					Cent20:  [5]int{},
					Cent10:  [5]int{0, 0, 0, 0, 0},
					Cent5:   [5]int{},
					Cent2:   [5]int{},
					Cent1:   [5]int{},
				},
				BoxValues:   BoxValues{},
				RollValues:  RollValues{},
				PayloadType: 1,
			},
			expected: ResponsePayload{
				ResponseValues: ResponseValues{
					TotalValue:      "100.00",
					DifferenceValue: "50.00",
				},
				PayloadType: 2,
			},
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := calculateTotalValue(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
