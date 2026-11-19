// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.

package tests

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/markoxley/mud"
	"github.com/markoxley/mud/utils"
	"github.com/markoxley/mud/where"
)

// whereCase describes a single where clause test. The `out` field uses "%p"
// as a sentinel for each successive placeholder; it is replaced with the
// manager-specific placeholder in sequence order before comparison.
// `args` may be nil, in which case the expected arguments are derived from
// input[1:] via utils.MakeArg.
type whereCase struct {
	mgr   mud.Manager
	name  string
	input []interface{}
	out   string
	args  []interface{}
}

// phSeq replaces each "%p" marker in sql with the next placeholder produced
// by the manager, in left-to-right order, matching Builder.Build's counter.
func phSeq(sql string, mgr mud.Manager) string {
	pos := 1
	for strings.Contains(sql, "%p") {
		sql = strings.Replace(sql, "%p", mgr.Placeholder(&pos), 1)
	}
	return sql
}

// makeArgs converts values to their expected bind-argument form.
func makeArgs(vals ...interface{}) []interface{} {
	res := make([]interface{}, 0, len(vals))
	for _, v := range vals {
		a, ok := utils.MakeArg(v)
		if ok {
			res = append(res, a)
		}
	}
	return res
}

// runWhereCase builds the clause and asserts the parameterized SQL and args.
func runWhereCase(t *testing.T, tst whereCase, b *where.Builder) {
	t.Helper()
	want := phSeq(tst.out, tst.mgr)
	wantArgs := tst.args
	if wantArgs == nil {
		wantArgs = makeArgs(tst.input[1:]...)
	}
	got, args := b.Build(tst.mgr.Operators(), tst.mgr.Placeholder)
	if got != want {
		t.Errorf("%v : Expected %v, got %v", tst.name, want, got)
	}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Errorf("%v : Expected args %v, got %v", tst.name, wantArgs, args)
	}
}

func TestWhereEqual(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Equal Bool", input: []interface{}{"valid", true}, out: "`valid` = %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Equal Int", input: []interface{}{"miles", int(829)}, out: "`miles` = %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` = %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Equal Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` = %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Equal Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` = %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Equal String", input: []interface{}{"name", "Sally"}, out: "`name` = %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Equal DateTime", input: []interface{}{"dob", tm1}, out: "`dob` = %p"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Equal Bool", input: []interface{}{"valid", true}, out: "[valid] = %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Equal Int", input: []interface{}{"miles", int(829)}, out: "[miles] = %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] = %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Equal Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] = %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Equal Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] = %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Equal String", input: []interface{}{"name", "Sally"}, out: "[name] = %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Equal DateTime", input: []interface{}{"dob", tm1}, out: "[dob] = %p"},

		{mgr: &mud.SqliteManager{}, name: "SQLite Equal Bool", input: []interface{}{"valid", true}, out: "\"valid\" = %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Equal Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" = %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" = %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Equal Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" = %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Equal Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" = %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Equal String", input: []interface{}{"name", "Sally"}, out: "\"name\" = %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Equal DateTime", input: []interface{}{"dob", tm1}, out: "\"dob\" = %p"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.Equal(tst.input[0].(string), tst.input[1]))
	}
}

func TestWhereGreater(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Greater Int", input: []interface{}{"miles", int(829)}, out: "`miles` > %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` > %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Greater Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` > %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Greater Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` > %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Greater DateTime", input: []interface{}{"dob", tm1}, out: "`dob` > %p"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Greater Int", input: []interface{}{"miles", int(829)}, out: "[miles] > %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] > %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Greater Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] > %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Greater Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] > %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Greater DateTime", input: []interface{}{"dob", tm1}, out: "[dob] > %p"},

		{mgr: &mud.SqliteManager{}, name: "SQLite Greater Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" > %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" > %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Greater Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" > %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Greater Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" > %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Greater DateTime", input: []interface{}{"dob", tm1}, out: "\"dob\" > %p"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.Greater(tst.input[0].(string), tst.input[1]))
	}
}

func TestWhereLess(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Less Int", input: []interface{}{"miles", int(829)}, out: "`miles` < %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Less Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` < %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Less Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` < %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Less Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` < %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Less DateTime", input: []interface{}{"dob", tm1}, out: "`dob` < %p"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Less Int", input: []interface{}{"miles", int(829)}, out: "[miles] < %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Less Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] < %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Less Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] < %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Less Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] < %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Less DateTime", input: []interface{}{"dob", tm1}, out: "[dob] < %p"},

		{mgr: &mud.SqliteManager{}, name: "SQLite Less Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" < %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Less Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" < %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Less Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" < %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Less Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" < %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Less DateTime", input: []interface{}{"dob", tm1}, out: "\"dob\" < %p"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.Less(tst.input[0].(string), tst.input[1]))
	}
}

func TestWhereIn(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL In Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "`miles` IN (%p,%p,%p)"},
		{mgr: &mud.MySQLManager{}, name: "MySQL In Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "`counters` IN (%p,%p)"},
		{mgr: &mud.MySQLManager{}, name: "MySQL In Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "`weight` IN (%p,%p,%p,%p)"},
		{mgr: &mud.MySQLManager{}, name: "MySQL In Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` IN (%p)"},
		{mgr: &mud.MySQLManager{}, name: "MySQL In String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "`name` IN (%p,%p,%p,%p,%p)"},
		{mgr: &mud.MySQLManager{}, name: "MySQL In DateTime", input: []interface{}{"dob", tm1}, out: "`dob` IN (%p)"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL In Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "[miles] IN (%p,%p,%p)"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL In Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "[counters] IN (%p,%p)"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL In Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "[weight] IN (%p,%p,%p,%p)"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL In Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] IN (%p)"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL In String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "[name] IN (%p,%p,%p,%p,%p)"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL In DateTime", input: []interface{}{"dob", tm1}, out: "[dob] IN (%p)"},

		{mgr: &mud.SqliteManager{}, name: "SQLite In Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "\"miles\" IN (%p,%p,%p)"},
		{mgr: &mud.SqliteManager{}, name: "SQLite In Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "\"counters\" IN (%p,%p)"},
		{mgr: &mud.SqliteManager{}, name: "SQLite In Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "\"weight\" IN (%p,%p,%p,%p)"},
		{mgr: &mud.SqliteManager{}, name: "SQLite In Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" IN (%p)"},
		{mgr: &mud.SqliteManager{}, name: "SQLite In String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "\"name\" IN (%p,%p,%p,%p,%p)"},
		{mgr: &mud.SqliteManager{}, name: "SQLite In DateTime", input: []interface{}{"dob", tm1}, out: "\"dob\" IN (%p)"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.In(tst.input[0].(string), tst.input[1:]))
	}
}

func TestWhereLike(t *testing.T) {
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Like", input: []interface{}{"name", "%ma%"}, out: "`name` LIKE %p"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Like", input: []interface{}{"name", "%ma%"}, out: "[name] LIKE %p"},

		{mgr: &mud.SqliteManager{}, name: "SQLite Like", input: []interface{}{"name", "%ma%"}, out: "\"name\" LIKE %p"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.Like(tst.input[0].(string), tst.input[1].(string)))
	}
}

func TestWhereStartsWith(t *testing.T) {
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Starts With", input: []interface{}{"name", "ma"}, out: "`name` LIKE %p", args: []interface{}{"ma%"}},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Starts With", input: []interface{}{"name", "ma"}, out: "[name] LIKE %p", args: []interface{}{"ma%"}},

		{mgr: &mud.SqliteManager{}, name: "SQLite Starts With", input: []interface{}{"name", "ma"}, out: "\"name\" LIKE %p", args: []interface{}{"ma%"}},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.StartsWith(tst.input[0].(string), tst.input[1].(string)))
	}
}

func TestWhereEndsWith(t *testing.T) {
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Ends With", input: []interface{}{"name", "ma"}, out: "`name` LIKE %p", args: []interface{}{"%ma"}},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Ends With", input: []interface{}{"name", "ma"}, out: "[name] LIKE %p", args: []interface{}{"%ma"}},

		{mgr: &mud.SqliteManager{}, name: "SQLite Ends With", input: []interface{}{"name", "ma"}, out: "\"name\" LIKE %p", args: []interface{}{"%ma"}},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.EndsWith(tst.input[0].(string), tst.input[1].(string)))
	}
}

func TestWhereContains(t *testing.T) {
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Contains", input: []interface{}{"name", "ma"}, out: "`name` LIKE %p", args: []interface{}{"%ma%"}},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Contains", input: []interface{}{"name", "ma"}, out: "[name] LIKE %p", args: []interface{}{"%ma%"}},

		{mgr: &mud.SqliteManager{}, name: "SQLite Contains", input: []interface{}{"name", "ma"}, out: "\"name\" LIKE %p", args: []interface{}{"%ma%"}},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.Contains(tst.input[0].(string), tst.input[1].(string)))
	}
}

func TestWhereBetween(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tm2 := time.Date(2020, 2, 7, 22, 0, 0, 0, time.UTC)
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL between Int", input: []interface{}{"amount", 7, 3}, out: "`amount` BETWEEN %p AND %p", args: []interface{}{3, 7}},
		{mgr: &mud.MySQLManager{}, name: "MySQL between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "`miles` BETWEEN %p AND %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "`height` BETWEEN %p AND %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "`weight` BETWEEN %p AND %p", args: []interface{}{float64(73.3212), float64(983.23)}},
		{mgr: &mud.MySQLManager{}, name: "MySQL between DateTime", input: []interface{}{"dob", tm1, tm2}, out: "`dob` BETWEEN %p AND %p"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL between Int", input: []interface{}{"amount", 7, 3}, out: "[amount] BETWEEN %p AND %p", args: []interface{}{3, 7}},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "[miles] BETWEEN %p AND %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "[height] BETWEEN %p AND %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "[weight] BETWEEN %p AND %p", args: []interface{}{float64(73.3212), float64(983.23)}},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL between DateTime", input: []interface{}{"dob", tm1, tm2}, out: "[dob] BETWEEN %p AND %p"},

		{mgr: &mud.SqliteManager{}, name: "SQLite between Int", input: []interface{}{"amount", 7, 3}, out: "\"amount\" BETWEEN %p AND %p", args: []interface{}{3, 7}},
		{mgr: &mud.SqliteManager{}, name: "SQLite between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "\"miles\" BETWEEN %p AND %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "\"height\" BETWEEN %p AND %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "\"weight\" BETWEEN %p AND %p", args: []interface{}{float64(73.3212), float64(983.23)}},
		{mgr: &mud.SqliteManager{}, name: "SQLite between DateTime", input: []interface{}{"dob", tm1, tm2}, out: "\"dob\" BETWEEN %p AND %p"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.Between(tst.input[0].(string), tst.input[1], tst.input[2]))
	}
}

func TestWhereNotBetween(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tm2 := time.Date(2020, 2, 7, 22, 0, 0, 0, time.UTC)
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Not between Int", input: []interface{}{"amount", 7, 3}, out: "`amount` NOT BETWEEN %p AND %p", args: []interface{}{3, 7}},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "`miles` NOT BETWEEN %p AND %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "`height` NOT BETWEEN %p AND %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "`weight` NOT BETWEEN %p AND %p", args: []interface{}{float64(73.3212), float64(983.23)}},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not between DateTime", input: []interface{}{"dob", tm1, tm2}, out: "`dob` NOT BETWEEN %p AND %p"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not between Int", input: []interface{}{"amount", 7, 3}, out: "[amount] NOT BETWEEN %p AND %p", args: []interface{}{3, 7}},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "[miles] NOT BETWEEN %p AND %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "[height] NOT BETWEEN %p AND %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "[weight] NOT BETWEEN %p AND %p", args: []interface{}{float64(73.3212), float64(983.23)}},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not between DateTime", input: []interface{}{"dob", tm1, tm2}, out: "[dob] NOT BETWEEN %p AND %p"},

		{mgr: &mud.SqliteManager{}, name: "SQLite Not between Int", input: []interface{}{"amount", 7, 3}, out: "\"amount\" NOT BETWEEN %p AND %p", args: []interface{}{3, 7}},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "\"miles\" NOT BETWEEN %p AND %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "\"height\" NOT BETWEEN %p AND %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "\"weight\" NOT BETWEEN %p AND %p", args: []interface{}{float64(73.3212), float64(983.23)}},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not between DateTime", input: []interface{}{"dob", tm1, tm2}, out: "\"dob\" NOT BETWEEN %p AND %p"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.NotBetween(tst.input[0].(string), tst.input[1], tst.input[2]))
	}
}

func TestWhereNotEqual(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Not equal Bool", input: []interface{}{"valid", true}, out: "`valid` <> %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not equal Int", input: []interface{}{"miles", int(829)}, out: "`miles` <> %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` <> %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not equal Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` <> %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not equal Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` <> %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not equal String", input: []interface{}{"name", "Sally"}, out: "`name` <> %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not equal DateTime", input: []interface{}{"dob", tm1}, out: "`dob` <> %p"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not equal Bool", input: []interface{}{"valid", true}, out: "[valid] <> %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not equal Int", input: []interface{}{"miles", int(829)}, out: "[miles] <> %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] <> %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not equal Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] <> %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not equal Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] <> %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not equal String", input: []interface{}{"name", "Sally"}, out: "[name] <> %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not equal DateTime", input: []interface{}{"dob", tm1}, out: "[dob] <> %p"},

		{mgr: &mud.SqliteManager{}, name: "SQLite Not equal Bool", input: []interface{}{"valid", true}, out: "\"valid\" <> %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not equal Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" <> %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" <> %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not equal Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" <> %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not equal Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" <> %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not equal String", input: []interface{}{"name", "Sally"}, out: "\"name\" <> %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not equal DateTime", input: []interface{}{"dob", tm1}, out: "\"dob\" <> %p"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.NotEqual(tst.input[0].(string), tst.input[1]))
	}
}

func TestWhereNotGreater(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Not greater Int", input: []interface{}{"miles", int(829)}, out: "`miles` <= %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` <= %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not greater Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` <= %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not greater Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` <= %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not greater DateTime", input: []interface{}{"dob", tm1}, out: "`dob` <= %p"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not greater Int", input: []interface{}{"miles", int(829)}, out: "[miles] <= %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] <= %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not greater Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] <= %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not greater Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] <= %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not greater DateTime", input: []interface{}{"dob", tm1}, out: "[dob] <= %p"},

		{mgr: &mud.SqliteManager{}, name: "SQLite Not greater Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" <= %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" <= %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not greater Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" <= %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not greater Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" <= %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not greater DateTime", input: []interface{}{"dob", tm1}, out: "\"dob\" <= %p"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.NotGreater(tst.input[0].(string), tst.input[1]))
	}
}

func TestWhereNotLess(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Not less Int", input: []interface{}{"miles", int(829)}, out: "`miles` >= %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not less Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` >= %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not less Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` >= %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not less Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` >= %p"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not less DateTime", input: []interface{}{"dob", tm1}, out: "`dob` >= %p"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not less Int", input: []interface{}{"miles", int(829)}, out: "[miles] >= %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not less Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] >= %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not less Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] >= %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not less Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] >= %p"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not less DateTime", input: []interface{}{"dob", tm1}, out: "[dob] >= %p"},

		{mgr: &mud.SqliteManager{}, name: "SQLite Not less Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" >= %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not less Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" >= %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not less Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" >= %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not less Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" >= %p"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not less DateTime", input: []interface{}{"dob", tm1}, out: "\"dob\" >= %p"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.NotLess(tst.input[0].(string), tst.input[1]))
	}
}

func TestWhereNotIn(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Not in Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "`miles` NOT IN (%p,%p,%p)"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not in Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "`counters` NOT IN (%p,%p)"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not in Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "`weight` NOT IN (%p,%p,%p,%p)"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not in Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` NOT IN (%p)"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not in String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "`name` NOT IN (%p,%p,%p,%p,%p)"},
		{mgr: &mud.MySQLManager{}, name: "MySQL Not in DateTime", input: []interface{}{"dob", tm1}, out: "`dob` NOT IN (%p)"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not in Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "[miles] NOT IN (%p,%p,%p)"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not in Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "[counters] NOT IN (%p,%p)"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not in Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "[weight] NOT IN (%p,%p,%p,%p)"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not in Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] NOT IN (%p)"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not in String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "[name] NOT IN (%p,%p,%p,%p,%p)"},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not in DateTime", input: []interface{}{"dob", tm1}, out: "[dob] NOT IN (%p)"},

		{mgr: &mud.SqliteManager{}, name: "SQLite Not in Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "\"miles\" NOT IN (%p,%p,%p)"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not in Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "\"counters\" NOT IN (%p,%p)"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not in Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "\"weight\" NOT IN (%p,%p,%p,%p)"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not in Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" NOT IN (%p)"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not in String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "\"name\" NOT IN (%p,%p,%p,%p,%p)"},
		{mgr: &mud.SqliteManager{}, name: "SQLite Not in DateTime", input: []interface{}{"dob", tm1}, out: "\"dob\" NOT IN (%p)"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.NotIn(tst.input[0].(string), tst.input[1:]))
	}
}

func TestWhereNotLike(t *testing.T) {
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Not Like", input: []interface{}{"name", "%ma%"}, out: "`name` NOT LIKE %p"},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not Like", input: []interface{}{"name", "%ma%"}, out: "[name] NOT LIKE %p"},

		{mgr: &mud.SqliteManager{}, name: "SQLite Not Like", input: []interface{}{"name", "%ma%"}, out: "\"name\" NOT LIKE %p"},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.NotLike(tst.input[0].(string), tst.input[1].(string)))
	}
}

func TestWhereNotStartsWith(t *testing.T) {
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Not Starts With", input: []interface{}{"name", "ma"}, out: "`name` NOT LIKE %p", args: []interface{}{"ma%"}},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not Starts With", input: []interface{}{"name", "ma"}, out: "[name] NOT LIKE %p", args: []interface{}{"ma%"}},

		{mgr: &mud.SqliteManager{}, name: "SQLite Not Starts With", input: []interface{}{"name", "ma"}, out: "\"name\" NOT LIKE %p", args: []interface{}{"ma%"}},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.NotStartsWith(tst.input[0].(string), tst.input[1].(string)))
	}
}

func TestWhereNotEndsWith(t *testing.T) {
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Not Ends With", input: []interface{}{"name", "ma"}, out: "`name` NOT LIKE %p", args: []interface{}{"%ma"}},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not Ends With", input: []interface{}{"name", "ma"}, out: "[name] NOT LIKE %p", args: []interface{}{"%ma"}},

		{mgr: &mud.SqliteManager{}, name: "SQLite Not Ends With", input: []interface{}{"name", "ma"}, out: "\"name\" NOT LIKE %p", args: []interface{}{"%ma"}},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.NotEndsWith(tst.input[0].(string), tst.input[1].(string)))
	}
}

func TestWhereNotContains(t *testing.T) {
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Not Contains", input: []interface{}{"name", "ma"}, out: "`name` NOT LIKE %p", args: []interface{}{"%ma%"}},

		{mgr: &mud.MSSQLManager{}, name: "MSSQL Not Contains", input: []interface{}{"name", "ma"}, out: "[name] NOT LIKE %p", args: []interface{}{"%ma%"}},

		{mgr: &mud.SqliteManager{}, name: "SQLite Not Contains", input: []interface{}{"name", "ma"}, out: "\"name\" NOT LIKE %p", args: []interface{}{"%ma%"}},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.NotContains(tst.input[0].(string), tst.input[1].(string)))
	}
}

func TestWhereIsNull(t *testing.T) {
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Is Null", input: []interface{}{"name"}, out: "`name` IS NULL", args: []interface{}{}},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Is Null", input: []interface{}{"name"}, out: "[name] IS NULL", args: []interface{}{}},
		{mgr: &mud.SqliteManager{}, name: "SQLite Is Null", input: []interface{}{"name"}, out: "\"name\" IS NULL", args: []interface{}{}},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.IsNull(tst.input[0].(string)))
	}
}

func TestWhereIsNotNull(t *testing.T) {
	tests := []whereCase{
		{mgr: &mud.MySQLManager{}, name: "MySQL Is Not Null", input: []interface{}{"name"}, out: "`name` IS NOT NULL", args: []interface{}{}},
		{mgr: &mud.MSSQLManager{}, name: "MSSQL Is Not Null", input: []interface{}{"name"}, out: "[name] IS NOT NULL", args: []interface{}{}},
		{mgr: &mud.SqliteManager{}, name: "SQLite Is Not Null", input: []interface{}{"name"}, out: "\"name\" IS NOT NULL", args: []interface{}{}},
	}
	for _, tst := range tests {
		runWhereCase(t, tst, where.IsNotNull(tst.input[0].(string)))
	}
}

func TestWhereAndConjunction(t *testing.T) {
	mgr := &mud.MySQLManager{}
	e := "`Age` = ? AND `Name` = ?"
	r, args := where.Equal("Age", 12).AndEqual("Name", "Alex").Build(mgr.Operators(), mgr.Placeholder)
	if r != e {
		t.Errorf("Expected %v, got %v", e, r)
	}
	if !reflect.DeepEqual(args, []interface{}{12, "Alex"}) {
		t.Errorf("Expected args %v, got %v", []interface{}{12, "Alex"}, args)
	}
}

func TestWhereOrConjunction(t *testing.T) {
	mgr := &mud.MySQLManager{}
	e := "`Age` = ? OR `Name` <> ?"
	r, args := where.Equal("Age", 12).OrNotEqual("Name", "Alex").Build(mgr.Operators(), mgr.Placeholder)
	if r != e {
		t.Errorf("Expected %v, got %v", e, r)
	}
	if !reflect.DeepEqual(args, []interface{}{12, "Alex"}) {
		t.Errorf("Expected args %v, got %v", []interface{}{12, "Alex"}, args)
	}
}

func TestWhereMSSQLNumberedPlaceholders(t *testing.T) {
	mgr := &mud.MSSQLManager{}
	e := "[Age] = @p1 OR [Name] <> @p2"
	r, args := where.Equal("Age", 12).OrNotEqual("Name", "Alex").Build(mgr.Operators(), mgr.Placeholder)
	if r != e {
		t.Errorf("Expected %v, got %v", e, r)
	}
	if !reflect.DeepEqual(args, []interface{}{12, "Alex"}) {
		t.Errorf("Expected args %v, got %v", []interface{}{12, "Alex"}, args)
	}
}

func TestWhereIn2(t *testing.T) {
	mgr := &mud.MySQLManager{}
	type inCase struct {
		out  string
		args []interface{}
	}
	expected := []inCase{
		{"`ID` IN (?,?,?,?)", []interface{}{1, 2, 3, 4}},
		{"`Name` IN (?,?,?)", []interface{}{"Mark", "Sally", "Oliver"}},
		{"`Age` IN (?)", []interface{}{42}},
		{"`Colour` IN (?)", []interface{}{"RED"}},
	}
	builders := []*where.Builder{
		where.In("ID", []int{1, 2, 3, 4}),
		where.In("Name", []string{"Mark", "Sally", "Oliver"}),
		where.In("Age", 42),
		where.In("Colour", "RED"),
	}
	for i, b := range builders {
		result, args := b.Build(mgr.Operators(), mgr.Placeholder)
		if result != expected[i].out {
			t.Errorf("expecting '%s' got '%s'", expected[i].out, result)
		}
		if !reflect.DeepEqual(args, expected[i].args) {
			t.Errorf("expecting args %v got %v", expected[i].args, args)
		}
	}
}

func TestWhereNotIn2(t *testing.T) {
	mgr := &mud.MySQLManager{}
	type inCase struct {
		out  string
		args []interface{}
	}
	expected := []inCase{
		{"`ID` NOT IN (?,?,?,?)", []interface{}{1, 2, 3, 4}},
		{"`Name` NOT IN (?,?,?)", []interface{}{"Mark", "Sally", "Oliver"}},
	}
	builders := []*where.Builder{
		where.NotIn("ID", []int{1, 2, 3, 4}),
		where.NotIn("Name", []string{"Mark", "Sally", "Oliver"}),
	}
	for i, b := range builders {
		result, args := b.Build(mgr.Operators(), mgr.Placeholder)
		if result != expected[i].out {
			t.Errorf("expecting '%s' got '%s'", expected[i].out, result)
		}
		if !reflect.DeepEqual(args, expected[i].args) {
			t.Errorf("expecting args %v got %v", expected[i].args, args)
		}
	}
}

func TestWhereAndIn2(t *testing.T) {
	mgr := &mud.MySQLManager{}
	type inCase struct {
		out  string
		args []interface{}
	}
	expected := []inCase{
		{"`ID` = ? AND `Size` IN (?,?,?)", []interface{}{2, 2, 4, 6}},
		{"`ID` = ? AND `Name` IN (?,?,?)", []interface{}{3, "Mark", "Sally", "Oliver"}},
	}
	builders := []*where.Builder{
		where.Equal("ID", 2).AndIn("Size", []int{2, 4, 6}),
		where.Equal("ID", 3).AndIn("Name", []string{"Mark", "Sally", "Oliver"}),
	}
	for i, b := range builders {
		result, args := b.Build(mgr.Operators(), mgr.Placeholder)
		if result != expected[i].out {
			t.Errorf("expecting '%s' got '%s'", expected[i].out, result)
		}
		if !reflect.DeepEqual(args, expected[i].args) {
			t.Errorf("expecting args %v got %v", expected[i].args, args)
		}
	}
}

func TestWhereNestedBuilder(t *testing.T) {
	mgr := &mud.SqliteManager{}
	e := "\"Status\" = ? AND \"Age\" > ? OR (\"Role\" = ?)"
	r, args := where.Equal("Status", "active").
		AndGreater("Age", 21).
		OrSub(where.Equal("Role", "admin")).
		Build(mgr.Operators(), mgr.Placeholder)
	if r != e {
		t.Errorf("Expected %v, got %v", e, r)
	}
	wantArgs := []interface{}{"active", 21, "admin"}
	if !reflect.DeepEqual(args, wantArgs) {
		t.Errorf("Expected args %v, got %v", wantArgs, args)
	}
}
