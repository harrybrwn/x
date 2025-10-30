package xmath

import (
	"iter"
	"math"
	"slices"

	"github.com/harrybrwn/x/types"
)

type StreamAvg[T types.Number] struct {
	count uint64
	avg   float64
}

func (ma *StreamAvg[T]) Value() float64                  { return ma.avg }
func (ma *StreamAvg[T]) Count() uint64                   { return ma.count }
func (ma *StreamAvg[T]) AppendSlice(s []T) *StreamAvg[T] { return ma.AppendIter(slices.Values(s)) }

func (ma *StreamAvg[T]) AppendIter(i iter.Seq[T]) *StreamAvg[T] {
	for v := range i {
		ma.Append(v)
	}
	return ma
}

func (ma *StreamAvg[T]) Append(v T) *StreamAvg[T] {
	ma.count++
	count := float64(ma.count)
	val := float64(v)
	ma.avg = (ma.avg*(count-1) + val) / count
	return ma
}

func Round(val float64, decimals int) float64 {
	pow := math.Pow10(decimals)
	return (math.Round(val) * pow) / pow
}
