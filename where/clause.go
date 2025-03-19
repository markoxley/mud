package where

import (
	"fmt"
	"strings"

	"github.com/markoxley/dtorm/utils"
)

// clause
type clause struct {
	conjunction conjunction
	field       string
	not         bool
	op          operator
	values      []interface{}
}

// String
// @receiver c
// @return string
func (c clause) String(operators []string) string {
	opCode := int(c.op)
	if c.not {
		opCode += len(operators) / 2
	}
	fieldCount := 1
	switch c.op {
	case opBetween:
		fieldCount = 2
	case opIn:
		fieldCount = len(c.values)
		if fieldCount == 0 {
			return ""
		}
	case opIsNull:
		fieldCount = 0
	}
	if len(c.values) < fieldCount {
		return ""
	}
	vls := make([]string, fieldCount)
	for i := 0; i < fieldCount; i++ {
		f, ok := utils.MakeValue(c.values[i])
		if !ok {
			return ""
		}
		vls[i] = f
	}

	switch c.op {
	case opIn:
		return fmt.Sprintf(operators[opCode], c.field, strings.Join(vls, ","))
	case opBetween:
		v1 := vls[0]
		v2 := vls[1]
		if v1 > v2 {
			v1 = vls[1]
			v2 = vls[0]
		}
		return fmt.Sprintf(operators[opCode], c.field, v1, v2)
	case opIsNull:
		return fmt.Sprintf(operators[opCode], c.field)
	default:
		return fmt.Sprintf(operators[opCode], c.field, vls[0])
	}
}

// getConjunction
// @receiver c
// @return conjunction
func (c *clause) getConjunction() conjunction {
	return c.conjunction
}

// newClause
// @param c
// @param f
// @param o
// @param n
// @param v
// @return *clause
func newClause(c conjunction, f string, o operator, n bool, v ...interface{}) *clause {
	return &clause{
		conjunction: c,
		field:       f,
		not:         n,
		op:          o,
		values:      v,
	}
}
