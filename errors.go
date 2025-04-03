// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package dtorm provides a database ORM (Object-Relational Mapping) implementation
// with support for SQLite, MySQL, and SQL Server databases.
package dtorm

import "fmt"

// ErrNoResults represents an error that occurs when a database query returns no results.
// This error is typically returned when a query expects at least one result but finds none.
type ErrNoResults struct {
	// Err contains the underlying error that caused the no results condition
	Err error
}

// NoResults creates a new ErrNoResults error with the specified message.
// This function is used to wrap errors when a database query returns no results.
// Parameters:
//
//	msg: The error message to include in the error
//
// Returns:
//
//	A new ErrNoResults error
func NoResults(msg string) ErrNoResults {
	return ErrNoResults{Err: fmt.Errorf("%s", msg)}
}

// Error returns the error message for ErrNoResults.
// This method implements the error interface for ErrNoResults.
// Returns:
//
//	The error message as a string
func (e ErrNoResults) Error() string {
	return e.Err.Error()
}
