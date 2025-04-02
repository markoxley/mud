// Package utils provides utility functions for working with pointers
package utils

// Ptr returns a pointer to the provided value
// This is a generic helper function that works with any type
//
// @param value The value to create a pointer to
// @return A pointer to the provided value
func Ptr[T any](value T) *T {
	return &value
}
