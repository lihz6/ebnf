package main

import (
	"fmt"
	"io"
	"text/scanner"
)

func Parse(filename string, src io.Reader) map[string]*Production {
	return new(parser).parse(filename, src)
}

type parser struct {
	scanner scanner.Scanner
	pos     scanner.Position
	lit     string
	tok     rune
}

func (p *parser) parse(filename string, src io.Reader) map[string]*Production {
	p.scanner.Init(src)
	p.scanner.Filename = filename
	p.scanner.Mode ^= scanner.SkipComments
	p.next()

	grammar := make(map[string]*Production)

	for p.tok != scanner.EOF {
		prod := p.parseProduction()
		if _, found := grammar[prod.Identifier.Name]; found {
			panic(fmt.Sprintf("%s %s declared already", prod.Identifier.Position, prod.Identifier.Name))
		}
		grammar[prod.Identifier.Name] = prod
	}

	return grammar
}

func (p *parser) next() {
	p.tok = p.scanner.Scan()
	p.pos = p.scanner.Position
	p.lit = p.scanner.TokenText()
}

func (p *parser) parseProduction() *Production {
	Header := &Identifier{p.pos, p.lit}
	p.expect(scanner.Ident)
	p.expect('=')
	var Content Expression
	if p.tok != '.' {
		Content = p.parseAlternative()
		if grouping, ok := Content.(*Grouping); ok {
			Content = grouping.Content
		}
	}
	p.expect('.')
	return &Production{Header, Content}
}

func (p *parser) expect(tok rune) {
	if p.tok != tok {
		p.expectError(p.pos, scanner.TokenString(tok))
	}
	p.next()
}

func (p *parser) expectError(pos scanner.Position, expected string) {
	msg := `expected "` + expected + `"`
	if pos.Offset == p.pos.Offset {
		msg += ", found " + scanner.TokenString(p.tok)
		if p.tok < 0 {
			msg += " " + p.lit
		}
	}
	panic(fmt.Sprintf("%s %s", pos, msg))
}

func (p *parser) parseAlternative() (exp Expression) {
	var Content Alternative
	for {
		Content = append(Content, p.parseSequence())
		if p.tok != '|' {
			break
		}
		p.next()
	}
	switch len(Content) {
	case 0:
		p.expectError(p.pos, "alternative")
	case 1:
		exp = Content[0]
	default:
		exp = &Content
	}
	return
}

func (p *parser) parseSequence() (exp Expression) {
	var Content Sequence
	for x := p.parseExpression(); x != nil; x = p.parseExpression() {
		Content = append(Content, x)
	}
	switch len(Content) {
	case 0:
		p.expectError(p.pos, "sequence")
	case 1:
		exp = Content[0]
	default:
		exp = &Content
	}
	return
}

func (p *parser) parseExpression() (exp Expression) {
	Position := p.pos
	switch p.tok {
	case '(':
		p.next()
		exp = p.parseAlternative()
		if _, ok := exp.(*Alternative); ok {
			exp = &Grouping{Position, exp}
		}
		p.expect(')')
	case '[':
		p.next()
		exp = p.parseAlternative()
		switch e := exp.(type) {
		case *Repetition:
			exp = e.Content
		case *Optional:
			exp = e.Content
		case *Grouping:
			exp = e.Content
		}
		exp = &Optional{Position, exp}
		p.expect(']')
	case '{':
		p.next()
		exp = p.parseAlternative()
		switch e := exp.(type) {
		case *Repetition:
			exp = e.Content
		case *Optional:
			exp = e.Content
		case *Grouping:
			exp = e.Content
		}
		exp = &Repetition{Position, exp}
		p.expect('}')
	case scanner.String, scanner.RawString:
		exp = p.parseToken()
		if p.tok == '…' || p.tok == '-' {
			p.next()
			exp = &Range{exp.(*Token), p.parseToken()}
		}
	case scanner.Ident:
		exp = &Identifier{p.pos, p.lit}
		p.next()
	case scanner.Comment:
		exp = &Token{p.pos, p.lit}
		p.next()
	case scanner.Char, scanner.Int, scanner.Float:
		p.expectError(p.pos, "expression")
	}
	return
}

func (p *parser) parseToken() (token *Token) {
	if p.tok != scanner.String && p.tok != scanner.RawString {
		p.expect(scanner.String)
	}
	token = &Token{p.pos, p.lit}
	p.next()
	return
}
