// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.

package tests

import (
	"reflect"
	"testing"
	"time"

	"github.com/markoxley/mud/utils"
)

func TestMakeArg(t *testing.T) {
	tm := time.Date(2025, 4, 3, 15, 23, 0, 0, time.UTC)
	tests := []struct {
		name   string
		input  interface{}
		want   interface{}
		wantOK bool
	}{
		{"float32", float32(1.5), float64(1.5), true},
		{"float64", float64(2.5), float64(2.5), true},
		{"int", int(42), int(42), true},
		{"int8", int8(8), int8(8), true},
		{"int16", int16(16), int16(16), true},
		{"int32", int32(32), int32(32), true},
		{"int64", int64(64), int64(64), true},
		{"uint", uint(1), uint(1), true},
		{"uint8", uint8(2), uint8(2), true},
		{"uint16", uint16(3), uint16(3), true},
		{"uint32", uint32(4), uint32(4), true},
		{"uint64", uint64(5), uint64(5), true},
		{"bool true", true, 1, true},
		{"bool false", false, 0, true},
		{"string", "hello", "hello", true},
		{"time", tm, utils.TimeToSQL(tm), true},
		{"nil", nil, nil, false},
		{"unsupported slice", []byte{1}, nil, false},
		{"unsupported struct", struct{ X int }{1}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := utils.MakeArg(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("MakeArg(%v) ok = %v, want %v", tt.input, ok, tt.wantOK)
			}
			if ok && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeArg(%v) = %v (%T), want %v (%T)", tt.input, got, got, tt.want, tt.want)
			}
		})
	}
}

// Named types used to exercise the kind-based fallback in MakeArg/MakeValue.
type namedInt int
type namedUint8 uint8
type namedBool bool
type namedString string
type namedFloat64 float64

func TestMakeArgNamedTypes(t *testing.T) {
	tests := []struct {
		name   string
		input  interface{}
		want   interface{}
		wantOK bool
	}{
		{"named int", namedInt(7), int64(7), true},
		{"named uint8", namedUint8(2), int64(2), true},
		{"named bool true", namedBool(true), 1, true},
		{"named bool false", namedBool(false), 0, true},
		{"named string", namedString("hi"), "hi", true},
		{"named float64", namedFloat64(1.5), float64(1.5), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := utils.MakeArg(tt.input)
			if ok != tt.wantOK {
				t.Fatalf("MakeArg(%v) ok = %v, want %v", tt.input, ok, tt.wantOK)
			}
			if ok && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeArg(%v) = %v (%T), want %v (%T)", tt.input, got, got, tt.want, tt.want)
			}
		})
	}
}

func TestSQLToTimeEdges(t *testing.T) {
	valid := []struct {
		name  string
		input string
		want  time.Time
	}{
		{"T separator", "2024-05-06T07:08:09", time.Date(2024, 5, 6, 7, 8, 9, 0, time.UTC)},
		{"date only", "2024-05-06", time.Date(2024, 5, 6, 0, 0, 0, 0, time.UTC)},
		{"with nanos", "2024-05-06 07:08:09.123", time.Date(2024, 5, 6, 7, 8, 9, 123, time.UTC)},
		{"trailing junk in seconds", "2024-05-06 07:08:09.1.2", time.Date(2024, 5, 6, 7, 8, 9, 1, time.UTC)},
	}
	for _, tt := range valid {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := utils.SQLToTime(tt.input)
			if !ok {
				t.Fatalf("SQLToTime(%q) failed", tt.input)
			}
			if !got.Equal(tt.want) {
				t.Errorf("SQLToTime(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}

	invalid := []string{
		"",                    // no date parts
		"2024/05/06",          // wrong date separator
		"abc-05-06",           // bad year
		"2024-xx-06",          // bad month
		"2024-05-yy",          // bad day
		"2024-05-06 07:08",    // incomplete time
		"2024-05-06 aa:08:09", // bad hour
		"2024-05-06 07:bb:09", // bad minute
		"2024-05-06 07:08:cc", // bad second
	}
	for _, in := range invalid {
		t.Run("invalid_"+in, func(t *testing.T) {
			if got, ok := utils.SQLToTime(in); ok {
				t.Errorf("SQLToTime(%q) = %v, expected failure", in, got)
			}
		})
	}
}
