// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
package mud

import (
	"strings"
	"sync"
	"testing"
	"time"
)

// ddlTestModel exercises every field type used in dialect DDL assertions.
type ddlTestModel struct {
	Model
	Name   string    `mud:"size:50"`
	Count  int       `mud:""`
	Sized  int       `mud:"size:10"`
	UCount uint      `mud:""`
	ULong  uint64    `mud:""`
	Price  float64   `mud:"type:decimal,size:10,2"`
	Active bool      `mud:""`
	When   time.Time `mud:""`
}

func newDDLTestDB(t *testing.T, dialect string) *DB {
	t.Helper()
	cfg := &Config{Type: dialect}
	mgr, err := GetManager(cfg)
	if err != nil {
		t.Fatalf("GetManager(%q): %v", dialect, err)
	}
	return &DB{cfg: cfg, mgr: mgr, mu: &sync.RWMutex{}, createMu: &sync.Mutex{}, tableDef: make(map[string][]field)}
}

func TestTableDefinitionDialects(t *testing.T) {
	cases := []struct {
		name     string
		dialect  string
		contains []string
		excludes []string
	}{
		{
			name:    "sqlite",
			dialect: "sqlite",
			contains: []string{
				`"ID" VARCHAR(36)`,
				`"Name" VARCHAR(50)`,
				`"Count" INTEGER`,
				`"UCount" INTEGER`,
				`"ULong" INTEGER`,
				`"Price" DECIMAL(10,2)`,
				`"Active" INTEGER`,
				`"When" DATETIME`,
				` PRIMARY KEY`,
			},
			excludes: []string{"UNSIGNED"},
		},
		{
			name:    "mysql",
			dialect: "mysql",
			contains: []string{
				"`ID` CHAR(36)",
				"`Name` VARCHAR(50)",
				"`Count` INT NOT NULL",
				"`Sized` INT(10)",
				"`UCount` INT UNSIGNED",
				"`ULong` BIGINT UNSIGNED",
				"`Price` DECIMAL(10,2)",
				"`Active` TINYINT(1)",
				"`When` DATETIME",
			},
		},
		{
			name:    "mssql",
			dialect: "mssql",
			contains: []string{
				"[ID] VARCHAR(36)",
				"[Name] NVARCHAR(50)",
				"[Count] INT NOT NULL",
				"[UCount] BIGINT",       // unsigned int widened
				"[ULong] DECIMAL(20,0)", // unsigned bigint widened
				"[Price] DECIMAL(10,2)",
				"[Active] BIT",
				"DATETIME2",
			},
			excludes: []string{"UNSIGNED", "DOUBLE", "INT("},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := newDDLTestDB(t, tc.dialect)
			sql, reqd := db.tableDefinition(&ddlTestModel{})
			if !reqd {
				t.Fatal("tableDefinition returned not required")
			}
			joined := strings.Join(sql, "\n")
			for _, want := range tc.contains {
				if !strings.Contains(joined, want) {
					t.Errorf("DDL missing %q:\n%s", want, joined)
				}
			}
			for _, bad := range tc.excludes {
				if strings.Contains(joined, bad) {
					t.Errorf("DDL unexpectedly contains %q:\n%s", bad, joined)
				}
			}
		})
	}
}
