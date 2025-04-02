// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package dtorm provides interfaces and utilities for database model management.
package dtorm

import (
	"reflect"
	"strconv"
	"strings"
)

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

	// Disable marks the model as deleted (soft deletion).
	Disable()
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
		sv := v.Field(i)
		null := false

		// Handle pointer fields
		if sv.Kind() == reflect.Ptr {
			null = true
			sv = sv.Elem()
		}

		// Handle nested structs (except time.Time)
		if sv.Kind() == reflect.Struct && sv.Type().Name() != "Time" {
			if subf := getDefs(sv.Interface(), false); len(subf) > 0 {
				res = append(res, subf...)
			}
		} else {
			// Process mxorm tags for field configuration
			if tg, ok := st.Tag.Lookup("mxorm"); ok {
				nm := st.Name
				szMj := 0  // Major size (e.g., varchar length)
				szMn := 0  // Minor size (e.g., decimal places)
				id := false // Is identity field
				key := false // Is key field
				uns := false // Is unsigned
				fld := tString // Default field type

				// Find matching field type from reflection Kind
			FieldSearchLoop:
				for k, v := range fieldTrans {
					for _, v2 := range v {
						if v2 == sv.Kind() {
							fld = k
							for _, sn := range fieldUnsigned {
								if sn == sv.Kind() {
									uns = true
								}
							}
							break FieldSearchLoop
						}
					}
				}

				// Parse field tags
				if tg != "" {
					tgs := strings.Split(tg, ",")
					for _, t := range tgs {
						pts := strings.Split(t, ":")
						if len(pts) == 2 {
							switch pts[0] {
							case "type":
								typeKey := pts[1]
								if typeKey == "time" {
									typeKey = "struct"
								}
								if v, ok := fieldNames[typeKey]; ok {
									fld = v
								}
							case "size":
								szPt := strings.Split(pts[1], ",")
								if v, err := strconv.ParseInt(szPt[0], 10, 64); err == nil {
									szMj = int(v)
									if len(szPt) > 1 {
										if v, err = strconv.ParseInt(szPt[1], 10, 64); err == nil {
											szMn = int(v)
										}
									}
								}
							case "identity":
								id = pts[1] == "true"
							case "key":
								key = pts[1] == "true"
							case "unsigned":
								uns = pts[1] == "true"
							}
						}
					}
				}
				res = append(res, newField(nm, fld, szMj, szMn, id, key, uns, null))
			}
		}
	}
	return res
}

