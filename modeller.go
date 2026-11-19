// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package mud provides interfaces and utilities for database model management.
package mud

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

// timeType is the reflect type of time.Time, used to distinguish datetime
// fields from nested model structs.
var timeType = reflect.TypeOf(time.Time{})

// Modeller defines the interface for database model objects.
// Any struct that implements this interface can be used as a database model.
type Modeller interface {
	// StandingData returns the standing data for the model.
	// This is used to provide seed or default data for the model.
	StandingData() []Modeller

	// GetID returns the ID of the model.
	// Returns nil if the model hasn't been saved to the database.
	GetID() *string

	// IsNew returns true if the model has yet to be saved to the database.
	IsNew() bool

	// IsDeleted returns true if the model has been marked as deleted.
	// This is used for soft deletion support.
	IsDeleted() bool

	// Note: Disable is intentionally not part of this interface. It is
	// defined on *Model, and requiring it here would exclude models used
	// by value (model values are created internally during Fetch).
}

// getDefs extracts field definitions from a struct using reflection.
// It processes struct tags and builds field metadata for database operations.
//
// Parameters:
//   - t: The struct to analyze
//   - first: If true, includes standard model fields (ID, CreateDate, etc.)
//
// Returns a slice of field definitions containing metadata about each field.
func getDefs(t interface{}, first bool) []field {
	res := make([]field, 0, 10)

	// Add standard model fields if this is the top-level struct
	if first {
		res = append(res, field{
			name:     "ID",
			fType:    tUUID,
			identity: true,
			key:      true,
		})
		res = append(res, field{
			name:  "CreateDate",
			fType: tDateTime,
			key:   true,
		})
		res = append(res, field{
			name:  "LastUpdate",
			fType: tDateTime,
			key:   true,
		})
		res = append(res, field{
			name:      "DeleteDate",
			fType:     tDateTime,
			allowNull: true,
		})
	}

	// Use reflection to analyze the struct fields
	v := reflect.ValueOf(t)
	ft := reflect.TypeOf(t)
	nf := v.NumField()
	for i := 0; i < nf; i++ {
		st := ft.Field(i)
		// Recurse into anonymous embedded structs before the exported check:
		// an embedded field of an unexported struct type is itself unexported,
		// but its exported fields are promoted and remain settable through
		// FieldByName on the parent model.
		if st.Anonymous {
			et := st.Type
			if et.Kind() == reflect.Ptr {
				et = et.Elem()
			}
			if et.Kind() == reflect.Struct && et != timeType {
				if subf := getDefs(reflect.New(et).Elem().Interface(), false); len(subf) > 0 {
					res = append(res, subf...)
				}
				continue
			}
		}
		// Skip unexported fields: they cannot be populated by reflection
		// and sv.Interface() would panic on them.
		if !st.IsExported() {
			continue
		}
		sv := v.Field(i)
		null := false

		// Handle pointer fields. When the pointer is nil, use a zero value
		// of the declared element type so the kind and any nested-struct
		// recursion still work.
		if sv.Kind() == reflect.Ptr {
			null = true
			if sv.IsNil() {
				sv = reflect.New(sv.Type().Elem()).Elem()
			} else {
				sv = sv.Elem()
			}
		}
		kind := sv.Kind()

		// Handle nested structs (except time.Time)
		if kind == reflect.Struct && sv.Type() != timeType {
			if subf := getDefs(sv.Interface(), false); len(subf) > 0 {
				res = append(res, subf...)
			}
		} else {
			// Process mud tags for field configuration
			if tg, ok := st.Tag.Lookup("mud"); ok {
				nm := st.Name
				szMj := 0      // Major size (e.g., varchar length)
				szMn := 0      // Minor size (e.g., decimal places)
				id := false    // Is identity field
				key := false   // Is key field
				uns := false   // Is unsigned
				fld := tString // Default field type

				// Find matching field type from reflection Kind
			FieldSearchLoop:
				for k, v := range fieldTrans {
					for _, v2 := range v {
						if v2 == kind {
							fld = k
							for _, sn := range fieldUnsigned {
								if sn == kind {
									uns = true
								}
							}
							break FieldSearchLoop
						}
					}
				}

				// Parse field tags. Note the size decimal component may
				// appear as a separate token because the tag is comma
				// separated (e.g. `size:10,2`), so look ahead for it.
				if tg != "" {
					tgs := strings.Split(tg, ",")
					for j := 0; j < len(tgs); j++ {
						pts := strings.SplitN(strings.TrimSpace(tgs[j]), ":", 2)
						if len(pts) == 0 {
							continue
						}

						switch pts[0] {
						case "type":
							if len(pts) < 2 {
								continue
							}
							typeKey := pts[1]
							if typeKey == "time" {
								typeKey = sStruct
							}
							if v, ok := fieldNames[typeKey]; ok {
								fld = v
							}
						case "size":
							if len(pts) < 2 {
								continue
							}
							if v, err := strconv.ParseInt(pts[1], 10, 64); err == nil {
								szMj = int(v)
								// Decimal places may have been split into
								// the next comma-separated token.
								if j+1 < len(tgs) && isDigits(tgs[j+1]) {
									j++
									if v, err = strconv.ParseInt(tgs[j], 10, 64); err == nil {
										szMn = int(v)
									}
								}
							}
						case "identity":
							id = true
						case "key":
							key = true
						case "unsigned":
							uns = true
						}

					}
				}
				res = append(res, newField(nm, fld, szMj, szMn, id, key, uns, null))
			}
		}
	}
	return res
}

// isDigits reports whether s consists only of ASCII digits.
func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
