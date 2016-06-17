// raku

/* Raku is an easy to use scripting language that can also be used easily interactively

Syntax (verified LL(1) )

PROGRAM -> STATEMENTS .
STATEMENTS -> STATEMENT STATEMENTS | .
STATEMENT -> EXPRESSION | BLOCK | EMPTY_LINE | comment .
EXPRESSION -> VALUE PARAMETERS NL.
PARAMETERS_NONEMPTY -> PARAMETER PARAMETERS.
PARAMETERS-> PARAMETERS_NONEMPTY | .
PARAMETER -> BLOCK | VALUE .
EMPTY_LINE -> NL .
BLOCK -> ob STATEMENTS cb | op STATEMENTS cp | oa STATEMENTS ca.
NL -> nl | semicolon .
VALUE -> string | float | integer | symbol .

Lexer:


*/
package raku

import (
	"fmt"
	"io"
)

type Value string
type TokenType int

type Position struct {
	Index  int
	Row    int
	Column int
}

const (
	TokenError TokenType = iota
	TokenEOF
)

type Token struct {
	TokenType
	Value
	Position
}

func (me Token) String() string {
	return fmt.Sprintf("Token: %d >%s< %d %d %d.", me.TokenType, string(me.Value), me.Index, me.Row, me.Column)
}

type TokenChannel chan Token

type Lexer struct {
	Reader  io.Reader
	Current Position
	Last    Position
	Token   Token
	rule    LexerRule
	Output  TokenChannel
	buffer  []byte
}

type LexerRule func(lexer *Lexer) LexerRule

func (lexer *Lexer) Emit(t TokenType, v Value) {
	tok := Token{t, v, lexer.Current}
	lexer.Output <- tok
}

func (lexer *Lexer) Error(message string, args ...interface{}) {
	value := fmt.Sprintf(message, args...)
	lexer.Emit(TokenError, Value(value))
}

func LexError(lexer *Lexer) LexerRule {
	lexer.Error("Error")
	return nil
}

func LexNormal(lexer *Lexer) LexerRule {
	return LexError
}

func OpenLexer(reader io.Reader) *Lexer {
	lexer := &Lexer{}
	lexer.Reader = reader
	lexer.Output = make(TokenChannel)
	// lexer.buffer = new(byte[1024])
	return lexer
}

func (me *Lexer) ReadReader() (bool, error) {
	buffer := make([]byte, 1024)

	n, err := me.Reader.Read(buffer)
	if n > 0 {
		me.buffer = append(me.buffer, buffer...)
	}
	if err == io.EOF {
		me.Emit(TokenEOF, "")
		return true, nil
	} else if err != nil {
		me.Error("Error reading from reader: %s", err)
		return true, err
	}
	return false, nil
}

func (me *Lexer) Start() {
	more, err := me.ReadReader()
	for err == nil && more {
		more, err = me.ReadReader()
	}

	if err != nil {
		return
	}

	rule := LexNormal
	for rule != nil {
		rule = rule(me)
	}

	close(me.Output)
}

/*
func (me *Lexer) TryLexing() {
	go {
		me.Start()
	}

	for token := range me.Output {
		fmt.Println("Token %s", token)
	}
}
*/

type Parser struct {
	Lexer
}

type Environment struct {
	Parent *Environment
}

func main() {
	fmt.Println("Hello World!")
}
