package lexer

import "vnh1/types"

type Lexer struct {
	input        string
	position     int  // aktuelle Position im Eingabetext (auf das aktuelle Zeichen)
	readPosition int  // aktuelle Leseposition im Eingabetext (nach dem aktuellen Zeichen)
	ch           byte // aktuelles Zeichen
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0 // EOF repräsentieren
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readLineComment() string {
	// Überspringe die beiden Schrägstriche "//"
	l.readChar() // Aktuelles Zeichen ist '/', gehe zum nächsten Zeichen
	l.readChar() // Gehe über den zweiten '/' hinaus zum Beginn des Kommentartextes

	startPosition := l.position
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
	// Extrahiere den Kommentartext ohne die Schrägstriche und ohne das Zeilenumbruchzeichen
	return l.input[startPosition:l.position]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readBlockComment() string {
	// Überspringe die Zeichen '/*' am Anfang des Blockkommentars
	l.readChar() // Aktuelles Zeichen ist '*', gehe zum nächsten Zeichen
	l.readChar() // Gehe über '*' hinaus, um den Inhalt des Kommentars zu erreichen

	startPosition := l.position
	for !(l.ch == '*' && l.peekChar() == '/') {
		l.readChar()
		if l.ch == 0 {
			// EOF erreicht, ohne das Ende des Kommentars zu finden, potenzieller Fehler
			return "Unvollständiger Blockkommentar"
		}
	}
	// Speichere die Endposition des Inhalts, bevor '*/' übersprungen wird
	endPosition := l.position

	// Überspringe die Zeichen '*/' am Ende des Blockkommentars
	l.readChar() // Überspringe '*'
	l.readChar() // Überspringe '/', um nach dem Kommentar fortzufahren

	// Gib den Text des Kommentars ohne die umgebenden '/*' und '*/' zurück
	return l.input[startPosition:endPosition]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readString() string {
	startPosition := l.position + 1 // Überspringe das Anfangsanführungszeichen
	for {
		l.readChar()
		// Beende die Schleife, wenn ein schließendes Anführungszeichen oder das Dateiende erreicht wird
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	// Gibt den String ohne die Anführungszeichen zurück
	return l.input[startPosition:l.position]
}

func (l *Lexer) NextToken() types.Token {
	var tok types.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = types.Token{Type: types.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(types.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()                         // Gehe zum '=' Zeichen
			literal := string(ch) + string(l.ch) // Kombiniere '!' und '=' zu "!="
			tok = types.Token{Type: types.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(types.ILLEGAL, l.ch) // oder behandele '!' als eigenständiges types.Token, falls erforderlich
		}
	case ';':
		tok = newToken(types.SEMICOLON, l.ch)
	case ',':
		tok = newToken(types.COMMA, l.ch)
	case '(':
		tok = newToken(types.LPAREN, l.ch)
	case ')':
		tok = newToken(types.RPAREN, l.ch)
	case '&':
		tok = newToken(types.AND, l.ch)
	case '{':
		tok = newToken(types.LBRACE, l.ch)
	case '}':
		tok = newToken(types.RBRACE, l.ch)
	case '<':
		tok = newToken(types.LT, l.ch)
	case '>':
		tok = newToken(types.GT, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = types.EOF
	case '/':
		if l.peekChar() == '/' {
			tok.Literal = l.readLineComment()
			tok.Type = types.COMMENT
		} else if l.peekChar() == '*' {
			tok.Literal = l.readBlockComment()
			tok.Type = types.COMMENT
		} else {
			tok = newToken(types.ILLEGAL, l.ch)
		}
	case '.':
		tok = newToken(types.PERIOD, l.ch)
	case '"':
		tok.Type = types.STRING
		tok.Literal = l.readString()
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()                         // Gehe zum '=' Zeichen
			literal := string(ch) + string(l.ch) // Kombiniere ':' und '=' zu ":="
			tok = types.Token{Type: types.ASSIGN_INIT, Literal: literal}
		} else {
			tok = newToken(types.ILLEGAL, l.ch)
		}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = types.INT
			return tok
		} else {
			tok = newToken(types.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) LexTokenList() []types.Token {
	retriveList := make([]types.Token, 0)
	for tok := l.NextToken(); tok.Type != types.EOF; tok = l.NextToken() {
		retriveList = append(retriveList, tok)
	}
	return retriveList
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func LookupIdent(ident string) types.TokenType {
	if tok, ok := types.Keywords[ident]; ok {
		return tok
	}
	return types.IDENT
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tokenType types.TokenType, ch byte) types.Token {
	return types.Token{Type: tokenType, Literal: string(ch)}
}
