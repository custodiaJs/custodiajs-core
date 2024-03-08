package main

import (
	"fmt"
)

// Hier sollten alle zuvor definierten Typen und Methoden eingefügt werden
// Token, TokenType, Lexer, newToken, isLetter, isDigit,
// skipWhitespace, readIdentifier, readNumber, readLineComment, readBlockComment,
// und alle anderen benötigten Funktionen.

func main() {
	// Beispieltext, der analysiert werden soll
	_ = `
rblockcall ("server uri", {}, <userPub:=userPub, host:=store.host, agent:=store.agent>) {}
catch(error) {}
final(result) {}
`
	input := `
rblockcall ("server uri", {}, <userPub:=userPub, host:=store.host, agent:=store.agent>) {
	final;
}
catch(error) {

}
final(result) {

}
`

	// Initialisiere den Lexer mit dem Eingabetext
	lexer := NewLexer(input)

	// Iteriere durch die Token, bis das Ende der Eingabe erreicht ist
	for tok := lexer.NextToken(); tok.Type != EOF; tok = lexer.NextToken() {
		fmt.Println(tok)
	}
}
