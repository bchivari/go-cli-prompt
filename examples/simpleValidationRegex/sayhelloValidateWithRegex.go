package main

import (
	"fmt"
	"go-cli-prompt/prompt"
	"regexp"
)

func main() {
	namePrompt := prompt.CliPrompt{
		PromptMessage:        "Enter Name",
		InputValidatorRegex: regexp.MustCompile(`^[a-zA-Z]*$`),
		AllowNil: true,
	}

	// Try entering a number
	ret := namePrompt.Display()
	if ret != nil {
		fmt.Printf("Hello %v!", ret.(string))
	}

}
