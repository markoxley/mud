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

	"github.com/markoxley/mud/order"
)

func TestNewOrder(t *testing.T) {
	tests := []struct {
		name      string
		field     string
		ascending bool
		want      string
	}{
		{
			name:      "ascending order",
			field:     "name",
			ascending: true,
			want:      "`name` asc",
		},
		{
			name:      "descending order",
			field:     "age",
			ascending: false,
			want:      "`age` desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := order.Asc(tt.field)
			if !tt.ascending {
				o = order.Desc(tt.field)
			}
			if got := o.String(); got != tt.want {
				t.Errorf("order.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuilder(t *testing.T) {
	tests := []struct {
		name string
		fn   func() *order.Builder
		want string
	}{
		{
			name: "single ascending",
			fn: func() *order.Builder {
				return order.Asc("name")
			},
			want: "`name` asc",
		},
		{
			name: "single descending",
			fn: func() *order.Builder {
				return order.Desc("age")
			},
			want: "`age` desc",
		},
		{
			name: "multiple mixed orders",
			fn: func() *order.Builder {
				return order.Asc("name").Desc("age").Asc("created_at")
			},
			want: "`name` asc, `age` desc, `created_at` asc",
		},
		{
			name: "chained descending orders",
			fn: func() *order.Builder {
				return order.Desc("score").Desc("timestamp")
			},
			want: "`score` desc, `timestamp` desc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.fn()
			if got := b.String(); got != tt.want {
				t.Errorf("Builder.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
