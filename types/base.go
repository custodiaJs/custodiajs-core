package types

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
	NOT_EQ      TokenType = "!="
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
	SWITCH     TokenType = "SWITCH"
	CASE       TokenType = "CASE"
	DEFAULT    TokenType = "DEFAULT"
	READSTORE  TokenType = "READSTORE"

	// Kommentare
	COMMENT TokenType = "COMMENT"
)

var Keywords = map[string]TokenType{
	"const":      CONST,
	"isnull":     ISNULL,
	"thirpf":     THIRPF,
	"rblockcall": RBLOCKCALL,
	"catch":      CATCH,
	"switch":     SWITCH,
	"case":       CASE,
	"default":    DEFAULT,
	"readStore":  READSTORE,
	"if":         IF,
	"else":       ELSE,
}
