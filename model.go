// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package mud provides base model functionality for database entities.
package mud

import (
	"reflect"
	"time"

	"github.com/markoxley/mud/utils"
)

// Model represents the base structure for all database entities.
// It provides common fields and functionality for tracking creation,
// updates, and soft deletion of records.
type Model struct {
	// Unique identifier for the record
	ID *string
	// Timestamp when the record was created
	CreateDate time.Time
	// Timestamp of the last update
	LastUpdate time.Time
	// Timestamp when the record was soft deleted (nil if active)
	DeleteDate *time.Time
	// Optional custom table name override
	tableName *string
}

// CreateModel initializes a new Model instance with current timestamps.
// This should be called when creating new database entities.
func CreateModel() Model {
	return Model{
		CreateDate: time.Now(),
		LastUpdate: time.Now(),
	}
}

// StandingData returns a list of default records for the model.
// This can be overridden by implementing models to provide seed data.
func (m Model) StandingData() []Modeller {
	return nil
}

// GetID returns the unique identifier of the model.
func (m Model) GetID() *string {
	return m.ID
}

// IsNew checks if the model is a new record (has not been saved to database).
// Returns true if the ID is nil, indicating the record hasn't been assigned an ID.
func (m Model) IsNew() bool {
	return m.ID == nil
}

// IsDeleted checks if the model has been soft deleted.
// Returns true if DeleteDate is not nil, indicating the record is deleted.
func (m Model) IsDeleted() bool {
	return m.DeleteDate != nil
}

// Disable marks the model as soft deleted by setting its DeleteDate to current time.
func (m *Model) Disable() {
	m.DeleteDate = utils.Ptr(time.Now())
}

// TableName returns the custom table name set with SetTableName, or an
// empty string when none has been set. Embedding models may override this
// method to provide a static table name.
func (m Model) TableName() string {
	if m.tableName == nil {
		return ""
	}
	return *m.tableName
}

// SetTableName sets a custom table name for the model instance.
func (m *Model) SetTableName(name string) {
	m.tableName = utils.Ptr(name)
}

// GetTableName determines the database table name for a model.
// A model can override the derived struct type name by implementing the
// Tabler interface (a TableName() string method) or calling SetTableName.
func GetTableName(m Modeller) string {
	typ := reflect.TypeOf(m)
	if typ.Kind() == reflect.Pointer {
		// Honour an override on the instance itself.
		if t, ok := m.(Tabler); ok {
			if name := t.TableName(); name != "" {
				return name
			}
		}
		typ = typ.Elem()
	} else {
		// Honour value-receiver overrides carried by the value, then
		// pointer-receiver overrides via a fresh *T instance.
		if t, ok := m.(Tabler); ok {
			if name := t.TableName(); name != "" {
				return name
			}
		}
		if t, ok := reflect.New(typ).Interface().(Tabler); ok {
			if name := t.TableName(); name != "" {
				return name
			}
		}
	}
	return typ.Name()
}
