// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package utils provides utility functions for converting values to SQL-compatible strings
package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// MakeValue converts a value to a SQL-compatible string representation
// Supports various types including bool, string, and time.Time
// Returns the SQL string and a boolean indicating success
//
// @param value The value to convert
// @return A SQL-compatible string and a boolean indicating success
func MakeValue(value interface{}) (string, bool) {
	switch v := value.(type) {
	case float32:
		// Convert float32 to SQL-compatible string, removing trailing zeros
		result := fmt.Sprintf("%f", v)
		return result[:len(result)-2], true
	case float64:
		// Convert float64 to SQL-compatible string
		return fmt.Sprintf("%f", v), true
	case int, int8, int16, int32, int64:
		// Convert integer to SQL-compatible string
		return fmt.Sprintf("%d", v), true
	case bool:
		// Convert boolean to SQL-compatible string (1 for true, 0 for false)
		if v == true {
			return "1", true
		}
		return "0", true
	case string:
		// Convert string, escaping single quotes for SQL
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''")), true
	case time.Time:
		// Convert time.Time to SQL datetime string
		return fmt.Sprintf("'%s'", TimeToSQL(v)), true
	}
	// Named types (e.g. type Status int) miss every case in the type
	// switch; fall back to their underlying kind so they convert like
	// the equivalent builtin.
	if arg, ok := MakeValueByKind(value); ok {
		return arg, true
	}
	// Return empty string and false for unsupported types
	return "", false
}

// MakeValueByKind converts a value to a SQL-compatible string by its
// reflect.Kind rather than its exact type, so named types (enums and
// similar aliases with methods) are handled like the builtins they are
// based on.
func MakeValueByKind(value interface{}) (string, bool) {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return "", false
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", rv.Uint()), true
	case reflect.Bool:
		if rv.Bool() {
			return "1", true
		}
		return "0", true
	case reflect.String:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(rv.String(), "'", "''")), true
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", rv.Float()), true
	}
	return "", false
}

// MakeArg converts a value to a database-bindable argument suitable for use
// with parameterized queries. It mirrors the type support of MakeValue but
// returns values in a form the database drivers can bind safely, rather than
// as SQL literal text (so no quoting or escaping is required).
// Returns the argument and a boolean indicating success.
//
// @param value The value to convert
// @return A bindable argument and a boolean indicating success
func MakeArg(value interface{}) (interface{}, bool) {
	switch v := value.(type) {
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v, true
	case bool:
		// Match MakeValue semantics: booleans are stored as 1 or 0
		if v {
			return 1, true
		}
		return 0, true
	case string:
		return v, true
	case time.Time:
		// Bind as the same text representation previously interpolated into
		// queries so populateModel's SQLToTime conversion keeps working
		return TimeToSQL(v), true
	}
	// Named types (e.g. type Viewed int) miss every case in the type
	// switch; fall back to their underlying kind so they bind like the
	// equivalent builtin.
	return makeArgByKind(value)
}

// makeArgByKind converts a value to a bindable argument by its
// reflect.Kind rather than its exact type. Integer kinds are widened to
// int64 so database drivers receive a concrete builtin type.
func makeArgByKind(value interface{}) (interface{}, bool) {
	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return nil, false
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int(), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(rv.Uint()), true
	case reflect.Bool:
		if rv.Bool() {
			return 1, true
		}
		return 0, true
	case reflect.String:
		return rv.String(), true
	case reflect.Float32, reflect.Float64:
		return rv.Float(), true
	}
	return nil, false
}
