package ast

import (
	"vnh1/types"
)

type Parser struct {
	tokens  []types.Token
	current int // Aktuelle Position im tokens slice
}

func NewParser(tokens []types.Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) nextToken() types.Token {
	if p.current < len(p.tokens) {
		tok := p.tokens[p.current]
		p.current++
		return tok
	}
	return types.Token{Type: "EOF", Literal: ""}
}

func (p *Parser) currentToken() types.Token {
	if p.current < len(p.tokens) {
		return p.tokens[p.current]
	}
	return types.Token{Type: "EOF", Literal: ""}
}

func (p *Parser) ParseStatement() types.Statement {
	currentToken := p.currentToken()
	switch currentToken.Type {
	case types.RBLOCKCALL:
		p.parseRBlockCallStatement()
		return nil
	default:
		return nil
	}
}

func (p *Parser) ParseProgram() *types.Program {
	program := &types.Program{}
	for p.currentToken().Type != "EOF" {
		stmt := p.ParseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) expectPeek(t types.TokenType) bool {
	if p.current < len(p.tokens) {
		tok := p.tokens[p.current+1]
		if tok.Type == t {
			p.nextToken()
			return true
		} else {
			return false
		}
	}
	return false
}

func (p *Parser) currentTokenIs(t types.TokenType) bool {
	return p.currentToken().Type == t
}

func (p *Parser) currentTokenIsAndNext(t types.TokenType) bool {
	if p.currentToken().Type == t {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) currentTokenAndNext() types.Token {
	cToken := p.currentToken()
	p.nextToken()
	return cToken
}
