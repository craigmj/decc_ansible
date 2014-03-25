package decc

import (
	"testing"
)

func TestGems(t *testing.T) {
	a := &Gem{ "one", "3.5.1", "" }
	b := &Gem{ "two", "3.5.1pre", ""}
	c := &Gem{ "three","3.5.2", ""}
	d := &Gem{ "four", "0.54.1", ""}
	e := &Gem{ "five", "0.71.2", ""}

	if !a.EarlierThan(b) {
		t.Errorf("%s not earlier than %s", *a, *b)
	}
	if !b.EarlierThan(c) {
		t.Errorf("%s not earlier than %s", *b, *c)
	}
	if !c.EarlierThan(d) {
		t.Errorf("%s not earlier than %s", *c, *d)
	}
	if !d.EarlierThan(e) {
		t.Errorf("%s not earlier than %s", *d, *e)
	}
	if e.EarlierThan(a) {
		t.Errorf("%s earlier than %s", e, a)
	}

}