package main

import (
	"fmt"
	"io"
	"strings"
)

func Print(grammar map[string]*Production, start string, dst io.Writer) {
	new(printer).print(grammar, start, dst)
}

type printer struct {
	writer  io.Writer
	reached map[string]bool
	grammar map[string]*Production
}

func (p *printer) print(grammar map[string]*Production, start string, dst io.Writer) {
	p.writer = dst
	p.reached = make(map[string]bool)
	p.grammar = grammar
	p.printProduction(start, 0)
}

func (p *printer) printProduction(name string, depth int) {
	if p.reached[name] {
		return
	}
	p.reached[name] = true
	if prod, ok := p.grammar[name]; ok {
		str := fmt.Sprintf("%s%s\n", strings.Join(make([]string, depth), "  "), prod)
		p.writer.Write([]byte(str))
		p.printExpression(prod.Content, depth+1)
	}
}

func (v *printer) printExpression(expr Expression, depth int) {
	switch x := expr.(type) {
	case nil:
	case *Token:
	case *Range:
	case *Identifier:
		v.printProduction(x.Name, depth)
	case *Alternative:
		for _, e := range *x {
			v.printExpression(e, depth)
		}
	case *Sequence:
		for _, e := range *x {
			v.printExpression(e, depth)
		}
	case *Grouping:
		v.printExpression(x.Content, depth)
	case *Optional:
		v.printExpression(x.Content, depth)
	case *Repetition:
		v.printExpression(x.Content, depth)
	default:
		panic(fmt.Sprintf("internal error: unexpected type %T", expr))
	}
}
