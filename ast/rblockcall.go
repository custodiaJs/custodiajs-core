package ast

import (
	"fmt"
	"vnh1/static"
)

func (p *Parser) parseRBlockCallStatementParmOptions() map[string]interface{} {
	options := make(map[string]interface{})

	// Erwarte, dass das aktuelle Token eine öffnende geschweifte Klammer '{' ist
	if !p.currentTokenIs(static.LBRACE) {
		return nil
	}
	p.nextToken() // Gehe zur nächsten Token

	// Verarbeite die Token, bis eine schließende geschweifte Klammer '}' gefunden wird
	for !p.currentTokenIs(static.RBRACE) {
		if p.currentTokenIs(static.STRING) {
			key := p.currentToken().Literal
			// Erwarte ASSIGN_INIT Token (:=)
			if !p.expectPeek(static.ASSIGN_INIT) {
				return nil
			}
			p.nextToken() // Zum Wert-Token gehen
			// Hier verarbeiten wir den Wert. Wir unterstützen Strings und geschachtelte Objekte
			var value interface{}
			if p.currentTokenIs(static.STRING) {
				value = p.currentToken().Literal
			} else if p.currentTokenIs(static.LBRACE) { // Start eines geschachtelten Objekts
				value = p.parseRBlockCallStatementParmOptions()
				if value == nil {
					return nil
				}
			} else {
				// Behandle andere Typen oder gib einen Fehler zurück
				return nil
			}
			options[key] = value
			p.nextToken() // Weiter zum nächsten Token
		}

		// Wenn ein Komma gefunden wird, überspringe es und mache weiter (Optionale Logik, basierend auf deiner Syntax)
		if p.currentTokenIs(static.COMMA) {
			p.nextToken()
		}
	}

	// Überspringe die schließende geschweifte Klammer '}'
	p.nextToken()
	return options
}

func (p *Parser) parseRBlockCallStatementPassedParms() []*static.RBlockCallPassParms {
	// Prüfung auf Beginn der Argumentliste '('
	if !p.currentTokenIs(static.LBRACE) {
		return nil
	}

	if !p.expectPeek(static.RBRACE) {
		return nil
	}

	p.nextToken()
	return []*static.RBlockCallPassParms{}
}

func (p *Parser) parseRBlockCallStatementParameterParents() (string, map[string]interface{}, []*static.RBlockCallPassParms, bool) {
	// Prüfung auf Beginn der Argumentliste '('
	if !p.currentTokenIsAndNext(static.LPAREN) {
		fmt.Println("HERE", p.currentToken())
		return "", nil, nil, false
	}

	// Parsen der URI als STRING
	if !p.currentTokenIs(static.STRING) {
		fmt.Println("HERE 1", p.currentToken())
		return "", nil, nil, false
	}

	// Die URL wird ausgelesen
	uri := p.currentTokenAndNext().Literal

	// Gehe zum nächsten Token, das ein ',' sein sollte
	if !p.currentTokenIs(static.COMMA) {
		fmt.Println("HERE 2")
		return "", nil, nil, false
	}

	// Gehe zum nächsten Token, das ein '{' sein sollte
	if !p.expectPeek(static.LBRACE) {
		fmt.Println("HERE 3")
		return "", nil, nil, false
	}

	// Die Optionen werden ausgelesen
	options := p.parseRBlockCallStatementParmOptions()
	if options == nil {
		fmt.Println("HERE 4")
		return "", nil, nil, false
	}

	// Prüfe ob es sich um ein Komma handelt
	if !p.currentTokenIsAndNext(static.RPAREN) {
		fmt.Println("HERE 5", p.currentToken().Literal)
		return "", nil, nil, false
	}

	// Die Passed Parms werden eingelesen
	currentPassedParms := []*static.RBlockCallPassParms{}

	// Es wird geprüft ob als nächstes eine Zulässige Kette vorhanden ist, wenn ja wird das Token entfernt
	if p.expectNextTokenChain(static.AND, static.LPAREN) {
		// Nächster Token
		p.nextToken()

		// Die Vorhandenen Parameter werden eigneleesn
		currentPassedParms = append(currentPassedParms, p.parseRBlockCallStatementPassedParms()...)
	}

	// Die Passed Parms werden ermittelt
	return uri, options, currentPassedParms, true
}

func (p *Parser) parseRBlockCallStatement() *static.RBlockCallStatement {
	// Initialisiere ein neues RBlockCallStatement
	statement := &static.RBlockCallStatement{}

	// Überspringe das rblockcall-Token
	p.nextToken()

	// Die Parameter werden eingelesen
	p.parseRBlockCallStatementParameterParents()

	// Parsen der URI als STRING
	if !p.currentTokenIs(static.STRING) {
		return nil
	}
	statement.URI = p.currentToken().Literal

	// Gehe zum nächsten Token, das entweder ',' oder ')' sein sollte
	p.nextToken()

	// Parsen weiterer Argumente oder schließen der Argumentliste
	// Dies würde erfordern, dass du durch die Token iterierst, bis du ')' findest
	// Für dieses Beispiel überspringen wir die Details der Argument-Parsing-Logik

	// Suche nach dem Beginn des Körpers '{'
	if !p.expectPeek(static.LBRACE) {
		return nil
	}

	// Stelle sicher, dass der gesamte rblockcall korrekt geparst wurde,
	// und kehre dann das Statement zurück
	return statement
}
