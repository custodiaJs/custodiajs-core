package ast

import "vnh1/static"

func (p *Parser) expectPeek(t static.TokenType) bool {
	if p.current+1 < len(p.tokens) {
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

func (p *Parser) currentTokenIs(t static.TokenType) bool {
	return p.currentToken().Type == t
}

func (p *Parser) currentTokenIsAndNext(t static.TokenType) bool {
	if p.currentToken().Type == t {
		p.nextToken()
		return true
	} else {
		return false
	}
}

func (p *Parser) currentTokenAndNext() static.Token {
	cToken := p.currentToken()
	p.nextToken()
	return cToken
}

func (p *Parser) expectNextTokenChain(types ...static.TokenType) bool {
	tempPosition := p.current // Speichere die aktuelle Position, um keine Änderungen am Parser-Zustand vorzunehmen

	for _, t := range types {
		tempPosition++ // Bewege die temporäre Position vorwärts
		// Stelle sicher, dass die temporäre Position nicht über die Länge der Token-Liste hinausgeht
		if tempPosition >= len(p.tokens) || p.tokens[tempPosition].Type != t {
			return false // Frühe Rückkehr, falls ein Token nicht wie erwartet ist
		}
	}

	// Alle erwarteten Token-Typen wurden in der Sequenz gefunden
	return true
}

func (p *Parser) matchAndUpdateForTokenChain(types ...static.TokenType) bool {
	tempPosition := p.current // Speichere die aktuelle Position, um keine Änderungen am Parser-Zustand vorzunehmen

	for _, t := range types {
		tempPosition++ // Bewege die temporäre Position vorwärts
		// Stelle sicher, dass die temporäre Position nicht über die Länge der Token-Liste hinausgeht
		if tempPosition >= len(p.tokens) || p.tokens[tempPosition].Type != t {
			return false // Frühe Rückkehr, falls ein Token nicht wie erwartet ist
		}
	}

	// Alle erwarteten Token-Typen wurden in der Sequenz gefunden,
	// aktualisiere die aktuelle Position im Parser
	p.current = tempPosition
	return true
}
