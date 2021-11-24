package main

import (
	"context"
	"fmt"
	"go-cli-prompt/prompt"
	"time"
)

func main() {
	namePrompt := prompt.CliPrompt{
		PromptMessage:        "Enter Name",
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second * 5)

	ret, err := namePrompt.DisplayWithContext(ctx)

	if err != nil {
		fmt.Printf("Timed out!")
	} else {
		fmt.Printf("Hello %v!", ret.(string))
	}
}
