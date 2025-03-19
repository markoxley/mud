package dtorm

import (
	"fmt"
	"reflect"
)

const (
	majorVersion   = 0
	minorVersion   = 1
	releaseVersion = 0
)

var (
	fieldTrans = map[fieldType][]reflect.Kind{
		tBool:     {reflect.Bool},
		tDateTime: {reflect.Struct},
		tDouble:   {reflect.Float64},
		tFloat:    {reflect.Float32},
		tInt:      {reflect.Int8, reflect.Uint8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint16, reflect.Uint32},
		tLong:     {reflect.Int64, reflect.Uint64},
		tString:   {reflect.String},
	}
	fieldUnsigned = []reflect.Kind{
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
	}
)

func Version() string {
	return fmt.Sprintf("DTORM v%d.%d.%d ©2023 I Have a Hat", majorVersion, minorVersion, releaseVersion)
}
