package lexer

import (
	"github.com/ibzsy/cardboard/lexer/token"
	"unicode"
)

type Lexer struct {
	data    string
	curPos  int
	nextPos int
	char    byte
}

func CreateLexer(inputData string) *Lexer {
	lexer := Lexer{data: inputData}
	lexer.readChar()
	return &lexer
}

func (lex *Lexer) NextToken() token.Token {
	var curToken token.Token

	lex.eatWhiteSpace()

	switch lex.char {

	// BRACES AND DELIMITERS
	case '(':
		curToken = token.NewToken(token.LPAREN, "(")
	case ')':
		curToken = token.NewToken(token.RPAREN, ")")
	case '{':
		curToken = token.NewToken(token.LCURLY, "{")
	case '}':
		curToken = token.NewToken(token.RCURLY, "}")
	case ',':
		curToken = token.NewToken(token.COMMA, ",")
	case ';':
		curToken = token.NewToken(token.SCOLON, ";")

	// Arithmetic Operators
	case '+':
		curToken = token.NewToken(token.ADD, "+")
	case '-':
		curToken = token.NewToken(token.SUB, "-")
	case '=':
		curToken = token.NewToken(token.ASSIGN, "=")

	// EOF
	case 0:
		curToken = token.NewToken(token.EOF, "")

	default:
		// isInteger, isLetter return their tokens because
		// theres no need to move the lexer char pointer forwards!
		if isInteger(lex.char) {
			readInteger := lex.readInteger()
			return token.NewToken(token.INT, readInteger)
		} else if isLetter(lex.char) {
			readIdentifier := lex.readIdentifier()
			return token.NewToken(token.GetIdentifierType(readIdentifier), readIdentifier)
		} else {
			// Unknown token
			curToken = token.NewToken(token.UNKNOWN, string(lex.char))
		}
	}
	lex.readChar()
	return curToken
}

func (lex *Lexer) readChar() {
	if lex.nextPos >= len(lex.data) {
		lex.char = 0
	} else {
		lex.char = lex.data[lex.nextPos]
	}
	lex.curPos = lex.nextPos
	lex.nextPos++
}

func (lex *Lexer) readIdentifier() string {
	startPos := lex.curPos
	for isLetter(lex.char) {
		lex.readChar()
	}
	return string(lex.data[startPos:lex.curPos])
}

func (lex *Lexer) readInteger() string {
	startPos := lex.curPos
	for isInteger(lex.char) {
		lex.readChar()
	}
	return string(lex.data[startPos:lex.curPos])
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch))
}

func isInteger(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

func (lex *Lexer) eatWhiteSpace() {
	for unicode.IsSpace(rune(lex.char)) {
		lex.readChar()
	}
}
