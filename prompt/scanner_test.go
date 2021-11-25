package prompt

import (
	"fmt"
	"testing"
)

func TestMockScanner_Err(t *testing.T) {
	type fields struct {
		returnTextFifo []string
		returnError    error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Empty FIFO - Error",
			fields: fields{
				returnTextFifo: nil,
				returnError:    fmt.Errorf("some error"),
			},
			wantErr: true,
		},
		{
			name: "Non-Empty FIFO - No Error",
			fields: fields{
				returnTextFifo: []string{"text"},
				returnError:    fmt.Errorf("some error"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockScanner{
				returnTextFifo: tt.fields.returnTextFifo,
				returnError:    tt.fields.returnError,
			}
			if err := s.Err(); (err != nil) != tt.wantErr {
				t.Errorf("Err() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMockScanner_Scan(t *testing.T) {
	type fields struct {
		returnTextFifo []string
		returnError    error
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Gets true",
			fields: fields{
				returnTextFifo: []string{"text"},
				returnError:    fmt.Errorf("some error"),
			},
			want: true,
		},
		{
			name: "Gets false",
			fields: fields{
				returnTextFifo: []string{},
				returnError:    fmt.Errorf("some error"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockScanner{
				returnTextFifo: tt.fields.returnTextFifo,
				returnError:    tt.fields.returnError,
			}
			if got := s.Scan(); got != tt.want {
				t.Errorf("Scan() = %v, wantText %v", got, tt.want)
			}
		})
	}
}

func TestMockScanner_Text(t *testing.T) {
	type fields struct {
		returnTextFifo []string
		returnError    error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Gets text",
			fields: fields{
				returnTextFifo: []string{"text"},
				returnError:    fmt.Errorf("some error"),
			},
			want: "text",
		},
		{
			name: "Gets first text",
			fields: fields{
				returnTextFifo: []string{"text1", "text2"},
				returnError:    fmt.Errorf("some error"),
			},
			want: "text1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockScanner{
				returnTextFifo: tt.fields.returnTextFifo,
				returnError:    tt.fields.returnError,
			}
			if got := s.Text(); got != tt.want {
				t.Errorf("Text() = %v, wantText %v", got, tt.want)
			}
		})
	}
}

func TestMockScanner_TextMultiCall(t *testing.T) {
	var (
		strings     = []string{"string1", "string2"}
		wantStrings []string
		err         = fmt.Errorf("some error")
	)
	copy(wantStrings, strings)

	type fields struct {
		returnTextFifo []string
		returnError    error
	}
	tests := []struct {
		name     string
		fields   fields
		wantText []string
		wantScan []bool
		wantErr  []error
	}{
		{
			name: "MultiCall Text",
			fields: fields{
				returnTextFifo: strings,
				returnError:    err,
			},
			wantText: wantStrings,
			wantScan: []bool{true, true, false},
			wantErr:  []error{nil, nil, err},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &mockScanner{
				returnTextFifo: tt.fields.returnTextFifo,
				returnError:    tt.fields.returnError,
			}
			for i := 0; i < len(tt.wantText); i++ {
				if got := s.Scan(); got != tt.wantScan[i] {
					t.Errorf("Scan() = %v, wantScan %v", got, tt.wantScan[i])
				}
				if got := s.Text(); got != tt.wantText[i] {
					t.Errorf("Text() = %v, wantText %v", got, tt.wantText[i])
				}
				if got := s.Err(); got != tt.wantErr[i] {
					t.Errorf("Err() = %v, wantErr %v", got, tt.wantErr[i])
				}
			}
		})
	}
}
