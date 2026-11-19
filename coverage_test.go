// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
package mud

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/markoxley/mud/order"
	"github.com/markoxley/mud/utils"
	"github.com/markoxley/mud/where"
)

// covEmbed is anonymously embedded to exercise getDefs struct recursion.
type covEmbed struct {
	EmbeddedField string `mud:""`
}

// covModel exercises every field type and both nullable/non-nullable paths.
type covModel struct {
	Model
	covEmbed
	Name    string     `mud:""`
	Age     int        `mud:""`
	Big     int64      `mud:""`
	UInt    uint       `mud:""`
	ULong   uint64     `mud:""`
	Score   float32    `mud:""`
	Price   float64    `mud:"type:decimal,size:10,2"`
	Active  bool       `mud:""`
	Initial string     `mud:"type:char"`
	When    time.Time  `mud:""`
	NName   *string    `mud:""`
	NAge    *int       `mud:""`
	NUInt   *uint      `mud:""`
	NFloat  *float64   `mud:""`
	NBool   *bool      `mud:""`
	NTime   *time.Time `mud:""`
	Blob    *[]byte    `mud:""`
	Skipped string
	hidden  int
}

// Restore implements Restorer so the doRestore positive path is covered.
func (m covModel) Restore(mgr Manager) error { return nil }

// seedModel returns standing data so the tableTest seeding path is covered.
type seedModel struct {
	Model
	Name string `mud:""`
}

func (m seedModel) StandingData() []Modeller {
	return []Modeller{&seedModel{Name: "seeded"}}
}

// errModel implements Updater with a failing update hook.
type errModel struct {
	Model
	Name string `mud:""`
}

func (m *errModel) Update(mgr Manager) error { return errors.New("update failed") }

// updModel implements Updater successfully.
type updModel struct {
	Model
	Name    string `mud:""`
	updated bool
}

func (m *updModel) Update(mgr Manager) error { m.updated = true; return nil }

// stringCrit implements fmt.Stringer to cover the Stringer criteria branch.
type stringCrit string

func (s stringCrit) String() string { return string(s) }

func covDB(t *testing.T) *DB {
	t.Helper()
	db, err := New(&Config{Type: "sqlite", Database: filepath.Join(t.TempDir(), "test.db")})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func newCovModel() *covModel {
	return &covModel{
		covEmbed: covEmbed{EmbeddedField: "emb"},
		Name:     "cover",
		Age:      30,
		Big:      9999999,
		UInt:     7,
		ULong:    12345678901,
		Score:    1.5,
		Price:    19.99,
		Active:   true,
		Initial:  "Q",
		When:     time.Date(2024, 3, 2, 1, 2, 3, 0, time.UTC),
		NName:    utils.Ptr("nullable"),
		NAge:     utils.Ptr(11),
		NUInt:    utils.Ptr(uint(13)),
		NFloat:   utils.Ptr(2.25),
		NBool:    utils.Ptr(true),
		Blob:     utils.Ptr([]byte{1, 2, 3}),
		// NTime deliberately left nil to cover the nil-pointer insert skip.
	}
}

func TestCoverageSaveFetchUpdate(t *testing.T) {
	db := covDB(t)

	m := newCovModel()
	if err := db.Save(m); err != nil {
		t.Fatalf("Save insert: %v", err)
	}
	if m.GetID() == nil {
		t.Fatal("expected ID to be assigned on insert")
	}

	// Fetch and verify all field types round-tripped.
	res, err := db.Fetch(covModel{}, where.Equal("Name", "cover"))
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	f, ok := res[0].(covModel)
	if !ok {
		t.Fatalf("unexpected model type %T", res[0])
	}
	if f.Age != 30 || f.Big != 9999999 || f.UInt != 7 || f.ULong != 12345678901 {
		t.Errorf("integer fields wrong: %+v", f)
	}
	if !f.Active || f.Initial != "Q" || f.EmbeddedField != "emb" {
		t.Errorf("char/bool/embed fields wrong: %+v", f)
	}
	if f.When.IsZero() {
		t.Errorf("time field not populated: %+v", f)
	}
	if f.NName == nil || *f.NName != "nullable" || f.NAge == nil || *f.NAge != 11 ||
		f.NUInt == nil || *f.NUInt != 13 || f.NFloat == nil || *f.NFloat != 2.25 ||
		f.NBool == nil || !*f.NBool {
		t.Errorf("nullable fields wrong: %+v", f)
	}
	if f.NTime != nil {
		t.Errorf("NTime should be nil")
	}
	if f.Blob != nil {
		t.Errorf("Blob should be nil (unsupported type not bound)")
	}

	// Update: null out NName (covers "field = null" assignment) and change Name.
	m.NName = nil
	m.Name = "updated"
	m.NAge = utils.Ptr(42)
	if err := db.Save(m); err != nil {
		t.Fatalf("Save update: %v", err)
	}
	res, err = db.Fetch(covModel{}, where.Equal("Name", "updated"))
	if err != nil || len(res) != 1 {
		t.Fatalf("Fetch after update: %v len=%d", err, len(res))
	}
	fu := res[0].(covModel)
	if fu.NName != nil {
		t.Errorf("expected NName to be NULL after update")
	}
	if fu.NAge == nil || *fu.NAge != 42 {
		t.Errorf("expected NAge=42, got %+v", fu.NAge)
	}
}

func TestCoverageRangeAndFirst(t *testing.T) {
	db := covDB(t)
	for i := 0; i < 3; i++ {
		m := newCovModel()
		m.Name = "range-test"
		if err := db.Save(m); err != nil {
			t.Fatalf("Save: %v", err)
		}
	}

	// db.Range covers populateSingle.
	count := 0
	for r := range db.Range(covModel{}, where.Equal("Name", "range-test")) {
		if _, ok := r.(covModel); !ok {
			t.Fatalf("unexpected type %T", r)
		}
		count++
	}
	if count != 3 {
		t.Fatalf("expected 3 range results, got %d", count)
	}

	// Early exit covers yield==false branch.
	for range db.Range(covModel{}) {
		break
	}

	// Generic Range.
	count = 0
	for range Range[covModel](db) {
		count++
	}
	if count != 3 {
		t.Fatalf("expected 3 generic range results, got %d", count)
	}

	// Generic First and FromID.
	saved := newCovModel()
	if err := db.Save(saved); err != nil {
		t.Fatalf("Save: %v", err)
	}
	one, err := First[covModel](db, where.Equal("ID", *saved.GetID()))
	if err != nil || one == nil {
		t.Fatalf("First: %v", err)
	}
	byID, err := FromID[covModel](db, *saved.GetID())
	if err != nil || byID == nil {
		t.Fatalf("FromID: %v", err)
	}

	// No-results path covers db.First via Refresh.
	missing := covModel{}
	id := "00000000-0000-0000-0000-000000000000"
	missing.ID = &id
	if err := db.Refresh(missing); err == nil {
		t.Error("expected Refresh to fail for missing model")
	}
	if err := db.Refresh(*saved); err != nil {
		t.Errorf("Refresh saved model: %v", err)
	}
	empty := covModel{}
	if err := db.Refresh(empty); err == nil {
		t.Error("expected Refresh with nil ID to fail")
	}
}

func TestCoverageGetCriteria(t *testing.T) {
	db := covDB(t)
	m := newCovModel()
	m.Name = "crit"
	m.Age = 50
	if err := db.Save(m); err != nil {
		t.Fatalf("Save: %v", err)
	}

	checks := []struct {
		name string
		crit []interface{}
	}{
		{"nil", nil},
		{"nil entry", []interface{}{nil}},
		{"*Criteria", []interface{}{&Criteria{Where: where.Equal("Name", "crit")}}},
		{"Criteria", []interface{}{Criteria{Where: where.Equal("Name", "crit")}}},
		{"*where.Builder", []interface{}{where.Equal("Name", "crit")}},
		{"where.Builder", []interface{}{*where.Equal("Name", "crit")}},
		{"*order.Builder", []interface{}{order.Desc("Age")}},
		{"order.Builder", []interface{}{*order.Desc("Age")}},
		{"uuid string", []interface{}{*m.GetID()}},
		{"raw string", []interface{}{"Age > 40"}},
		{"Stringer", []interface{}{stringCrit("Age > 40")}},
	}
	for _, c := range checks {
		res, err := db.Fetch(covModel{}, c.crit...)
		if err != nil {
			t.Errorf("%s: unexpected error %v", c.name, err)
		}
		_ = res
	}

	if _, err := db.Fetch(covModel{}, 3.14); err == nil {
		t.Error("expected error for invalid criteria type")
	}
	if got := db.Count(&covModel{}, 3.14); got != -1 {
		t.Errorf("expected Count=-1 for invalid criteria, got %d", got)
	}
}

func TestCoverageRawAndTransactions(t *testing.T) {
	db := covDB(t)
	m := newCovModel()
	if err := db.Save(m); err != nil {
		t.Fatalf("Save: %v", err)
	}

	rows, err := db.RawSelect(`SELECT Name, Age FROM covModel`, nil)
	if err != nil || len(rows) != 1 {
		t.Fatalf("RawSelect: %v rows=%d", err, len(rows))
	}
	val, ok := db.RawScalar("SELECT COUNT(*) FROM covModel", nil)
	if !ok || val == nil {
		t.Fatalf("RawScalar failed: %v %v", val, ok)
	}
	if err := db.RawExecute("DELETE FROM covModel WHERE 1 = 0", nil); err != nil {
		t.Fatalf("RawExecute: %v", err)
	}

	// Error paths.
	if _, err := db.RawSelect("SELECT * FROM missing_table", nil); err == nil {
		t.Error("expected RawSelect error for missing table")
	}
	if _, ok := db.RawScalar("SELECT * FROM missing_table", nil); ok {
		t.Error("expected RawScalar failure for missing table")
	}
	if err := db.RawExecute("THIS IS NOT SQL", nil); err == nil {
		t.Error("expected RawExecute error for bad SQL")
	}

	// Transactions are nil-ops on sqlite but still exercise the methods.
	tx, err := db.BeginTransaction()
	if err != nil {
		t.Fatalf("BeginTransaction: %v", err)
	}
	if err := db.CommitTransaction(tx); err != nil {
		t.Fatalf("CommitTransaction: %v", err)
	}
	if err := db.RollbackTransaction(tx); err != nil {
		t.Fatalf("RollbackTransaction: %v", err)
	}
}

func TestCoverageRemovePaths(t *testing.T) {
	db := covDB(t)

	// Remove with nil ID returns nil without touching the database.
	if err := db.Remove(&covModel{}); err != nil {
		t.Fatalf("Remove nil ID: %v", err)
	}

	m := newCovModel()
	if err := db.Save(m); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Soft-delete (Deletable=false) covers massDisable via Remove.
	if err := db.Remove(m); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	// RemoveMany on a table that does not exist short-circuits.
	if n, err := db.RemoveMany(&updModel{}, nil); err != nil || n != 0 {
		t.Fatalf("RemoveMany missing table: %d %v", n, err)
	}

	// RemoveMany with zero matching rows.
	if n, err := db.RemoveMany(&covModel{}, &Criteria{Where: where.Equal("Name", "nobody")}); err != nil || n != 0 {
		t.Fatalf("RemoveMany zero rows: %d %v", n, err)
	}

	// RemoveMany matching rows (soft delete path).
	m2 := newCovModel()
	m2.Name = "batch"
	if err := db.Save(m2); err != nil {
		t.Fatalf("Save: %v", err)
	}
	n, err := db.RemoveMany(&covModel{}, &Criteria{Where: where.Equal("Name", "batch")})
	if err != nil || n != 1 {
		t.Fatalf("RemoveMany: %d %v", n, err)
	}
}

func TestCoverageRemoveDeletable(t *testing.T) {
	db, err := New(&Config{
		Type:      "sqlite",
		Database:  filepath.Join(t.TempDir(), "test.db"),
		Deletable: true,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	m := newCovModel()
	if err := db.Save(m); err != nil {
		t.Fatalf("Save: %v", err)
	}
	// Hard delete covers massDelete via Remove.
	if err := db.Remove(m); err != nil {
		t.Fatalf("Remove hard delete: %v", err)
	}

	m2 := newCovModel()
	m2.Name = "batch"
	if err := db.Save(m2); err != nil {
		t.Fatalf("Save: %v", err)
	}
	n, err := db.RemoveMany(&covModel{}, &Criteria{Where: where.Equal("Name", "batch")})
	if err != nil || n != 1 {
		t.Fatalf("RemoveMany hard delete: %d %v", n, err)
	}
}

func TestCoverageStandingData(t *testing.T) {
	db := covDB(t)
	m := &seedModel{Name: "main"}
	if err := db.Save(m); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if got := db.Count(&seedModel{}); got != 2 {
		t.Fatalf("expected standing data + saved model = 2, got %d", got)
	}
	res, err := db.Fetch(seedModel{}, where.Equal("Name", "seeded"))
	if err != nil || len(res) != 1 {
		t.Fatalf("seeded fetch: %v len=%d", err, len(res))
	}
}

func TestCoverageSaveHooks(t *testing.T) {
	db := covDB(t)

	// Updater returning an error aborts Save.
	if err := db.Save(&errModel{}); err == nil {
		t.Error("expected Save to propagate Updater error")
	}

	// Successful Updater runs before insert.
	u := &updModel{Name: "u"}
	if err := db.Save(u); err != nil {
		t.Fatalf("Save with Updater: %v", err)
	}
	if !u.updated {
		t.Error("expected Updater.Update to be called")
	}
}

func TestCoverageNoResults(t *testing.T) {
	db := covDB(t)
	// Ensure the table exists so the query succeeds with zero rows.
	db.tableTest(&covModel{})

	_, err := First[covModel](db, where.Equal("Name", "missing"))
	if err == nil {
		t.Fatal("expected NoResults error")
	}
	var nr ErrNoResults
	if !errors.As(err, &nr) {
		t.Error("expected errors.As to find ErrNoResults")
	}
	if nr.Error() == "" {
		t.Error("expected non-empty error message")
	}
	if nr.Unwrap() == nil {
		t.Error("expected Unwrap to return underlying error")
	}
}

func TestCoverageModelHelpers(t *testing.T) {
	m := CreateModel()
	m.SetTableName("custom_table")
	if got := m.TableName(); got != "custom_table" {
		t.Errorf("expected custom_table, got %q", got)
	}
	if got := GetTableName(&covModel{Model: m}); got != "custom_table" {
		t.Errorf("expected GetTableName override, got %q", got)
	}

	if v := Version(); v == "" {
		t.Error("expected non-empty Version")
	}
	if got := GetTypeName(nil); got != "" {
		t.Errorf("expected empty type name for nil, got %q", got)
	}
	if got := GetTypeName(&SqliteManager{}); got == "" {
		t.Error("expected non-empty type name")
	}

	if !unsignedCapable(tInt) || unsignedCapable(tString) {
		t.Error("unsignedCapable wrong")
	}
	if isDigits("") || isDigits("12a") || !isDigits("123") {
		t.Error("isDigits wrong")
	}
}

func TestCoverageNewErrors(t *testing.T) {
	if _, err := New(nil); err == nil {
		t.Error("expected error for nil config")
	}
	if _, err := New(&Config{Type: "sqlite"}); err == nil {
		t.Error("expected error for empty database path")
	}
	if _, err := New(&Config{Type: "bogus", Database: "x"}); err == nil {
		t.Error("expected error for invalid type")
	}

	// connect returns the sql.Open error for an unknown driver.
	d := &DB{dbtype: "no-such-driver", connectionString: "x"}
	if _, err := d.connect(); err == nil {
		t.Error("expected connect error for unknown driver")
	}
}

func TestCoverageSetFieldValue(t *testing.T) {
	e := reflect.ValueOf(&covModel{}).Elem()

	cases := []struct {
		name string
		fld  field
		s    string
	}{
		{"missing field", field{name: "Nope", fType: tInt}, "1"},
		{"int bad", field{name: "Age", fType: tInt}, "abc"},
		{"int ok", field{name: "Age", fType: tInt}, "7"},
		{"uint bad", field{name: "UInt", fType: tInt, unsigned: true}, "-1"},
		{"uint ok", field{name: "UInt", fType: tInt, unsigned: true}, "9"},
		{"bool int", field{name: "Active", fType: tBool}, "1"},
		{"bool text", field{name: "Active", fType: tBool}, "false"},
		{"bool bad", field{name: "Active", fType: tBool}, "junk"},
		{"float bad", field{name: "Score", fType: tFloat}, "xx"},
		{"float ok", field{name: "Score", fType: tFloat}, "3.5"},
		{"time ok", field{name: "When", fType: tDateTime}, "2024-01-02 03:04:05"},
		{"time bad", field{name: "When", fType: tDateTime}, "bad"},
		{"char empty", field{name: "Initial", fType: tChar}, ""},
		{"char ok", field{name: "Initial", fType: tChar}, "xyz"},
		{"nullable int", field{name: "NAge", fType: tInt, allowNull: true}, "5"},
		{"nullable int again", field{name: "NAge", fType: tInt, allowNull: true}, "6"},
		{"nullable uint", field{name: "NUInt", fType: tInt, unsigned: true, allowNull: true}, "8"},
		{"nullable float", field{name: "NFloat", fType: tDouble, allowNull: true}, "1.5"},
		{"nullable bool", field{name: "NBool", fType: tBool, allowNull: true}, "true"},
		{"nullable time", field{name: "NTime", fType: tDateTime, allowNull: true}, "2024-01-02 03:04:05"},
		{"nullable string", field{name: "NName", fType: tString, allowNull: true}, "v"},
		{"nullable char", field{name: "NName", fType: tChar, allowNull: true}, "xy"},
	}
	for _, c := range cases {
		setFieldValue(e, c.fld, c.s) // must not panic
	}

	m := e.Addr().Interface().(*covModel)
	if m.Age != 7 || m.UInt != 9 || m.Active || m.Score != 3.5 || m.Initial != "x" {
		t.Errorf("scalar assignments wrong: %+v", m)
	}
	if m.NAge == nil || *m.NAge != 6 || m.NUInt == nil || *m.NUInt != 8 ||
		m.NFloat == nil || *m.NFloat != 1.5 || m.NBool == nil || !*m.NBool ||
		m.NTime == nil || m.NName == nil || *m.NName != "x" {
		t.Errorf("nullable assignments wrong: %+v", m)
	}
}

func TestCoverageSeededPos(t *testing.T) {
	// whereClause with an explicit position covers the seeded placeholder path
	// used by massDisable on SQL Server.
	pos := 5
	c := &Criteria{Where: where.Equal("A", 1)}
	mgr := &MSSQLManager{}
	w, args := c.whereClause(mgr, &pos)
	if len(args) != 1 {
		t.Fatalf("expected 1 arg, got %v", args)
	}
	if w == "" || pos != 6 {
		t.Errorf("expected WHERE clause and pos=6, got %q pos=%d", w, pos)
	}

	// Criteria may also hold a where.Builder value (not pointer).
	wb := where.Equal("B", 2)
	c2 := &Criteria{Where: *wb}
	mgr2 := &SqliteManager{}
	w2, args2 := c2.whereClause(mgr2, nil)
	if w2 == "" || len(args2) != 1 {
		t.Errorf("expected WHERE clause with 1 arg, got %q args=%v", w2, args2)
	}
}

func TestCoverageMassDisableMSSQL(t *testing.T) {
	// The sqlserver branch of massDisable seeds the placeholder counter at 2
	// because @p1 is consumed by the SET clause.
	d := &DB{dbtype: "sqlserver", mgr: &MSSQLManager{}}
	c := &Criteria{Where: where.Equal("A", 1)}
	s, args := d.massDisable(&covModel{}, c)
	if !strings.Contains(s, "@p1") || !strings.Contains(s, "@p2") {
		t.Errorf("expected @p1 and @p2 placeholders, got %q", s)
	}
	if len(args) != 2 {
		t.Errorf("expected 2 args, got %v", args)
	}
}

// dupA and dupB each contribute a column named "Same"; embedding both produces
// an invalid CREATE TABLE so every tableTest error path can be exercised.
type dupA struct {
	Same string `mud:""`
}

type dupB struct {
	Same string `mud:""`
}

type dupModel struct {
	Model
	dupA
	dupB
}

func TestCoverageDupColumns(t *testing.T) {
	db := covDB(t)

	// insertCommand -> tableTest -> CREATE fails on duplicate column name.
	if err := db.Save(&dupModel{}); err == nil {
		t.Error("expected Save to fail on duplicate columns")
	}

	// Second attempt: definition is cached so the INSERT itself fails.
	if err := db.Save(&dupModel{}); err == nil {
		t.Error("expected second Save to fail on missing table")
	}

	// updateCommand error path (non-new model).
	dm := &dupModel{}
	id := "x"
	dm.ID = &id
	if err := db.Save(dm); err == nil {
		t.Error("expected update Save to fail")
	}

	// On a fresh database the CREATE fails inside tableTest; Count must run
	// first because tableDefinition caches the (bad) def, after which later
	// calls skip CREATE and fail at query time instead.
	db2 := covDB(t)
	if got := db2.Count(&dupModel{}); got != -1 {
		t.Errorf("expected Count=-1, got %d", got)
	}
	if _, err := db2.Fetch(dupModel{}); err == nil {
		t.Error("expected Fetch to fail")
	}

	// On the original database the definition is cached, so the query itself
	// fails and Count reports 0.
	if _, err := db.Fetch(dupModel{}); err == nil {
		t.Error("expected Fetch to fail on missing table")
	}
	if got := db.Count(&dupModel{}); got != 0 {
		t.Errorf("expected Count=0 for query failure, got %d", got)
	}
}

// failSeedModel's standing data cannot be created, exercising the
// propagation of seed errors out of tableTest.
type failSeedModel struct {
	Model
	Name string `mud:""`
}

func (m failSeedModel) StandingData() []Modeller {
	return []Modeller{&dupModel{}}
}

func TestCoverageStandingDataFail(t *testing.T) {
	db := covDB(t)
	if err := db.Save(&failSeedModel{}); err == nil {
		t.Error("expected Save to propagate standing-data failure")
	}
}

func TestCoverageRangeErrors(t *testing.T) {
	db := covDB(t)
	m := newCovModel()
	if err := db.Save(m); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Invalid criteria aborts the iterator.
	for range db.Range(covModel{}, 42) {
		t.Error("should not yield on invalid criteria")
	}

	// Malformed raw SQL errors the query.
	for range db.Range(covModel{}, "WHERE ((") {
		t.Error("should not yield on bad SQL")
	}

	// tableTest failure aborts the iterator.
	for range db.Range(dupModel{}) {
		t.Error("should not yield when the table cannot be created")
	}
}

func TestCoveragePopulateMiss(t *testing.T) {
	db := covDB(t)
	rows, err := db.db.Query("SELECT 1 AS A, 2 AS B")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	defer rows.Close()

	// updModel has no cached table definition in this database.
	if _, ok := db.populateSingle(updModel{}, rows); ok {
		t.Error("expected populateSingle to fail without a table definition")
	}
	if _, ok := db.populateModel(updModel{}, rows); ok {
		t.Error("expected populateModel to fail without a table definition")
	}
}

func TestCoverageRealTransactions(t *testing.T) {
	db := covDB(t)

	tx, err := db.db.Begin()
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if err := db.CommitTransaction(tx); err != nil {
		t.Errorf("CommitTransaction: %v", err)
	}

	tx2, err := db.db.Begin()
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}
	if err := db.RollbackTransaction(tx2); err != nil {
		t.Errorf("RollbackTransaction: %v", err)
	}
}

func TestCoverageSoftDeletedFetch(t *testing.T) {
	db := covDB(t)
	m := newCovModel()
	if err := db.Save(m); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if err := db.Remove(m); err != nil {
		t.Fatalf("Remove: %v", err)
	}

	// The default criteria excludes soft-deleted rows (DeleteDate IS NULL);
	// IncDeleted includes them so the DeleteDate populate branch fires.
	res, err := db.Fetch(covModel{}, &Criteria{IncDeleted: true})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	found := false
	for _, r := range res {
		f := r.(covModel)
		if f.DeleteDate != nil {
			found = true
		}
	}
	if !found {
		t.Error("expected a fetched row with DeleteDate set")
	}
}

type covInt int

// taggy exercises the tag-parsing edge paths in getDefs.
type taggy struct {
	A       string    `mud:"type:bogus"`  // unknown type override is ignored
	B       string    `mud:"size:x"`      // non-numeric size is ignored
	C       string    `mud:"size:10,key"` // size followed by a non-digit token
	D       int       `mud:"identity,key,unsigned"`
	E       int       `mud:"type:time"` // type:time maps to the datetime type
	F       *covEmbed `mud:""`          // nil pointer to struct is materialised
	*covInt           // anonymous non-struct falls through
}

func TestCoverageGetDefs(t *testing.T) {
	byName := make(map[string]field)
	for _, f := range getDefs(taggy{}, false) {
		byName[f.name] = f
	}
	if _, ok := byName["A"]; !ok {
		t.Error("expected field A")
	}
	if byName["B"].size.Size != 0 {
		t.Error("expected size 0 for invalid size tag")
	}
	if byName["C"].size.Size != 10 || byName["C"].size.Decimal != 0 {
		t.Error("expected size 10/0 for size:10,key")
	}
	if !byName["D"].identity || !byName["D"].key || !byName["D"].unsigned {
		t.Error("expected identity/key/unsigned flags on D")
	}
	if byName["E"].fType != tDateTime {
		t.Error("expected type:time to produce tDateTime")
	}
	if _, ok := byName["EmbeddedField"]; !ok {
		t.Error("expected nested field from *covEmbed")
	}
	if _, ok := byName["covInt"]; ok {
		t.Error("anonymous non-struct field should be skipped")
	}

	// Non-nil pointer to struct follows the else branch.
	defs := getDefs(taggy{F: &covEmbed{}}, false)
	found := false
	for _, f := range defs {
		if f.name == "EmbeddedField" {
			found = true
		}
	}
	if !found {
		t.Error("expected nested defs for non-nil embedded pointer")
	}
}

func TestCoverageMisc(t *testing.T) {
	if got := GetTableName(covModel{}); got != "covModel" {
		t.Errorf("expected covModel, got %q", got)
	}

	db := covDB(t)
	// db.First error path via invalid criteria.
	if _, err := db.First(covModel{}, 42); err == nil {
		t.Error("expected db.First to fail on invalid criteria")
	}
	// Generic First error propagation.
	if _, err := First[covModel](db, 42); err == nil {
		t.Error("expected generic First to propagate error")
	}
	// Generic Range early exit.
	for range Range[covModel](db) {
		break
	}
}
