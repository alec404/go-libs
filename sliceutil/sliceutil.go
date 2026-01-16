package sliceutil

// UniqueSlice removes duplicate elements from a slice while preserving order.
// It uses a map to track seen elements, ensuring O(n) time complexity.
// Example: UniqueSlice([]int{1, 2, 2, 3, 1}) -> [1, 2, 3]
func UniqueSlice[K comparable](slice []K) []K {
	result := make([]K, 0, len(slice))
	seen := map[K]struct{}{}
	for _, item := range slice {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// Remove removes the first element from the slice that satisfies the predicate function.
// The slice is modified in-place by swapping the found element with the last element.
// Returns the modified slice with length reduced by 1.
// Example: Remove([]interface{}{1, 2, 3}, func(v interface{}) bool { return v == 2 }) -> [1, 3]
func Remove(sl []interface{}, f func(v interface{}) bool) []interface{} {
	for k, v := range sl {
		if f(v) {
			sl[k] = sl[len(sl)-1]
			sl = sl[:len(sl)-1]
			return sl
		}
	}
	return sl
}

// RemoveSlice removes the first occurrence of a specific element from the slice.
// It maintains the order of remaining elements by shifting elements left.
// Returns a new slice with the element removed, or the original slice if not found.
// Example: RemoveSlice([]int{1, 2, 3, 2}, 2) -> [1, 3, 2]
func RemoveSlice[K comparable](src []K, target K) []K {
	for k, v := range src {
		if v == target {
			copy(src[k:], src[k+1:])
			return src[:len(src)-1]
		}
	}
	return src
}

// DifferenceSlice returns elements that are in s2 but not in s1.
// It creates a map from s1 for O(1) lookup, then iterates through s2.
// Example:
//
//	s1 := []int{1, 2, 3, 4, 5}
//	s2 := []int{4, 5, 6, 7, 8}
//	DifferenceSlice(s1, s2) -> [6, 7, 8]
func DifferenceSlice[T comparable](s1, s2 []T) []T {
	m := make(map[T]bool)
	for _, item := range s1 {
		m[item] = true
	}

	var diff []T
	for _, item := range s2 {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

// InSlice checks if a target element exists in a slice.
// Returns true if found, false otherwise.
// Example: InSlice(2, []int{1, 2, 3}) -> true
func InSlice[T comparable](target T, slice []T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}
