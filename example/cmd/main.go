package main

import (
	"context"
	"errors"
	"fmt"
	"go-cli-prompt/prompt"
	"go-cli-prompt/validation"
	"net"
	"regexp"
	"strconv"
	"time"
)

func main() {
	namePrompt := prompt.CliPrompt{
		PromptMessage:        "Full Name",
		AllowNil:             false,
		IsPassword:           false,
		InvalidInputMessage:  "",
		DefaultAsString:      "",
		InputValidatorFunc:   validation.MakeInputValidatorChain(validation.MinLen(2), validation.MaxLen(4)),
		OutputSerializerFunc: nil,
		MapKey:               "Name",
	}

	agePrompt := prompt.CliPrompt{
		PromptMessage:       "Age",
		AllowNil:            false,
		IsPassword:          false,
		InvalidInputMessage: "",
		DefaultAsString:     "",
		InputValidatorRegex: regexp.MustCompile(`^[1-9]+[0-9]*$`),
		OutputSerializerFunc: func(s string) (interface{}, error) {
			return strconv.Atoi(s)
		},
		MapKey: "Age",
	}

	ipAddressPrompt := prompt.CliPrompt{
		PromptMessage:       "Ip",
		AllowNil:            false,
		IsPassword:          false,
		InvalidInputMessage: "",
		DefaultAsString:     "192.168.1.1",
		InputValidatorFunc:  validation.IpAddress,
		OutputSerializerFunc: func(s string) (interface{}, error) {
			ip := net.ParseIP(s)
			if ip == nil {
				return nil, errors.New("some ip error")
			}
			return ip, nil
		},
		MapKey: "Ip",
	}

	pwQ := prompt.CliPrompt{
		PromptMessage:        "Password",
		AllowNil:             false,
		IsPassword:           true,
		InvalidInputMessage:  "",
		DefaultAsString:      "",
		InputValidatorFunc:   nil,
		OutputSerializerFunc: nil,
		MapKey:               "Password",
	}

	r1, _ := namePrompt.DisplayWithContext(context.Background())
	fmt.Printf("GOT: %#v\n", r1)

	timeout, _ := context.WithTimeout(context.Background(), time.Second*5)

	r2, err := agePrompt.DisplayWithContext(timeout)
	fmt.Printf("Ret: %#v Err: %#v\n", r2, err)

	responses, _ := prompt.MakePromptChain(namePrompt, agePrompt, ipAddressPrompt, pwQ).Display()

	fmt.Printf("%#v", responses)
}
