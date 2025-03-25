package array

import (
	"fmt"
	"iter"
	"slices"
)

// Append will merge all elements of a number of slices.
//
// Deprecated: Use [slices.Concat] instead, its also more efficient.
func Append[T any, S ~[]T](arrs ...S) S {
	res := make([]T, 0, len(arrs))
	for _, arr := range arrs {
		res = append(res, arr...)
	}
	return res
}

// Reverse a slice.
//
// Deprecated: Use [slices.Reverse] instead.q
func Reverse[T any, S ~[]T](s S) {
	l := len(s)
	m := l / 2
	j := 0
	for i := 0; i < m; i++ {
		j = l - i - 1
		s[i], s[j] = s[j], s[i]
	}
}

// Map will run a function on each element of a slice and collect the output.
func Map[I, O any, In ~[]I](arr In, fn func(I) O) []O {
	res := make([]O, len(arr))
	for i, v := range arr {
		res[i] = fn(v)
	}
	return res
}

// MapErr is the same as [Map] but will handle errors returned by the mapping
// function.
func MapErr[I, O any, S ~[]I](s S, fn func(I) (O, error)) (res []O, err error) {
	res = make([]O, len(s))
	for i, v := range s {
		res[i], err = fn(v)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// ForEach is like [Map] but it doesn't collect any function output.
func ForEach[I any, S ~[]I](s S, fn func(I)) {
	for _, v := range s {
		fn(v)
	}
}

// MapStringers will convert a slice of [fmt.Stringer] to a slice of strings.
func MapStringers[Slice ~[]T, T fmt.Stringer](s Slice) []string {
	return Map(s, func(e T) string { return e.String() })
}

func ToString[T fmt.Stringer](v T) string { return v.String() }

func FilterMap[I, O any](in iter.Seq[I], fn func(*I) (O, bool)) iter.Seq[O] {
	return func(yield func(O) bool) {
		for item := range in {
			o, ok := fn(&item)
			if ok && !yield(o) {
				return
			}
		}
	}
}

func Iter[T any, S ~[]T](s S) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range s {
			if !yield(v) {
				return
			}
		}
	}
}

func IterRef[T any, S ~[]T](s S) iter.Seq[*T] {
	return func(yield func(*T) bool) {
		for i := range s {
			if !yield(&s[i]) {
				return
			}
		}
	}
}

func IterDeref[T any, S ~[]T](s []*T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range s {
			v := *s[i]
			if !yield(v) {
				return
			}
		}
	}
}

func Move[T any, S ~[]T](array S, srcIndex int, dstIndex int) []T {
	value := array[srcIndex]
	return slices.Insert(remove(array, srcIndex), dstIndex, value)
}

func MoveToFront[T any, S ~[]T](arr S, fn func(T) bool) {
	for i := 0; i < len(arr); i++ {
		if fn(arr[i]) {
			Move(arr, i, 0)
			return
		}
	}
}

func remove[T any, A ~[]T](array A, index int) []T {
	return append(array[:index], array[index+1:]...)
}

func Remove[T any, A ~[]T](arr A, index int) A {
	l := len(arr) - 1
	arr[index] = arr[l]
	return arr[:l]
}

// ToAny will convert a slice of some type T to a slice of type any.
func ToAny[T any, S ~[]T](s S) []any {
	res := make([]any, len(s))
	for i, v := range s {
		res[i] = v
	}
	return res
}

// All will return true if all the values in the slice are true.
func All(vals []bool) bool {
	for _, v := range vals {
		if !v {
			return false
		}
	}
	return true
}

// Contains will check a slice of slices and return true if it contains a given
// slice.
//
// This is a brute force implementation with O(n^2) runtime.
func Contains[T comparable, S ~[]T](list []S, s S) bool {
	for _, sub := range list {
		if slices.Equal(sub, s) {
			return true
		}
	}
	return false
}
