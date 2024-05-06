package main

import (
	"testing"
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
	requestValues := RequestValues{
		Euro200: [5]int{1,1,1,1,1},
		Euro100: [5]int{1,1,1,1,1},
		Euro50:  [5]int{1,1,1,1,1},
		Euro20:  [5]int{1,1,1,1,1},
		Euro10:  [5]int{1,1,1,1,1},
		Euro5:   [5]int{1,1,1,1,1},
		Euro2:   [5]int{1,1,1,1,1},
		Euro1:   [5]int{1,1,1,1,1},
		Cent50:  [5]int{1,1,1,1,1},
		Cent20:  [5]int{1,1,1,1,1},
		Cent10:  [5]int{1,1,1,1,1},
		Cent5:   [5]int{1,1,1,1,1},
		Cent2:   [5]int{1,1,1,1,1},
		Cent1:   [5]int{1,1,1,1,1},
	}
	got := calculateDailyValues(requestValues)
	want := 1944.3999999999999
	if got != want {
		t.Errorf("calculateDailyValues(%v) = %f; want %f", requestValues, got, want)
	}
}
