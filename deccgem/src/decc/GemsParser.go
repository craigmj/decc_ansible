package decc

import (
	"fmt"
	"io"
	"lexer"
)

type ParserInput struct {
	C lexer.TokenChan
	T *lexer.Token
}

var workingGem *Gem

type GemRead struct {
	G   *Gem
	Err error
}

type ParserState func(in *ParserInput, out GemReadC) ParserState

type GemReadC chan GemRead

func (C GemReadC) Error(err error) ParserState {
	C <- GemRead{nil, err}
	return nil
}

func newParserInput(in io.Reader) *ParserInput {
	return &ParserInput{
		C: make(lexer.TokenChan),
		T: nil,
	}
}

func NewGemReadChan(in io.Reader) GemReadC {
	pi := &ParserInput{
		C: make(lexer.TokenChan),
		T: nil,
	}
	out := make(GemReadC)
	go lexer.ParseTokens(in, pi.C)
	go func() {
		defer close(out)
		for state := readGemName; nil != state; state = state(pi, out) {
		}
	}()
	return out
}

func (in *ParserInput) Read() (*lexer.Token, error) {
	if nil != in.T {
		t := in.T
		in.T = nil
		return t.OkError()
	}
	t, ok := <-in.C
	if !ok {
		return nil, nil // Channel is closed => EOF
	}
	return t.OkError()
}

func (in *ParserInput) Push(t *lexer.Token) {
	if nil != in.T {
		panic("Trying to push into ParserInput, but it already has a peek'd value")
	}
	in.T = t
}

func (in *ParserInput) Peek() (*lexer.Token, error) {
	if nil != in.T {
		return in.T.Ok(), in.T.Error()
	}
	var ok bool
	in.T, ok = <-in.C
	if !ok {
		return nil, nil // EOF
	}
	return in.T, nil
}

func readGemName(in *ParserInput, out GemReadC) ParserState {
	workingGem = new(Gem)
	t, err := in.Read()
	if nil != err {
		return out.Error(err)
	}
	if nil == t { // EOF
		return nil
	}
	if t.Type != lexer.GEMNAME {
		return out.Error(fmt.Errorf("Unexpected token %s", t.String()))
	}
	workingGem.Name = t.Value
	return readGemVersion
}

func readGemVersion(in *ParserInput, out GemReadC) ParserState {
	t, err := in.Read()
	if nil != err {
		return out.Error(err)
	}
	if nil == t {
		return out.Error(fmt.Errorf("Unexpected EOF after GemName %s", workingGem.Name))
	}
	if t.Type == lexer.GEMNAME {
		if 0 == len(workingGem.Version) {
			return out.Error(fmt.Errorf("GemName %s found but no version", workingGem.Name))
		}
		in.Push(t)
		return readGemName
	}
	if t.Type == lexer.GEMVERSION {
		workingGem = &Gem{Name: workingGem.Name, Version: t.Value}
		return readGemLang
	}
	return out.Error(fmt.Errorf("Unexpected token %s reading gem %s", t.String(), workingGem.Name))
}

func readGemLang(in *ParserInput, out GemReadC) ParserState {
	t, err := in.Read()
	if nil != err {
		return out.Error(err)
	}
	if nil == t { // EOF
		out <- GemRead{workingGem, nil}
		return nil
	}
	switch t.Type {
	case lexer.GEMNAME:
		out <- GemRead{workingGem, nil}
		in.Push(t)
		return readGemName
	case lexer.GEMLANG:
		workingGem.Lang = t.Value
		return readGemBinary
	case lexer.GEMBINARY:
		in.Push(t)
		return readGemBinary
	case lexer.GEMVERSION:
		out <- GemRead{workingGem, nil}
		in.Push(t)
		return readGemVersion
	}
	return out.Error(fmt.Errorf("Unexpected token type %s", t.String()))
}

func readGemBinary(in *ParserInput, out GemReadC) ParserState {
	t, err := in.Read()
	if nil != err {
		return out.Error(err)
	}
	if nil == t { // EOF
		out <- GemRead{workingGem, nil}
		return nil
	}
	switch t.Type {
	case lexer.GEMNAME:
		out <- GemRead{workingGem, nil}
		in.Push(t)
		return readGemName
	case lexer.GEMVERSION:
		out <- GemRead{workingGem, nil}
		in.Push(t)
		return readGemVersion
	case lexer.GEMLANG:
		return out.Error(fmt.Errorf("Unexpected LANG token %s reading Gem %s", t.String(), workingGem.Name))
	case lexer.GEMBINARY:
		workingGem.Binary = t.Value
		out <- GemRead{workingGem, nil}
		return readGemVersion
	}
	return out.Error(fmt.Errorf("Unexpected token type %s", t.String()))
}
