// Package where provides functionality for building SQL WHERE clauses
package where

// conjunction represents a logical operator used to combine conditions in SQL WHERE clauses
type conjunction string

// Conjunction constants define the supported logical operators for combining conditions
const (
	conAnd conjunction = " AND " // Logical AND operator for combining conditions
	conOr  conjunction = " OR "  // Logical OR operator for combining conditions
)
