// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package dtorm provides a database ORM (Object-Relational Mapping) implementation
// with support for SQLite, MySQL, and SQL Server databases.
package dtorm

import (
	"fmt"
	"reflect"
)

// Version constants define the current version of the dtorm package.
const (
	// majorVersion represents the major version number
	majorVersion = 0
	// minorVersion represents the minor version number
	minorVersion = 1
	// releaseVersion represents the patch version number
	releaseVersion = 0
)

// Version returns the current version of the dtorm package as a string.
// The version follows semantic versioning format (MAJOR.MINOR.PATCH).
// Returns:
//   A string representing the current version
func Version() string {
	return fmt.Sprintf("%d.%d.%d", majorVersion, minorVersion, releaseVersion)
}

// fieldTrans maps database field types to their corresponding Go reflect kinds.
// This is used for type conversion between database values and Go types.
var (
	fieldTrans = map[fieldType][]reflect.Kind{
		// Boolean fields map to Go bool type
		tBool:     {reflect.Bool},
		// DateTime fields map to Go struct (time.Time)
		tDateTime: {reflect.Struct},
		// Double fields map to Go float64 type
		tDouble:   {reflect.Float64},
		// Float fields map to Go float32 type
		tFloat:    {reflect.Float32},
		// Integer fields map to various Go integer types
		tInt:      {reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint16, reflect.Uint32},
		// Long fields map to Go int64 and uint64 types
		tLong:     {reflect.Int64, reflect.Uint64},
		// String fields map to Go string type
		tString:   {reflect.String},
	}
)

// fieldUnsigned contains the Go reflect kinds that represent unsigned integer types.
// This is used to determine if a field should be treated as unsigned in the database.
var fieldUnsigned = []reflect.Kind{
	reflect.Uint,
	reflect.Uint8,
	reflect.Uint16,
	reflect.Uint32,
	reflect.Uint64,
}

