// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package where provides functionality for building SQL WHERE clauses
package where

import (
	"fmt"
	"strings"

	"github.com/markoxley/mud/utils"
)

// clause represents a single SQL WHERE clause condition
// It contains the operator, values, and whether the condition is negated
// Implements the clauser interface for generating SQL clause strings
type clause struct {
	conjunction conjunction
	field       string
	not         bool
	op          operator
	values      []interface{}
}

// build generates the SQL clause string representation using placeholder
// markers instead of literal values. This method implements the clauser
// interface.
//
// @receiver c The clause instance
// @param operators List of operator strings to use for generating the clause
// @param ph Function that returns the next placeholder marker
// @param pos Counter used by numbered placeholders (e.g. @pN)
// @return SQL clause string (or empty string if invalid) and its bind arguments
func (c clause) build(operators []string, ph func(*int) string, pos *int) (string, []interface{}) {
	// Calculate operator index, adjusting for NOT conditions
	opCode := int(c.op)
	if c.not {
		opCode += len(operators) / 2
	}

	// Determine number of fields based on operator type
	fieldCount := 1
	switch c.op {
	case opBetween:
		fieldCount = 2
	case opIn:
		fieldCount = len(c.values)
		if fieldCount == 0 {
			return "", nil
		}
	case opIsNull:
		fieldCount = 0
	}

	// Validate values count
	if len(c.values) < fieldCount {
		return "", nil
	}

	// Convert values to bindable arguments
	args := make([]interface{}, fieldCount)
	for i := 0; i < fieldCount; i++ {
		a, ok := utils.MakeArg(c.values[i])
		if !ok {
			return "", nil
		}
		args[i] = a
	}

	// Construct and return the SQL clause
	switch c.op {
	case opIn:
		phs := make([]string, fieldCount)
		for i := range phs {
			phs[i] = ph(pos)
		}
		return fmt.Sprintf(operators[opCode], c.field, strings.Join(phs, ",")), args
	case opBetween:
		// Preserve previous behaviour of ordering the range values, using
		// their SQL string representation for comparison as before
		v1, _ := utils.MakeValue(c.values[0])
		v2, _ := utils.MakeValue(c.values[1])
		if v1 > v2 {
			args[0], args[1] = args[1], args[0]
		}
		return fmt.Sprintf(operators[opCode], c.field, ph(pos), ph(pos)), args
	case opIsNull:
		return fmt.Sprintf(operators[opCode], c.field), args
	default:
		return fmt.Sprintf(operators[opCode], c.field, ph(pos)), args
	}
}

// getConjunction returns the logical conjunction used to combine this clause
// with other clauses in a WHERE condition
//
// @receiver c The clause instance
// @return The conjunction used for combining clauses
func (c *clause) getConjunction() conjunction {
	return c.conjunction
}

// newClause creates a new clause instance
//
// @param c The logical conjunction for combining clauses
// @param f The field name for the clause
// @param o The operator for the clause
// @param n Whether the clause is negated
// @param v The values for the clause
// @return A new clause instance
func newClause(c conjunction, f string, o operator, n bool, v ...interface{}) *clause {
	return &clause{
		conjunction: c,
		field:       f,
		not:         n,
		op:          o,
		values:      v,
	}
}
