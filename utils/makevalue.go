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
