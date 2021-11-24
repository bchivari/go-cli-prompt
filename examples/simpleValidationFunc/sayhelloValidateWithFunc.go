package main

import (
	"fmt"
	"go-cli-prompt/prompt"
	"strings"
)

func main() {
	namePrompt := prompt.CliPrompt{
		PromptMessage:        "Enter Name",
		InputValidatorFunc: func(s string) bool {
			return !strings.Contains(s, "!")
		},
		InvalidInputMessage: "Don't take that tone with me!",
	}

	// Try entering a string with a "!"
	ret := namePrompt.Display()

	fmt.Printf("Hello %v!", ret.(string))
}
