package dtorm

type fieldType int
type fieldName string

const (
	tInt fieldType = iota
	tLong
	tBool
	tDecimal
	tFloat
	tDouble
	tDateTime
	tChar
	tString
	tUUID
)

const (
	sInt      = "int"
	sLong     = "long"
	sBool     = "bool"
	sDecimal  = "decimal"
	sFloat    = "float"
	sDouble   = "double"
	sDateTime = "struct"
	sChar     = "char"
	sString   = "string"
	sUUID     = "uuid"
)

var (
	fieldNames = map[string]fieldType{
		sInt:      tInt,
		sLong:     tLong,
		sBool:     tBool,
		sDecimal:  tDecimal,
		sFloat:    tFloat,
		sDouble:   tDouble,
		sDateTime: tDateTime,
		sChar:     tChar,
		sString:   tString,
		sUUID:     tUUID,
	}
)

var (
	pgFieldNames = map[fieldType]string{
		tInt:      "INT",
		tLong:     "BIGINT",
		tBool:     "SMALLINT",
		tDecimal:  "DECIMAL",
		tFloat:    "REAL",
		tDouble:   "DOUBLE",
		tDateTime: "DATETIME",
		tChar:     "VARCHAR(1)",
		tString:   "VARCHAR",
		tUUID:     "VARCHAR(36)",
	}
)

type field struct {
	name      string
	fType     fieldType
	size      fieldSize
	identity  bool
	key       bool
	unsigned  bool
	allowNull bool
}

func newField(nm string, tp fieldType, sz, dec int, id, ky, us bool, nl bool) field {
	return field{
		name:      nm,
		fType:     tp,
		size:      newSize(sz, dec),
		identity:  id,
		key:       ky,
		unsigned:  us,
		allowNull: nl,
	}
}
