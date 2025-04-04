// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.

// Package mud provides a database ORM (Object-Relational Mapping) implementation
// with support for SQLite, MySQL, and SQL Server databases.
package mud

// Config represents the configuration for a database connection.
// It contains all the necessary information to establish and configure a database connection.
type Config struct {
	// Type specifies the type of database (sqlite, mysql, sqlserver)
	Type string `json:"type"`
	// Host specifies the hostname or IP address of the database server
	// This is optional for SQLite databases
	Host string `json:"host,omitzero"`
	// Database specifies the name of the database to connect to
	Database string `json:"database"`
	// User specifies the username for database authentication
	// This is optional for SQLite databases
	User string `json:"user,omitzero"`
	// Password specifies the password for database authentication
	// This is optional for SQLite databases
	Password string `json:"password,omitzero"`
	// Deletable indicates whether records can be permanently deleted
	// When false, records are soft-deleted (marked with delete date)
	Deletable bool `json:"deletable,omitzero"`
	// DisabledTransactions indicates whether database transactions should be disabled
	// When true, each operation will be executed independently
	DisabledTransactions bool `json:"disabledTransactions,omitzero"`
}
