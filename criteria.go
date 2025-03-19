package dtorm

import "fmt"

// Criteria is used to safely build your criteria for searches
type Criteria struct {
	Where      interface{}
	Order      interface{}
	Limit      int
	Offset     int
	IncDeleted bool
}

// WhereString returns the where condition in SQL format
// @receiver c
// @return string
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

// OrderString returns the orderering in SQL format
// @receiver c
// @return string
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
