package dtorm

import "fmt"

type fieldSize struct {
	size    int
	decimal int
}

// newSize
// @param sz
// @param dec
// @return fieldSize
func newSize(sz, dec int) fieldSize {
	return fieldSize{
		size:    sz,
		decimal: dec,
	}
}

// String
// @receiver s
// @return string
func (s fieldSize) String() string {
	if s.decimal > 0 {
		return fmt.Sprintf("%d,%d", s.size, s.decimal)
	}
	return fmt.Sprintf("%d", s.size)
}
