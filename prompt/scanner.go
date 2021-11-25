package prompt

import "bufio"

// scanner interface to enable mocking
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

// MockScanner to enable testing. Emits elements in Fifo until empty then returns error
type MockScanner struct {
	returnTextFifo []string
	returnError    error
}

func (s *MockScanner) popText() string {
	if len(s.returnTextFifo) > 0 {
		ret := s.returnTextFifo[0]
		s.returnTextFifo[0] = ""
		s.returnTextFifo = s.returnTextFifo[1:]
		return ret
	}
	return ""
}

func (s *MockScanner) Scan() bool {
	return len(s.returnTextFifo) > 0
}

func (s *MockScanner) Err() error {
	if len(s.returnTextFifo) > 0 {
		return nil
	}
	return s.returnError
}

func (s *MockScanner) Text() string {
	return s.popText()
}
