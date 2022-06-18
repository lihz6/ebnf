package main

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

func Verify(grammar map[string]*Production, start string) {
	new(verifier).verify(grammar, start)
}

type verifier struct {
	worklist []*Production
	reached  map[string]*Production
	grammar  map[string]*Production
}

func (v *verifier) verify(grammar map[string]*Production, start string) {
	root, found := grammar[start]
	if !found {
		panic(fmt.Sprintf("no start production %s", start))
	}

	v.worklist = v.worklist[0:0]
	v.reached = make(map[string]*Production)
	v.grammar = grammar

	v.push(root)
	for {
		n := len(v.worklist) - 1
		if n < 0 {
			break
		}
		prod := v.worklist[n]
		v.worklist = v.worklist[0:n]
		v.verifyExpr(prod.Content, isLexical(prod.Identifier.Name))
	}

	if len(v.reached) < len(v.grammar) {
		for name, prod := range v.grammar {
			if _, found := v.reached[name]; !found {
				panic(fmt.Sprintf("%s %s is unreachable", prod.Identifier.Position, name))
			}
		}
	}
}

func isLexical(name string) bool {
	ch, _ := utf8.DecodeRuneInString(name)
	return !unicode.IsUpper(ch)
}

func (v *verifier) push(prod *Production) {
	name := prod.Identifier.Name
	if _, found := v.reached[name]; !found {
		v.worklist = append(v.worklist, prod)
		v.reached[name] = prod
	}
}

func (v *verifier) verifyExpr(expr Expression, lexical bool) {
	switch x := expr.(type) {
	case nil:
	case *Alternative:
		for _, e := range *x {
			v.verifyExpr(e, lexical)
		}
	case *Sequence:
		for _, e := range *x {
			v.verifyExpr(e, lexical)
		}
	case *Identifier:
		if prod, found := v.grammar[x.Name]; found {
			v.push(prod)
		} else {
			panic(fmt.Sprintf("%s missing production %s", x.Position, x.Name))
		}
		if lexical && !isLexical(x.Name) {
			panic(fmt.Sprintf("%s reference to non-lexical production %s", x.Position, x.Name))
		}
	case *Token:
	case *Range:
		i := v.verifyChar(x.Begin)
		j := v.verifyChar(x.End)
		if i >= j {
			panic(fmt.Sprintf("%s decreasing character range", x.Begin.Position))
		}
	case *Grouping:
		v.verifyExpr(x.Content, lexical)
	case *Optional:
		v.verifyExpr(x.Content, lexical)
	case *Repetition:
		v.verifyExpr(x.Content, lexical)
	default:
		panic(fmt.Sprintf("internal error: unexpected type %T", expr))
	}
}

func (v *verifier) verifyChar(x *Token) rune {
	s, _ := strconv.Unquote(x.Token)
	if utf8.RuneCountInString(s) != 1 {
		panic(fmt.Sprintf("%s single char expected, found %s", x.Position, s))
	}
	ch, _ := utf8.DecodeRuneInString(s)
	return ch
}
