// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.

// MIT License
//
// Copyright (c) 2025 DaggerTech
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package tests

import (
	"fmt"
	"testing"

	"github.com/markoxley/mud"
	"github.com/markoxley/mud/where"
)

// mockManager implements the Manager interface for testing
type mockManager struct {
	db *mud.DB
}

func (m *mockManager) SetDB(db *mud.DB) {
	m.db = db
}

func (m *mockManager) GetDB() *mud.DB {
	return m.db
}

func (m *mockManager) ConnectionString(cfg *mud.Config) (string, error) {
	return "mock_connection", nil
}

func (m *mockManager) Operators() []string {
	return []string{
		"%s = %s",                  // Equal
		"%s > %s",                  // Greater than
		"%s < %s",                  // Less than
		"%s LIKE %s",               // Pattern matching
		"%s IN (%s)",               // In list
		"%s BETWEEN %s AND %s",     // Between range
		"%s IS NULL",               // Is null check
		"%s <> %s",                 // Not equal
		"%s <= %s",                 // Less than or equal
		"%s >= %s",                 // Greater than or equal
		"%s NOT LIKE %s",           // Not like pattern
		"%s NOT IN (%s)",           // Not in list
		"%s NOT BETWEEN %s AND %s", // Not between range
		"%s IS NOT NULL",           // Is not null check
	}
}

func (m *mockManager) LimitString(c *mud.Criteria) string {
	if c == nil || c.Limit < 1 {
		return ""
	}
	return fmt.Sprintf(" LIMIT %d", c.Limit)
}

func (m *mockManager) OffsetString(c *mud.Criteria) string {
	if c == nil || c.Offset < 1 {
		return ""
	}
	return fmt.Sprintf(" OFFSET %d", c.Offset)
}

func (m *mockManager) IdentityString(f string) string {
	return f
}

func (m *mockManager) BuildQuery(where, order, limit, offset string) string {
	result := ""
	if where != "" {
		result += where
	}
	if order != "" {
		result += order
	}
	if limit != "" {
		result += limit
	}
	if offset != "" {
		result += offset
	}
	return result
}

func (m *mockManager) TableExistsQuery(name string) string {
	return "SELECT 1"
}

func (m *mockManager) TableCreate() string {
	return "CREATE TABLE"
}

func (m *mockManager) IndexCreate() string {
	return "CREATE INDEX"
}

func TestCriteriaWhereString(t *testing.T) {
	mgr := &mockManager{}

	tests := []struct {
		name      string
		criteria  mud.Criteria
		wantWhere string
	}{
		{
			name: "nil where",
			criteria: mud.Criteria{
				Where:      nil,
				IncDeleted: false,
			},
			wantWhere: "WHERE DeleteDate IS NULL",
		},
		{
			name: "string where",
			criteria: mud.Criteria{
				Where:      "name = 'test'",
				IncDeleted: false,
			},
			wantWhere: " WHERE name = 'test' AND DeleteDate IS NULL",
		},
		{
			name: "where builder pointer",
			criteria: mud.Criteria{
				Where:      where.Equal("name", "test"),
				IncDeleted: false,
			},
			wantWhere: " WHERE name = 'test' AND DeleteDate IS NULL",
		},
		{
			name: "include deleted",
			criteria: mud.Criteria{
				Where:      where.Equal("name", "test"),
				IncDeleted: true,
			},
			wantWhere: " WHERE name = 'test'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.criteria.WhereString(mgr)
			if got != tt.wantWhere {
				t.Errorf("Criteria.WhereString() = [%v], want [%v]", got, tt.wantWhere)
			}
		})
	}
}

type orderStringer struct {
	order string
}

func (o orderStringer) String() string {
	return o.order
}

func TestCriteriaOrderString(t *testing.T) {
	mgr := &mockManager{}

	tests := []struct {
		name      string
		criteria  mud.Criteria
		wantOrder string
	}{
		{
			name: "nil order",
			criteria: mud.Criteria{
				Order: nil,
			},
			wantOrder: "",
		},
		{
			name: "string order",
			criteria: mud.Criteria{
				Order: "name ASC",
			},
			wantOrder: " ORDER BY name ASC",
		},
		{
			name: "stringer order",
			criteria: mud.Criteria{
				Order: orderStringer{order: "name DESC"},
			},
			wantOrder: " ORDER BY name DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.criteria.OrderString(mgr)
			if got != tt.wantOrder {
				t.Errorf("Criteria.OrderString() = [%v], want [%v]", got, tt.wantOrder)
			}
		})
	}
}

func TestCriteriaLimitOffset(t *testing.T) {
	mgr := &mockManager{}

	tests := []struct {
		name       string
		criteria   mud.Criteria
		wantLimit  string
		wantOffset string
	}{
		{
			name: "no limit or offset",
			criteria: mud.Criteria{
				Limit:  0,
				Offset: 0,
			},
			wantLimit:  "",
			wantOffset: "",
		},
		{
			name: "with limit",
			criteria: mud.Criteria{
				Limit:  10,
				Offset: 0,
			},
			wantLimit:  " LIMIT 10",
			wantOffset: "",
		},
		{
			name: "with offset",
			criteria: mud.Criteria{
				Limit:  0,
				Offset: 5,
			},
			wantLimit:  "",
			wantOffset: " OFFSET 5",
		},
		{
			name: "with both",
			criteria: mud.Criteria{
				Limit:  10,
				Offset: 5,
			},
			wantLimit:  " LIMIT 10",
			wantOffset: " OFFSET 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLimit := tt.criteria.LimitString(mgr)
			if gotLimit != tt.wantLimit {
				t.Errorf("Criteria.LimitString() = %v, want %v", gotLimit, tt.wantLimit)
			}

			gotOffset := tt.criteria.OffsetString(mgr)
			if gotOffset != tt.wantOffset {
				t.Errorf("Criteria.OffsetString() = %v, want %v", gotOffset, tt.wantOffset)
			}
		})
	}
}

func TestCriteriaString(t *testing.T) {
	mgr := &mockManager{}

	tests := []struct {
		name     string
		criteria mud.Criteria
		want     string
	}{
		{
			name: "complete query",
			criteria: mud.Criteria{
				Where:      where.Equal("name", "test"),
				Order:      "name ASC",
				Limit:      10,
				Offset:     5,
				IncDeleted: false,
			},
			want: " WHERE name = 'test' AND DeleteDate IS NULL ORDER BY name ASC LIMIT 10 OFFSET 5",
		},
		{
			name: "where only",
			criteria: mud.Criteria{
				Where:      where.Equal("name", "test"),
				IncDeleted: false,
			},
			want: " WHERE name = 'test' AND DeleteDate IS NULL",
		},
		{
			name: "order only",
			criteria: mud.Criteria{
				Order:      "name ASC",
				IncDeleted: false,
			},
			want: "WHERE DeleteDate IS NULL ORDER BY name ASC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.criteria.String(mgr)
			if got != tt.want {
				t.Errorf("Criteria.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
