package slices

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// Unique omits duplicate entries form slice s.
func Unique[T constraints.Ordered](slice []T) []T {
	if len(slice) < 1 {
		return slice
	}

	// sort slice
	sort.SliceStable(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})

	// shift unique to the left
	prev := 1
	for curr := 1; curr < len(slice); curr++ {
		if slice[curr-1] != slice[curr] {
			slice[prev] = slice[curr]
			prev++
		}
	}

	// return left of the slice containing unique entries only
	return slice[:prev]
}
