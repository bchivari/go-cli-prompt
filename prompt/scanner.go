package prompt

import "bufio"

type scanner interface {
	Scan() bool
	Err() error
	Text() string
}

type defaultScanner struct {
	scanner *bufio.Scanner
}

func newDefaultScanner(s *bufio.Scanner) *defaultScanner {
	i := new(defaultScanner)
	i.scanner = s
	return i
}

func (s *defaultScanner) Scan() bool {
	return s.scanner.Scan()
}

func (s *defaultScanner) Err() error {
	return s.scanner.Err()
}

func (s *defaultScanner) Text() string {
	return s.scanner.Text()
}
