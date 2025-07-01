package funcs

// Map returns a slice by performing transformFunc with each collection element
func Map[T any, R any](collection []T, transformFunc func(t T) R) []R {
	result := make([]R, 0, len(collection))
	for _, item := range collection {
		result = append(result, transformFunc(item))
	}
	return result
}

// Filter returns a slice contains values that passed the predicate
func Filter[V any](collection []V, predicate func(item V, index int) bool) []V {
	result := make([]V, 0, len(collection))

	for i, item := range collection {
		if predicate(item, i) {
			result = append(result, item)
		}
	}

	return result
}

// GroupBy returns a map composed of keys
// generated from the results of running each element of collection through iteratee.
func GroupBy[T any, U comparable](collection []T, iteratee func(item T) U) map[U][]T {
	result := map[U][]T{}

	for _, item := range collection {
		key := iteratee(item)
		result[key] = append(result[key], item)
	}

	return result
}

// AssociateBy returns a map with key provided by iteratee applied to elements of the given slice.
// If any of two pairs would have the same key the last one gets added to the map.
func AssociateBy[T any, K comparable](collection []T, iteratee func(item T) K) map[K]T {
	result := make(map[K]T, len(collection))

	for _, t := range collection {
		k := iteratee(t)
		result[k] = t
	}

	return result
}

// Chunk split a slice into groups the length of size, if slice can't be split evenly,
// the final chunk will be the remaining elements.
func Chunk[T any](collection []T, size int) [][]T {
	if size <= 0 {
		size = 1
	}

	chunksNum := len(collection) / size
	if len(collection)%size != 0 {
		chunksNum += 1
	}

	result := make([][]T, 0, chunksNum)

	for i := 0; i < chunksNum; i++ {
		last := (i + 1) * size
		if last > len(collection) {
			last = len(collection)
		}
		result = append(result, collection[i*size:last])
	}

	return result
}

// Reduce reduces collection to a value with accumulator
func Reduce[T any, R any](collection []T, accumulator func(agg R, item T, index int) R, initial R) R {
	for i, item := range collection {
		initial = accumulator(initial, item, i)
	}

	return initial
}

// UniqueElements returns a slice with unique elements
func UniqueElements[T comparable](inputSlices ...[]T) []T {
	numOfElements := 0
	for _, sl := range inputSlices {
		numOfElements += len(sl)
	}
	uniqueSlice := make([]T, 0, numOfElements)
	seen := make(map[T]bool, numOfElements)
	for _, sl := range inputSlices {
		for _, element := range sl {
			if !seen[element] {
				uniqueSlice = append(uniqueSlice, element)
				seen[element] = true
			}
		}
	}
	return uniqueSlice
}

// AnyEqual returns true if any element in collection is equal to value
func AnyEqual[T comparable](collection []T, value T) bool {
	for _, item := range collection {
		if item == value {
			return true
		}
	}
	return false
}
