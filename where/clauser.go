package where

// clauser
type clauser interface {
	String([]string) string
	getConjunction() conjunction
}
