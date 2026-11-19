// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package utils provides utility functions for converting values to SQL-compatible strings
package utils

import (
	"fmt"
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
	// Return empty string and false for unsupported types
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
	// Return nil and false for unsupported types
	return nil, false
}
