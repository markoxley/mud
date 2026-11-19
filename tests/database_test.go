// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
package tests

import (
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/markoxley/mud"
	"github.com/markoxley/mud/where"
	"github.com/stretchr/testify/assert"
)

func getConfig(dbType string) *mud.Config {
	switch strings.ToLower(dbType) {
	case "sqlite", "sqlite3":
		return &mud.Config{
			Type:     "sqlite",
			Host:     "localhost:1433",
			User:     "sa",
			Password: "Dantooine2020!",
			Database: "mud_test.db",
		}
	case "mssql", "sqlserver":
		return &mud.Config{
			Type:     "mssql",
			Host:     "localhost:1433",
			User:     "sa",
			Password: "Dantooine2020!",
			Database: "mud_test",
		}
	case "mysql":
		return &mud.Config{
			Type:     "mysql",
			Host:     "localhost:3306",
			User:     "root",
			Password: "Dantooine2020!",
			Database: "mud_test",
		}
	}
	panic("Invalid database type")
}

func getDB(t *testing.T, dbType string) *mud.DB {
	t.Helper()
	config := getConfig(dbType)
	db, err := mud.New(config)
	if err != nil {
		t.Skipf("database %s unavailable: %v", dbType, err)
	}
	if err := db.SqlDB().Ping(); err != nil {
		db.Close()
		t.Skipf("database %s unreachable: %v", dbType, err)
	}
	db.RawExecute("Delete from test_models", nil)
	return db
}

// TestModel is a test model for database operations
type TestModel struct {
	mud.Model
	Name string `mud:"size:255"`
	Age  int    `mud:""`
}

func (t *TestModel) GetCreated() time.Time  { return t.CreateDate }
func (t *TestModel) GetUpdated() time.Time  { return t.LastUpdate }
func (t *TestModel) GetDeleted() *time.Time { return t.DeleteDate }
func (t *TestModel) TableName() string      { return "test_models" }

// ConcurrentModel is a separate model type used by the concurrency test so
// that its table definition is cold at the start of the test, forcing
// concurrent goroutines to race on the shared table-definition cache.
type ConcurrentModel struct {
	mud.Model
	Name string `mud:"size:255"`
	N    int    `mud:""`
}

func (t *ConcurrentModel) TableName() string { return "concurrent_models" }

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *mud.Config
		wantErr bool
	}{
		{
			name:    "Valid MySQL Config",
			config:  getConfig("mysql"),
			wantErr: false,
		},
		{
			name:    "Valid SQL Server Config",
			config:  getConfig("mssql"),
			wantErr: false,
		},
		{
			name:    "Valid SQLite Config",
			config:  getConfig("sqlite"),
			wantErr: false,
		},
		{
			name: "Invalid DB Type",
			config: &mud.Config{
				Type:     "invalid",
				Database: "test.db",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := mud.New(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
				db.Close()
			}
		})
	}
}

func TestSaveAndFetchSQLite(t *testing.T) {
	db := getDB(t, "sqlite")
	defer db.Close()
	// Create test model
	model := &TestModel{
		Name: "Test User",
		Age:  25,
	}

	// Test Save
	err := db.Save(model)
	assert.Nil(t, err)
	assert.NotEmpty(t, model.ID)
	assert.False(t, model.CreateDate.IsZero())
	assert.False(t, model.LastUpdate.IsZero())
	assert.Nil(t, model.DeleteDate)

	// Test Fetch
	//fetchedModels := make([]TestModel, 0)
	fetchedModels, err := mud.Fetch[TestModel](db, where.Equal("ID", *model.ID))
	assert.NoError(t, err)
	assert.Len(t, fetchedModels, 1)

	fetchedModel := fetchedModels[0]
	assert.Equal(t, model.ID, fetchedModel.ID)
	assert.Equal(t, model.Name, fetchedModel.Name)
	assert.Equal(t, model.Age, fetchedModel.Age)
}

func TestSaveAndFetchMSSQL(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL test on Linux")
	}
	db := getDB(t, "mssql")
	defer db.Close()
	// Create test model
	model := &TestModel{
		Name: "Test User",
		Age:  25,
	}

	// Test Save
	err := db.Save(model)
	assert.Nil(t, err)
	assert.NotEmpty(t, model.ID)
	assert.False(t, model.CreateDate.IsZero())
	assert.False(t, model.LastUpdate.IsZero())
	assert.Nil(t, model.DeleteDate)

	// Test Fetch
	//fetchedModels := make([]TestModel, 0)
	fetchedModels, err := mud.Fetch[TestModel](db, where.Equal("ID", *model.ID))
	assert.NoError(t, err)
	assert.Len(t, fetchedModels, 1)

	fetchedModel := fetchedModels[0]
	assert.Equal(t, model.ID, fetchedModel.ID)
	assert.Equal(t, model.Name, fetchedModel.Name)
	assert.Equal(t, model.Age, fetchedModel.Age)
}
func TestSaveAndFetchMySQL(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping MySQL test on Windows")
	}
	db := getDB(t, "mysql")
	defer db.Close()
	// Create test model
	model := &TestModel{
		Name: "Test User",
		Age:  25,
	}

	// Test Save
	err := db.Save(model)
	assert.Nil(t, err)
	assert.NotEmpty(t, model.ID)
	assert.False(t, model.CreateDate.IsZero())
	assert.False(t, model.LastUpdate.IsZero())
	assert.Nil(t, model.DeleteDate)

	// Test Fetch
	//fetchedModels := make([]TestModel, 0)
	fetchedModels, err := mud.Fetch[TestModel](db, where.Equal("ID", *model.ID))
	assert.NoError(t, err)
	assert.Len(t, fetchedModels, 1)

	fetchedModel := fetchedModels[0]
	assert.Equal(t, model.ID, fetchedModel.ID)
	assert.Equal(t, model.Name, fetchedModel.Name)
	assert.Equal(t, model.Age, fetchedModel.Age)
}
func TestRemoveSQLite(t *testing.T) {
	db := getDB(t, "sqlite")
	defer db.Close()

	// Create and save test model
	model := &TestModel{
		Name: "To Delete",
		Age:  30,
	}
	err := db.Save(model)
	assert.Nil(t, err)
	id := *model.ID
	// Test Remove
	err = db.Remove(model)
	assert.Nil(t, err)

	//Verify model is not fetchable
	fetched, err := mud.Fetch[TestModel](db, where.Equal("ID", id))
	assert.NoError(t, err)
	assert.Len(t, fetched, 0)
	db.Close()
}

func TestRemoveMSSQL(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL test on Linux")
	}
	db := getDB(t, "mssql")
	defer db.Close()

	// Create and save test model
	model := &TestModel{
		Name: "To Delete",
		Age:  30,
	}
	err := db.Save(model)
	assert.Nil(t, err)
	id := *model.ID
	// Test Remove
	err = db.Remove(model)
	assert.Nil(t, err)

	//Verify model is not fetchable
	fetched, err := mud.Fetch[TestModel](db, where.Equal("ID", id))
	assert.NoError(t, err)
	assert.Len(t, fetched, 0)
	db.Close()
}
func TestRemoveMySQL(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping MySQL test on Windows")
	}
	db := getDB(t, "mysql")
	defer db.Close()

	// Create and save test model
	model := &TestModel{
		Name: "To Delete",
		Age:  30,
	}
	err := db.Save(model)
	assert.Nil(t, err)
	id := *model.ID
	// Test Remove
	err = db.Remove(model)
	assert.Nil(t, err)

	//Verify model is not fetchable
	fetched, err := mud.Fetch[TestModel](db, where.Equal("ID", id))
	assert.NoError(t, err)
	assert.Len(t, fetched, 0)
	db.Close()
}
func TestCountSQLite(t *testing.T) {
	db := getDB(t, "sqlite")
	defer db.Close()

	// Create test models
	models := []*TestModel{
		{Name: "User 1", Age: 20},
		{Name: "User 2", Age: 25},
		{Name: "User 3", Age: 25},
	}

	for _, m := range models {
		err := db.Save(m)
		assert.Nil(t, err)
	}

	// Test Count
	count := db.Count(&TestModel{})
	assert.Equal(t, 3, count)

	// Test Count with criteria
	count = db.Count(&TestModel{}, where.Equal("Age", 25))
	assert.Equal(t, 2, count)

	// Test Count
	count = db.Count(&TestModel{})
	assert.Equal(t, 3, count)

	// Test Count with criteria
	count = db.Count(&TestModel{}, where.Equal("Age", 25))
	assert.Equal(t, 2, count)
}

func TestCountMSSQL(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping MSSQL test on Linux")
	}
	db := getDB(t, "mssql")
	defer db.Close()

	// Create test models
	models := []*TestModel{
		{Name: "User 1", Age: 20},
		{Name: "User 2", Age: 25},
		{Name: "User 3", Age: 25},
	}

	for _, m := range models {
		err := db.Save(m)
		assert.Nil(t, err)
	}

	// Test Count
	count := db.Count(&TestModel{})
	assert.Equal(t, 3, count)

	// Test Count with criteria
	count = db.Count(&TestModel{}, where.Equal("Age", 25))
	assert.Equal(t, 2, count)

	// Test Count
	count = db.Count(&TestModel{})
	assert.Equal(t, 3, count)

	// Test Count with criteria
	count = db.Count(&TestModel{}, where.Equal("Age", 25))
	assert.Equal(t, 2, count)
}

func TestCountMySQL(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping MySQL test on Windows")
	}
	db := getDB(t, "mysql")
	defer db.Close()

	// Create test models
	models := []*TestModel{
		{Name: "User 1", Age: 20},
		{Name: "User 2", Age: 25},
		{Name: "User 3", Age: 25},
	}

	for _, m := range models {
		err := db.Save(m)
		assert.Nil(t, err)
	}

	// Test Count
	count := db.Count(&TestModel{})
	assert.Equal(t, 3, count)

	// Test Count with criteria
	count = db.Count(&TestModel{}, where.Equal("Age", 25))
	assert.Equal(t, 2, count)

	// Test Count
	count = db.Count(&TestModel{})
	assert.Equal(t, 3, count)

	// Test Count with criteria
	count = db.Count(&TestModel{}, where.Equal("Age", 25))
	assert.Equal(t, 2, count)
}

func TestFirstSQLite(t *testing.T) {
	db := getDB(t, "sqlite")
	defer db.Close()

	// Create test models
	models := []*TestModel{
		{Name: "First User", Age: 52},
		{Name: "Second User", Age: 53},
	}
	for _, m := range models {
		err := db.Save(m)
		assert.Nil(t, err)
	}

	// Test First
	result, err := mud.First[TestModel](db, where.Equal("Age", 53))
	assert.NoError(t, err)
	assert.Equal(t, "Second User", result.Name)
	assert.Equal(t, 53, result.Age)

	// Test First with non-existent criteria
	result, err = mud.First[TestModel](db, where.Equal("Age", 999))
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

// TestConcurrentFetchSaveSQLite exercises the race between Fetch and Save.
// Both operations touch the shared table-definition cache and known-tables
// slice on the *mud.DB. Without synchronization this triggers concurrent
// map read/write (which panics) and concurrent slice appends. Run with
// -race to verify the fix.
func TestConcurrentFetchSaveSQLite(t *testing.T) {
	cfg := getConfig("sqlite")
	// Add a busy timeout so concurrent SQLite writes wait instead of failing
	// with SQLITE_BUSY. This isolates the test from DB-level locking so it
	// exercises the library's in-memory shared state under real concurrency.
	cfg.Database = "mud_test.db?_pragma=busy_timeout(5000)"
	db, err := mud.New(cfg)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	// Serialize the connection pool. This avoids SQLITE_BUSY write-write
	// deadlocks (orthogonal to the library race) while still allowing the
	// Go-level table-definition cache access to run concurrently across
	// goroutines, which is what this test is verifying.
	db.Connection().SetMaxOpenConns(1)
	db.RawExecute("DROP TABLE IF EXISTS concurrent_models", nil)

	// Use a cold model type so its table definition is not yet cached. All
	// goroutines start simultaneously via the barrier, forcing them to race
	// on the shared tableDef map and knownTables slice.
	var wg sync.WaitGroup
	const goroutines = 8
	const iterations = 50
	barrier := make(chan struct{})

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-barrier
			for j := 0; j < iterations; j++ {
				m := &ConcurrentModel{Name: "Concurrent", N: id*100 + j}
				if err := db.Save(m); err != nil {
					t.Errorf("save: %v", err)
					return
				}
				if _, err := mud.Fetch[ConcurrentModel](db, where.Equal("ID", *m.ID)); err != nil {
					t.Errorf("fetch: %v", err)
					return
				}
			}
		}(i)
	}
	close(barrier)
	wg.Wait()
}

func TestFirstMSSQL(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping MSSQL test on Windows")
	}
	db := getDB(t, "mssql")
	defer db.Close()

	// Create test models
	models := []*TestModel{
		{Name: "First User", Age: 52},
		{Name: "Second User", Age: 53},
	}
	for _, m := range models {
		err := db.Save(m)
		assert.Nil(t, err)
	}

	// Test First
	result, err := mud.First[TestModel](db, where.Equal("Age", 53))
	assert.NoError(t, err)
	assert.Equal(t, "Second User", result.Name)
	assert.Equal(t, 53, result.Age)

	// Test First with non-existent criteria
	result, err = mud.First[TestModel](db, where.Equal("Age", 999))
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestFirstMySQL(t *testing.T) {

	db := getDB(t, "mysql")
	defer db.Close()

	// Create test models
	models := []*TestModel{
		{Name: "First User", Age: 52},
		{Name: "Second User", Age: 53},
	}
	for _, m := range models {
		err := db.Save(m)
		assert.Nil(t, err)
	}

	// Test First
	result, err := mud.First[TestModel](db, where.Equal("Age", 53))
	assert.NoError(t, err)
	assert.Equal(t, "Second User", result.Name)
	assert.Equal(t, 53, result.Age)

	// Test First with non-existent criteria
	result, err = mud.First[TestModel](db, where.Equal("Age", 999))
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
