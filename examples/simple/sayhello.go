package main

import (
	"fmt"
	"go-cli-prompt/prompt"
)

func main() {
	namePrompt := prompt.Prompt{
		PromptMessage: "Enter Name",
		AllowNil:      false,
	}

	// Blocks & re-prompts until valid input is received
	name, err := namePrompt.Show()
	if err != nil {
		fmt.Printf("Hello %v!", name.(string))
	}
}
