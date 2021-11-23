package validation

import (
	"net"
	"regexp"
)

var (
	lettersOnly  = regexp.MustCompile(`^[a-zA-Z]*$`)
	numbersOnly  = regexp.MustCompile(`^[1-9]*[0-9]+$`)
	alphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
	singleDigit  = regexp.MustCompile(`^[0-9]$`)
)

type InputValidatorChain func(...InputValidator) InputValidator

func MakeInputValidatorChain(validators ...InputValidator) InputValidator {
	chain := func(s string) bool {
		for _, v := range validators {
			if !v(s) {
				return false
			}
		}
		return true
	}
	return chain
}

type InputValidator func(string) bool

func MinLen(min int) InputValidator {
	return InputValidator(
		func(s string) bool {
			return len([]rune(s)) >= min
		})
}

func MaxLen(max int) InputValidator {
	return InputValidator(
		func(s string) bool {
			return len([]rune(s)) <= max
		})
}

func IpAddress(s string) bool {
	ip := net.ParseIP(s)
	return ip != nil
}

func AlphaNumericString(s string) bool {
	return alphaNumeric.MatchString(s)
}
