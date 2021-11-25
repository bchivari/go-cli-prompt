package prompt

import (
	"io"
)

// Opt accepts a reference to a Prompt and is used to modify its state
type Opt func(*Prompt) error

// SetOptions will iterate over all provided Opt objects and call them
func (h *Prompt) SetOptions(opts ...Opt) error {
	for _, o := range opts {
		if err := o(h); err != nil {
			return err
		}
	}
	return nil
}

// WithWriter returns an option func which sets a customized (non stdout) io.Writer
func WithWriter(w io.Writer) Opt {
	return func(p *Prompt) error {
		p.outputWriter = w
		return nil
	}
}

// WithReader returns an option func which sets a customized (non stdout) io.Reader
func WithReader(r io.Reader) Opt {
	return func(p *Prompt) error {
		p.inputReader = r
		return nil
	}
}
