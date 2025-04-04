// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package dtorm provides a simple ORM implementation for MySQL databases.
package dtorm

import (
	"fmt"
)

// MySQLManager implements the database management interface for MySQL.
// It handles MySQL specific query generation and database operations.
type MySQLManager struct {
	// dbName is the name of the connected database.
	dbName string
	// db is a reference to the database connection.
	db *DB
}

// SetDB assigns a database connection to the manager.
func (m *MySQLManager) SetDB(db *DB) {
	m.db = db
}

// GetDB returns the current database connection.
func (m *MySQLManager) GetDB() *DB {
	return m.db
}

// ConnectionString generates a MySQL connection string from the provided configuration.
// Returns an error if the configuration is nil or missing required fields.
func (m *MySQLManager) ConnectionString(cfg *Config) (string, error) {
	if cfg == nil {
		return "", fmt.Errorf("no config provided")
	}
	if cfg.User == "" || cfg.Password == "" || cfg.Host == "" || cfg.Database == "" {
		return "", fmt.Errorf("invalid config provided")
	}
	m.dbName = cfg.Database
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Database), nil
}

// LimitString generates the MySQL specific LIMIT clause for result limiting.
// Returns an empty string if criteria is nil or limit is less than 1.
func (m *MySQLManager) LimitString(c *Criteria) string {
	if c == nil || c.Limit < 1 {
		return ""
	}
	return fmt.Sprintf(" LIMIT %d", c.Limit)
}

// OffsetString generates the MySQL specific OFFSET clause.
// Returns an empty string if criteria is nil or offset is less than 1.
func (m *MySQLManager) OffsetString(c *Criteria) string {
	if c == nil || c.Offset < 1 {
		return ""
	}
	return fmt.Sprintf(" OFFSET %d", c.Offset)
}

// IdentityString wraps a field name in backticks for MySQL identifier escaping.
func (m *MySQLManager) IdentityString(f string) string {
	return fmt.Sprintf("`%s`", f)
}

// TableCreate returns the MySQL table creation query template.
// The template includes IF NOT EXISTS to prevent duplicate table creation errors.
func (m *MySQLManager) TableCreate() string {
	return "CREATE TABLE IF NOT EXISTS `%s` (%s);"
}

// IndexCreate returns the MySQL index creation query template.
func (m *MySQLManager) IndexCreate() string {
	return "CREATE INDEX `%s_%s_Idx` ON %s(`%s`);"
}

// BuildQuery combines WHERE, ORDER BY, LIMIT, and OFFSET clauses into a complete query string.
func (m *MySQLManager) BuildQuery(where string, order string, limit string, offset string) string {
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
func (m *MySQLManager) TableExistsQuery(name string) string {
	return fmt.Sprintf("SHOW TABLES WHERE Tables_in_%s = '%s'", m.dbName, name)
}

// Operators returns a list of MySQL compatible operator formats for query building.
// These formats include comparison, LIKE, IN, BETWEEN, and NULL check operators.
func (m *MySQLManager) Operators() []string {
	return []string{
		"`%s` = %s",                  // Equal
		"`%s` > %s",                  // Greater than
		"`%s` < %s",                  // Less than
		"`%s` LIKE %s",               // Pattern matching
		"`%s` IN (%s)",               // In list
		"`%s` BETWEEN %s AND %s",     // Between range
		"`%s` IS NULL",               // Is null check
		"`%s` <> %s",                 // Not equal
		"`%s` <= %s",                 // Less than or equal
		"`%s` >= %s",                 // Greater than or equal
		"`%s` NOT LIKE %s",           // Not like pattern
		"`%s` NOT IN (%s)",           // Not in list
		"`%s` NOT BETWEEN %s AND %s", // Not between range
		"`%s` IS NOT NULL",           // Is not null check
	}
}
