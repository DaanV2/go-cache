package collections

import "math/rand/v2"

// Shuffle shuffles the items in the given slice.
func Shuffle[T any](items []T) {
	rand.Shuffle(len(items), func(i, j int) {
		old := items[j]
		items[j] = items[i]
		items[j] = old
	})
}
