package dtormtests

import (
	"fmt"
	"testing"
	"time"

	"github.com/markoxley/dtorm"
	"github.com/markoxley/dtorm/utils"
	"github.com/markoxley/dtorm/where"
)

func TestWhereBetween(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tm2 := time.Date(2020, 2, 7, 22, 0, 0, 0, time.UTC)
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Between Int", input: []interface{}{"amount", 7, 3}, out: "`amount` BETWEEN 3 AND 7"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "`miles` BETWEEN 12987 AND 898989"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "`height` BETWEEN 13.2300 AND 45.4000"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "`weight` BETWEEN 73.321200 AND 983.230000"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Between DateTime", input: []interface{}{"dob", tm1, tm2}, out: fmt.Sprintf("`dob` BETWEEN '%s' AND '%s'", utils.TimeToSQL(tm1), utils.TimeToSQL(tm2))},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Between Int", input: []interface{}{"amount", 7, 3}, out: "[amount] BETWEEN 3 AND 7"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "[miles] BETWEEN 12987 AND 898989"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "[height] BETWEEN 13.2300 AND 45.4000"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "[weight] BETWEEN 73.321200 AND 983.230000"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Between DateTime", input: []interface{}{"dob", tm1, tm2}, out: fmt.Sprintf("[dob] BETWEEN '%s' AND '%s'", utils.TimeToSQL(tm1), utils.TimeToSQL(tm2))},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Between Int", input: []interface{}{"amount", 7, 3}, out: "\"amount\" BETWEEN 3 AND 7"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "\"miles\" BETWEEN 12987 AND 898989"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "\"height\" BETWEEN 13.2300 AND 45.4000"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "\"weight\" BETWEEN 73.321200 AND 983.230000"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Between DateTime", input: []interface{}{"dob", tm1, tm2}, out: fmt.Sprintf("\"dob\" BETWEEN '%s' AND '%s'", utils.TimeToSQL(tm1), utils.TimeToSQL(tm2))},
	}

	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.Between(tst.input[0].(string), tst.input[1], tst.input[2]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereEqual(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Equal Bool", input: []interface{}{"valid", true}, out: "`valid` = 1"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Equal Int", input: []interface{}{"miles", int(829)}, out: "`miles` = 829"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` = 1003322443"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Equal Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` = 73.1200"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Equal Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` = 432.543300"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Equal String", input: []interface{}{"name", "Sally"}, out: "`name` = 'Sally'"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Equal DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` = '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Equal Bool", input: []interface{}{"valid", true}, out: "[valid] = 1"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Equal Int", input: []interface{}{"miles", int(829)}, out: "[miles] = 829"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] = 1003322443"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Equal Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] = 73.1200"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Equal Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] = 432.543300"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Equal String", input: []interface{}{"name", "Sally"}, out: "[name] = 'Sally'"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Equal DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("[dob] = '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Equal Bool", input: []interface{}{"valid", true}, out: "\"valid\" = 1"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Equal Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" = 829"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" = 1003322443"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Equal Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" = 73.1200"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Equal Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" = 432.543300"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Equal String", input: []interface{}{"name", "Sally"}, out: "\"name\" = 'Sally'"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Equal DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("\"dob\" = '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.Equal(tst.input[0].(string), tst.input[1]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereGreater(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Greater Int", input: []interface{}{"miles", int(829)}, out: "`miles` > 829"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` > 1003322443"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Greater Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` > 73.1200"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Greater Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` > 432.543300"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Greater DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` > '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Greater Int", input: []interface{}{"miles", int(829)}, out: "[miles] > 829"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] > 1003322443"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Greater Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] > 73.1200"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Greater Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] > 432.543300"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Greater DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("[dob] > '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Greater Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" > 829"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" > 1003322443"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Greater Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" > 73.1200"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Greater Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" > 432.543300"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Greater DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("\"dob\" > '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.Greater(tst.input[0].(string), tst.input[1]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereLess(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Less Int", input: []interface{}{"miles", int(829)}, out: "`miles` < 829"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Less Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` < 1003322443"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Less Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` < 73.1200"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Less Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` < 432.543300"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Less DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` < '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Less Int", input: []interface{}{"miles", int(829)}, out: "[miles] < 829"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Less Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] < 1003322443"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Less Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] < 73.1200"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Less Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] < 432.543300"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Less DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("[dob] < '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Less Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" < 829"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Less Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" < 1003322443"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Less Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" < 73.1200"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Less Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" < 432.543300"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Less DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("\"dob\" < '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.Less(tst.input[0].(string), tst.input[1]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereIn(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL In Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "`miles` IN (829,21,1)"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL In Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "`counters` IN (1003322443,437216784)"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL In Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "`weight` IN (73.1200,1.4300,0.7600,32.2000)"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL In Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` IN (432.543300)"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL In String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "`name` IN ('Sally','Mark','Jane','Sam','Jack')"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL In DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` IN ('%s')", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL In Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "[miles] IN (829,21,1)"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL In Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "[counters] IN (1003322443,437216784)"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL In Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "[weight] IN (73.1200,1.4300,0.7600,32.2000)"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL In Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] IN (432.543300)"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL In String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "[name] IN ('Sally','Mark','Jane','Sam','Jack')"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL In DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("[dob] IN ('%s')", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite In Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "\"miles\" IN (829,21,1)"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite In Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "\"counters\" IN (1003322443,437216784)"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite In Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "\"weight\" IN (73.1200,1.4300,0.7600,32.2000)"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite In Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" IN (432.543300)"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite In String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "\"name\" IN ('Sally','Mark','Jane','Sam','Jack')"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite In DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("\"dob\" IN ('%s')", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.In(tst.input[0].(string), tst.input[1:]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereLike(t *testing.T) {
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Like", input: []interface{}{"name", "%ma%"}, out: "`name` LIKE '%ma%'"},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Like", input: []interface{}{"name", "%ma%"}, out: "[name] LIKE '%ma%'"},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Like", input: []interface{}{"name", "%ma%"}, out: "\"name\" LIKE '%ma%'"},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.Like(tst.input[0].(string), tst.input[1].(string)).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereStartsWith(t *testing.T) {
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Starts With", input: []interface{}{"name", "ma"}, out: "`name` LIKE 'ma%'"},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Starts With", input: []interface{}{"name", "ma"}, out: "[name] LIKE 'ma%'"},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Starts With", input: []interface{}{"name", "ma"}, out: "\"name\" LIKE 'ma%'"},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.StartsWith(tst.input[0].(string), tst.input[1].(string)).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereEndsWith(t *testing.T) {
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Ends With", input: []interface{}{"name", "ma"}, out: "`name` LIKE '%ma'"},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Ends With", input: []interface{}{"name", "ma"}, out: "[name] LIKE '%ma'"},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Ends With", input: []interface{}{"name", "ma"}, out: "\"name\" LIKE '%ma'"},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.EndsWith(tst.input[0].(string), tst.input[1].(string)).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereContains(t *testing.T) {
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Contains", input: []interface{}{"name", "ma"}, out: "`name` LIKE '%ma%'"},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Contains", input: []interface{}{"name", "ma"}, out: "[name] LIKE '%ma%'"},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Contains", input: []interface{}{"name", "ma"}, out: "\"name\" LIKE '%ma%'"},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.Contains(tst.input[0].(string), tst.input[1].(string)).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotBetween(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tm2 := time.Date(2020, 2, 7, 22, 0, 0, 0, time.UTC)
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not between Int", input: []interface{}{"amount", 7, 3}, out: "`amount` NOT BETWEEN 3 AND 7"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "`miles` NOT BETWEEN 12987 AND 898989"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "`height` NOT BETWEEN 13.2300 AND 45.4000"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "`weight` NOT BETWEEN 73.321200 AND 983.230000"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not between DateTime", input: []interface{}{"dob", tm1, tm2}, out: fmt.Sprintf("`dob` NOT BETWEEN '%s' AND '%s'", utils.TimeToSQL(tm1), utils.TimeToSQL(tm2))},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not between Int", input: []interface{}{"amount", 7, 3}, out: "[amount] NOT BETWEEN 3 AND 7"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "[miles] NOT BETWEEN 12987 AND 898989"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "[height] NOT BETWEEN 13.2300 AND 45.4000"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "[weight] NOT BETWEEN 73.321200 AND 983.230000"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not between DateTime", input: []interface{}{"dob", tm1, tm2}, out: fmt.Sprintf("[dob] NOT BETWEEN '%s' AND '%s'", utils.TimeToSQL(tm1), utils.TimeToSQL(tm2))},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not between Int", input: []interface{}{"amount", 7, 3}, out: "\"amount\" NOT BETWEEN 3 AND 7"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not between Long", input: []interface{}{"miles", int64(12987), int64(898989)}, out: "\"miles\" NOT BETWEEN 12987 AND 898989"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not between Float", input: []interface{}{"height", float32(13.23), float32(45.4)}, out: "\"height\" NOT BETWEEN 13.2300 AND 45.4000"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not between Double", input: []interface{}{"weight", float64(983.23), float64(73.3212)}, out: "\"weight\" NOT BETWEEN 73.321200 AND 983.230000"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not between DateTime", input: []interface{}{"dob", tm1, tm2}, out: fmt.Sprintf("\"dob\" NOT BETWEEN '%s' AND '%s'", utils.TimeToSQL(tm1), utils.TimeToSQL(tm2))},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.NotBetween(tst.input[0].(string), tst.input[1], tst.input[2]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotEqual(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not equal Bool", input: []interface{}{"valid", true}, out: "`valid` <> 1"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not equal Int", input: []interface{}{"miles", int(829)}, out: "`miles` <> 829"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` <> 1003322443"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not equal Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` <> 73.1200"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not equal Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` <> 432.543300"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not equal String", input: []interface{}{"name", "Sally"}, out: "`name` <> 'Sally'"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not equal DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` <> '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not equal Bool", input: []interface{}{"valid", true}, out: "[valid] <> 1"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not equal Int", input: []interface{}{"miles", int(829)}, out: "[miles] <> 829"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] <> 1003322443"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not equal Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] <> 73.1200"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not equal Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] <> 432.543300"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not equal String", input: []interface{}{"name", "Sally"}, out: "[name] <> 'Sally'"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not equal DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("[dob] <> '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not equal Bool", input: []interface{}{"valid", true}, out: "\"valid\" <> 1"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not equal Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" <> 829"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not equal Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" <> 1003322443"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not equal Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" <> 73.1200"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not equal Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" <> 432.543300"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not equal String", input: []interface{}{"name", "Sally"}, out: "\"name\" <> 'Sally'"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not equal DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("\"dob\" <> '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.NotEqual(tst.input[0].(string), tst.input[1]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotGreater(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not greater Int", input: []interface{}{"miles", int(829)}, out: "`miles` <= 829"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` <= 1003322443"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not greater Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` <= 73.1200"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not greater Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` <= 432.543300"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not greater DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` <= '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not greater Int", input: []interface{}{"miles", int(829)}, out: "[miles] <= 829"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] <= 1003322443"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not greater Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] <= 73.1200"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not greater Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] <= 432.543300"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not greater DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("[dob] <= '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not greater Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" <= 829"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not greater Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" <= 1003322443"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not greater Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" <= 73.1200"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not greater Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" <= 432.543300"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not greater DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("\"dob\" <= '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.NotGreater(tst.input[0].(string), tst.input[1]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotLess(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not less Int", input: []interface{}{"miles", int(829)}, out: "`miles` >= 829"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not less Long", input: []interface{}{"counters", int64(1003322443)}, out: "`counters` >= 1003322443"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not less Float", input: []interface{}{"weight", float32(73.12)}, out: "`weight` >= 73.1200"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not less Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` >= 432.543300"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not less DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` >= '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not less Int", input: []interface{}{"miles", int(829)}, out: "[miles] >= 829"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not less Long", input: []interface{}{"counters", int64(1003322443)}, out: "[counters] >= 1003322443"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not less Float", input: []interface{}{"weight", float32(73.12)}, out: "[weight] >= 73.1200"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not less Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] >= 432.543300"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not less DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("[dob] >= '%s'", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not less Int", input: []interface{}{"miles", int(829)}, out: "\"miles\" >= 829"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not less Long", input: []interface{}{"counters", int64(1003322443)}, out: "\"counters\" >= 1003322443"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not less Float", input: []interface{}{"weight", float32(73.12)}, out: "\"weight\" >= 73.1200"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not less Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" >= 432.543300"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not less DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("\"dob\" >= '%s'", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.NotLess(tst.input[0].(string), tst.input[1]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotIn(t *testing.T) {
	tm1 := time.Date(1971, 11, 15, 22, 30, 0, 12, time.UTC)
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not in Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "`miles` NOT IN (829,21,1)"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not in Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "`counters` NOT IN (1003322443,437216784)"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not in Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "`weight` NOT IN (73.1200,1.4300,0.7600,32.2000)"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not in Double", input: []interface{}{"height", float64(432.5433)}, out: "`height` NOT IN (432.543300)"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not in String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "`name` NOT IN ('Sally','Mark','Jane','Sam','Jack')"},
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Not in DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("`dob` NOT IN ('%s')", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not in Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "[miles] NOT IN (829,21,1)"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not in Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "[counters] NOT IN (1003322443,437216784)"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not in Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "[weight] NOT IN (73.1200,1.4300,0.7600,32.2000)"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not in Double", input: []interface{}{"height", float64(432.5433)}, out: "[height] NOT IN (432.543300)"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not in String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "[name] NOT IN ('Sally','Mark','Jane','Sam','Jack')"},
		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Not in DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("[dob] NOT IN ('%s')", utils.TimeToSQL(tm1))},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not in Int", input: []interface{}{"miles", int(829), int(21), int(1)}, out: "\"miles\" NOT IN (829,21,1)"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not in Long", input: []interface{}{"counters", int64(1003322443), int64(437216784)}, out: "\"counters\" NOT IN (1003322443,437216784)"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not in Float", input: []interface{}{"weight", float32(73.12), float32(1.43), float32(0.76), float32(32.2)}, out: "\"weight\" NOT IN (73.1200,1.4300,0.7600,32.2000)"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not in Double", input: []interface{}{"height", float64(432.5433)}, out: "\"height\" NOT IN (432.543300)"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not in String", input: []interface{}{"name", "Sally", "Mark", "Jane", "Sam", "Jack"}, out: "\"name\" NOT IN ('Sally','Mark','Jane','Sam','Jack')"},
		{mgr: &dtorm.SqliteManager{}, name: "SQLite Not in DateTime", input: []interface{}{"dob", tm1}, out: fmt.Sprintf("\"dob\" NOT IN ('%s')", utils.TimeToSQL(tm1))},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.NotIn(tst.input[0].(string), tst.input[1:]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}
func TestWhereNotLike(t *testing.T) {
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Like", input: []interface{}{"name", "%ma%"}, out: "`name` NOT LIKE '%ma%'"},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Like", input: []interface{}{"name", "%ma%"}, out: "[name] NOT LIKE '%ma%'"},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Like", input: []interface{}{"name", "%ma%"}, out: "\"name\" NOT LIKE '%ma%'"},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.NotLike(tst.input[0].(string), tst.input[1].(string)).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotStartsWith(t *testing.T) {
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Starts With", input: []interface{}{"name", "ma"}, out: "`name` NOT LIKE 'ma%'"},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Starts With", input: []interface{}{"name", "ma"}, out: "[name] NOT LIKE 'ma%'"},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Starts With", input: []interface{}{"name", "ma"}, out: "\"name\" NOT LIKE 'ma%'"},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.NotStartsWith(tst.input[0].(string), tst.input[1].(string)).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotEndsWith(t *testing.T) {
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Ends With", input: []interface{}{"name", "ma"}, out: "`name` NOT LIKE '%ma'"},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Ends With", input: []interface{}{"name", "ma"}, out: "[name] NOT LIKE '%ma'"},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Ends With", input: []interface{}{"name", "ma"}, out: "\"name\" NOT LIKE '%ma'"},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.NotEndsWith(tst.input[0].(string), tst.input[1].(string)).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereNotContains(t *testing.T) {
	tests := []struct {
		mgr   dtorm.Manager
		name  string
		input []interface{}
		out   string
	}{
		{mgr: &dtorm.MySQLManager{}, name: "MySQL Contains", input: []interface{}{"name", "ma"}, out: "`name` NOT LIKE '%ma%'"},

		{mgr: &dtorm.MSSQLManager{}, name: "MSSQL Contains", input: []interface{}{"name", "ma"}, out: "[name] NOT LIKE '%ma%'"},

		{mgr: &dtorm.SqliteManager{}, name: "SQLite Contains", input: []interface{}{"name", "ma"}, out: "\"name\" NOT LIKE '%ma%'"},
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.NotContains(tst.input[0].(string), tst.input[1].(string)).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

func TestWhereAndConjunction(t *testing.T) {
	mgr := &dtorm.MySQLManager{}
	e := "`Age` = 12 AND `Name` = 'Alex'"
	if r := where.Equal("Age", 12).AndEqual("Name", "Alex").String(mgr.Operators()); r != e {
		t.Errorf("Expected %v, got %v", e, r)
	}
}

func TestWhereOrConjunction(t *testing.T) {
	mgr := &dtorm.MySQLManager{}
	e := "`Age` = 12 OR `Name` <> 'Alex'"
	if r := where.Equal("Age", 12).OrNotEqual("Name", "Alex").String(mgr.Operators()); r != e {
		t.Errorf("Expected %v, got %v", e, r)
	}
}

func TestWhereIn2(t *testing.T) {
	mgr := &dtorm.MySQLManager{}
	ops := mgr.Operators()
	expected := []string{
		"`ID` IN (1,2,3,4)",
		"`Name` IN ('Mark','Sally','Oliver')",
		"`Age` IN (42)",
		"`Colour` IN ('RED')",
	}
	d := []int{1, 2, 3, 4}
	s := []string{"Mark", "Sally", "Oliver"}
	sd := 42
	ss := "RED"
	result := where.In("ID", d).String(ops)
	if result != expected[0] {
		t.Errorf("expecting '%s' got '%s'", expected[0], result)
	}
	result = where.In("Name", s).String(ops)
	if result != expected[1] {
		t.Errorf("expecting '%s' got '%s'", expected[1], result)
	}
	result = where.In("Age", sd).String(ops)
	if result != expected[2] {
		t.Errorf("expecting '%s' got '%s'", expected[2], result)
	}
	result = where.In("Colour", ss).String(ops)
	if result != expected[3] {
		t.Errorf("expecting '%s' got '%s'", expected[3], result)
	}

}

func TestWhereNotIn2(t *testing.T) {
	mgr := &dtorm.MySQLManager{}
	ops := mgr.Operators()
	expected := []string{
		"`ID` NOT IN (1,2,3,4)",
		"`Name` NOT IN ('Mark','Sally','Oliver')",
	}
	d := []int{1, 2, 3, 4}
	s := []string{"Mark", "Sally", "Oliver"}
	result := where.NotIn("ID", d).String(ops)
	if result != expected[0] {
		t.Errorf("expecting '%s' got '%s'", expected[0], result)
	}
	result = where.NotIn("Name", s).String(ops)
	if result != expected[1] {
		t.Errorf("expecting '%s' got '%s'", expected[1], result)
	}
}

func TestWhereAndIn2(t *testing.T) {
	mgr := &dtorm.MySQLManager{}
	ops := mgr.Operators()
	expected := []string{
		"`ID` = 2 AND `Size` IN (2,4,6)",
		"`ID` = 3 AND `Name` IN ('Mark','Sally','Oliver')",
	}
	d := []int{2, 4, 6}
	s := []string{"Mark", "Sally", "Oliver"}
	result := where.Equal("ID", 2).AndIn("Size", d).String(ops)
	if result != expected[0] {
		t.Errorf("expecting '%s' got '%s'", expected[0], result)
	}
	result = where.Equal("ID", 3).AndIn("Name", s).String(ops)
	if result != expected[1] {
		t.Errorf("expecting '%s' got '%s'", expected[1], result)
	}
}
