package ast

import (
	"bytes"
	"github.com/ibzsy/cardboard/lexer/token"
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

func (program *Program) TokenLiteral() string { return "" }

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
func (es *ExpressionStatement) String() string       { return es.Expression.String() }

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

// Infix Expression
type InfixExpression struct {
	NodeToken token.Token
	Left      Expression
	Operator  string
	Right     Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.NodeToken.TokenLiteral }
func (ie *InfixExpression) String() string {
	var outputString bytes.Buffer
	outputString.WriteString("(")
	outputString.WriteString(ie.Left.String())
	outputString.WriteString(ie.Operator)
	outputString.WriteString(ie.Right.String())
	outputString.WriteString(")")
	return outputString.String()
}

// Box Statement -> box <identifier> <parameter list> <block statement>
type BlockStatement struct {
	NodeToken  token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.NodeToken.TokenLiteral }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")
	return out.String()
}

type BoxExpression struct {
	NodeToken     token.Token
	ParameterList []*Identifier
	Body          *BlockStatement
}

func (box *BoxExpression) expressionNode()      {}
func (box *BoxExpression) TokenLiteral() string { return box.NodeToken.TokenLiteral }
func (box *BoxExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	for _, arg := range box.ParameterList {
		out.WriteString(arg.String() + ",")
	}
	out.WriteString(")")

	out.WriteString("{")
	for _, s := range box.Body.Statements {
		out.WriteString(s.String())
	}
	out.WriteString("}")
	return out.String()
}

type CallExpression struct {
	NodeToken token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.NodeToken.TokenLiteral }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	for _, arg := range ce.Arguments {
		out.WriteString(arg.String())
	}
	out.WriteString(")")
	out.WriteString(ce.Function.String())
	return out.String()
}
