package babelcoin

import (
	"strings"
	"errors"
)

type Symbol struct {
	s string
	parts []string
}

func ParseSymbol(s string) (Symbol) {
	return Symbol{s, strings.Split(s, "/")}
}

func (s *Symbol) Exchange() string {
	return s.parts[0]
}

func (s *Symbol) Pair() (string, error) {
	if len(s.parts) < 2 {
		return "", errors.New("Failed to parse pair from symbol "+s.s)
	} else {
		return s.parts[1], nil
	}
}
