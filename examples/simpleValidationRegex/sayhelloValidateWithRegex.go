package main

import (
	"fmt"
	"github.com/bchivari/go-cli-prompt/prompt"
	"regexp"
)

func main() {
	namePrompt := prompt.Prompt{
		PromptMessage:       "Enter Name",
		InputValidatorRegex: regexp.MustCompile(`^[a-zA-Z]*$`),
		AllowNil:            true,
	}

	// Try entering a number
	ret, _ := namePrompt.Show()
	if ret != nil {
		fmt.Printf("Hello %v!", ret.(string))
	}

}
