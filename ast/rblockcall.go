package ast

import (
	"fmt"
	"vnh1/types"
)

func (p *Parser) parseRBlockCallStatementParmOptions() *types.RBlockCallPassParms {
	// Prüfung auf Beginn der Argumentliste '('
	if !p.currentTokenIs(types.LBRACE) {
		return nil
	}

	if !p.expectPeek(types.RBRACE) {
		return nil
	}

	p.nextToken()
	return &types.RBlockCallPassParms{}
}

func (p *Parser) parseRBlockCallStatementPassedParms() *types.RBlockCallPassParms {
	// Prüfung auf Beginn der Argumentliste '('
	if !p.currentTokenIs(types.LBRACE) {
		return nil
	}

	if !p.expectPeek(types.RBRACE) {
		return nil
	}

	p.nextToken()
	return &types.RBlockCallPassParms{}
}

func (p *Parser) parseRBlockCallStatementParameterParents() (string, *types.RBlockCallOptions, *types.RBlockCallPassParms, bool) {
	// Prüfung auf Beginn der Argumentliste '('
	if !p.currentTokenIsAndNext(types.LPAREN) {
		fmt.Println("HERE", p.currentToken())
		return "", nil, nil, false
	}

	// Parsen der URI als STRING
	if !p.currentTokenIs(types.STRING) {
		fmt.Println("HERE 1", p.currentToken())
		return "", nil, nil, false
	}

	// Die URL wird ausgelesen
	uri := p.currentTokenAndNext().Literal

	// Gehe zum nächsten Token, das ein ',' sein sollte
	if !p.currentTokenIs(types.COMMA) {
		fmt.Println("HERE 2")
		return "", nil, nil, false
	}

	// Gehe zum nächsten Token, das ein '{' sein sollte
	if !p.expectPeek(types.LBRACE) {
		fmt.Println("HERE 3")
		return "", nil, nil, false
	}

	// Die Optionen werden ausgelesen
	options := p.parseRBlockCallStatementParmOptions()
	if options == nil {
		fmt.Println("HERE 4")
		return "", nil, nil, false
	}

	fmt.Println(p.currentToken())

	// Prüfe ob es sich um ein Komma handelt
	if !p.currentTokenIsAndNext(types.COMMA) {
		fmt.Println("HERE 5", p.currentToken())
		return "", nil, nil, false
	}

	fmt.Println(p.currentToken())

	// Gehe zum nächsten Token, das ein < sein sollte
	if !p.currentTokenIs(types.LT) {
		fmt.Println("HERE 6", p.currentToken())
		return "", nil, nil, false
	}

	// Die Passed Parms werden ermittelt

	return uri, nil, nil, true
}

func (p *Parser) parseRBlockCallStatement() *types.RBlockCallStatement {
	// Initialisiere ein neues RBlockCallStatement
	statement := &types.RBlockCallStatement{}

	// Überspringe das rblockcall-Token
	p.nextToken()

	// Die Parameter werden eingelesen
	p.parseRBlockCallStatementParameterParents()

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
