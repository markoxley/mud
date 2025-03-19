package dtorm

import "fmt"

type ErrNoResults struct {
	Err error
}

func NoResults(msg string) ErrNoResults {
	return ErrNoResults{Err: fmt.Errorf(msg)}
}

func (e ErrNoResults) Error() string {
	return e.Err.Error()
}
