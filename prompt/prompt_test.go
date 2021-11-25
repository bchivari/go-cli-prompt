package prompt

import (
	"bytes"
	"context"
	"fmt"
	"go-cli-prompt/serialization"
	"go-cli-prompt/validation"
	"io"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestCliPromptChain_Display(t *testing.T) {
	promptMessage := "Enter Name"
	responseText1 := "Bobby"
	responseText2 := "Sarah"
	key1 := "name1"
	key2 := "name2"

	tests := []struct {
		name    string
		p1      Prompt
		p2      Prompt
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "Working Chain",
			p1: Prompt{
				PromptMessage: promptMessage,
				MapKey:        key1,
				outputWriter:  new(bytes.Buffer),
				inputReader:   bytes.NewBufferString(fmt.Sprintf("%v\n", responseText1)),
			},
			p2: Prompt{
				PromptMessage: promptMessage,
				MapKey:        key2,
				outputWriter:  new(bytes.Buffer),
				inputReader:   bytes.NewBufferString(fmt.Sprintf("%v\n", responseText2)),
			},
			want:    map[string]interface{}{key1: responseText1, key2: responseText2},
			wantErr: false,
		},
		{
			name: "Broken Chain",
			p1: Prompt{
				PromptMessage: promptMessage,
				MapKey:        key1,
				outputWriter:  new(bytes.Buffer),
				inputReader:   bytes.NewBufferString(fmt.Sprintf("%v\n", responseText1)),
			},
			p2: Prompt{
				PromptMessage: promptMessage,
				outputWriter:  new(bytes.Buffer),
				inputReader:   bytes.NewBufferString(fmt.Sprintf("%v\n", responseText2)),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := makeChain(tt.p1, tt.p2)
			got, err := chain.Show()
			if (err != nil) != tt.wantErr {
				t.Errorf("Show() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Show() got = %v, wantText %v", got, tt.want)
			}
		})
	}
}

func TestCliPromptChain_DisplayWithContext(t *testing.T) {
	var (
		promptMessage     = "Enter Name"
		responseText1     = "Bobby"
		responseText2     = "Sarah"
		key1              = "name1"
		key2              = "name2"
		ctxWithTimeout, _ = context.WithTimeout(context.Background(), time.Millisecond*250)
		ctxWithoutTimeout = context.Background()
	)
	tests := []struct {
		name    string
		ctx     context.Context
		p1      Prompt
		p2      Prompt
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "With Timeout - Times out",
			p1: Prompt{
				PromptMessage: promptMessage,
				MapKey:        key1,
				outputWriter:  new(bytes.Buffer),
				inputReader:   bytes.NewBufferString(fmt.Sprintf("%v\n", responseText1)),
			},
			p2: Prompt{
				PromptMessage: promptMessage,
				MapKey:        key2,
				outputWriter:  new(bytes.Buffer),
				inputReader:   bytes.NewBufferString(""),
			},
			ctx:     ctxWithTimeout,
			want:    nil,
			wantErr: true,
		},
		{
			name: "Without Timeout",
			p1: Prompt{
				PromptMessage: promptMessage,
				MapKey:        key1,
				outputWriter:  new(bytes.Buffer),
				inputReader:   bytes.NewBufferString(fmt.Sprintf("%v\n", responseText1)),
			},
			p2: Prompt{
				PromptMessage: promptMessage,
				MapKey:        key2,
				outputWriter:  new(bytes.Buffer),
				inputReader:   bytes.NewBufferString(fmt.Sprintf("%v\n", responseText2)),
			},
			ctx:     ctxWithoutTimeout,
			want:    map[string]interface{}{key1: responseText1, key2: responseText2},
			wantErr: false,
		},
		{
			name: "Without Timeout - Missing Key Error",
			p1: Prompt{
				PromptMessage: promptMessage,
				MapKey:        key1,
				outputWriter:  new(bytes.Buffer),
				inputReader:   bytes.NewBufferString(fmt.Sprintf("%v\n", responseText1)),
			},
			p2: Prompt{
				PromptMessage: promptMessage,
				outputWriter:  new(bytes.Buffer),
				inputReader:   bytes.NewBufferString(fmt.Sprintf("%v\n", responseText2)),
			},
			ctx:     ctxWithoutTimeout,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chain := makeChain(tt.p1, tt.p2)
			got, err := chain.ShowWithContext(tt.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShowWithContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShowWithContext() got = %v, wantText %v", got, tt.want)
			}
		})
	}
}

func makeChain(prompts ...Prompt) *PromptList {
	return MakePromptList(prompts...)
}

func TestPrompt_RealStdin(t *testing.T) {
	var (
		promptMessage = "Enter Name1"
		writer        = new(bytes.Buffer)
		p             = &Prompt{
			PromptMessage: promptMessage,
			AllowNil:      true,
			outputWriter:  writer,
		}
	)

	got, _ := p.Show()

	time.Sleep(time.Millisecond * 100)
	if !reflect.DeepEqual(got, nil) {
		t.Errorf("ShowWithContext() got = %v, wantText %v", got, nil)
	}
}

func TestPrompt_Display(t *testing.T) {
	var (
		promptMessage1              = "Enter Name1"
		responseText                = "Bobby"
		invalidInputMessage         = "badinput"
		mockScanner         scanner = &MockScanner{
			returnTextFifo: nil,
			returnError:    fmt.Errorf("some scanner error"),
		}
	)

	type fields struct {
		PromptMessage              string
		AllowNil                   bool
		IsPassword                 bool
		InvalidInputMessage        string
		DefaultAsString            string
		InputValidatorFunc         validation.InputValidator
		InputValidatorRegex        *regexp.Regexp
		OutputSerializerFunc       serialization.OutputSerializer
		MapKey                     string
		PromptMessageDelim         string
		SuppressTrimWhitespace     bool
		SuppressEchoInputOnInvalid bool
		outputWriter               io.Writer
		inputReader                io.Reader
		scanner                    scanner
	}
	tests := []struct {
		name                  string
		fields                fields
		want                  interface{}
		wantErr               bool
		wantOutputWriterRegex *regexp.Regexp
	}{
		{
			name: "Test Prompt Received",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("%v\n", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: nil,
		},
		{
			name: "Scanner Errors out",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("%v\n", responseText)),
				scanner:                    mockScanner,
			},
			want:                  nil,
			wantErr:               true,
			wantOutputWriterRegex: nil,
		},
		{
			name: "Test using default Stdout",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               nil,
				inputReader:                bytes.NewBufferString(fmt.Sprintf("%v\n", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: nil,
		},
		{
			name: "Test AllowNil = false && InvalidInputMessage printed",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        invalidInputMessage,
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("\n%v", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: regexp.MustCompile(`(?s)^.*` + invalidInputMessage + `.*$`),
		},
		{
			name: "Test IsPassword = true",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 true,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("%v\n", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: nil,
		},
		{
			name: "Test Default Used",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            responseText,
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString("\n"),
			},
			want:                  responseText,
			wantOutputWriterRegex: nil,
		},
		{
			name: "InputValidatorFunc fail/pass",
			fields: fields{
				PromptMessage:       promptMessage1,
				AllowNil:            false,
				IsPassword:          false,
				InvalidInputMessage: invalidInputMessage,
				DefaultAsString:     "",
				InputValidatorFunc: func(s string) bool {
					if s == responseText {
						return true
					}
					return false
				},
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("Billy\n%v\n", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: regexp.MustCompile(`(?s)^.*` + invalidInputMessage + `.*$`),
		},
		{
			name: "InputValidatorRegex fail/pass",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        invalidInputMessage,
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        regexp.MustCompile("^" + responseText + "$"),
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("Billy\n%v\n", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: regexp.MustCompile(`(?s)^.*` + invalidInputMessage + `.*$`),
		},
		{
			name: "OutputSerializerFunc",
			fields: fields{
				PromptMessage:       promptMessage1,
				AllowNil:            false,
				IsPassword:          false,
				InvalidInputMessage: "",
				DefaultAsString:     "",
				InputValidatorFunc:  nil,
				InputValidatorRegex: nil,
				OutputSerializerFunc: func(s string) (interface{}, error) {
					return fmt.Sprintf("***%v", s), nil
				},
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("%v\n", responseText)),
			},
			want:                  "***" + responseText,
			wantOutputWriterRegex: nil,
		},
		{
			name: "OutputSerializerFunc return error",
			fields: fields{
				PromptMessage:       promptMessage1,
				AllowNil:            false,
				IsPassword:          false,
				InvalidInputMessage: "",
				DefaultAsString:     "Reggie",
				InputValidatorFunc:  nil,
				InputValidatorRegex: nil,
				OutputSerializerFunc: func(s string) (interface{}, error) {
					if s != responseText {
						return nil, fmt.Errorf("some error")
					}
					return responseText, nil
				},
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("\n%v\n", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: regexp.MustCompile(`(?s)^.*some error.*$`),
		},
		{
			name: "PromptMessageDelim",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "<> ",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("%v\n", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: regexp.MustCompile(`(?s)^.*` + promptMessage1 + `<> .*$`),
		},
		{
			name: "SuppressTrimWhitespace = false",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "<> ",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("  %v  \n", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: nil,
		},
		{
			name: "SuppressTrimWhitespace = true",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "<> ",
				SuppressTrimWhitespace:     true,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("  %v  \n", responseText)),
			},
			want:                  "  " + responseText + "  ",
			wantOutputWriterRegex: nil,
		},
		{
			name: "SuppressEchoInputOnInvalid = false",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        regexp.MustCompile(`(?s)^Bobby$`),
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("Craig\n%v\n", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: regexp.MustCompile(`(?s)^.*Craig.*$`),
		},
		{
			name: "SuppressEchoInputOnInvalid = true",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        regexp.MustCompile(`(?s)^Bobby$`),
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: true,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("zzzz\n%v\n", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: regexp.MustCompile(`(?s)^[^z]*$`),
		},
		{
			name: "SuppressEchoInputOnInvalid = true, empty string, no default",
			fields: fields{
				PromptMessage:              promptMessage1,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: true,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("\n%v\n", responseText)),
			},
			want:                  responseText,
			wantOutputWriterRegex: regexp.MustCompile(`(?s)^[^\[]*$`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Prompt{
				PromptMessage:              tt.fields.PromptMessage,
				AllowNil:                   tt.fields.AllowNil,
				IsPassword:                 tt.fields.IsPassword,
				InvalidInputMessage:        tt.fields.InvalidInputMessage,
				DefaultAsString:            tt.fields.DefaultAsString,
				InputValidatorFunc:         tt.fields.InputValidatorFunc,
				InputValidatorRegex:        tt.fields.InputValidatorRegex,
				OutputSerializerFunc:       tt.fields.OutputSerializerFunc,
				MapKey:                     tt.fields.MapKey,
				PromptMessageDelim:         tt.fields.PromptMessageDelim,
				SuppressTrimWhitespace:     tt.fields.SuppressTrimWhitespace,
				SuppressEchoInputOnInvalid: tt.fields.SuppressEchoInputOnInvalid,
				outputWriter:               tt.fields.outputWriter,
				inputReader:                tt.fields.inputReader,
				scanner:                    tt.fields.scanner,
			}

			got, err := h.Show()

			if (err != nil) != tt.wantErr {
				t.Errorf("ShowWithContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Show() = %v, wantText %v", got, tt.want)
			}
			if tt.fields.outputWriter != nil && !strings.HasPrefix(tt.fields.outputWriter.(*bytes.Buffer).String(), tt.fields.PromptMessage) {
				t.Errorf("Show() wantText stdout prefix = %v, got stdout prefix = %v", tt.fields.PromptMessage, tt.fields.outputWriter.(*bytes.Buffer).String())
			}
			if tt.wantOutputWriterRegex != nil && !tt.wantOutputWriterRegex.MatchString(tt.fields.outputWriter.(*bytes.Buffer).String()) {
				t.Errorf("Show() wantText outputwriter match on = %v, got content = %v", tt.wantOutputWriterRegex.String(), tt.fields.outputWriter.(*bytes.Buffer).String())
			}
			if tt.fields.outputWriter != nil {
				fmt.Println(tt.fields.outputWriter.(*bytes.Buffer).String())
			}

		})
	}
}

func TestPrompt_DisplayWithContext(t *testing.T) {
	var (
		promptMessage             = "Enter Name"
		responseText              = "Bobby"
		ctxWithTimeout, _         = context.WithTimeout(context.Background(), time.Millisecond*250)
		ctxWithoutTimeout         = context.Background()
		mockScanner       scanner = &MockScanner{
			returnTextFifo: nil,
			returnError:    fmt.Errorf("some scanner error"),
		}
	)

	type fields struct {
		PromptMessage              string
		AllowNil                   bool
		IsPassword                 bool
		InvalidInputMessage        string
		DefaultAsString            string
		InputValidatorFunc         validation.InputValidator
		InputValidatorRegex        *regexp.Regexp
		OutputSerializerFunc       serialization.OutputSerializer
		MapKey                     string
		PromptMessageDelim         string
		SuppressTrimWhitespace     bool
		SuppressEchoInputOnInvalid bool
		outputWriter               io.Writer
		inputReader                io.Reader
		scanner                    scanner
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "Without Timeout",
			fields: fields{
				PromptMessage:              promptMessage,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("%v\n", responseText)),
			},
			args: args{
				ctx: ctxWithoutTimeout,
			},
			want:    responseText,
			wantErr: false,
		},
		{
			name: "Without Timeout - input error",
			fields: fields{
				PromptMessage:              promptMessage,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(fmt.Sprintf("%v\n", responseText)),
				scanner:                    mockScanner,
			},
			args: args{
				ctx: ctxWithoutTimeout,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "With Timeout",
			fields: fields{
				PromptMessage:              promptMessage,
				AllowNil:                   false,
				IsPassword:                 false,
				InvalidInputMessage:        "",
				DefaultAsString:            "",
				InputValidatorFunc:         nil,
				InputValidatorRegex:        nil,
				OutputSerializerFunc:       nil,
				MapKey:                     "",
				PromptMessageDelim:         "",
				SuppressTrimWhitespace:     false,
				SuppressEchoInputOnInvalid: false,
				outputWriter:               new(bytes.Buffer),
				inputReader:                bytes.NewBufferString(""),
			},
			args: args{
				ctx: ctxWithTimeout,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Prompt{
				PromptMessage:              tt.fields.PromptMessage,
				AllowNil:                   tt.fields.AllowNil,
				IsPassword:                 tt.fields.IsPassword,
				InvalidInputMessage:        tt.fields.InvalidInputMessage,
				DefaultAsString:            tt.fields.DefaultAsString,
				InputValidatorFunc:         tt.fields.InputValidatorFunc,
				InputValidatorRegex:        tt.fields.InputValidatorRegex,
				OutputSerializerFunc:       tt.fields.OutputSerializerFunc,
				MapKey:                     tt.fields.MapKey,
				PromptMessageDelim:         tt.fields.PromptMessageDelim,
				SuppressTrimWhitespace:     tt.fields.SuppressTrimWhitespace,
				SuppressEchoInputOnInvalid: tt.fields.SuppressEchoInputOnInvalid,
				outputWriter:               tt.fields.outputWriter,
				inputReader:                tt.fields.inputReader,
				scanner:                    tt.fields.scanner,
			}
			got, err := h.DisplayWithContext(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShowWithContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShowWithContext() got = %v, wantText %v", got, tt.want)
			}
		})
	}
}

func TestMakePromptChain(t *testing.T) {
	type args struct {
		prompts []Prompt
	}
	tests := []struct {
		name string
		args args
		want *PromptList
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakePromptList(tt.args.prompts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakePromptList() = %v, wantText %v", got, tt.want)
			}
		})
	}
}
