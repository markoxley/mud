// Copyright (c) 2025 DaggerTech. All rights reserved.
// Use of this source code is governed by an MIT license that can be
// found in the LICENSE file.
// Package mud provides a database ORM (Object-Relational Mapping) implementation
// with support for SQLite, MySQL, and SQL Server databases.
package mud

import (
	"fmt"

	"github.com/markoxley/mud/order"
	"github.com/markoxley/mud/where"
)

// Criteria is used to safely build search criteria for database queries.
// It provides a structured way to define WHERE, ORDER BY, LIMIT, and OFFSET conditions.
type Criteria struct {
	// Where defines the WHERE condition for the query
	Where interface{}
	// Order defines the ORDER BY condition for the query
	Order interface{}
	// Limit specifies the maximum number of rows to return
	Limit int
	// Offset specifies the number of rows to skip
	Offset int
	// IncDeleted indicates whether to include soft-deleted records
	IncDeleted bool
}

// WhereClause returns the WHERE condition in parameterized SQL format along
// with the ordered list of arguments to bind to its placeholders.
// Where conditions built with a where.Builder emit placeholder markers; a
// Where supplied as a raw string is passed through verbatim as an explicit
// escape hatch.
// Parameters:
//
//	mgr: The database manager used to format the WHERE condition
//
// Returns:
//
//	A string containing the SQL WHERE clause and its bind arguments
func (c Criteria) WhereClause(mgr Manager) (string, []interface{}) {
	return c.whereClause(mgr, nil)
}

// whereClause is WhereClause with an optional placeholder position counter.
// A nil counter starts numbering at 1; pass a seeded counter when earlier
// placeholders in the same statement have already been consumed.
func (c Criteria) whereClause(mgr Manager, pos *int) (string, []interface{}) {
	wh := ""
	var args []interface{}
	whereDone := false
	if c.Where != nil {
		if b, ok := c.Where.(*where.Builder); ok {
			wh, args = b.BuildAt(mgr.Operators(), mgr.Placeholder, seededPos(pos))
		} else if b, ok := c.Where.(where.Builder); ok {
			wh, args = b.BuildAt(mgr.Operators(), mgr.Placeholder, seededPos(pos))
		} else if w, ok := c.Where.(string); ok {
			wh = w
		}
	}

	if wh != "" {
		wh = fmt.Sprintf(" WHERE %s", wh)
		whereDone = true
	}
	if !c.IncDeleted {
		if whereDone {
			wh += " AND"
		} else {
			wh += "WHERE"
		}
		wh += fmt.Sprintf(" %s IS NULL", mgr.IdentityString("DeleteDate"))
	}
	return wh, args
}

// OrderString returns the ORDER BY condition in SQL format.
// It converts the criteria's Order condition into a properly formatted SQL ORDER BY clause.
// Parameters:
//
//	mgr: The database manager used to format the ORDER BY condition
//
// Returns:
//
//	A string containing the SQL ORDER BY clause
func (c Criteria) OrderString(mgr Manager) string {
	if c.Order == nil {
		return ""
	}
	orderStr := ""
	if s, ok := c.Order.(string); ok {
		orderStr = s
	} else if ob, ok := c.Order.(*order.Builder); ok {
		orderStr = ob.String(mgr.IdentityString)
	} else if str, ok := c.Order.(interface{ String() string }); ok {
		orderStr = str.String()
	}

	if orderStr != "" {
		orderStr = fmt.Sprintf(" ORDER BY %s", orderStr)
	}
	return orderStr
}

// LimitString returns the limiter in SQL format
// @receiver c
// @return string
func (c Criteria) LimitString(mgr Manager) string {
	return mgr.LimitString(&c)
}

// OffsetString returns the offset in SQL format
// @receiver c
// @return string
func (c Criteria) OffsetString(mgr Manager) string {
	return mgr.OffsetString(&c)
}

// Build returns the full criteria in parameterized SQL format along with the
// ordered list of arguments to bind to the WHERE clause placeholders.
// @receiver c
// @return string
// @return []interface{}
func (c Criteria) Build(mgr Manager) (string, []interface{}) {
	where, args := c.WhereClause(mgr)
	return mgr.BuildQuery(where, c.OrderString(mgr), c.LimitString(mgr), c.OffsetString(mgr)), args
}

// seededPos returns pos as-is when provided, or a fresh counter starting at 1.
func seededPos(pos *int) *int {
	if pos != nil {
		return pos
	}
	p := 1
	return &p
}
