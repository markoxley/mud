// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package where provides functionality for building SQL WHERE clauses
package where

import "time"

// DBInt represents all integer types supported by the database
// This includes both signed and unsigned integers of various sizes
type DBInt interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

// DBFloat represents floating-point number types supported by the database
type DBFloat interface {
	float32 | float64
}

// DBNumeric is a union type that includes all numeric types
// (both integers and floating-point numbers) supported by the database
type DBNumeric interface {
	DBInt | DBFloat
}

// DBField represents all valid database field types
// This includes numeric types, boolean, string, and time.Time
type DBField interface {
	DBNumeric | bool | string | time.Time
}

// Database type constants used for type checking and validation
const (
	dBool = 1 << iota // Boolean type
	dDate             // Date/time type
	dFloat            // Single-precision floating point
	dDouble           // Double-precision floating point
	dInt              // Integer types
	dLong             // Long integer types
	dText             // Text/string type
)

