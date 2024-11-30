package parser

import (
	"github.com/ibzsy/cardboard/lexer"
	"github.com/ibzsy/cardboard/lexer/token"
	"github.com/ibzsy/cardboard/parser/ast"
	"fmt"
	"strconv"
)

type Parser struct {
	lexer       *lexer.Lexer
	curToken    token.Token
	peekToken   token.Token
	errors      []string
	prefixFuncs map[token.TokenType]prefixFunc
	infixFuncs  map[token.TokenType]infixFunc
}

// Precedence For Operators < Not all supported right now! >
const (
	_ int = iota
	LOWEST
	EQUALS  // == LESSGREATER // > or <
	SUM     //+
	PRODUCT //*
	PREFIX  //-X or !X
	CALL    // myFunction(X)
)

// Precedence mapping
var precedence = map[token.TokenType]int{
	token.ADD:    SUM,
	token.SUB:    SUM,
	token.LPAREN: CALL,
}

func CreateParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}

	// Need to initialize both Tokens Pointers
	p.nextToken()
	p.nextToken()

	// Instantiate Mapping
	// Prefix
	p.prefixFuncs = make(map[token.TokenType]prefixFunc)
	p.setPrefixFunction(token.IDENTIFIER, p.parseIdentifier)
	p.setPrefixFunction(token.INT, p.parseIntegerLiteral)
	p.setPrefixFunction(token.SUB, p.parsePrefixExpression)
	p.setPrefixFunction(token.ADD, p.parsePrefixExpression)
	p.setPrefixFunction(token.LPAREN, p.parseGroupedExpression)
	p.setPrefixFunction(token.BOX, p.parseBoxStatement)

	// Infix
	p.infixFuncs = make(map[token.TokenType]infixFunc)
	p.setInfixFunction(token.ADD, p.parseInfixExpression)
	p.setInfixFunction(token.SUB, p.parseInfixExpression)
	p.setInfixFunction(token.LPAREN, p.parseCallExpression)

	return p
}

// Parses Cardboard Program
func (p *Parser) ParseCardBoard() *ast.Program {
	program := ast.Program{}

	for !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.UNKNOWN) {
			p.addError(fmt.Sprintf("Unknown Token: %s", p.curToken.TokenLiteral))
			return &ast.Program{}
		}
		stmt := p.parseStatement()
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}
	return &program
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.TokenType {
	case token.PUT:
		return p.parsePutStatement()
	case token.UNBOX:
		return p.parseUnboxStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parsePutStatement() *ast.PutStatement {
	putStmt := &ast.PutStatement{}

	// Parse Put Token
	putStmt.NodeToken = p.curToken

	// Parse Identifier
	if !p.expectPeek(token.IDENTIFIER) {
		p.typeError(token.IDENTIFIER, p.peekToken.TokenType)
		return nil
	}

	putStmt.NodeIdentifier = ast.Identifier{NodeToken: p.curToken, Value: p.curToken.TokenLiteral}

	// Ensure Next Token is Assign
	if !p.expectPeek(token.ASSIGN) {
		p.typeError(token.ASSIGN, p.peekToken.TokenType)
		return nil
	}

	p.nextToken()
	putStmt.NodeExpression = p.parseExpression(LOWEST)

	// At this point peek token should be semi colon!
	for !p.expectPeek(token.SCOLON) {
		p.addError("Error. Expected <;> at the end of the box statement.")
		return nil
	}

	return putStmt
}

func (p *Parser) parseUnboxStatement() *ast.UnboxStatement {
	unboxStmt := &ast.UnboxStatement{}

	// Parse Unbox Token
	unboxStmt.NodeToken = p.curToken

	p.nextToken()
	unboxStmt.NodeExpression = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SCOLON) {
		p.addError("Error. Expected <;> at the end of UNBOX statement.")
		return nil
	}

	return unboxStmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expStmt := &ast.ExpressionStatement{NodeToken: p.curToken}
	expStmt.Expression = p.parseExpression(LOWEST)

	// Optional Semi-colon
	if p.peekTokenIs(token.SCOLON) {
		p.nextToken()
	}
	return expStmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixFuncs[p.curToken.TokenType]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("Unknown character: <%s>", p.curToken.TokenLiteral))
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SCOLON) && precedence < p.peekPrecedence() {
		infix := p.infixFuncs[p.peekToken.TokenType]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		NodeToken: p.curToken,
		Value:     p.curToken.TokenLiteral,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	val, err := strconv.ParseInt(p.curToken.TokenLiteral, 10, 0)
	if err != nil {
		error := fmt.Sprintf("Error. Couldn't Parse Integer From String = <%s>", p.curToken.TokenLiteral)
		p.errors = append(p.errors, error)
		return nil
	}
	return &ast.IntegerLiteral{NodeToken: p.curToken, Value: val}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		NodeToken: p.curToken,
		Operator:  p.curToken.TokenLiteral,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		NodeToken: p.curToken,
		Left:      left,
		Operator:  p.curToken.TokenLiteral,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		err := fmt.Sprintf("Error. Expected <)>. Found <%s> instead.", p.peekToken.TokenLiteral)
		p.addError(err)
		return nil
	}
	return expr
}

func (p *Parser) parseBoxStatement() ast.Expression {
	box := &ast.BoxExpression{NodeToken: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		p.addError(fmt.Sprintf("Error. Expected parameter list after function name. Found <%s>", p.peekToken.TokenType))
		return nil
	}

	box.ParameterList = p.parseFunctionParameters()

	if !p.expectPeek(token.LCURLY) {
		p.addError("Error. Expected function block statement after parameter list.")
		return nil
	}

	box.Body = p.parseBlockStatement()

	if !p.curTokenIs(token.RCURLY) {
		p.addError(fmt.Sprintf("Error. Expected block statement closure! Got <%s>", p.peekToken.TokenType))
		return nil
	}

	return box
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.NodeToken = p.curToken

	p.nextToken()

	for !p.curTokenIs(token.RCURLY) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		block.Statements = append(block.Statements, stmt)
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	list := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return list
	}

	p.nextToken()

	list = append(list, &ast.Identifier{NodeToken: p.curToken, Value: p.curToken.TokenLiteral})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, &ast.Identifier{NodeToken: p.curToken, Value: p.curToken.TokenLiteral})
	}

	if !p.expectPeek(token.RPAREN) {
		p.addError(fmt.Sprintf("Error. Expected parameter list closure. Got <%s>", p.peekToken.TokenLiteral))
		return nil
	}

	return list
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{NodeToken: p.curToken, Function: function}
	expr.Arguments = p.parseCallArguments()
	return expr
}

func (p *Parser) parseCallArguments() []ast.Expression {
	arguments := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return arguments
	}

	p.nextToken()
	arguments = append(arguments, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		arguments = append(arguments, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		p.addError(fmt.Sprintf("Error. No closure to call expression argument list. Got <%s>", p.peekToken.TokenLiteral))
		return nil
	}
	return arguments
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedence[p.peekToken.TokenType]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if precedence, ok := precedence[p.curToken.TokenType]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.TokenType == t
}
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.TokenType == t
}

func (p *Parser) GetErrors() []string {
	return p.errors
}

// In the case where the statement is invalid, we'll
// need to skip it!
func (p *Parser) skipStatement() {
	for !p.curTokenIs(token.SCOLON) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}
}

func (p *Parser) typeError(expectedType token.TokenType, gotType token.TokenType) {
	p.errors = append(p.errors,
		fmt.Sprintf("Error. Expected Token Type <%s>. Got Token Type <%s>.\n", expectedType, gotType))
	p.skipStatement()
}

// Parser Functions -> Certain token types are associated to infix and prefix operations.
type (
	infixFunc  func(ast.Expression) ast.Expression
	prefixFunc func() ast.Expression
)

func (p *Parser) setInfixFunction(token token.TokenType, function infixFunc) {
	p.infixFuncs[token] = function
}

func (p *Parser) setPrefixFunction(token token.TokenType, function prefixFunc) {
	p.prefixFuncs[token] = function
}

// Helper function to help to register errors
func (p *Parser) addError(err string) {
	p.errors = append(p.errors, err)
}
