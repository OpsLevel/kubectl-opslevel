package common

import (
	"math/rand"
	"slices"
)

// GetSample returns a random selection of N items in the slice in the original order.
// The elements are copied using assignment, so this is a shallow clone.
func GetSample[T any](sampleCount int, data []T) []T {
	var (
		keys = make([]int, len(data))
		copy []T
	)
	if sampleCount < 1 || sampleCount >= len(data) {
		return slices.Clone(data)
	}
	for i := range keys {
		keys[i] = i
	}
	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})
	keys = keys[:sampleCount]
	slices.Sort(keys)
	copy = make([]T, sampleCount)
	for i := range keys {
		copy[i] = data[keys[i]]
	}
	return copy
}
