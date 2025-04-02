// Package dtorm provides a database ORM (Object-Relational Mapping) implementation
// with support for SQLite, MySQL, and SQL Server databases.
package dtorm

import "fmt"

// fieldSize represents the size and precision of a database field.
// It is used to store the maximum length and decimal places for numeric types.
type fieldSize struct {
	// size represents the maximum length of the field
	size int
	// decimal represents the number of decimal places for numeric types
	decimal int
}

// newSize creates a new fieldSize with the specified size and decimal places.
// This function is used to define the size constraints for database fields.
// Parameters:
//   sz: The maximum length of the field
//   dec: The number of decimal places (for numeric types)
// Returns:
//   A new fieldSize instance
func newSize(sz, dec int) fieldSize {
	return fieldSize{
		size:    sz,
		decimal: dec,
	}
}

// String returns a string representation of the field size.
// For numeric types with decimal places, it returns "size,decimal".
// For other types, it returns just the size.
// This method is used for generating SQL schema definitions.
// Returns:
//   A string representation of the field size
func (s fieldSize) String() string {
	if s.decimal > 0 {
		return fmt.Sprintf("%d,%d", s.size, s.decimal)
	}
	return fmt.Sprintf("%d", s.size)
}
