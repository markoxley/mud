// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package dtorm provides database management interfaces and implementations.
package dtorm

import (
	"fmt"
)

// Manager defines the interface for database-specific operations.
// Each supported database type (SQLite, MySQL, SQL Server) implements this interface
// to provide its specific SQL syntax and behavior.
type Manager interface {
	// SetDB assigns a database connection to the manager
	SetDB(db *DB)

	// GetDB returns the current database connection
	GetDB() *DB

	// ConnectionString generates a database-specific connection string from the configuration
	ConnectionString(cfg *Config) (string, error)

	// LimitString generates the database-specific LIMIT clause
	LimitString(c *Criteria) string

	// OffsetString generates the database-specific OFFSET clause
	OffsetString(c *Criteria) string

	// IdentityString wraps field names with database-specific identifier quotes
	IdentityString(f string) string

	// BuildQuery combines WHERE, ORDER BY, LIMIT, and OFFSET clauses into a query
	BuildQuery(where string, order string, limit string, offset string) string

	// TableExistsQuery generates a query to check if a table exists
	TableExistsQuery(name string) string

	// Operators returns a list of database-specific operator formats
	Operators() []string

	// TableCreate returns the database-specific table creation template
	TableCreate() string

	// IndexCreate returns the database-specific index creation template
	IndexCreate() string
}

// GetManager creates and returns a database-specific Manager implementation based on the configuration.
// Supported database types are: "sqlite", "mysql", and "sqlserver".
// Returns an error if an unsupported database type is specified or if config is nil.
func GetManager(config *Config) (Manager, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	switch config.Type {
	case "sqlite", "sqlite3":
		return &SqliteManager{}, nil
	case "mysql", "mariadb":
		return &MySQLManager{}, nil
	case "sqlserver", "mssql":
		return &MSSQLManager{}, nil
	default:
		return nil, fmt.Errorf("invalid database type: %s", config.Type)
	}
}

// GetTypeName returns the type name of the manager implementation
func GetTypeName(m Manager) string {
	if m == nil {
		return ""
	}
	return fmt.Sprintf("%T", m)
}
