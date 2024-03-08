package main

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers + literals
	IDENT    TokenType = "IDENT"  // add, foobar, x, y, ...
	INT      TokenType = "INT"    // 123456
	STRING   TokenType = "STRING" // "foobar"
	TRUE     TokenType = "TRUE"   // True
	FALSE    TokenType = "FALSE"  // False
	NULL     TokenType = "NULL"   // Null
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	// Operatoren
	ASSIGN      TokenType = "="
	EQ          TokenType = "=="
	ASSIGN_INIT TokenType = ":="

	// Delimiters
	PERIOD    TokenType = "."
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	LPAREN    TokenType = "("
	RPAREN    TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
	LT        TokenType = "<"
	GT        TokenType = ">"

	// Schlüsselwörter
	FUNCTION   TokenType = "FUNCTION"
	LET        TokenType = "LET"
	IF         TokenType = "IF"
	ELSE       TokenType = "ELSE"
	RETURN     TokenType = "RETURN"
	CONST      TokenType = "CONST"
	ISNULL     TokenType = "ISNULL"
	THIRPF     TokenType = "THIRPF"
	RBLOCKCALL TokenType = "RBLOCKCALL"
	CATCH      TokenType = "CATCH"
	FINAL      TokenType = "FINAL"
	SWITCH     TokenType = "SWITCH"
	CASE       TokenType = "CASE"
	DEFAULT    TokenType = "DEFAULT"
	READSTORE  TokenType = "READSTORE"
	PACKAGE    TokenType = "PACKAGE"

	// Kommentare
	COMMENT TokenType = "COMMENT"
)

var keywords = map[string]TokenType{
	"const":      CONST,
	"isnull":     ISNULL,
	"thirpf":     THIRPF,
	"rblockcall": RBLOCKCALL,
	"catch":      CATCH,
	"final":      FINAL,
	"switch":     SWITCH,
	"case":       CASE,
	"default":    DEFAULT,
	"readStore":  READSTORE,
	"package":    PACKAGE,
}

type Token struct {
	Type    TokenType
	Literal string
}

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

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ASSIGN, l.ch)
		}
	case ';':
		tok = newToken(SEMICOLON, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case '{':
		tok = newToken(LBRACE, l.ch)
	case '}':
		tok = newToken(RBRACE, l.ch)
	case '<':
		tok = newToken(LT, l.ch)
	case '>':
		tok = newToken(GT, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	case '/':
		if l.peekChar() == '/' {
			tok.Literal = l.readLineComment()
			tok.Type = COMMENT
		} else if l.peekChar() == '*' {
			tok.Literal = l.readBlockComment()
			tok.Type = COMMENT
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	case '.':
		tok = newToken(PERIOD, l.ch)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case ':':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()                         // Gehe zum '=' Zeichen
			literal := string(ch) + string(l.ch) // Kombiniere ':' und '=' zu ":="
			tok = Token{Type: ASSIGN_INIT, Literal: literal}
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = INT
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}
