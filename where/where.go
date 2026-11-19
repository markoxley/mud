// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package clause This is a very basic, and yet versatile ORM package.
// At present, this is only for postgres
package where

import (
	"fmt"
)

// Builder is the main clause builder mechanism used for dagger
type Builder struct {
	conjunction conjunction
	children    []clauser
}

// NewBuilder
// @param c
// @return *Builder
func NewBuilder(c ...conjunction) *Builder {
	conj := conAnd
	if len(c) > 0 {
		conj = c[0]
	}
	return &Builder{
		conjunction: conj,
		children:    make([]clauser, 0, 5),
	}
}

func (c *Builder) Count() int {
	return len(c.children)
}

// Build returns the parameterized version of the clause.
// All values are emitted as placeholder markers via the ph function and
// returned separately in the args slice so they can be bound safely by the
// database driver instead of being interpolated into the SQL text.
//
// @receiver c
// @param operators List of operator format strings to use in the clause
// @param ph Function that returns the next placeholder marker
// @return string The SQL clause with placeholder markers
// @return []interface{} The ordered arguments to bind to the placeholders
func (c *Builder) Build(operators []string, ph func(*int) string) (string, []interface{}) {
	pos := 1
	return c.BuildAt(operators, ph, &pos)
}

// BuildAt is like Build but lets the caller seed and own the placeholder
// position counter. This is needed when the clause forms part of a larger
// statement that already consumed earlier positions (e.g. @pN in SQL Server).
//
// @receiver c
// @param operators List of operator format strings to use in the clause
// @param ph Function that returns the next placeholder marker
// @param pos Counter used by numbered placeholders
// @return string The SQL clause with placeholder markers
// @return []interface{} The ordered arguments to bind to the placeholders
func (c *Builder) BuildAt(operators []string, ph func(*int) string, pos *int) (string, []interface{}) {
	result := ""
	args := make([]interface{}, 0, len(c.children))
	for _, child := range c.children {
		if result != "" {
			result += string(child.getConjunction())
		}

		value, childArgs := child.build(operators, ph, pos)
		if _, ok := child.(*Builder); ok {
			value = fmt.Sprintf("(%s)", value)
		}
		if value == "" {
			return "", nil
		}
		result += value
		args = append(args, childArgs...)
	}
	return result, args
}

// build implements the clauser interface so a Builder can be nested inside
// another Builder as a sub-clause.
//
// @receiver c
// @param operators List of operator format strings to use in the clause
// @param ph Function that returns the next placeholder marker
// @param pos Counter used by numbered placeholders
// @return string The SQL clause with placeholder markers
// @return []interface{} The ordered arguments to bind to the placeholders
func (c *Builder) build(operators []string, ph func(*int) string, pos *int) (string, []interface{}) {
	return c.BuildAt(operators, ph, pos)
}

// getConjunction
// @receiver c
// @return conjunction
func (c *Builder) getConjunction() conjunction {
	return c.conjunction
}

// Sub creates a new Builder with the first clause being a sub clause
// @param c
// @return *Builder
func Sub(c *Builder) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, c)
	return n
}

// Equal creates a new Builder with the first clause being an equal clause
// @param f
// @param v
// @return *Builder
func Equal(field string, value interface{}) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opEqual, false, value))
	return n
}

// Greater creates a new Builder with the first clause being a greater than clause
// @param f
// @param v
// @return *Builder
func Greater(field string, value interface{}) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opGreater, false, value))
	return n
}

// Less creates a new Builder with the first clause being a less than clause
// @param f
// @param v
// @return *Builder
func Less(field string, value interface{}) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opLess, false, value))
	return n
}

// Like creates a new Builder with the first clause being a like clause
// @param f
// @param v
// @return *Builder
func Like(field string, value string) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opLike, false, value))
	return n
}

// StartsWith creates a new Builder with the first clause being a starts with clause
// @param f
// @param v
// @return *Builder
func StartsWith(field string, value string) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opLike, false, value+"%"))
	return n
}

// EndsWith creates a new Builder with the first clause being an ends with clause
// @param f
// @param v
// @return *Builder
func EndsWith(field string, value string) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opLike, false, "%"+value))
	return n
}

// Contains creates a new Builder with the first clause being a contains clause
// @param f
// @param v
// @return *Builder
func Contains(field string, value string) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opLike, false, "%"+value+"%"))
	return n
}

// In creates a new Builder with the first clause being an in clause
// @param f
// @param v
// @return *Builder
func In(field string, value interface{}) *Builder {
	values := convertToInterfaceArray(value)
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opIn, false, values...))
	return n
}

// Between creates a new Builder with the first clause being a between clause
// @param f
// @param v1
// @param v2
// @return *Builder
func Between(field string, value1 interface{}, value2 interface{}) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opBetween, false, value1, value2))
	return n
}

// IsNull creates a new Builder with the first clause being an is null clause
// @param f
// @return *Builder
func IsNull(field string) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opIsNull, false))
	return n
}

// IsNotNull creates a new Builder with the first clause being an is null clause
// @param f
// @return *Builder
func IsNotNull(field string) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opIsNull, true))
	return n
}

// NotEqual creates a new Builder with the first clause being a not equal clause
// @param f
// @param v
// @return *Builder
func NotEqual(field string, value interface{}) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opEqual, true, value))
	return n
}

// NotGreater creates a new Builder with the first clause being a not greater than clause
// @param f
// @param v
// @return *Builder
func NotGreater(field string, value interface{}) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opGreater, true, value))
	return n
}

// NotLess creates a new Builder with the first clause being a not less than clause
// @param f
// @param v
// @return *Builder
func NotLess(field string, value interface{}) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opLess, true, value))
	return n
}

// NotLike creates a new Builder with the first clause being a not like clause
// @param f
// @param v
// @return *Builder
func NotLike(field string, value string) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opLike, true, value))
	return n
}

// NotStartsWith creates a new Builder with the first clause being a not starts with clause
// @param f
// @param v
// @return *Builder
func NotStartsWith(field string, value string) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opLike, true, value+"%"))
	return n
}

// NotEndsWith creates a new Builder with the first clause being a not ends with clause
// @param f
// @param v
// @return *Builder
func NotEndsWith(field string, value string) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opLike, true, "%"+value))
	return n
}

// NotContains creates a new Builder with the first clause being a not contains clause
// @param f
// @param v
// @return *Builder
func NotContains(field string, value string) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opLike, true, "%"+value+"%"))
	return n
}

// NotIn creates a new Builder with the first clause being a not in clause
// @param f
// @param v
// @return *Builder
func NotIn(field string, value interface{}) *Builder {
	values := convertToInterfaceArray(value)
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opIn, true, values...))
	return n
}

// NotBetween creates a new Builder with the first clause being a not between clause
// @param f
// @param v1
// @param v2
// @return *Builder
func NotBetween(field string, value1 interface{}, value2 interface{}) *Builder {
	n := NewBuilder(conAnd)
	n.children = append(n.children, newClause(conAnd, field, opBetween, true, value1, value2))
	return n
}

// AndSub add and existing subclause to the clause with an AND conjunction
// @receiver c
// @param n
// @return *Builder
func (c *Builder) AndSub(n *Builder) *Builder {
	n.conjunction = conAnd
	c.children = append(c.children, n)
	return c
}

// AndEqual add an equal clause to the clause with an AND conjunction
// @receiver c
// @param f
// @param v
// @return *Builder
func (c *Builder) AndEqual(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opEqual, false, value))
	return c
}

// AndGreater add a greater than clause to the clause with an AND conjunction
// @receiver c
// @param f
// @param v
// @return *Builder
func (c *Builder) AndGreater(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opGreater, false, value))
	return c
}

// AndLess add a less than clause to the clause with an AND conjunction
func (c *Builder) AndLess(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opLess, false, value))
	return c
}

// AndLike add a like clause to the clause with an AND conjunction
func (c *Builder) AndLike(field string, value string) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opLike, false, value))
	return c
}

// AndStartsWith add a starts with clause to the clause with an AND conjunction
func (c *Builder) AndStartsWith(field string, value string) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opLike, false, value+"%"))
	return c
}

// AndEndsWith add a ends with clause to the clause with an AND conjunction
func (c *Builder) AndEndsWith(field string, value string) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opLike, false, "%"+value))
	return c
}

// AndContains add a contains clause to the clause with an AND conjunction
func (c *Builder) AndContains(field string, value string) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opLike, false, "%"+value+"%"))
	return c
}

// AndIn add an in clause to the clause with an AND conjunction
func (c *Builder) AndIn(field string, value interface{}) *Builder {
	values := convertToInterfaceArray(value)
	c.children = append(c.children, newClause(conAnd, field, opIn, false, values...))
	return c
}

// AndBetween add a between clause to the clause with an AND conjunction
func (c *Builder) AndBetween(field string, value1 interface{}, value2 interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opBetween, false, value1, value2))
	return c
}

// AndIsNull adds an is null clause to the clause with an AND conjunction
func (c *Builder) AndIsNull(field string) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opIsNull, false))
	return c
}

// AndNotIsNull adds a not is null clause to the clause with an AND conjunction
func (c *Builder) AndNotIsNull(field string) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opIsNull, true))
	return c
}

// AndNotEqual add a not equal clause to the clause with an AND conjunction
func (c *Builder) AndNotEqual(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opEqual, true, value))
	return c
}

// AndNotGreater add a not greater than clause to the clause with an AND conjunction
func (c *Builder) AndNotGreater(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opGreater, true, value))
	return c
}

// AndNotLess add a not less than clause to the clause with an AND conjunction
func (c *Builder) AndNotLess(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opLess, true, value))
	return c
}

// AndNotLike add a not like clause to the clause with an AND conjunction
func (c *Builder) AndNotLike(field string, value string) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opLike, true, value))
	return c
}

// AndNotStartsWith add a starts with clause to the clause with an AND conjunction
func (c *Builder) AndNotStartsWith(field string, value string) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opLike, true, value+"%"))
	return c
}

// AndNotEndsWith add a ends with clause to the clause with an AND conjunction
func (c *Builder) AndNotEndsWith(field string, value string) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opLike, true, "%"+value))
	return c
}

// AndNotContains add a contains clause to the clause with an AND conjunction
func (c *Builder) AndNotContains(field string, value string) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opLike, true, "%"+value+"%"))
	return c
}

// AndNotIn add a not in clause to the clause with an AND conjunction
func (c *Builder) AndNotIn(field string, value interface{}) *Builder {
	values := convertToInterfaceArray(value)
	c.children = append(c.children, newClause(conAnd, field, opIn, true, values...))
	return c
}

// AndNotBetween add a not between clause to the clause with an AND conjunction
func (c *Builder) AndNotBetween(field string, value1 interface{}, value2 interface{}) *Builder {
	c.children = append(c.children, newClause(conAnd, field, opBetween, true, value1, value2))
	return c
}

// OrSub add and existing subclause to the clause with an OR conjunction
func (c *Builder) OrSub(n *Builder) *Builder {
	n.conjunction = conOr
	c.children = append(c.children, n)
	return c
}

// OrEqual add an equal clause to the clause with an OR conjunction
func (c *Builder) OrEqual(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, field, opEqual, false, value))
	return c
}

// OrGreater add a greater than clause to the clause with an OR conjunction
func (c *Builder) OrGreater(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, field, opGreater, false, value))
	return c
}

// OrLess add a less than clause to the clause with an OR conjunction
func (c *Builder) OrLess(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, field, opLess, false, value))
	return c
}

// OrLike add a like clause to the clause with an OR conjunction
func (c *Builder) OrLike(field string, value string) *Builder {
	c.children = append(c.children, newClause(conOr, field, opLike, false, value))
	return c
}

// OrStartsWith add a starts with clause to the clause with an OR conjunction
func (c *Builder) OrStartsWith(field string, value string) *Builder {
	c.children = append(c.children, newClause(conOr, field, opLike, false, value+"%"))
	return c
}

// OrEndsWith add a ends with clause to the clause with an OR conjunction
func (c *Builder) OrEndsWith(field string, value string) *Builder {
	c.children = append(c.children, newClause(conOr, field, opLike, false, "%"+value))
	return c
}

// OrContains add a contains clause to the clause with an OR conjunction
func (c *Builder) OrContains(field string, value string) *Builder {
	c.children = append(c.children, newClause(conOr, field, opLike, false, "%"+value+"%"))
	return c
}

// OrIn add an in clause to the clause with an OR conjunction
func (c *Builder) OrIn(field string, value interface{}) *Builder {
	values := convertToInterfaceArray(value)
	c.children = append(c.children, newClause(conOr, field, opIn, false, values...))
	return c
}

// OrBetween add a between clause to the clause with an OR conjunction
func (c *Builder) OrBetween(field string, value1 interface{}, value2 interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, field, opBetween, false, value1, value2))
	return c
}

// OrIsNull adds an is null clause to the clause with an OR conjunction
func (c *Builder) OrIsNull(field string) *Builder {
	c.children = append(c.children, newClause(conOr, field, opIsNull, false))
	return c
}

// OrNotIsNull adds a not is null clause to the clause with an OR conjunction
func (c *Builder) OrNotIsNull(field string) *Builder {
	c.children = append(c.children, newClause(conOr, field, opIsNull, true))
	return c
}

// OrNotEqual add a not equal clause to the clause with an OR conjunction
func (c *Builder) OrNotEqual(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, field, opEqual, true, value))
	return c
}

// OrNotGreater add a not greater than clause to the clause with an OR conjunction
func (c *Builder) OrNotGreater(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, field, opGreater, true, value))
	return c
}

// OrNotLess add a not less than clause to the clause with an OR conjunction
func (c *Builder) OrNotLess(field string, value interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, field, opLess, true, value))
	return c
}

// OrNotLike add a not like clause to the clause with an OR conjunction
func (c *Builder) OrNotLike(field string, value string) *Builder {
	c.children = append(c.children, newClause(conOr, field, opLike, true, value))
	return c
}

// OrNotStartsWith add a not starts with clause to the clause with an OR conjunction
func (c *Builder) OrNotStartsWith(field string, value string) *Builder {
	c.children = append(c.children, newClause(conOr, field, opLike, true, value+"%"))
	return c
}

// OrNotEndsWith add a not ends with clause to the clause with an OR conjunction
func (c *Builder) OrNotEndsWith(field string, value string) *Builder {
	c.children = append(c.children, newClause(conOr, field, opLike, true, "%"+value))
	return c
}

// OrNotContains add a not contains clause to the clause with an OR conjunction
func (c *Builder) OrNotContains(field string, value string) *Builder {
	c.children = append(c.children, newClause(conOr, field, opLike, true, "%"+string(value)+"%"))
	return c
}

// OrNotIn add a not in clause to the clause with an OR conjunction
func (c *Builder) OrNotIn(field string, value interface{}) *Builder {
	values := convertToInterfaceArray(value)
	c.children = append(c.children, newClause(conOr, field, opIn, true, values...))
	return c
}

// OrNotBetween add a not between clause to the clause with an OR conjunction
func (c *Builder) OrNotBetween(field string, value1 interface{}, value2 interface{}) *Builder {
	c.children = append(c.children, newClause(conOr, field, opBetween, true, value1, value2))
	return c
}
