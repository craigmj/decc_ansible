package lexer

import (
	// "fmt"
	"io"
	"bytes"
	"bufio"
	"unicode"
)


// Input is an input source, and handles the conversion of incoming data into tokens.
type Input struct {
	in *bufio.Reader

	buf bytes.Buffer
}

func NewInput(in io.Reader) *Input {
	return &Input{
		in: bufio.NewReader(in),
	}
}

// peek peeks at the next rune in the input
func (i *Input) peek() (rune, error) {
	r, _, err := i.in.ReadRune()
	if nil==err {
		i.in.UnreadRune()
		return r, nil
	}
	return 0, err
}

// reads a single rune
func (i *Input) read() (rune, error) {
	r, _, err := i.in.ReadRune()
	if nil==err {
		i.buf.WriteRune(r)
	}
	return r, err
}

// clear clears the buffer
func (i *Input) clear() {
	// fmt.Println("    BUFFER = ", i.buf.String())
	i.buf.Reset()
	// fmt.Println("RST BUFFER = ", i.buf.String())
}

// readUntil reads until the rune function returns true, or EOF, adding data to the buffer
func (i *Input) readUntil(until func(r rune) bool) error {
	for {
		r, _, err := i.in.ReadRune()
		if io.EOF==err {
			return nil
		}
		if nil!=err {
			return err
		}
		if (until(r)) {
			i.in.UnreadRune()
			return nil
		}
		i.buf.WriteRune(r)
	}
	return nil
}

func (i *Input) readUntilSpace() error {
	return i.readUntil(unicode.IsSpace)
}

// readWhile reads until the rune function returns false, or EOF, adding data to the buffer
func (i *Input) readWhile(while func(r rune) bool) error {
	return i.readUntil(func (r rune) bool { return !while(r) })
}

// skipSpace skips all whitespace
func (i *Input) skipSpace() error {
	for {
		r, _, err := i.in.ReadRune()
		if nil!=err {
			return err
		}
		if !unicode.IsSpace(r) {
			i.in.UnreadRune()
			return nil
		}
	}
}
