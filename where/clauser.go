// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package where provides functionality for building SQL WHERE clauses
package where

// clauser is an interface that defines the contract for SQL clause builders
// Implementations of this interface are responsible for generating SQL clauses
// and managing logical conjunctions between conditions
type clauser interface {
	// build returns the SQL clause string representation with placeholder
	// markers instead of literal values, plus the ordered list of arguments
	// to bind to those placeholders.
	//
	// @param operators List of operator format strings to use in the clause
	// @param ph Function that returns the next placeholder marker
	// @param pos Counter used by numbered placeholders (e.g. @pN)
	// @return SQL clause string and its bind arguments
	build([]string, func(*int) string, *int) (string, []interface{})

	// getConjunction returns the logical conjunction used to combine this clause
	// with other clauses in a WHERE condition
	//
	// @return The conjunction used for combining clauses
	getConjunction() conjunction
}
