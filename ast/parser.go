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
		return p.parseRBlockCallStatement()
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
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		// Fehlerbehandlung: Ungültiges Token
		return false
	}
}

func (p *Parser) currentTokenIs(t types.TokenType) bool {
	return p.currentToken().Type == t
}

func (p *Parser) parseRBlockCallStatement() *types.RBlockCallStatement {
	// Initialisiere ein neues RBlockCallStatement
	statement := &types.RBlockCallStatement{}

	// Überspringe das rblockcall-Token
	p.nextToken()

	// Prüfung auf Beginn der Argumentliste '('
	if !p.expectPeek(types.LPAREN) {
		return nil
	}

	// Überspringe '('
	p.nextToken()

	// Parsen der URI als STRING
	if !p.currentTokenIs(types.STRING) {
		return nil // Oder Fehlerbehandlung
	}
	statement.URI = p.currentToken().Literal

	// Gehe zum nächsten Token, das entweder ',' oder ')' sein sollte
	p.nextToken()

	// Parsen weiterer Argumente oder schließen der Argumentliste
	// Dies würde erfordern, dass du durch die Token iterierst, bis du ')' findest
	// Für dieses Beispiel überspringen wir die Details der Argument-Parsing-Logik

	// Suche nach dem Beginn des Körpers '{'
	if !p.expectPeek(types.LBRACE) {
		return nil
	}

	// Parsen des Körpers des rblockcall
	// Auch hier überspringen wir die Details des Körper-Parsings für dieses Beispiel

	// Suche nach dem Ende des Körpers '}'
	// Dies setzt voraus, dass der Körper korrekt geparst wurde

	// Stelle sicher, dass der gesamte rblockcall korrekt geparst wurde,
	// und kehre dann das Statement zurück
	return statement
}
