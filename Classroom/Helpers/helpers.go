package main

import "fmt"

func main() {
	// Example usage of Map
	nums := []int{1, 2, 3, 4}
	squared := Map(nums, func(n int) int {
		return n * n
	})
	fmt.Println(squared) // Output: [1 4 9 16]

	// Example usage of Filter
	nums2 := []int{1, 2, 3, 4, 5, 6}
	even := Filter(nums2, func(n int) bool {
		return n%2 == 0
	})
	fmt.Println(even) // Output: [2 4 6]
}

// Map applies a function to each item in a slice and returns a new slice with the results.
func Map[T any, U any](s []T, f func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

// Filter returns a new slice holding only the elements of s that satisfy f()
func Filter[T any](s []T, f func(T) bool) []T {
	var result []T
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
