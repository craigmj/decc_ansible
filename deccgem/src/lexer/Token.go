package lexer

import (
	"bytes"
	"errors"
)

// Token is a single lexed token that is emitted by the lexer.
type Token struct {
	Type int
	Value string
}

func (t *Token) Error() error {
	if t.Type==ERROR {
		return errors.New(t.Value)
	}
	return nil
}
func (t *Token) Ok() *Token {
	if t.Type==ERROR {
		return nil
	}
	return t
}
func (t *Token) OkError() (*Token, error) {
	if t.Type==ERROR {
		if t.Value=="EOF" {
			return nil, nil
		}
		return nil, errors.New(t.Value)
	}
	return t, nil
}

func (t *Token) String() string {
	var out bytes.Buffer
	switch t.Type {
	case ERROR:
		out.Write([]byte("ERROR: ["))
	case GEMNAME:
		out.Write([]byte("NAME : ["))
	case GEMVERSION:
		out.Write([]byte("VER  : ["))
	case GEMLANG:
		out.Write([]byte("LANG : ["))
	case GEMBINARY:
		out.Write([]byte("BINRY: ["))
	}
	out.Write([]byte(t.Value))
	out.Write([]byte("]"))
	return out.String()
}

const (
	ERROR = iota
	GEMNAME
	GEMVERSION
	GEMLANG
	GEMBINARY
)

type TokenChan chan*Token

func (t TokenChan) Error(err error) {
	t <- &Token{ ERROR, err.Error() }
}

func (t TokenChan) Emit(tok int, v string) {
	t <- &Token{ tok, v }
}