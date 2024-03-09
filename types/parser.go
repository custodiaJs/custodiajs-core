package types

type Token struct {
	Type    TokenType
	Literal string
}

type Statement interface {
	statementNode()
	String() string
}

type Expression interface {
	expressionNode()
	String() string
}

type Program struct {
	Statements []Statement
}
