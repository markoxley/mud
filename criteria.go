// Package dtorm provides a database ORM (Object-Relational Mapping) implementation
// with support for SQLite, MySQL, and SQL Server databases.
package dtorm

import "fmt"

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

// WhereString returns the WHERE condition in SQL format.
// It converts the criteria's Where condition into a properly formatted SQL WHERE clause.
// Parameters:
//   mgr: The database manager used to format the WHERE condition
// Returns:
//   A string containing the SQL WHERE clause
func (c Criteria) WhereString(mgr Manager) string {
	if c.Where == nil {
		return ""
	}
	where := ""
	whereDone := false
	switch c.Where.(type) {
	case string:
		where, _ = c.Where.(string)
	case fmt.Stringer:
		st, _ := c.Where.(fmt.Stringer)
		where = st.String()
	}
	if where != "" {
		where = fmt.Sprintf(" WHERE %s", where)
		whereDone = true
	}
	if !c.IncDeleted {
		if whereDone {
			where += " AND"
		} else {
			where += "WHERE"
		}
		where += fmt.Sprintf(" %s IS NULL", mgr.IdentityString("DeleteDate"))
	}
	return where
}

// OrderString returns the ORDER BY condition in SQL format.
// It converts the criteria's Order condition into a properly formatted SQL ORDER BY clause.
// Parameters:
//   mgr: The database manager used to format the ORDER BY condition
// Returns:
//   A string containing the SQL ORDER BY clause
func (c Criteria) OrderString(mgr Manager) string {
	if c.Order == nil {
		return ""
	}
	order := ""
	switch c.Order.(type) {
	case string:
		order, _ = c.Order.(string)
	case fmt.Stringer:
		st, _ := c.Order.(fmt.Stringer)
		order = st.String()
	}

	if order != "" {
		order = fmt.Sprintf(" ORDER BY %s", order)
	}
	return order
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

// String returns the full criteria in SQL format
// @receiver c
// @return string
func (c Criteria) String(mgr Manager) string {
	return mgr.BuildQuery(c.WhereString(mgr), c.OrderString(mgr), c.LimitString(mgr), c.OffsetString(mgr))
}
