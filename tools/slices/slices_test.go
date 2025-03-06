package slices

import (
	"math"
	"testing"
)

var equalIntTests = []struct {
	s1, s2 []int
	want   bool
}{
	{
		[]int{1},
		nil,
		false,
	},
	{
		[]int{},
		nil,
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3},
		true,
	},
	{
		[]int{1, 2, 3},
		[]int{1, 2, 3, 4},
		false,
	},
}

var equalFloatTests = []struct {
	s1, s2       []float64
	wantEqual    bool
	wantEqualNaN bool
}{
	{
		[]float64{1, 2},
		[]float64{1, 2},
		true,
		true,
	},
	{
		[]float64{1, 2, math.NaN()},
		[]float64{1, 2, math.NaN()},
		false,
		true,
	},
}

var filterTests = []struct {
	s    []int
	fn   func(int) bool
	want []int
}{
	{
		nil,
		func(int) bool { return false },
		nil,
	},
	{
		[]int{1, 2, 3},
		func(int) bool { return false },
		nil,
	},
	{
		[]int{1, 2, 3},
		func(int) bool { return true },
		[]int{1, 2, 3},
	},
	{
		[]int{1, 2, 3},
		func(i int) bool { return i <= 2 },
		[]int{1, 2},
	},
	{
		[]int{1, 2, 3},
		func(i int) bool { return i >= 2 },
		[]int{2, 3},
	},
	{
		[]int{10, 2, 30},
		func(i int) bool { return i < 10 },
		[]int{2},
	},
}

func TestClone(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := Clone(s1)
	if !Equal(s1, s2) {
		t.Errorf("Clone(%v) = %v, want %v", s1, s2, s1)
	}
	s1[0] = 4
	want := []int{1, 2, 3}
	if !Equal(s2, want) {
		t.Errorf("Clone(%v) changed unexpectedly to %v", want, s2)
	}
	if got := Clone([]int(nil)); got != nil {
		t.Errorf("Clone(nil) = %#v, want nil", got)
	}
	if got := Clone(s1[:0]); got == nil || len(got) != 0 {
		t.Errorf("Clone(%v) = %#v, want %#v", s1[:0], got, s1[:0])
	}
}

func TestEqual(t *testing.T) {
	for _, test := range equalIntTests {
		if got := Equal(test.s1, test.s2); got != test.want {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, got, test.want)
		}
	}
	for _, test := range equalFloatTests {
		if got := Equal(test.s1, test.s2); got != test.wantEqual {
			t.Errorf("Equal(%v, %v) = %t, want %t", test.s1, test.s2, got, test.wantEqual)
		}
	}
}

func TestFilter(t *testing.T) {
	for i, test := range filterTests {
		copied := Clone(test.s)
		if got := Filter(copied, test.fn); !Equal(got, test.want) {
			t.Errorf("Filter case %d: got %v, want %v", i, got, test.want)
		}
	}
}
