package main

import (
	"context"
	"fmt"
	"github.com/bchivari/go-cli-prompt/prompt"
	"time"
)

func main() {
	namePrompt := prompt.Prompt{
		PromptMessage: "Enter Name",
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)

	ret, err := namePrompt.DisplayWithContext(ctx)

	if err != nil {
		fmt.Printf("Timed out!")
	} else {
		fmt.Printf("Hello %v!", ret.(string))
	}
}
