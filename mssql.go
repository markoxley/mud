// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package mud provides a simple ORM implementation for MS SQL Server databases.
package mud

import "fmt"

// MSSQLManager implements the database management interface for Microsoft SQL Server.
// It handles SQL Server specific query generation and database operations.
type MSSQLManager struct {
	// dbName is the name of the connected database
	dbName string
	// db is a reference to the database connection
	db *DB
}

// SetDB assigns a database connection to the manager.
func (m *MSSQLManager) SetDB(db *DB) {
	m.db = db
}

// GetDB returns the current database connection.
func (m *MSSQLManager) GetDB() *DB {
	return m.db
}

// ConnectionString generates a SQL Server connection string from the provided configuration.
// Returns an error if the configuration is nil or missing required fields.
func (m *MSSQLManager) ConnectionString(cfg *Config) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("no config provided")
	}
	if cfg.User == "" || cfg.Password == "" || cfg.Host == "" || cfg.Database == "" {
		return "", fmt.Errorf("invalid config provided")
	}
	m.dbName = cfg.Database
	return fmt.Sprintf("sqlserver://%s:%s@%s?database=%s", cfg.User, cfg.Password, cfg.Host, cfg.Database), nil
}

// LimitString generates the SQL Server specific FETCH NEXT clause for result limiting.
// Returns an empty string if criteria is nil or limit is less than 1.
func (m *MSSQLManager) LimitString(c *Criteria) string {
	if c == nil || c.Limit < 1 {
		return ""
	}
	res := ""
	if c.Order == nil {
		res = fmt.Sprintf(" ORDER BY [ID] ")
	}
	if c.Offset < 1 {
		return fmt.Sprintf("%s OFFSET 0 ROWS FETCH NEXT %d ROWS ONLY", res, c.Limit)
	}
	return fmt.Sprintf("%s OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", res, c.Offset, c.Limit)
}

// OffsetString generates the SQL Server specific OFFSET clause.
// Returns an empty string if criteria is nil or offset is less than 1.
func (m *MSSQLManager) OffsetString(c *Criteria) string {
	if c == nil || c.Offset < 1 {
		return ""
	}
	res := ""
	if c.Order == nil {
		res = fmt.Sprintf(" ORDER BY [ID] ")
	}
	return fmt.Sprintf("%s OFFSET %d ROWS", res, c.Offset)
}

// IdentityString wraps a field name in square brackets for SQL Server identifier escaping.
func (m *MSSQLManager) IdentityString(f string) string {
	return fmt.Sprintf("[%s]", f)
}

// BuildQuery combines WHERE, ORDER BY, LIMIT, and OFFSET clauses into a complete query string.
func (m *MSSQLManager) BuildQuery(where string, order string, limit string, offset string) string {
	res := ""
	if where != "" {
		res += fmt.Sprintf(" %s", where)
	}
	if order != "" {
		res += fmt.Sprintf(" %s", order)
	}
	return res + offset + limit
}

// TableExistsQuery generates a query to check if a table exists in the database.
func (m *MSSQLManager) TableExistsQuery(name string) string {
	return fmt.Sprintf("SELECT [Name] FROM [sys].[tables] WHERE [Name] = '%s'", name)
}

// Operators returns a list of SQL Server compatible operator formats for query building.
// These formats include comparison, LIKE, IN, BETWEEN, and NULL check operators.
func (m *MSSQLManager) Operators() []string {
	return []string{
		"[%s] = %s",                  // Equal
		"[%s] > %s",                  // Greater than
		"[%s] < %s",                  // Less than
		"[%s] LIKE %s",               // Pattern matching
		"[%s] IN (%s)",               // In list
		"[%s] BETWEEN %s AND %s",     // Between range
		"[%s] IS NULL",               // Is null check
		"[%s] <> %s",                 // Not equal
		"[%s] <= %s",                 // Less than or equal
		"[%s] >= %s",                 // Greater than or equal
		"[%s] NOT LIKE %s",           // Not like pattern
		"[%s] NOT IN (%s)",           // Not in list
		"[%s] NOT BETWEEN %s AND %s", // Not between range
		"[%s] IS NOT NULL",           // Is not null check
	}
}

// TableCreate returns the SQL Server table creation query template.
// The template includes a check to prevent creating duplicate tables.
func (m *MSSQLManager) TableCreate() string {
	return "IF OBJECT_ID(N'dbo.%[1]s', N'U') IS NULL BEGIN CREATE TABLE dbo.[%[1]s] (%[2]s); END;"
}

// IndexCreate returns the SQL Server index creation query template.
func (m *MSSQLManager) IndexCreate() string {
	return "CREATE INDEX [%s_%s_Idx] ON [%s]([%s]);"
}
