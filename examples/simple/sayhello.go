package main

import (
	"fmt"
	"go-cli-prompt/prompt"
)

func main() {
	namePrompt := prompt.CliPrompt{
		PromptMessage: "Enter Name",
	}

	// Blocks until valid input is received
	ret := namePrompt.Display()

	fmt.Printf("Hello %v!", ret.(string))
}
