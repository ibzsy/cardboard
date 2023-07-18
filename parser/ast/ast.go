package ast

import (
	"bytes"
	"cardboard/lexer/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Every Cardboard Program is a list of statements,
// therefore the AST Root Node is the list of statements of the program
type Program struct {
	Statements []Statement
}

func (program *Program) String() string {
	var outputString bytes.Buffer
	for _, stm := range program.Statements {
		outputString.WriteString(stm.String())
	}
	return outputString.String()
}

// Identifiers are Expressions.
type Identifier struct {
	NodeToken token.Token
	Value     string
}

func (ident *Identifier) expressionNode()      {}
func (ident *Identifier) TokenLiteral() string { return ident.Value }
func (ident *Identifier) String() string       { return ident.Value }

// 'put' statement
// put <identifier> = <expression>
type PutStatement struct {
	NodeToken      token.Token
	NodeIdentifier Identifier
	NodeExpression Expression
}

func (p *PutStatement) statementNode()       {}
func (p *PutStatement) TokenLiteral() string { return p.NodeToken.TokenLiteral }

// Helps during debugging to observe what the Node represents
func (p *PutStatement) String() string {
	var outputString bytes.Buffer
	outputString.WriteString(p.TokenLiteral() + " " + p.NodeIdentifier.Value + " = ")
	if p.NodeExpression != nil {
		outputString.WriteString(p.NodeExpression.String())
	}
	outputString.WriteString(";")
	return outputString.String()
}

// 'unbox' statement. Basically return statement!
// unbox <expression>;
type UnboxStatement struct {
	NodeToken      token.Token
	NodeExpression Expression
}

func (u *UnboxStatement) statementNode()       {}
func (u *UnboxStatement) TokenLiteral() string { return u.NodeToken.TokenLiteral }
func (u *UnboxStatement) String() string {
	var outputString bytes.Buffer
	outputString.WriteString(u.NodeToken.TokenLiteral + " ")
	if u.NodeExpression != nil {
		outputString.WriteString(u.NodeExpression.String())
	}
	outputString.WriteString(";")
	return outputString.String()
}

// Expression Statements link identifiers to expressions

type ExpressionStatement struct {
	NodeToken  token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.NodeToken.TokenLiteral }
func (es *ExpressionStatement) String() string       { return es.NodeToken.TokenLiteral }

// Integers are just expressions as well

type IntegerLiteral struct {
	NodeToken token.Token
	Value     int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.NodeToken.TokenLiteral }
func (il *IntegerLiteral) String() string       { return il.NodeToken.TokenLiteral }

// Prefix Expressions

type PrefixExpression struct {
	NodeToken token.Token
	Operator  string
	Right     Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.NodeToken.TokenLiteral }
func (pe *PrefixExpression) String() string {
	var outputString bytes.Buffer
	outputString.WriteString("(")
	outputString.WriteString(pe.Operator)
	outputString.WriteString(pe.Right.String())
	outputString.WriteString(")")
	return outputString.String()
}
