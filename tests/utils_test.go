// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.

// MIT License
//
// Copyright (c) 2025 DaggerTech
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package tests

import (
	"testing"
	"time"

	"github.com/markoxley/dtorm/utils"
)

func TestMakeValue(t *testing.T) {
	// Define test time
	testTime := time.Date(2025, 4, 3, 15, 23, 0, 0, time.UTC)

	tests := []struct {
		name     string
		input    interface{}
		want     string
		wantBool bool
	}{
		{
			name:     "float32 value",
			input:    float32(123.44),
			want:     "123.4400",
			wantBool: true,
		},
		{
			name:     "float64 value",
			input:    float64(123.456789),
			want:     "123.456789",
			wantBool: true,
		},
		{
			name:     "integer value",
			input:    42,
			want:     "42",
			wantBool: true,
		},
		{
			name:     "int64 value",
			input:    int64(9223372036854775807),
			want:     "9223372036854775807",
			wantBool: true,
		},
		{
			name:     "boolean true",
			input:    true,
			want:     "1",
			wantBool: true,
		},
		{
			name:     "boolean false",
			input:    false,
			want:     "0",
			wantBool: true,
		},
		{
			name:     "simple string",
			input:    "hello",
			want:     "'hello'",
			wantBool: true,
		},
		{
			name:     "string with single quotes",
			input:    "O'Connor",
			want:     "'O''Connor'",
			wantBool: true,
		},
		{
			name:     "empty string",
			input:    "",
			want:     "''",
			wantBool: true,
		},
		{
			name:     "time value",
			input:    testTime,
			want:     "'2025-04-03 15:23:00.0'", // Assuming TimeToSQL formats like this
			wantBool: true,
		},
		{
			name:     "unsupported type (struct)",
			input:    struct{}{},
			want:     "",
			wantBool: false,
		},
		{
			name:     "nil value",
			input:    nil,
			want:     "",
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotBool := utils.MakeValue(tt.input)
			if got != tt.want {
				t.Errorf("MakeValue() got = %v, want %v", got, tt.want)
			}
			if gotBool != tt.wantBool {
				t.Errorf("MakeValue() gotBool = %v, want %v", gotBool, tt.wantBool)
			}
		})
	}
}

func TestPtr(t *testing.T) {
	t.Run("string pointer", func(t *testing.T) {
		value := "test"
		ptr := utils.Ptr(value)
		if ptr == nil {
			t.Error("Expected non-nil pointer")
		}
		if *ptr != value {
			t.Errorf("Expected %v, got %v", value, *ptr)
		}
	})

	t.Run("integer pointer", func(t *testing.T) {
		value := 42
		ptr := utils.Ptr(value)
		if ptr == nil {
			t.Error("Expected non-nil pointer")
		}
		if *ptr != value {
			t.Errorf("Expected %v, got %v", value, *ptr)
		}
	})

	t.Run("boolean pointer", func(t *testing.T) {
		value := true
		ptr := utils.Ptr(value)
		if ptr == nil {
			t.Error("Expected non-nil pointer")
		}
		if *ptr != value {
			t.Errorf("Expected %v, got %v", value, *ptr)
		}
	})

	t.Run("struct pointer", func(t *testing.T) {
		type testStruct struct {
			field string
		}
		value := testStruct{field: "test"}
		ptr := utils.Ptr(value)
		if ptr == nil {
			t.Error("Expected non-nil pointer")
		}
		if ptr.field != value.field {
			t.Errorf("Expected %v, got %v", value, *ptr)
		}
	})
}

func TestTimeToSQL(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "zero time",
			input:    time.Date(0, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: "0000-01-01 00:00:00.0",
		},
		{
			name:     "standard date time",
			input:    time.Date(2025, 4, 3, 15, 26, 27, 123456789, time.UTC),
			expected: "2025-04-03 15:26:27.123",
		},
		{
			name:     "date with single digit values",
			input:    time.Date(2025, 1, 5, 9, 5, 7, 123000000, time.UTC),
			expected: "2025-01-05 09:05:07.123",
		},
		{
			name:     "max year",
			input:    time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC),
			expected: "9999-12-31 23:59:59.999",
		},
		{
			name:     "nanoseconds truncation",
			input:    time.Date(2025, 6, 15, 12, 30, 45, 123456789, time.UTC),
			expected: "2025-06-15 12:30:45.123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.TimeToSQL(tt.input)
			if result != tt.expected {
				t.Errorf("TimeToSQL(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
