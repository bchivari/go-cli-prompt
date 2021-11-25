package main

import (
	"fmt"
	"github.com/bchivari/go-cli-prompt/prompt"
	"strconv"
)

func main() {
	namePrompt := prompt.Prompt{
		PromptMessage: "Name",
		MapKey:        "name",
	}
	age := prompt.Prompt{
		PromptMessage: "Age",
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
		MapKey:              "age",
	}

	var (
		ret map[string]interface{}
		err error
	)

	ret, err = prompt.MakePromptList(namePrompt, age).Show()
	if err != nil {
		return
	}

	fmt.Printf("Got: %#v\n", ret)
	fmt.Printf("Hello %v, you are %v years old!", ret["name"].(string), ret["age"].(string))
}
