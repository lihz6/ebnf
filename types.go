package main

import (
	"fmt"
	"strings"
	"text/scanner"
)

type Expression interface {
	String() string
}



type Production struct {
	Identifier *Identifier
	Content    Expression
} // a = b .

func (p *Production) String() string {
	if p.Content == nil {
		return fmt.Sprintf("%s = .", p.Identifier)
	}
	return fmt.Sprintf("%s = %s .", p.Identifier, p.Content)
}

type Identifier struct {
	Position scanner.Position
	Name     string
} // a

func (n *Identifier) String() string {
	return n.Name
}

type Alternative []Expression // a | b

func (a *Alternative) String() string {
	sb := new(strings.Builder)
	for i, e := range *a {
		if i > 0 {
			sb.WriteString(" | ")
		}
		sb.WriteString(e.String())
	}
	return sb.String()
}

type Sequence []Expression // a b

func (s *Sequence) String() string {
	sb := new(strings.Builder)
	for i, e := range *s {
		if i > 0 {
			sb.WriteRune(' ')
		}
		sb.WriteString(e.String())
	}
	return sb.String()
}

type Grouping struct {
	Position scanner.Position
	Content  Expression
} // ( a )

func (g *Grouping) String() string {

	return fmt.Sprintf("( %s )", g.Content)
}

type Optional struct {
	Position scanner.Position
	Content  Expression
} // [ a ]

func (o *Optional) String() string {
	return fmt.Sprintf("[ %s ]", o.Content)
}

type Repetition struct {
	Position scanner.Position
	Content  Expression
} // { a }

func (r *Repetition) String() string {
	return fmt.Sprintf("{ %s }", r.Content)
}

type Token struct {
	Position scanner.Position
	Token    string
} // "a" or `a` or /* a */

func (t *Token) String() string {
	return t.Token
}

type Range struct {
	Begin, End *Token
} // '0' … '9'

func (r *Range) String() string {
	return fmt.Sprintf("%s - %s", r.Begin, r.End)
}
