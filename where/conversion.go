// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package where provides functionality for building SQL WHERE clauses
package where

import "reflect"

// convertToInterfaceArray converts an input value to a slice of interface{}
// If the input is already a slice or array, it converts each element to interface{}
// If the input is a single value, it wraps it in a slice
// If the input is nil or an empty slice/array, it returns nil
//
// @param values The input value to convert
// @return A slice of interface{} containing the converted values
func convertToInterfaceArray(values interface{}) []interface{} {
	if values == nil {
		return nil
	}
	v := reflect.ValueOf(values)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		for v.Kind() == reflect.Pointer {
			newValue := getValue(values)
			if newValue == nil {
				return nil
			}
			return []interface{}{newValue}
		}
		return []interface{}{values}
	}
	if v.Len() == 0 {
		return nil
	}
	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}
	return result
}

// getValue retrieves the value pointed to by a pointer
// This is a helper function used by convertToInterfaceArray
//
// @param value The pointer value to dereference
// @return The dereferenced value or nil if the pointer is nil
func getValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		return v.Elem().Interface()
	}
	return value
}

