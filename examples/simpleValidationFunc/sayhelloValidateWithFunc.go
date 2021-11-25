package main

import (
	"fmt"
	"go-cli-prompt/prompt"
	"net"
)

func main() {
	ipPrompt := prompt.Prompt{
		PromptMessage:   "IP Address",
		DefaultAsString: "192.168.1.1",
		InputValidatorFunc: func(s string) bool {
			parsedIp := net.ParseIP(s)
			if parsedIp == nil {
				return false
			}
			return true
		},
		OutputSerializerFunc: func(s string) (interface{}, error) {
			parsedIp := net.ParseIP(s)
			if parsedIp == nil {
				return nil, fmt.Errorf("unparsable IP address") // This shouldn't happen if the validator is correct
			}
			return parsedIp, nil
		},
		InvalidInputMessage: "Not a valid IP address",
	}

	// Try entering a string with a "!"
	ret, _ := ipPrompt.Show()

	fmt.Printf("IP is %v!", ret.(net.IP))
}
