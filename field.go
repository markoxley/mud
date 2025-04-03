// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package dtorm provides a database ORM (Object-Relational Mapping) implementation
// with support for SQLite, MySQL, and SQL Server databases.
package dtorm

// fieldType represents the type of a database field.
// It maps to various database-specific types and is used for type conversion.
type fieldType int

// fieldName represents the name of a database field.
type fieldName string

// Field type constants define the available field types in the ORM.
// These correspond to common database data types and are used for type conversion.
const (
	// tInt represents an integer field type
	tInt fieldType = iota
	// tLong represents a long integer field type
	tLong
	// tBool represents a boolean field type
	tBool
	// tDecimal represents a decimal/numeric field type
	tDecimal
	// tFloat represents a floating-point field type
	tFloat
	// tDouble represents a double-precision floating-point field type
	tDouble
	// tDateTime represents a date/time field type
	tDateTime
	// tChar represents a single character field type
	tChar
	// tString represents a string/varchar field type
	tString
	// tUUID represents a UUID field type
	tUUID
)

// Field type string constants provide string representations of field types.
// These are used for configuration and type identification.
const (
	sInt      = "int"
	sLong     = "long"
	sBool     = "bool"
	sDecimal  = "decimal"
	sFloat    = "float"
	sDouble   = "double"
	sDateTime = "struct"
	sChar     = "char"
	sString   = "string"
	sUUID     = "uuid"
)

// fieldNames maps string representations of field types to their corresponding fieldType constants.
// This is used for type conversion and validation.
var (
	fieldNames = map[string]fieldType{
		sInt:      tInt,
		sLong:     tLong,
		sBool:     tBool,
		sDecimal:  tDecimal,
		sFloat:    tFloat,
		sDouble:   tDouble,
		sDateTime: tDateTime,
		sChar:     tChar,
		sString:   tString,
		sUUID:     tUUID,
	}
)

// pgFieldNames maps fieldType constants to their PostgreSQL-specific database types.
// This is used for generating correct database schema definitions for PostgreSQL.
var (
	pgFieldNames = map[fieldType]string{
		tInt:      "INT",
		tLong:     "BIGINT",
		tBool:     "SMALLINT",
		tDecimal:  "DECIMAL",
		tFloat:    "REAL",
		tDouble:   "DOUBLE",
		tDateTime: "DATETIME",
		tChar:     "VARCHAR(1)",
		tString:   "VARCHAR",
		tUUID:     "VARCHAR(36)",
	}
)

// field represents a database field definition.
// It contains all the necessary information to define and work with a database column.
type field struct {
	// name is the name of the field in the database
	name string
	// fType is the type of the field (int, string, etc.)
	fType fieldType
	// size contains the size and decimal places for numeric types
	size FieldSize
	// identity indicates if this field is an auto-incrementing primary key
	identity bool
	// key indicates if this field is part of a primary key
	key bool
	// unsigned indicates if this numeric field should be unsigned
	unsigned bool
	// allowNull indicates if NULL values are allowed for this field
	allowNull bool
}

// newField creates a new field definition with the specified properties.
// It is used internally by the ORM to define database schema.
// Parameters:
//
//	nm: Field name
//	tp: Field type
//	sz: Field size (for numeric types)
//	dec: Decimal places (for numeric types)
//	id: Whether this is an identity field
//	ky: Whether this is a key field
//	us: Whether this is an unsigned numeric field
//	nl: Whether NULL values are allowed
//
// Returns:
//
//	A new field definition
func newField(nm string, tp fieldType, sz, dec int, id, ky, us bool, nl bool) field {
	return field{
		name:      nm,
		fType:     tp,
		size:      NewSize(sz, dec),
		identity:  id,
		key:       ky,
		unsigned:  us,
		allowNull: nl,
	}
}
