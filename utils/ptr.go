// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
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

