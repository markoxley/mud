// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.

package tests

import (
	"reflect"
	"testing"

	"github.com/markoxley/mud"
	"github.com/markoxley/mud/utils"
	"github.com/markoxley/mud/where"
)

// fluentCase describes one fluent builder method applied to a builder that
// already contains a base clause ("z" = ?). `frag` is the expected SQL for the
// appended clause and `args` are its expected bind arguments.
type fluentCase struct {
	name string
	fn   func(*where.Builder) *where.Builder
	frag string
	args []interface{}
}

// runFluent builds "z" = ? plus the fluent clause and asserts the full SQL
// and ordered arguments for a SQLite manager.
func runFluent(t *testing.T, c fluentCase, conj string) {
	t.Helper()
	mgr := &mud.SqliteManager{}
	b := c.fn(where.Equal("z", 0))
	got, args := b.Build(mgr.Operators(), mgr.Placeholder)
	want := phSeq(`"z" = %p`+conj+c.frag, mgr)
	if got != want {
		t.Errorf("%s: expected %q, got %q", c.name, want, got)
	}
	wantArgs := append(makeArgs(0), c.args...)
	if !reflect.DeepEqual(args, wantArgs) {
		t.Errorf("%s: expected args %v, got %v", c.name, wantArgs, args)
	}
}

func TestWhereFluentAnd(t *testing.T) {
	cases := []fluentCase{
		{"AndEqual", func(b *where.Builder) *where.Builder { return b.AndEqual("x", 1) }, `"x" = %p`, makeArgs(1)},
		{"AndGreater", func(b *where.Builder) *where.Builder { return b.AndGreater("x", 2) }, `"x" > %p`, makeArgs(2)},
		{"AndLess", func(b *where.Builder) *where.Builder { return b.AndLess("x", 3) }, `"x" < %p`, makeArgs(3)},
		{"AndLike", func(b *where.Builder) *where.Builder { return b.AndLike("x", "a") }, `"x" LIKE %p`, makeArgs("a")},
		{"AndStartsWith", func(b *where.Builder) *where.Builder { return b.AndStartsWith("x", "a") }, `"x" LIKE %p`, makeArgs("a%")},
		{"AndEndsWith", func(b *where.Builder) *where.Builder { return b.AndEndsWith("x", "a") }, `"x" LIKE %p`, makeArgs("%a")},
		{"AndContains", func(b *where.Builder) *where.Builder { return b.AndContains("x", "a") }, `"x" LIKE %p`, makeArgs("%a%")},
		{"AndIn", func(b *where.Builder) *where.Builder { return b.AndIn("x", []int{1, 2}) }, `"x" IN (%p,%p)`, makeArgs(1, 2)},
		{"AndBetween", func(b *where.Builder) *where.Builder { return b.AndBetween("x", 1, 9) }, `"x" BETWEEN %p AND %p`, makeArgs(1, 9)},
		{"AndIsNull", func(b *where.Builder) *where.Builder { return b.AndIsNull("x") }, `"x" IS NULL`, nil},
		{"AndNotIsNull", func(b *where.Builder) *where.Builder { return b.AndNotIsNull("x") }, `"x" IS NOT NULL`, nil},
		{"AndNotEqual", func(b *where.Builder) *where.Builder { return b.AndNotEqual("x", 1) }, `"x" <> %p`, makeArgs(1)},
		{"AndNotGreater", func(b *where.Builder) *where.Builder { return b.AndNotGreater("x", 2) }, `"x" <= %p`, makeArgs(2)},
		{"AndNotLess", func(b *where.Builder) *where.Builder { return b.AndNotLess("x", 3) }, `"x" >= %p`, makeArgs(3)},
		{"AndNotLike", func(b *where.Builder) *where.Builder { return b.AndNotLike("x", "a") }, `"x" NOT LIKE %p`, makeArgs("a")},
		{"AndNotStartsWith", func(b *where.Builder) *where.Builder { return b.AndNotStartsWith("x", "a") }, `"x" NOT LIKE %p`, makeArgs("a%")},
		{"AndNotEndsWith", func(b *where.Builder) *where.Builder { return b.AndNotEndsWith("x", "a") }, `"x" NOT LIKE %p`, makeArgs("%a")},
		{"AndNotContains", func(b *where.Builder) *where.Builder { return b.AndNotContains("x", "a") }, `"x" NOT LIKE %p`, makeArgs("%a%")},
		{"AndNotIn", func(b *where.Builder) *where.Builder { return b.AndNotIn("x", []int{1, 2}) }, `"x" NOT IN (%p,%p)`, makeArgs(1, 2)},
		{"AndNotBetween", func(b *where.Builder) *where.Builder { return b.AndNotBetween("x", 1, 9) }, `"x" NOT BETWEEN %p AND %p`, makeArgs(1, 9)},
		{"AndSub", func(b *where.Builder) *where.Builder { return b.AndSub(where.Equal("x", 4)) }, `("x" = %p)`, makeArgs(4)},
	}
	for _, c := range cases {
		runFluent(t, c, " AND ")
	}
}

func TestWhereFluentOr(t *testing.T) {
	cases := []fluentCase{
		{"OrEqual", func(b *where.Builder) *where.Builder { return b.OrEqual("x", 1) }, `"x" = %p`, makeArgs(1)},
		{"OrGreater", func(b *where.Builder) *where.Builder { return b.OrGreater("x", 2) }, `"x" > %p`, makeArgs(2)},
		{"OrLess", func(b *where.Builder) *where.Builder { return b.OrLess("x", 3) }, `"x" < %p`, makeArgs(3)},
		{"OrLike", func(b *where.Builder) *where.Builder { return b.OrLike("x", "a") }, `"x" LIKE %p`, makeArgs("a")},
		{"OrStartsWith", func(b *where.Builder) *where.Builder { return b.OrStartsWith("x", "a") }, `"x" LIKE %p`, makeArgs("a%")},
		{"OrEndsWith", func(b *where.Builder) *where.Builder { return b.OrEndsWith("x", "a") }, `"x" LIKE %p`, makeArgs("%a")},
		{"OrContains", func(b *where.Builder) *where.Builder { return b.OrContains("x", "a") }, `"x" LIKE %p`, makeArgs("%a%")},
		{"OrIn", func(b *where.Builder) *where.Builder { return b.OrIn("x", []int{1, 2}) }, `"x" IN (%p,%p)`, makeArgs(1, 2)},
		{"OrBetween", func(b *where.Builder) *where.Builder { return b.OrBetween("x", 1, 9) }, `"x" BETWEEN %p AND %p`, makeArgs(1, 9)},
		{"OrIsNull", func(b *where.Builder) *where.Builder { return b.OrIsNull("x") }, `"x" IS NULL`, nil},
		{"OrNotIsNull", func(b *where.Builder) *where.Builder { return b.OrNotIsNull("x") }, `"x" IS NOT NULL`, nil},
		{"OrNotEqual", func(b *where.Builder) *where.Builder { return b.OrNotEqual("x", 1) }, `"x" <> %p`, makeArgs(1)},
		{"OrNotGreater", func(b *where.Builder) *where.Builder { return b.OrNotGreater("x", 2) }, `"x" <= %p`, makeArgs(2)},
		{"OrNotLess", func(b *where.Builder) *where.Builder { return b.OrNotLess("x", 3) }, `"x" >= %p`, makeArgs(3)},
		{"OrNotLike", func(b *where.Builder) *where.Builder { return b.OrNotLike("x", "a") }, `"x" NOT LIKE %p`, makeArgs("a")},
		{"OrNotStartsWith", func(b *where.Builder) *where.Builder { return b.OrNotStartsWith("x", "a") }, `"x" NOT LIKE %p`, makeArgs("a%")},
		{"OrNotEndsWith", func(b *where.Builder) *where.Builder { return b.OrNotEndsWith("x", "a") }, `"x" NOT LIKE %p`, makeArgs("%a")},
		{"OrNotContains", func(b *where.Builder) *where.Builder { return b.OrNotContains("x", "a") }, `"x" NOT LIKE %p`, makeArgs("%a%")},
		{"OrNotIn", func(b *where.Builder) *where.Builder { return b.OrNotIn("x", []int{1, 2}) }, `"x" NOT IN (%p,%p)`, makeArgs(1, 2)},
		{"OrNotBetween", func(b *where.Builder) *where.Builder { return b.OrNotBetween("x", 1, 9) }, `"x" NOT BETWEEN %p AND %p`, makeArgs(1, 9)},
		{"OrSub", func(b *where.Builder) *where.Builder { return b.OrSub(where.Equal("x", 4)) }, `("x" = %p)`, makeArgs(4)},
	}
	for _, c := range cases {
		runFluent(t, c, " OR ")
	}
}

func TestWhereSubAndCount(t *testing.T) {
	mgr := &mud.SqliteManager{}

	s := where.Sub(where.Equal("a", 1))
	got, args := s.Build(mgr.Operators(), mgr.Placeholder)
	if got != `("a" = ?)` {
		t.Errorf("expected (\"a\" = ?), got %q", got)
	}
	if !reflect.DeepEqual(args, makeArgs(1)) {
		t.Errorf("expected args [1], got %v", args)
	}

	b := where.Equal("a", 1)
	b.AndEqual("b", 2).OrEqual("c", 3)
	if b.Count() != 3 {
		t.Errorf("expected Count()=3, got %d", b.Count())
	}
}

func TestWhereInPointers(t *testing.T) {
	mgr := &mud.SqliteManager{}

	cases := []struct {
		name  string
		build func() *where.Builder
		out   string
		args  []interface{}
	}{
		{
			"slice pointer",
			func() *where.Builder { return where.In("x", utils.Ptr([]int{1, 2})) },
			`"x" IN (?,?)`,
			makeArgs(1, 2),
		},
		{
			"scalar pointer",
			func() *where.Builder { return where.In("x", utils.Ptr(5)) },
			`"x" IN (?)`,
			makeArgs(5),
		},
		{
			"double pointer",
			func() *where.Builder { return where.In("x", utils.Ptr(utils.Ptr(7))) },
			`"x" IN (?)`,
			makeArgs(7),
		},
		{
			"nil pointer",
			func() *where.Builder {
				var p *int
				return where.In("x", p)
			},
			``,
			nil,
		},
		{
			"array",
			func() *where.Builder { return where.In("x", [2]int{8, 9}) },
			`"x" IN (?,?)`,
			makeArgs(8, 9),
		},
		{
			"empty slice",
			func() *where.Builder { return where.In("x", []int{}) },
			``,
			nil,
		},
	}
	for _, c := range cases {
		got, args := c.build().Build(mgr.Operators(), mgr.Placeholder)
		if got != c.out {
			t.Errorf("%s: expected %q, got %q", c.name, c.out, got)
		}
		if !reflect.DeepEqual(args, c.args) {
			t.Errorf("%s: expected args %v, got %v", c.name, c.args, args)
		}
	}
}

func TestWhereBetweenReorders(t *testing.T) {
	// BETWEEN arguments are normalised into ascending order using their SQL
	// string representation.
	mgr := &mud.SqliteManager{}
	got, args := where.Between("x", "z", "a").Build(mgr.Operators(), mgr.Placeholder)
	if got != `"x" BETWEEN ? AND ?` {
		t.Errorf("unexpected SQL %q", got)
	}
	if !reflect.DeepEqual(args, makeArgs("a", "z")) {
		t.Errorf("expected reordered args [a z], got %v", args)
	}
}
