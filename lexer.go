package main

import (
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"
)


var EOF = errors.New("End of input reached")
var UnknownTokenError = errors.New("Unknown token!")

type tokenType int

const (
	let tokenType = iota
	identifier
	number
	str
	plus
	minus
	eq
	semicolon
)

type keyword string

const (
	letKeyword keyword = "let"
)

type token struct {
	value     string
	tokenType tokenType
}

type Lexer struct {
	Input string
	pos int
	width int
}

func (l *Lexer) next() (rune, error) {
	if l.pos >= len(l.Input) {
		return -1, EOF
	}
	r, width := utf8.DecodeRuneInString(l.Input[l.pos:])
	l.width = width
	l.pos += width
	return r, nil
}

func (l *Lexer) peek() (rune, error) {
	r, err := l.next()
	if err != nil {
		return r, err
	}
	l.backup()
	return r, nil
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) tokenize() ([]token, error) {
	tokens := make([]token, 0)
	for {
		if err := l.skipWhiteSpace(); err != nil {
			return tokens, err
		}
		r, err := l.peek()
		if err != nil {
			return tokens, err
		}

		switch {
		case r == '=':
			tokens = append(tokens, token{value: string(r), tokenType: eq})
			l.next()
		case r == '+':
			tokens = append(tokens, token{value: string(r), tokenType: plus})
			l.next()
		case r == '-':
			tokens = append(tokens, token{value: string(r), tokenType: minus})
			l.next()
		case r == ';':
			tokens = append(tokens, token{value: string(r), tokenType: semicolon})
			l.next()
		case r == '"':
			str, err := l.readString()
			if err != nil {
				return tokens, err
			}
			tokens = append(tokens, str)
		case unicode.IsDigit(r) || r == '.':
			num, err := l.readNum()
			if err != nil {
				return tokens, err
			}
			tokens = append(tokens, num)
		case unicode.IsLetter(r):
			ident, err := l.readIdentOrKeyword()
			if err != nil {
				return tokens, err
			}
			tokens = append(tokens, ident)
		default:
			return tokens, UnknownTokenError
		}
	}
}

func (l *Lexer) skipWhiteSpace() error {
	for {
		r, err := l.peek()
		if err != nil {
			return err
		}
		if !unicode.IsSpace(r) {
			return nil
		}
		l.next()
	}
}

func (l *Lexer) readString() (token, error) {
	// consume open quote
	if _, err := l.next(); err != nil {
		return token{}, err
	}
	start := l.pos
	for {
		r, err := l.next()
		if err != nil {
			return token{}, err
		}
		if r == '"' {
			break
		}
	}
	value := l.Input[start:l.pos-1]
	return token{value: value, tokenType: str}, nil
}

func (l *Lexer) readNum() (token, error) {
	start := l.pos
	for {
		r, err := l.next()
		if err != nil {
			return token{}, err
		}
		if !unicode.IsDigit(r) && r != '.' {
			break
		}
	}
	value := l.Input[start:l.pos]
	return token{value: value, tokenType: number}, nil
}

func (l *Lexer) readIdentOrKeyword() (token, error) {
	start := l.pos
	for {
		r, err := l.next()
		if err != nil {
			return token{}, err
		}
		if !unicode.IsLetter(r) && r != '_' {
			break
		}
	}
	value := l.Input[start:l.pos]
	switch {
	case value == string(letKeyword):
		return token{value: value, tokenType: let}, nil
	default:
		return token{value: value, tokenType: identifier}, nil
	}
}


func main() {
	lexer := Lexer{
		Input: `
			println + 420 69;
			let sayHello a b = printf "Hi, %s!" a;
			sayHello "world";
		`,
	}

	fmt.Println(lexer.tokenize())
}
