// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/markoxley/dtorm"
	"github.com/markoxley/dtorm/where"
	"github.com/stretchr/testify/assert"
)

func getConfig(dbType string) *dtorm.Config {
	switch strings.ToLower(dbType) {
	case "sqlite", "sqlite3":
		return &dtorm.Config{
			Type:     "sqlite",
			Host:     "localhost:1433",
			User:     "sa",
			Password: "Dantooine2020!",
			Database: "dtorm_test.db",
		}
	case "mssql", "sqlserver":
		return &dtorm.Config{
			Type:     "mssql",
			Host:     "localhost:1433",
			User:     "sa",
			Password: "Dantooine2020!",
			Database: "dtorm_test",
		}
	case "mysql":
		return &dtorm.Config{
			Type:     "mysql",
			Host:     "localhost:3306",
			User:     "root",
			Password: "Dantooine2020!",
			Database: "dtorm_test",
		}
	}
	panic("Invalid database type")
}

func getDB(dbType string) *dtorm.DB {
	config := getConfig(dbType)
	db, _ := dtorm.New(config)
	db.RawExecute("Delete from TestModel")
	return db
}

// TestModel is a test model for database operations
type TestModel struct {
	dtorm.Model
	Name string `dtorm:"size:255"`
	Age  int    `dtorm:""`
}

func (t *TestModel) GetCreated() time.Time  { return t.CreateDate }
func (t *TestModel) GetUpdated() time.Time  { return t.LastUpdate }
func (t *TestModel) GetDeleted() *time.Time { return t.DeleteDate }
func (t *TestModel) TableName() string      { return "test_models" }

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  *dtorm.Config
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
			config: &dtorm.Config{
				Type:     "invalid",
				Database: "test.db",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := dtorm.New(tt.config)
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
	db := getDB("sqlite")
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
	fetchedModels, err := dtorm.Fetch[TestModel](db, where.Equal("ID", *model.ID))
	assert.NoError(t, err)
	assert.Len(t, fetchedModels, 1)

	fetchedModel := fetchedModels[0]
	assert.Equal(t, model.ID, fetchedModel.ID)
	assert.Equal(t, model.Name, fetchedModel.Name)
	assert.Equal(t, model.Age, fetchedModel.Age)
}

func TestSaveAndFetchMSSQL(t *testing.T) {
	db := getDB("mssql")
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
	fetchedModels, err := dtorm.Fetch[TestModel](db, where.Equal("ID", *model.ID))
	assert.NoError(t, err)
	assert.Len(t, fetchedModels, 1)

	fetchedModel := fetchedModels[0]
	assert.Equal(t, model.ID, fetchedModel.ID)
	assert.Equal(t, model.Name, fetchedModel.Name)
	assert.Equal(t, model.Age, fetchedModel.Age)
}
func TestSaveAndFetchMySQL(t *testing.T) {
	db := getDB("mysql")
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
	fetchedModels, err := dtorm.Fetch[TestModel](db, where.Equal("ID", *model.ID))
	assert.NoError(t, err)
	assert.Len(t, fetchedModels, 1)

	fetchedModel := fetchedModels[0]
	assert.Equal(t, model.ID, fetchedModel.ID)
	assert.Equal(t, model.Name, fetchedModel.Name)
	assert.Equal(t, model.Age, fetchedModel.Age)
}
func TestRemoveSQLite(t *testing.T) {
	db := getDB("sqlite")
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
	fetched, err := dtorm.Fetch[TestModel](db, where.Equal("ID", id))
	assert.NoError(t, err)
	assert.Len(t, fetched, 0)
	db.Close()
}

func TestRemoveMSSQL(t *testing.T) {
	db := getDB("mssql")
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
	fetched, err := dtorm.Fetch[TestModel](db, where.Equal("ID", id))
	assert.NoError(t, err)
	assert.Len(t, fetched, 0)
	db.Close()
}
func TestRemoveMySQL(t *testing.T) {
	db := getDB("mysql")
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
	fetched, err := dtorm.Fetch[TestModel](db, where.Equal("ID", id))
	assert.NoError(t, err)
	assert.Len(t, fetched, 0)
	db.Close()
}
func TestCountSQLite(t *testing.T) {
	db := getDB("sqlite")
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
	db := getDB("mssql")
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
	db := getDB("sqlite")
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
	db := getDB("sqlite")
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
	result, err := dtorm.First[TestModel](db, where.Equal("Age", 53))
	assert.NoError(t, err)
	assert.Equal(t, "Second User", result.Name)
	assert.Equal(t, 53, result.Age)

	// Test First with non-existent criteria
	result, err = dtorm.First[TestModel](db, where.Equal("Age", 999))
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestFirstMSSQL(t *testing.T) {
	db := getDB("mssql")
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
	result, err := dtorm.First[TestModel](db, where.Equal("Age", 53))
	assert.NoError(t, err)
	assert.Equal(t, "Second User", result.Name)
	assert.Equal(t, 53, result.Age)

	// Test First with non-existent criteria
	result, err = dtorm.First[TestModel](db, where.Equal("Age", 999))
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
func TestFirstMySQL(t *testing.T) {
	db := getDB("mysql")
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
	result, err := dtorm.First[TestModel](db, where.Equal("Age", 53))
	assert.NoError(t, err)
	assert.Equal(t, "Second User", result.Name)
	assert.Equal(t, 53, result.Age)

	// Test First with non-existent criteria
	result, err = dtorm.First[TestModel](db, where.Equal("Age", 999))
	assert.NotNil(t, err)
	assert.Nil(t, result)
}
