package prompt

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go-cli-prompt/serialization"
	"go-cli-prompt/validation"
	"golang.org/x/term"
	"io"
	"os"
	"regexp"
	"strings"
	"syscall"
)

const (
	promptTemple               = "%v%v"
	promptWithDefaultTemplate  = "%v [%v]%v"
	errorTemplate              = "\n%v\n\n"
	errorEchoInputTemplate     = "\n%v [%v]\n\n"
	defaultInvalidInputMessage = "Invalid Input"
	defaultPromptMessageDelim  = ": "
	ioErrorMessage             = "I/O Error"
)

var (
	missingKeyError     = errors.New("'MapKey' field is missing from one more more CliPrompts")
	defaultOutputWriter = os.Stdout
	defaultInputReader  = os.Stdin
)

// CliPrompt The main struct that defines how a prompt is displayed and handled
// If both DefaultAsString and OutputSerializerFunc are set, the default string
// should be serializable via the OutputSerializerFunc
type CliPrompt struct {
	PromptMessage              string                         // The message prompt text displayed to user
	AllowNil                   bool                           // If this prompt accepts nil as allowable input
	IsPassword                 bool                           // If set, will suppress echoing of input to terminal
	InvalidInputMessage        string                         // Message displayed if InputValidatorFunc returns false, or nil is provided but not accepted by setting AllowNil
	DefaultAsString            string                         // The default value if the user just hits enter without providing input
	InputValidatorFunc         validation.InputValidator      // Function which validates the input string. If both InputValidatorFunc and InputValidatorRegex are provided both are tested, and both must pass for input to be valid
	InputValidatorRegex        *regexp.Regexp                 // Regex used to validate the input
	OutputSerializerFunc       serialization.OutputSerializer // Function which converts the input string into a desired type returned as interface{}
	MapKey                     string                         // If utilizing a CliPromptChain, this string is used as a key in the map[string]interface{} returned by Display()
	PromptMessageDelim         string                         // The string/character displayed after the PromptMessage. This will default to ": "
	SuppressTrimWhitespace     bool                           // By default, both leading and trailing whitespace are trimmed from all input strings, unless IsPassword is set, before validation. Setting SuppressTrimWhitespace will disable this behavior
	SuppressEchoInputOnInvalid bool                           // By default, input that fails validation will echod back as part of the error message. Setting SuppressEchoInputOnInvalid will disable this behavior

	outputWriter io.Writer // Advanced option so not exposed; Defaults to os.Stdout; Set with SetOption(WithWriter)
	inputReader  io.Reader // Advanced option so not exposed; Defaults to os.Stdin; Set with SetOption(WithReader)
	//scanner      *bufio.Scanner
	scanner scanner
}

// Display Displays a single CliPrompt and will return the supplied value. Blocks forever until valid input is received
func (h *CliPrompt) Display() interface{} {
	h.initializeScanner()
	for {
		h.displayPrompt()
		userInput, err := h.readInput()
		if err != nil {
			fmt.Fprintln(h.getOutputWriter(), ioErrorMessage)
			continue
		}

		// Got input
		if len(userInput) != 0 {
			if h.isValidInput(userInput) {
				serializedResp, err := h.serializeIfRequired(userInput)
				if err == nil && serializedResp != nil {
					return serializedResp
				}
			}
			h.displayInvalidInputMessage(userInput)
		} else {
			// Empty Input
			if h.hasDefault() {
				defaultSerialized, err := h.serializeIfRequired(h.DefaultAsString)
				if err != nil {
					fmt.Fprintf(h.getOutputWriter(), fmt.Sprintf("Default value cannot be serialized, This shouldn't happen. %v", err))
				} else {
					return defaultSerialized
				}
			}
			if h.AllowNil && !h.hasDefault() {
				return nil
			}
			if h.shouldEchoInput() {
				fmt.Fprintf(h.getOutputWriter(), errorEchoInputTemplate, h.getInvalidInputMessage(), "null")
			} else {
				fmt.Fprintf(h.getOutputWriter(), errorTemplate, h.getInvalidInputMessage())
			}
			// Loop until we get valid input
		}
	}
}

// DisplayWithContext - Same as Display but is context aware so can be canceled / timed out
func (h *CliPrompt) DisplayWithContext(ctx context.Context) (interface{}, error) {
	select {
	case ret := <-h.displayAsync():
		return ret, nil
	case <-ctx.Done():
		return nil, errors.New("call was canceled by context")
	}
}

func (h *CliPrompt) displayAsync() <-chan interface{} {
	resultChan := make(chan interface{})
	go func() {
		resultChan <- h.Display()
	}()
	return resultChan
}

// CliPromptChain represents a collection of CliPrompts; Used by (*CliPromptChain) Display for displaying prompts in series and collecting responses as a map
type CliPromptChain []CliPrompt

// MakePromptChain is a helper function used to assemble a CliPromptChain from individual CliPrompt instances
func MakePromptChain(prompts ...CliPrompt) *CliPromptChain {
	var chain CliPromptChain
	for _, p := range prompts {
		chain = append(chain, p)
	}
	return &chain
}

// Display displays all prompts in the CliPromptChain in succession and returns all responses as a map. Blocks forever until valid input is received for all prompts
func (c *CliPromptChain) Display() (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	for _, p := range *c {
		if p.MapKey == "" {
			return nil, missingKeyError
		}
		ret[p.MapKey] = p.Display()
	}
	return ret, nil
}

// DisplayWithContext - Same as Display but is context aware so can be canceled / timed out
func (c *CliPromptChain) DisplayWithContext(ctx context.Context) (map[string]interface{}, error) {
	resultChan, errChan := c.displayAsync()
	select {
	case ret := <-resultChan:
		return ret, nil
	case err := <-errChan:
		return nil, err
	case <-ctx.Done():
		return nil, errors.New("call was canceled by context")

	}
}

func (c *CliPromptChain) displayAsync() (<-chan map[string]interface{}, <-chan error) {
	resultChan := make(chan map[string]interface{})
	errChan := make(chan error)
	go func() {
		result, err := c.Display()
		if err != nil {
			errChan <- err
		}
		resultChan <- result
	}()
	return resultChan, errChan
}

func (h *CliPrompt) serializeIfRequired(input string) (interface{}, error) {
	if h.OutputSerializerFunc == nil {
		return input, nil
	}
	return h.OutputSerializerFunc(input)
}

func (h *CliPrompt) displayPrompt() {
	if h.hasDefault() {
		fmt.Fprintf(h.getOutputWriter(), promptWithDefaultTemplate, h.PromptMessage, h.DefaultAsString, h.getDelim())
	} else {
		fmt.Fprintf(h.getOutputWriter(), promptTemple, h.PromptMessage, h.getDelim())
	}
}

func (h *CliPrompt) getDelim() string {
	if h.PromptMessageDelim != "" {
		return h.PromptMessageDelim
	}
	return defaultPromptMessageDelim
}

func (h *CliPrompt) getOutputWriter() io.Writer {
	if h.outputWriter != nil {
		return h.outputWriter
	}
	return defaultOutputWriter
}

func (h *CliPrompt) getInputReader() io.Reader {
	if h.inputReader != nil {
		return h.inputReader
	}
	return defaultInputReader
}

func (h *CliPrompt) getInvalidInputMessage() string {
	if h.InvalidInputMessage != "" {
		return h.InvalidInputMessage
	}
	return defaultInvalidInputMessage
}

func (h *CliPrompt) shouldEchoInput() bool {
	if h.SuppressEchoInputOnInvalid {
		return false
	}
	return true
}

func (h *CliPrompt) hasDefault() bool {
	return h.DefaultAsString != ""
}

func (h *CliPrompt) isValidInput(s string) bool {
	if h.validateAgainstRegexIfProvided(s) && h.validateAgainstFuncIfProvided(s) {
		return true
	}
	return false
}

func (h *CliPrompt) validateAgainstRegexIfProvided(s string) bool {
	if h.InputValidatorRegex != nil {
		return h.InputValidatorRegex.MatchString(s)
	}
	return true
}

func (h *CliPrompt) validateAgainstFuncIfProvided(s string) bool {
	if h.InputValidatorFunc != nil {
		return h.InputValidatorFunc(s)
	}
	return true
}

func (h *CliPrompt) readInput() (string, error) {
	if !h.IsPassword {
		return h.readRegularInput()
	} else {
		defer func() {
			// Print blank line; Non-echoing password reader requires this
			fmt.Fprintln(h.getOutputWriter(), "")
		}()
		return h.readPasswordInput()
	}
}

func (h *CliPrompt) readPasswordInput() (string, error) {
	if h.inputReader != nil {
		return h.readPasswordFromIoReader()
	} else {
		return h.readPasswordFromStdin()
	}
}

func (h *CliPrompt) readPasswordFromStdin() (string, error) {
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(h.getOutputWriter(), "")
	if err != nil {
		return "", fmt.Errorf(errorTemplate, h.getInvalidInputMessage())
	}
	return string(bytePassword), nil
}

func (h *CliPrompt) readPasswordFromIoReader() (string, error) {
	var password []byte
	if file, ok := h.inputReader.(*os.File); ok && term.IsTerminal(int(file.Fd())) {
		password, err := term.ReadPassword(int(file.Fd()))
		if err != nil {
			return "", fmt.Errorf(errorTemplate, h.getInvalidInputMessage())
		}
		return string(password), nil
	}

	if _, err := fmt.Fscanf(h.inputReader, "%s\n", &password); err != nil {
		return "", fmt.Errorf(errorTemplate, h.getInvalidInputMessage())
	}
	return string(password), nil
}

func (h *CliPrompt) readRegularInput() (string, error) {
	h.scanner.Scan()
	if h.SuppressTrimWhitespace {
		return h.scanner.Text(), h.scanner.Err()
	}
	return strings.TrimSpace(h.scanner.Text()), h.scanner.Err()
}

func (h *CliPrompt) displayInvalidInputMessage(response string) {
	if h.shouldEchoInput() {
		fmt.Fprintf(h.getOutputWriter(), errorEchoInputTemplate, h.getInvalidInputMessage(), response)
		return
	}
	fmt.Fprintf(h.getOutputWriter(), errorTemplate, h.getInvalidInputMessage())
}

func (h *CliPrompt) initializeScanner() {
	if h.scanner == nil {
		h.scanner = newDefaultScanner(bufio.NewScanner(h.getInputReader()))
	}
}