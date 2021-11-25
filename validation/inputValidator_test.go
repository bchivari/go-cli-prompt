package validation

import (
	"fmt"
	"strings"
	"testing"
)

func TestMakeInputValidatorChain(t *testing.T) {
	type args struct {
		validators []InputValidator
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Make Chain",
			args: args{
				validators: []InputValidator{func(s string) bool {
					return strings.Contains(s, "test")
				}, func(s string) bool {
					return strings.Contains(s, "bob")
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := MakeInputValidatorChain(tt.args.validators...)
			if v("hello") {
				fmt.Errorf("improper validation")
			}
			if v("test") {
				fmt.Errorf("improper validation")
			}
			if v("bob") {
				fmt.Errorf("improper validation")
			}
			if !v("bob likes to test") {
				fmt.Errorf("improper validation")
			}
		})
	}
}
