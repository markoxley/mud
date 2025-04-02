// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package where provides functionality for building SQL WHERE clauses
package where

// clauser is an interface that defines the contract for SQL clause builders
// Implementations of this interface are responsible for generating SQL clauses
// and managing logical conjunctions between conditions
type clauser interface {
	// String returns the SQL clause string representation
	// The provided field names are used to construct the clause
	//
	// @param fields List of field names to use in the clause
	// @return SQL clause string
	String([]string) string

	// getConjunction returns the logical conjunction used to combine this clause
	// with other clauses in a WHERE condition
	//
	// @return The conjunction used for combining clauses
	getConjunction() conjunction
}

