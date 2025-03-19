package where

type operator uint8

const (
	opEqual operator = iota
	opGreater
	opLess
	opLike
	opIn
	opBetween
	opIsNull
)

var operatorType [7]int = [7]int{
	dBool & dDate & dFloat & dDouble & dInt & dLong & dText,
	dDate & dFloat & dDouble & dInt & dLong & dText,
	dDate & dFloat & dDouble & dInt & dLong & dText,
	dText,
	dDate & dFloat & dDouble & dInt & dLong & dText,
	dDate & dFloat & dDouble & dInt & dLong,
	dBool & dDate & dFloat & dDouble & dInt & dLong & dText,
}
