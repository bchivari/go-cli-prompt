package prompt

import (
	"io"
)

// Opt accepts a reference to a CliPrompt and is used to modify its state
type Opt func(*CliPrompt) error

// SetOptions will iterate over all provided Opt objects and call them
func (h *CliPrompt) SetOptions(opts ...Opt) error {
	for _, o := range opts {
		if err := o(h); err != nil {
			return err
		}
	}
	return nil
}

// WithWriter returns an option func which sets a customized (non stdout) io.Writer
func WithWriter(w io.Writer) Opt {
	return func(p *CliPrompt) error {
		p.outputWriter = w
		return nil
	}
}

// WithReader returns an option func which sets a customized (non stdout) io.Reader
func WithReader(r io.Reader) Opt {
	return func(p *CliPrompt) error {
		p.inputReader = r
		return nil
	}
}
