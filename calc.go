package main

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNumber
	TokenPlus
	TokenMinus
	TokenMultiply
	TokenDivide
	TokenLeftParen
	TokenRightParen
)

type Token struct {
	typ   TokenType
	value string
}

type Lexer struct {
	input string
	pos   int
	ch    rune
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, pos: 0}
	l.ch = l.nextChar()
	return l
}

func (l *Lexer) nextChar() rune {
	if l.pos >= len(l.input) {
		l.ch = 0
		return 0
	}
	ch := rune(l.input[l.pos])
	l.pos++
	return ch
}

func (l *Lexer) NextToken() Token {
	for unicode.IsSpace(l.ch) {
		l.ch = l.nextChar()
	}
	switch {
	case unicode.IsDigit(l.ch):
		start := l.pos - 1
		for unicode.IsDigit(l.ch) || l.ch == '.' {
			l.ch = l.nextChar()
		}
		return Token{typ: TokenNumber, value: l.input[start : l.pos-1]}
	case l.ch == '+':
		l.ch = l.nextChar()
		return Token{typ: TokenPlus, value: "+"}
	case l.ch == '-':
		l.ch = l.nextChar()
		return Token{typ: TokenMinus, value: "-"}
	case l.ch == '*':
		l.ch = l.nextChar()
		return Token{typ: TokenMultiply, value: "*"}
	case l.ch == '/':
		l.ch = l.nextChar()
		return Token{typ: TokenDivide, value: "/"}
	case l.ch == '(':
		l.ch = l.nextChar()
		return Token{typ: TokenLeftParen, value: "("}
	case l.ch == ')':
		l.ch = l.nextChar()
		return Token{typ: TokenRightParen, value: ")"}
	case l.ch == 0:
		return Token{typ: TokenEOF, value: ""}
	}

	return Token{typ: TokenEOF, value: ""}
}

type Parser struct {
	lexer   *Lexer
	current Token
}

func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.current = p.lexer.NextToken()
}

func (p *Parser) Parse() (float64, error) {
	result, err := p.expression()
	if err != nil {
		return 0, err
	}
	if p.current.typ != TokenEOF {
		return 0, errors.New("invalid expression")
	}
	return result, nil
}

func (p *Parser) expression() (float64, error) {
	result, err := p.term()
	if err != nil {
		return 0, err
	}

	for p.current.typ == TokenPlus || p.current.typ == TokenMinus {
		op := p.current
		p.nextToken()
		nextTerm, err := p.term()
		if err != nil {
			return 0, err
		}

		if op.typ == TokenPlus {
			result += nextTerm
		} else {
			result -= nextTerm
		}
	}

	return result, nil
}

func (p *Parser) term() (float64, error) {
	result, err := p.factor()
	if err != nil {
		return 0, err
	}

	for p.current.typ == TokenMultiply || p.current.typ == TokenDivide {
		op := p.current
		p.nextToken()
		nextFactor, err := p.factor()
		if err != nil {
			return 0, err
		}

		if op.typ == TokenMultiply {
			result *= nextFactor
		} else {
			if nextFactor == 0 {
				return 0, errors.New("division by zero")
			}
			result /= nextFactor
		}
	}

	return result, nil
}

func (p *Parser) factor() (float64, error) {
	var result float64
	var err error

	if p.current.typ == TokenNumber {
		result, err = strconv.ParseFloat(p.current.value, 64)
		if err != nil {
			return 0, err
		}
		p.nextToken()
	} else if p.current.typ == TokenLeftParen {
		p.nextToken()
		result, err = p.expression()
		if err != nil {
			return 0, err
		}
		if p.current.typ != TokenRightParen {
			return 0, errors.New("missing closing parenthesis")
		}
		p.nextToken()
	} else {
		return 0, errors.New("invalid expression")
	}

	return result, nil
}

func Calc(expression string) (float64, error) {
	lexer := NewLexer(expression)
	parser := NewParser(lexer)
	return parser.Parse()
}

func main() {
	expr := "3 + 5 * (2 - 8)"
	result, err := Calc(expr)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Result: %f\n", result)
	}
}
