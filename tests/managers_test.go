// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.

package tests

import (
	"strings"
	"testing"

	"github.com/markoxley/mud"
	"github.com/markoxley/mud/order"
)

func TestMySQLManagerQueries(t *testing.T) {
	mgr := &mud.MySQLManager{}
	mgr.SetDB(&mud.DB{})
	if mgr.GetDB() == nil {
		t.Error("expected GetDB to return the assigned DB")
	}

	if got := mgr.LimitString(nil); got != "" {
		t.Errorf("expected empty limit, got %q", got)
	}
	if got := mgr.LimitString(&mud.Criteria{Limit: 0}); got != "" {
		t.Errorf("expected empty limit for 0, got %q", got)
	}
	if got := mgr.LimitString(&mud.Criteria{Limit: 5}); got != " LIMIT 5" {
		t.Errorf("unexpected limit %q", got)
	}
	if got := mgr.OffsetString(nil); got != "" {
		t.Errorf("expected empty offset, got %q", got)
	}
	if got := mgr.OffsetString(&mud.Criteria{Offset: 0}); got != "" {
		t.Errorf("expected empty offset for 0, got %q", got)
	}
	if got := mgr.OffsetString(&mud.Criteria{Offset: 10}); got != " OFFSET 10" {
		t.Errorf("unexpected offset %q", got)
	}

	if got := mgr.BuildQuery("WHERE x", "ORDER BY y", " LIMIT 1", " OFFSET 2"); got != " WHERE x ORDER BY y LIMIT 1 OFFSET 2" {
		t.Errorf("unexpected BuildQuery %q", got)
	}
	if got := mgr.BuildQuery("", "", "", ""); got != "" {
		t.Errorf("expected empty BuildQuery, got %q", got)
	}

	q, args := mgr.TableExistsQuery("t")
	if !strings.Contains(q, "INFORMATION_SCHEMA") || len(args) != 2 {
		t.Errorf("unexpected TableExistsQuery %q args %v", q, args)
	}

	if _, err := mgr.ConnectionString(nil); err == nil {
		t.Error("expected error for nil config")
	}
	if _, err := mgr.ConnectionString(&mud.Config{User: "u"}); err == nil {
		t.Error("expected error for incomplete config")
	}
	cs, err := mgr.ConnectionString(&mud.Config{User: "u", Password: "p", Host: "h", Database: "d"})
	if err != nil || !strings.Contains(cs, "u:p@tcp(h)/d") {
		t.Errorf("unexpected connection string %q err %v", cs, err)
	}
}

func TestMSSQLManagerQueries(t *testing.T) {
	mgr := &mud.MSSQLManager{}
	mgr.SetDB(&mud.DB{})
	if mgr.GetDB() == nil {
		t.Error("expected GetDB to return the assigned DB")
	}

	if got := mgr.LimitString(nil); got != "" {
		t.Errorf("expected empty limit, got %q", got)
	}
	got := mgr.LimitString(&mud.Criteria{Limit: 5})
	if !strings.Contains(got, "ORDER BY [ID]") || !strings.Contains(got, "OFFSET 0 ROWS FETCH NEXT 5 ROWS ONLY") {
		t.Errorf("unexpected limit %q", got)
	}
	got = mgr.LimitString(&mud.Criteria{Limit: 5, Offset: 3})
	if !strings.Contains(got, "OFFSET 3 ROWS FETCH NEXT 5 ROWS ONLY") {
		t.Errorf("unexpected limit+offset %q", got)
	}
	// With an explicit order the default ORDER BY is omitted.
	got = mgr.LimitString(&mud.Criteria{Limit: 5, Order: order.Asc("Name")})
	if strings.Contains(got, "ORDER BY [ID]") {
		t.Errorf("expected no default order, got %q", got)
	}

	if got := mgr.OffsetString(nil); got != "" {
		t.Errorf("expected empty offset, got %q", got)
	}
	if got := mgr.OffsetString(&mud.Criteria{Offset: 0}); got != "" {
		t.Errorf("expected empty offset for 0, got %q", got)
	}
	got = mgr.OffsetString(&mud.Criteria{Offset: 7})
	if !strings.Contains(got, "ORDER BY [ID]") || !strings.Contains(got, "OFFSET 7 ROWS") {
		t.Errorf("unexpected offset %q", got)
	}
	got = mgr.OffsetString(&mud.Criteria{Offset: 7, Order: order.Asc("Name")})
	if strings.Contains(got, "ORDER BY [ID]") || !strings.Contains(got, "OFFSET 7 ROWS") {
		t.Errorf("unexpected offset with order %q", got)
	}

	if got := mgr.BuildQuery("WHERE x", "ORDER BY y", " FETCH", " OFFSET"); got != " WHERE x ORDER BY y OFFSET FETCH" {
		t.Errorf("unexpected BuildQuery %q", got)
	}

	pos := 1
	if p := mgr.Placeholder(&pos); p != "@p1" || pos != 2 {
		t.Errorf("unexpected placeholder %q pos %d", p, pos)
	}
	if p := mgr.Placeholder(nil); p != "@p1" {
		t.Errorf("expected @p1 for nil pos, got %q", p)
	}

	if _, err := mgr.ConnectionString(nil); err == nil {
		t.Error("expected error for nil config")
	}
	if _, err := mgr.ConnectionString(&mud.Config{User: "u"}); err == nil {
		t.Error("expected error for incomplete config")
	}
	cs, err := mgr.ConnectionString(&mud.Config{User: "u", Password: "p", Host: "h", Database: "d"})
	if err != nil || !strings.Contains(cs, "sqlserver://u:p@h?database=d") {
		t.Errorf("unexpected connection string %q err %v", cs, err)
	}
}

func TestSqliteManagerExtras(t *testing.T) {
	mgr := &mud.SqliteManager{}
	mgr.SetDB(&mud.DB{})
	if mgr.GetDB() == nil {
		t.Error("expected GetDB to return the assigned DB")
	}
	if got := mgr.OffsetString(&mud.Criteria{Offset: 0}); got != "" {
		t.Errorf("expected empty offset for 0, got %q", got)
	}
	if got := mgr.OffsetString(&mud.Criteria{Offset: 4}); got != " OFFSET 4" {
		t.Errorf("unexpected offset %q", got)
	}
	if got := mgr.LimitString(&mud.Criteria{Limit: 0}); got != "" {
		t.Errorf("expected empty limit for 0, got %q", got)
	}
	if got := mgr.BuildQuery("", "ORDER BY y", "", ""); got != " ORDER BY y" {
		t.Errorf("unexpected BuildQuery %q", got)
	}
	if _, err := mgr.ConnectionString(nil); err == nil {
		t.Error("expected error for nil config")
	}
	if _, err := mgr.ConnectionString(&mud.Config{}); err == nil {
		t.Error("expected error for empty database")
	}
}
