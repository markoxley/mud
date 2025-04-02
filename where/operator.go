// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package where provides functionality for building SQL WHERE clauses
package where

// operator represents the type of comparison operation in a WHERE clause
type operator uint8

// Operator constants define the supported comparison operations
const (
	opEqual   operator = iota // Equal (=) comparison
	opGreater               // Greater than (>) comparison
	opLess                  // Less than (<) comparison
	opLike                  // LIKE pattern matching
	opIn                    // IN clause for multiple values
	opBetween               // BETWEEN range comparison
	opIsNull                // IS NULL check
)

// operatorType defines the valid data types for each operator
// Each element is a bitmask of compatible data types
var operatorType [7]int = [7]int{
	dBool & dDate & dFloat & dDouble & dInt & dLong & dText,    // Equal supports all types
	dDate & dFloat & dDouble & dInt & dLong & dText,            // Greater than supports numeric and text types
	dDate & dFloat & dDouble & dInt & dLong & dText,            // Less than supports numeric and text types
	dText,                                                      // LIKE only supports text
	dDate & dFloat & dDouble & dInt & dLong & dText,            // IN supports all types except bool
	dDate & dFloat & dDouble & dInt & dLong,                    // BETWEEN supports numeric types
	dBool & dDate & dFloat & dDouble & dInt & dLong & dText,    // IS NULL supports all types
}

