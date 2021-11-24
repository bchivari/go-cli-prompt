package main

import (
	"fmt"
	"go-cli-prompt/prompt"
	"strconv"
)

func main() {
	namePrompt := prompt.CliPrompt{
		PromptMessage:        "Name",
		MapKey: "name",
	}
	age := prompt.CliPrompt{
		PromptMessage:        "Age",
		InputValidatorFunc: func(s string) bool {
			i, err := strconv.Atoi(s)
			if err != nil {
				return false
			}
			if i >= 0 && i <= 150 {
				return true
			}
			return false
		},
		InvalidInputMessage: "Age should be between 0 - 150",
		MapKey: "age",
	}


	var (
		ret map[string]interface{}
		err error
	)

	ret, err = prompt.MakePromptChain(namePrompt, age).Display()
	if err != nil {
		return
	}

	fmt.Printf("Hello %v, you are %v years old!", ret["name"].(string), ret["age"].(string))
}
