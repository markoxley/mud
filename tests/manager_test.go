// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.

package tests

import (
	"testing"

	"github.com/markoxley/mud"
)

func TestGetManager(t *testing.T) {
	tests := []struct {
		name       string
		config     *mud.Config
		wantType   string
		wantErr    bool
		errMessage string
	}{
		{
			name: "sqlite manager",
			config: &mud.Config{
				Type: "sqlite",
			},
			wantType: "*mud.SqliteManager",
			wantErr:  false,
		},
		{
			name: "sqlite3 manager",
			config: &mud.Config{
				Type: "sqlite3",
			},
			wantType: "*mud.SqliteManager",
			wantErr:  false,
		},
		{
			name: "mysql manager",
			config: &mud.Config{
				Type: "mysql",
			},
			wantType: "*mud.MySQLManager",
			wantErr:  false,
		},
		{
			name: "mariadb manager",
			config: &mud.Config{
				Type: "mariadb",
			},
			wantType: "*mud.MySQLManager",
			wantErr:  false,
		},
		{
			name: "sqlserver manager",
			config: &mud.Config{
				Type: "sqlserver",
			},
			wantType: "*mud.MSSQLManager",
			wantErr:  false,
		},
		{
			name: "mssql manager",
			config: &mud.Config{
				Type: "mssql",
			},
			wantType: "*mud.MSSQLManager",
			wantErr:  false,
		},
		{
			name: "invalid manager type",
			config: &mud.Config{
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
			got, err := mud.GetManager(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMessage {
				t.Errorf("GetManager() error message = %v, want %v", err.Error(), tt.errMessage)
				return
			}
			if !tt.wantErr {
				gotType := mud.GetTypeName(got)
				if gotType != tt.wantType {
					t.Errorf("GetManager() = %v, want %v", gotType, tt.wantType)
				}
			}
		})
	}
}
