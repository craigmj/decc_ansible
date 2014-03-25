package lexer

import (
	"fmt"
	"io"
	"unicode"
)

type State func(*Input, TokenChan) (State)

func errorState(out TokenChan, err error) State {
	out.Error(err)
	close(out)
	return nil
}

func DumpTokens(in io.Reader) {
	out := make(TokenChan)
	go ParseTokens(in, out)
	for t := range out {
		fmt.Println(t.String())
	}
}

func ParseTokens(inReader io.Reader, out TokenChan) {
	in := NewInput(inReader)
	var state State
	state = startLineState
	for nil!=state {
		state = state(in, out)
	}
}

func startLineState(in *Input, out TokenChan) (State) {
	err := in.skipSpace()
	if nil!=err {
		return errorState(out, err)
	}
	p, err := in.read()
	if io.EOF==err {
		close(out)
		return nil
	}
	if nil!=err {
		return errorState(out, err)
	}
	if '*'==p {
		return readCommentLine
	}
	err = in.readUntilSpace()
	if io.EOF==err {
		if 0==in.buf.Len() {
			close(out)
			return nil
		}
		return errorState(out, fmt.Errorf("Unexpected EOF while reading GEMNAME"))
	}
	if nil!=err {
		return errorState(out, err)
	}

	out.Emit(GEMNAME, in.buf.String())
	in.clear()
	err = in.skipSpace()
	p,err = in.read()
	if nil!=err {
		return errorState(out, fmt.Errorf("Error while looking for bracket: %s", err.Error()))
	}
	if '('!=p {
		return errorState(out, fmt.Errorf("Unexpected token '%c' while looking for bracket", p))
	}
	in.clear()
	return readGemVersion
}

func readCommentLine(in *Input, out TokenChan) State {
	err := in.readUntil(func (r rune) bool {
		return '\n'==r
	})
	if io.EOF==err {
		return nil
	}
	if nil!=err {
		return errorState(out, err)
	}
	in.clear()
	return startLineState
}

func spaceBracketOrComma(r rune) bool {
	return unicode.IsSpace(r) || ','==r || ')'==r
}

func readGemVersion(in *Input, out TokenChan) (State) {
	err := in.skipSpace()
	if nil!=err {
		return errorState(out, err)
	}
	if err = in.readUntil(spaceBracketOrComma); nil!=err {
		return errorState(out, err)
	}
	out.Emit(GEMVERSION, in.buf.String())
	in.clear()
	in.skipSpace()
	p, err := in.read()
	if nil!=err {
		return errorState(out, err)
	}
	if p==')' {
		in.clear()
		return startLineState
	}
	if p==',' {
		in.clear()
		return readGemVersion
	}
	// Any other case is either LANG or BINARY
	return readGemLang
}

func readGemLang(in *Input, out TokenChan) State {
	err := in.skipSpace()
	if nil!=err {
		return errorState(out, err)
	}
	if err = in.readUntil(spaceBracketOrComma); nil!=err {
		return errorState(out, err)
	}
	lang := in.buf.String()
	in.skipSpace()
	p, err := in.read()
	if nil!=err {
		return errorState(out, err)
	}
	if p==')' {
		in.clear()
		out.Emit(GEMBINARY, lang)
		return startLineState
	}
	if p==',' {
		in.clear()
		out.Emit(GEMBINARY, lang)
		return readGemVersion
	}
	out.Emit(GEMLANG, lang)
	return readGemBinary
}

func readGemBinary(in *Input, out TokenChan) State {
	err := in.skipSpace()
	if nil!=err {
		return errorState(out, err)
	}
	if err = in.readUntil(spaceBracketOrComma); nil!=err {
		return errorState(out, err)
	}
	binary := in.buf.String()	
	in.skipSpace()
	p, err := in.read()
	if nil!=err {
		return errorState(out, err)
	}
	if p==')' {
		in.clear()
		out.Emit(GEMBINARY, binary)
		return startLineState
	}
	if p==',' {
		in.clear()
		out.Emit(GEMBINARY, binary)
		return readGemVersion
	}

	out.Emit(GEMBINARY, binary)
	return errorState(out, fmt.Errorf("Too many tokens in GEM details"))
}
