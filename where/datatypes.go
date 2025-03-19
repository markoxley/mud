package where

import "time"

type DBInt interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type DBFloat interface {
	float32 | float64
}

type DBNumeric interface {
	DBInt | DBFloat
}

type DBField interface {
	DBNumeric | bool | string | time.Time
}

const (
	dBool = 1 << iota
	dDate
	dFloat
	dDouble
	dInt
	dLong
	dText
)
