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
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.NotIn(tst.input[0].(string), tst.input[1:]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
	for _, tst := range tests {
		ops := tst.mgr.Operators()
		if got := where.NotIn(tst.input[0].(string), tst.input[1:]).String(ops); got != tst.out {
			t.Errorf("%v : Expected %v, got %v", tst.name, tst.out, got)
		}
	}
}

// func TestWhereNotLike(t *testing.T) {
// 	cs := "`name` not like '%ma%'"
// 	if ts := where.NotLike("name", "%ma%").String(); ts != cs {
// 		t.Errorf("Expecting %v, got %v", cs, ts)
// 	}
// }

// func TestWhereNotStartsWith(t *testing.T) {
// 	cs := "`model` not like 'atar%'"
// 	if ts := where.NotStartsWith("model", "atar").String(); ts != cs {
// 		t.Errorf("Expecting %v, got %v", cs, ts)
// 	}
// }

// func TestWhereNotEndsWith(t *testing.T) {
// 	cs := "`product` not like 'ole%'"
// 	if ts := where.NotStartsWith("product", "ole").String(); ts != cs {
// 		t.Errorf("Expecting %v, got %v", cs, ts)
// 	}
// }

// func TestWhereNotContains(t *testing.T) {
// 	cs := "`breed` not like '%ige%'"
// 	if ts := where.NotContains("breed", "ige").String(); ts != cs {
// 		t.Errorf("Expecting %v, got %v", cs, ts)
// 	}
// }

// func TestWhereAndConjunction(t *testing.T) {
// 	e := "`Age` = 12 AND `Name` = 'Alex'"
// 	if r := where.Equal("Age", 12).AndEqual("Name", "Alex").String(); r != e {
// 		t.Errorf("Expected %v, got %v", e, r)
// 	}
// }

// func TestWhereOrConjunction(t *testing.T) {
// 	e := "`Age` = 12 OR `Name` <> 'Alex'"
// 	if r := where.Equal("Age", 12).OrNotEqual("Name", "Alex").String(); r != e {
// 		t.Errorf("Expected %v, got %v", e, r)
// 	}
// }

// func TestWhereIn2(t *testing.T) {
// 	expected := []string{
// 		"`ID` in (1,2,3,4)",
// 		"`Name` in ('Mark','Sally','Oliver')",
// 		"`Age` in (42)",
// 		"`Colour` in ('RED')",
// 	}
// 	d := []int{1, 2, 3, 4}
// 	s := []string{"Mark", "Sally", "Oliver"}
// 	sd := 42
// 	ss := "RED"
// 	result := where.In("ID", d).String()
// 	if result != expected[0] {
// 		t.Errorf("expecting '%s' got '%s'", expected[0], result)
// 	}
// 	result = where.In("Name", s).String()
// 	if result != expected[1] {
// 		t.Errorf("expecting '%s' got '%s'", expected[1], result)
// 	}
// 	result = where.In("Age", sd).String()
// 	if result != expected[2] {
// 		t.Errorf("expecting '%s' got '%s'", expected[2], result)
// 	}
// 	result = where.In("Colour", ss).String()
// 	if result != expected[3] {
// 		t.Errorf("expecting '%s' got '%s'", expected[3], result)
// 	}

// }

// func TestWhereNotIn2(t *testing.T) {
// 	expected := []string{
// 		"`ID` not in (1,2,3,4)",
// 		"`Name` not in ('Mark','Sally','Oliver')",
// 	}
// 	d := []int{1, 2, 3, 4}
// 	s := []string{"Mark", "Sally", "Oliver"}
// 	result := where.NotIn("ID", d).String()
// 	if result != expected[0] {
// 		t.Errorf("expecting '%s' got '%s'", expected[0], result)
// 	}
// 	result = where.NotIn("Name", s).String()
// 	if result != expected[1] {
// 		t.Errorf("expecting '%s' got '%s'", expected[1], result)
// 	}
// }

// func TestWhereAndIn2(t *testing.T) {
// 	expected := []string{
// 		"`ID` = 2 AND `Size` in (2,4,6)",
// 		"`ID` = 3 AND `Name` in ('Mark','Sally','Oliver')",
// 	}
// 	d := []int{2, 4, 6}
// 	s := []string{"Mark", "Sally", "Oliver"}
// 	result := where.Equal("ID", 2).AndIn("Size", d).String()
// 	if result != expected[0] {
// 		t.Errorf("expecting '%s' got '%s'", expected[0], result)
// 	}
// 	result = where.Equal("ID", 3).AndIn("Name", s).String()
// 	if result != expected[1] {
// 		t.Errorf("expecting '%s' got '%s'", expected[1], result)
// 	}
// }
