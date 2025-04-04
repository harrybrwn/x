package xiter

import (
	"slices"
	"testing"
)

func TestWindow(t *testing.T) {
	for _, tt := range []struct {
		size     int
		in       []int
		expected [][]int
	}{
		{
			3,
			[]int{7, 19, 3, 0, 14, 17, 10, 1, 2, 8, 5},
			[][]int{
				{7, 19, 3},
				{19, 3, 0},
				{3, 0, 14},
				{0, 14, 17},
				{14, 17, 10},
				{17, 10, 1},
				{10, 1, 2},
				{1, 2, 8},
				{2, 8, 5}},
		},
		{
			3,
			[]int{7, 19, 3, 0, 14, 17, 10, 1, 2, 8, 5, 4},
			[][]int{
				{7, 19, 3},
				{19, 3, 0},
				{3, 0, 14},
				{0, 14, 17},
				{14, 17, 10},
				{17, 10, 1},
				{10, 1, 2},
				{1, 2, 8},
				{2, 8, 5},
				{8, 5, 4}},
		},
		{
			4,
			[]int{7, 19, 3, 0, 14, 17, 10, 1, 2, 8, 5, 4},
			[][]int{
				{7, 19, 3, 0},
				{19, 3, 0, 14},
				{3, 0, 14, 17},
				{0, 14, 17, 10},
				{14, 17, 10, 1},
				{17, 10, 1, 2},
				{10, 1, 2, 8},
				{1, 2, 8, 5},
				{2, 8, 5, 4}},
		},
	} {
		windows := slices.Collect(Window(tt.size, Iter(tt.in)))
		for i, exp := range tt.expected {
			if !slices.Equal(exp, windows[i]) {
				t.Errorf("at %d: got %v, want %v", i, windows[i], exp)
			}
		}
	}
}

func TestChunk(t *testing.T) {
	for _, tt := range []struct {
		in  []int
		n   int
		exp [][]int
	}{
		{
			[]int{7, 19, 3, 0, 14, 17, 10, 1, 2, 8, 5, 4},
			3,
			[][]int{
				{7, 19, 3},
				{0, 14, 17},
				{10, 1, 2},
				{8, 5, 4},
			},
		},
		{
			[]int{7, 19, 3, 0, 14, 17, 10, 1, 2, 8, 5},
			3,
			[][]int{
				{7, 19, 3},
				{0, 14, 17},
				{10, 1, 2},
				{8, 5},
			},
		},
	} {
		chunks := slices.Collect(Chunk(Iter(tt.in), tt.n))
		for i, exp := range tt.exp {
			if !slices.Equal(exp, chunks[i]) {
				t.Errorf("at %d: got %v, want %v", i, chunks[i], exp)
			}
		}
	}
}
