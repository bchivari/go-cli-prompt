package prompt

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

func TestCliPrompt_SetOptions(t *testing.T) {
	type fields struct {
		PromptMessage string
		AllowNil      bool
	}
	type args struct {
		opts []Opt
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "No Error",
			fields: fields{
				PromptMessage: "",
				AllowNil:      false,
			},
			args: args{
				opts: []Opt{func(p *Prompt) error {
					p.IsPassword = true
					p.MapKey = "testkey"
					return nil
				}},
			},
			wantErr: false,
		},
		{
			name: "With Error",
			fields: fields{
				PromptMessage: "",
				AllowNil:      false,
			},
			args: args{
				opts: []Opt{func(p *Prompt) error {
					p.IsPassword = true
					p.MapKey = "testkey"
					return fmt.Errorf("some error")
				}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Prompt{
				PromptMessage: tt.fields.PromptMessage,
				AllowNil:      tt.fields.AllowNil,
			}
			if err := h.SetOptions(tt.args.opts...); (err != nil) != tt.wantErr {
				t.Errorf("SetOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWithWriter(t *testing.T) {
	p := &Prompt{PromptMessage: "My Message"}
	myWriter := bufio.NewWriter(bytes.NewBufferString("test"))
	p.SetOptions(WithWriter(myWriter))
	assertEqual(t, myWriter, p.outputWriter)
}

func TestWithReader(t *testing.T) {
	p := &Prompt{PromptMessage: "My Message"}
	myReader := bufio.NewReader(bytes.NewBufferString("test"))
	p.SetOptions(WithReader(myReader))
	assertEqual(t, myReader, p.inputReader)
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}
