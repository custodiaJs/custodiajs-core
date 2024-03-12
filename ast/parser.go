package ast

import "vnh1/static"

type Parser struct {
	tokens  []static.Token
	current int // Aktuelle Position im tokens slice
}

func NewParser(tokens []static.Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) nextToken() static.Token {
	if p.current < len(p.tokens) {
		tok := p.tokens[p.current]
		p.current++
		return tok
	}
	return static.Token{Type: "EOF", Literal: ""}
}

func (p *Parser) currentToken() static.Token {
	if p.current < len(p.tokens) {
		return p.tokens[p.current]
	}
	return static.Token{Type: "EOF", Literal: ""}
}

func (p *Parser) ParseStatement() static.Statement {
	currentToken := p.currentToken()
	switch currentToken.Type {
	case static.RBLOCKCALL:
		p.parseRBlockCallStatement()
		return nil
	default:
		return nil
	}
}

func (p *Parser) ParseProgram() *static.Program {
	program := &static.Program{}
	for p.currentToken().Type != "EOF" {
		stmt := p.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}
