// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package dtorm provides a database ORM (Object-Relational Mapping) implementation
// with support for SQLite, MySQL, and SQL Server databases.
package dtorm

import "fmt"

// FieldSize represents the size and precision of a database field.
// It is used to store the maximum length and decimal places for numeric types.
type FieldSize struct {
	// Size represents the maximum length of the field
	Size int
	// Decimal represents the number of decimal places for numeric types
	Decimal int
}

// NewSize creates a new FieldSize with the specified size and decimal places.
// This function is used to define the size constraints for database fields.
// Parameters:
//   sz: The maximum length of the field
//   dec: The number of decimal places (for numeric types)
// Returns:
//   A new FieldSize instance
func NewSize(sz, dec int) FieldSize {
	return FieldSize{
		Size:    sz,
		Decimal: dec,
	}
}

// String returns a string representation of the field size.
// For numeric types with decimal places, it returns "size,decimal".
// For other types, it returns just the size.
// This method is used for generating SQL schema definitions.
// Returns:
//   A string representation of the field size
func (s FieldSize) String() string {
	if s.Decimal > 0 {
		return fmt.Sprintf("%d,%d", s.Size, s.Decimal)
	}
	return fmt.Sprintf("%d", s.Size)
}
