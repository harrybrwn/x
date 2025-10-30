package xmath

import (
	"testing"

	"github.com/harrybrwn/x/random"
)

func TestAverage(t *testing.T) {
	for range 10 {
		arr := random.ArrayRange(256, 1, 1000)
		a := Round(Avg(arr), 6)
		b := Round(new(StreamAvg[int]).
			AppendSlice(arr).
			Value(), 6)
		if a != b {
			t.Errorf("%v != %v", a, b)
		}
	}
}

func TestStreamAvg_Count(t *testing.T) {
	var ma StreamAvg[int]
	for i := range 10 {
		ma.Append(i)
		if ma.Count() != uint64(i+1) {
			t.Errorf("expected %d", i+1)
		}
	}
}

func Avg[T random.Number](arr []T) float64 {
	var sum T
	for _, n := range arr {
		sum += n
	}
	return float64(sum) / float64(len(arr))
}
