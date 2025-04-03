// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.

package tests

import (
	"testing"

	"github.com/markoxley/dtorm"
)

func TestNewSize(t *testing.T) {
	tests := []struct {
		name         string
		size         int
		decimal      int
		wantSize     int
		wantDecimal  int
	}{
		{
			name:         "zero size and decimal",
			size:         0,
			decimal:      0,
			wantSize:     0,
			wantDecimal:  0,
		},
		{
			name:         "positive size no decimal",
			size:         10,
			decimal:      0,
			wantSize:     10,
			wantDecimal:  0,
		},
		{
			name:         "positive size with decimal",
			size:         10,
			decimal:      2,
			wantSize:     10,
			wantDecimal:  2,
		},
		{
			name:         "large size with large decimal",
			size:         255,
			decimal:      10,
			wantSize:     255,
			wantDecimal:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dtorm.NewSize(tt.size, tt.decimal)
			if got.Size != tt.wantSize {
				t.Errorf("NewSize() size = %v, want %v", got.Size, tt.wantSize)
			}
			if got.Decimal != tt.wantDecimal {
				t.Errorf("NewSize() decimal = %v, want %v", got.Decimal, tt.wantDecimal)
			}
		})
	}
}

func TestFieldSizeString(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		decimal int
		want    string
	}{
		{
			name:    "zero size no decimal",
			size:    0,
			decimal: 0,
			want:    "0",
		},
		{
			name:    "positive size no decimal",
			size:    10,
			decimal: 0,
			want:    "10",
		},
		{
			name:    "positive size with decimal",
			size:    10,
			decimal: 2,
			want:    "10,2",
		},
		{
			name:    "large size with large decimal",
			size:    255,
			decimal: 10,
			want:    "255,10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := dtorm.NewSize(tt.size, tt.decimal)
			if got := fs.String(); got != tt.want {
				t.Errorf("FieldSize.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
