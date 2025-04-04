// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.

package tests

import (
	"testing"
	"time"

	"github.com/markoxley/mud"
	"github.com/markoxley/mud/utils"
)

// ModelTest is a concrete implementation of Modeller for testing
type ModelTest struct {
	mud.Model
	Name string
}

func TestModelIsDeleted(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name       string
		deleteDate *time.Time
		want       bool
	}{
		{
			name:       "not deleted",
			deleteDate: nil,
			want:       false,
		},
		{
			name:       "deleted",
			deleteDate: &now,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModelTest{}
			m.DeleteDate = tt.deleteDate
			if got := m.IsDeleted(); got != tt.want {
				t.Errorf("Model.IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelDisable(t *testing.T) {
	m := &ModelTest{}
	if m.DeleteDate != nil {
		t.Error("DeleteDate should be nil before calling Disable()")
	}

	m.Disable()
	if m.DeleteDate == nil {
		t.Error("DeleteDate should not be nil after calling Disable()")
	}
}

func TestGetTableName(t *testing.T) {
	tests := []struct {
		name string
		m    mud.Modeller
		want string
	}{
		{
			name: "pointer to model",
			m:    &ModelTest{},
			want: "ModelTest",
		},
		{
			name: "direct model",
			m:    ModelTest{},
			want: "ModelTest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mud.GetTableName(tt.m)
			if got != tt.want {
				t.Errorf("getTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModelIsNew(t *testing.T) {
	tests := []struct {
		name       string
		id         *string
		createDate time.Time
		lastUpdate time.Time
		deleteDate *time.Time
		want       bool
	}{
		{
			name:       "new model",
			id:         nil,
			createDate: time.Time{},
			lastUpdate: time.Time{},
			deleteDate: nil,
			want:       true,
		},
		{
			name:       "existing model",
			id:         utils.Ptr("123"),
			createDate: time.Now(),
			lastUpdate: time.Now(),
			deleteDate: nil,
			want:       false,
		},
		{
			name:       "deleted model",
			id:         utils.Ptr("123"),
			createDate: time.Now(),
			lastUpdate: time.Now(),
			deleteDate: utils.Ptr(time.Now()),
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ModelTest{}
			m.ID = tt.id
			m.CreateDate = tt.createDate
			m.LastUpdate = tt.lastUpdate
			m.DeleteDate = tt.deleteDate

			if got := m.IsNew(); got != tt.want {
				t.Errorf("Model.IsNew() = %v, want %v", got, tt.want)
			}
		})
	}
}
