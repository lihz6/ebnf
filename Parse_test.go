package main_test

import (
	. "github.com/ngolin/ebnf"
	"strings"
	"testing"
)

func TestParse_good(t *testing.T) {
	tests := []struct{ prod, want string }{
		{`a=.`, `a = .`},
		{`a=b.`, `a = b .`},
		{`a="b".`, `a = "b" .`},
		{"a=`b`.", "a = `b` ."},
		{`a=/* b */.`, `a = /* b */ .`},
		{`a="0"…"9".`, `a = "0" - "9" .`},
		{`a=(b).`, `a = b .`},
		{`a=[b].`, `a = [ b ] .`},
		{`a={b}.`, `a = { b } .`},
		{`a=a|b.`, `a = a | b .`},
		{`a=(a|b).`, `a = a | b .`},
		{`a=(a|b).`, `a = a | b .`},
		{`a=[a|b].`, `a = [ a | b ] .`},
		{`a={a|b}.`, `a = { a | b } .`},
		{`a=([a|b]).`, `a = [ a | b ] .`},
		{`a=[(a|b)].`, `a = [ a | b ] .`},
		{`a=[(a|b)].`, `a = [ a | b ] .`},
		{`a=[[a|b]].`, `a = [ a | b ] .`},
		{`a=[{a|b}].`, `a = [ a | b ] .`},
		{`a=({a|b}).`, `a = { a | b } .`},
		{`a={(a|b)}.`, `a = { a | b } .`},
		{`a={(a|b)}.`, `a = { a | b } .`},
		{`a={[a|b]}.`, `a = { a | b } .`},
		{`a={{a|b}}.`, `a = { a | b } .`},
		{`a=[(a|b)]|{(c|d)}.`, `a = [ a | b ] | { c | d } .`},
	}
	for _, test := range tests {
		src := strings.NewReader(test.prod)
		for _, prod := range Parse("", src) {
			if prod.String() != test.want {
				t.Errorf("got %s, want %s", prod, test.want)
			}
		}
	}
}
