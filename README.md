# go-cli-prompt
[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)
![GitHub Workflow Status](https://github.com/bchivari/go-cli-prompt/workflows/CI/badge.svg)
[!codecov](https://codecov.io/gh/bchivari/go-cli-prompt)
[![Go Report Card](https://goreportcard.com/badge/github.com/bchivari/go-cli-prompt)](https://goreportcard.com/report/github.com/bchivari/go-cli-prompt)

**go-cli-prompt** is a simple go library that can be used to display a prompt (or series of prompts) for user
input and collect the result(s). It is highly configurable and provides a framework for supplying default values,
validating input, serializing output, and cancellation via `conetex.Context`.

By default, I/O takes place on `os.Stdin` / `os.Stdout`; However any implementation of `io.Reader` / `io.Writer` may be
supplied. 

## Getting Started

### Installing

* `go get github.com/bchivari/go-cli-prompt`

## Examples

### Simple Prompt

* A simple example of how to display a single prompt
* [See code...](http://www.google.ca)

*Code*
```golang
namePrompt := prompt.Prompt{
    PromptMessage: "Enter Name",
    AllowNil:      false,
}

// Blocks & re-prompts until valid input is received
name, err := namePrompt.Show()
if err != nil {
    fmt.Printf("Hello %v!", name.(string))
}
```

*Output*
```
Enter Name: Bob
Hello Bob!
```
### Multi-Prompt

* An example that demonstrates how to chain together multiple prompts
* [See code...](http://www.google.ca)

*Code*
```golang
namePrompt := prompt.CliPrompt{
    PromptMessage: "Name",
    MapKey:        "name",
}

agePrompt := prompt.CliPrompt{
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

ret, err = prompt.MakePromptChain(namePrompt, agePrompt).Display()
if err != nil {
    return
}

fmt.Printf("Got: %#v\n", ret)
fmt.Printf("Hello %v, you are %v years old!", ret["name"].(string), ret["age"].(string))
```

*Output*
```
Name: Bob
Age: 199

Age should be between 0 - 150 [199]

Age: 50
Got: map[string]interface {}{"age":"50", "name":"Bob"}
Hello Bob, you are 50 years old!

```

### Simple Prompt With Regexp Validation

* [See code...](http://www.google.ca)

*Code*
```golang
namePrompt := prompt.CliPrompt{
    PromptMessage: "Enter Name",
    InputValidatorRegex: regexp.MustCompile(`^[a-zA-Z]*$`),
    AllowNil: true,
}

// Try entering a number
name := namePrompt.Display()
if name != nil {
    fmt.Printf("Hello %v!", name.(string))
}
```

*Output*
```
Enter Name: 123

Invalid Input [123]

Enter Name: Bob
Hello Bob!
```

### Prompt With Default Value and Custom Validation

* [See code...](http://www.google.ca)

*Code*
```golang
ipPrompt := prompt.CliPrompt{
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

ret := ipPrompt.Display()

fmt.Printf("IP is %v!", ret.(net.IP))
```

*Output 1*
```
IP Address [192.168.1.1]: Not an IP

Not a valid IP address [Not an IP]

IP Address [192.168.1.1]: 10.0.0.1
IP is 10.0.0.1!
```

*Output 2 - Using default*
```
IP Address [192.168.1.1]: 
IP is 192.168.1.1!
```

### More examples
* [See code...](https://github.com/bchivari/go-cli-prompt/tree/master/examples)

## Full Documentation

[![Go Reference](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://pkg.go.dev/golang.org/x/example)


## Authors

Brad Chivari  
[@tracklessca](https://twitter.com/tracklessca)

## Version History

* 0.1
    * Initial Release

## License

This project is licensed under Apace License, Version 2.0 - see the LICENSE file for details