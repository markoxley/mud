// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.

package tests

import (
	"testing"

	"github.com/markoxley/dtorm"
)

func TestGetManager(t *testing.T) {
	tests := []struct {
		name       string
		config     *dtorm.Config
		wantType   string
		wantErr    bool
		errMessage string
	}{
		{
			name: "sqlite manager",
			config: &dtorm.Config{
				Type: "sqlite",
			},
			wantType: "*dtorm.SqliteManager",
			wantErr:  false,
		},
		{
			name: "sqlite3 manager",
			config: &dtorm.Config{
				Type: "sqlite3",
			},
			wantType: "*dtorm.SqliteManager",
			wantErr:  false,
		},
		{
			name: "mysql manager",
			config: &dtorm.Config{
				Type: "mysql",
			},
			wantType: "*dtorm.MySQLManager",
			wantErr:  false,
		},
		{
			name: "mariadb manager",
			config: &dtorm.Config{
				Type: "mariadb",
			},
			wantType: "*dtorm.MySQLManager",
			wantErr:  false,
		},
		{
			name: "sqlserver manager",
			config: &dtorm.Config{
				Type: "sqlserver",
			},
			wantType: "*dtorm.MSSQLManager",
			wantErr:  false,
		},
		{
			name: "mssql manager",
			config: &dtorm.Config{
				Type: "mssql",
			},
			wantType: "*dtorm.MSSQLManager",
			wantErr:  false,
		},
		{
			name: "invalid manager type",
			config: &dtorm.Config{
				Type: "invalid",
			},
			wantType:   "",
			wantErr:    true,
			errMessage: "invalid database type: invalid",
		},
		{
			name:       "nil config",
			config:     nil,
			wantType:   "",
			wantErr:    true,
			errMessage: "config cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dtorm.GetManager(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMessage {
				t.Errorf("GetManager() error message = %v, want %v", err.Error(), tt.errMessage)
				return
			}
			if !tt.wantErr {
				gotType := dtorm.GetTypeName(got)
				if gotType != tt.wantType {
					t.Errorf("GetManager() = %v, want %v", gotType, tt.wantType)
				}
			}
		})
	}
}
