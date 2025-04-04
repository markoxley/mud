// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package mud provides a simple ORM implementation for SQLite databases.
package mud

import "fmt"

// SqliteManager implements the database management interface for SQLite.
// It handles SQLite specific query generation and database operations.
type SqliteManager struct {
	db *DB // Reference to the database connection
}

// SetDB assigns a database connection to the manager.
func (m *SqliteManager) SetDB(db *DB) {
	m.db = db
}

// GetDB returns the current database connection.
func (m *SqliteManager) GetDB() *DB {
	return m.db
}

// ConnectionString generates a SQLite connection string from the provided configuration.
// Returns an error if the configuration is nil or missing required fields.
// Note that SQLite only requires the database path, unlike other SQL databases.
func (m *SqliteManager) ConnectionString(cfg *Config) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("no config provided")
	}
	if cfg.Database == "" {
		return "", fmt.Errorf("invalid config provided")
	}
	return cfg.Database, nil
}

// LimitString generates the SQLite specific LIMIT clause for result limiting.
// Returns an empty string if criteria is nil or limit is less than 1.
func (m *SqliteManager) LimitString(c *Criteria) string {
	if c == nil || c.Limit < 1 {
		return ""
	}
	return fmt.Sprintf(" LIMIT %d", c.Limit)
}

// OffsetString generates the SQLite specific OFFSET clause.
// Returns an empty string if criteria is nil or offset is less than 1.
func (m *SqliteManager) OffsetString(c *Criteria) string {
	if c == nil || c.Offset < 1 {
		return ""
	}
	return fmt.Sprintf(" OFFSET %d", c.Offset)
}

// IdentityString wraps a field name in double quotes for SQLite identifier escaping.
func (m *SqliteManager) IdentityString(f string) string {
	return fmt.Sprintf("\"%s\"", f)
}

// TableCreate returns the SQLite table creation query template.
// The template includes IF NOT EXISTS to prevent duplicate table creation errors.
func (m *SqliteManager) TableCreate() string {
	return "CREATE TABLE IF NOT EXISTS \"%s\" (%s);"
}

// IndexCreate returns the SQLite index creation query template.
func (m *SqliteManager) IndexCreate() string {
	return "CREATE INDEX \"%s_%s_Idx\" ON %s(\"%s\");"
}

// BuildQuery combines WHERE, ORDER BY, LIMIT, and OFFSET clauses into a complete query string.
func (m *SqliteManager) BuildQuery(where string, order string, limit string, offset string) string {
	res := ""
	if where != "" {
		res += fmt.Sprintf(" %s", where)
	}
	if order != "" {
		res += fmt.Sprintf(" %s", order)
	}
	return res + limit + offset
}

// TableExistsQuery generates a query to check if a table exists in the database.
// Uses the sqlite_master system table to check for table existence.
func (m *SqliteManager) TableExistsQuery(name string) string {
	return fmt.Sprintf("SELECT \"name\" FROM sqlite_master WHERE type='table' AND name='%s'", name)
}

// Operators returns a list of SQLite compatible operator formats for query building.
// These formats include comparison, LIKE, IN, BETWEEN, and NULL check operators.
func (m *SqliteManager) Operators() []string {
	return []string{
		"\"%s\" = %s",                  // Equal
		"\"%s\" > %s",                  // Greater than
		"\"%s\" < %s",                  // Less than
		"\"%s\" LIKE %s",               // Pattern matching
		"\"%s\" IN (%s)",               // In list
		"\"%s\" BETWEEN %s AND %s",     // Between range
		"\"%s\" IS NULL",               // Is null check
		"\"%s\" <> %s",                 // Not equal
		"\"%s\" <= %s",                 // Less than or equal
		"\"%s\" >= %s",                 // Greater than or equal
		"\"%s\" NOT LIKE %s",           // Not like pattern
		"\"%s\" NOT IN (%s)",           // Not in list
		"\"%s\" NOT BETWEEN %s AND %s", // Not between range
		"\"%s\" IS NOT NULL",           // Is not null check
	}
}
