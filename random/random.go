package random

import (
	"math/rand/v2"

	"github.com/harrybrwn/x/types"
)

type Number = types.Number

func Range[N types.Number](low, high N) N {
	l := int64(low)
	h := int64(high)
	return N(rand.Int64N(h-l) + l)
}

func Choice[T any, S ~[]T](s S) T {
	return s[rand.IntN(len(s))]
}

func Suffle[T any, S ~[]T](s S) {
	rand.Shuffle(
		len(s),
		func(i, j int) { s[i], s[j] = s[j], s[i] },
	)
}

func ArrayRange[T Number](size int, low, high T) []T {
	a := make([]T, size)
	for i := 0; i < size; i++ {
		a[i] = Range(low, high)
	}
	return a
}

func Array[T Number](size int) []T {
	a := make([]T, size)
	for i := 0; i < size; i++ {
		a[i] = T(rand.Int64())
	}
	return a
}
