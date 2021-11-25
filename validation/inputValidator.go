package validation

// InputValidator defines a function used to validate a user input string
type InputValidator func(string) bool

// MakeInputValidatorChain can wrap N InputValidator objects into a single InputValidator; Logical AND
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
