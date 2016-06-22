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
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Value string
type TokenType int

type Position struct {
	Index  int
	Row    int
	Column int
}

const (
	TokenEOS    TokenType = TokenType('.')
	TokenComma  TokenType = TokenType(',')
	TokenError  TokenType = -1
	TokenWord   TokenType = -2
	TokenEOL    TokenType = -3
	TokenEOF    TokenType = -4
	TokenNumber TokenType = -5
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
	runes   []rune
}

type LexerRule func(lexer *Lexer) LexerRule

func (me *Lexer) Emit(t TokenType, v Value) {
	tok := Token{t, v, me.Current}
	me.Output <- tok
}

func (me *Lexer) Error(message string, args ...interface{}) {
	value := fmt.Sprintf(message, args...)
	me.Emit(TokenError, Value(value))
}

func LexError(me *Lexer) LexerRule {
	me.Error("Error")
	return nil
}

func (me *Lexer) SkipComment() bool {
	if me.Peek() == '#' {
		if me.Next() == '(' {
			return me.SkipNotIn(")")
		} else {
			return me.SkipNotIn("\r\n")
		}
	}
	return true
}

func LexWord(me *Lexer) LexerRule {
	me.Found(TokenWord)
	return LexNormal
}

func LexNumber(me *Lexer) LexerRule {
	me.Found(TokenNumber)
	return LexNormal
}

func LexComment(me *Lexer) LexerRule {
	if !me.SkipComment() {
		me.Error("Unterminated comment")
		return LexError
	}
	me.Advance()
	return LexNormal
}

func LexEOS(me *Lexer) LexerRule {
	me.Found(TokenEOS)
	return LexNormal
}

func LexEOL(me *Lexer) LexerRule {
	me.Found(TokenEOL)
	return LexNormal
}

func LexNormal(me *Lexer) LexerRule {
	me.SkipWhitespace()
	peek := me.Peek()
	if peek == '#' {
		return LexComment
	} else if peek == '.' {
		return LexEOS
	} else if peek == '\n' || peek == '\r' {
		return LexEOL
	} else if unicode.IsLetter(me.Peek()) {
		return LexWord
	} else if unicode.IsDigit(me.Peek()) {
		return LexNumber
	}

	return nil
}

func OpenLexer(reader io.Reader) *Lexer {
	lexer := &Lexer{}
	lexer.Reader = reader
	lexer.Output = make(TokenChannel)
	// lexer.buffer = new(byte[1024])
	return lexer
}

func (me *Lexer) ReadReaderOnce() (bool, error) {
	buffer := make([]byte, 1024)

	n, err := me.Reader.Read(buffer)
	fmt.Printf("read %v %d %v\n", buffer[:n], n, err)
	if n > 0 {
		me.buffer = append(me.buffer, buffer[:n]...)
		fmt.Printf("append  %s", me.buffer)
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

func (me *Lexer) ReadReader() bool {
	me.buffer = make([]byte, 0)
	more, err := me.ReadReaderOnce()
	for err == nil && more {
		more, err = me.ReadReaderOnce()
	}
	me.runes = bytes.Runes(me.buffer)

	return err != nil && err != io.EOF
}

func (me *Lexer) Peek() rune {
	return me.runes[me.Current.Index]
}

func (me *Lexer) PeekNext() rune {
	if (me.Current.Index) >= len(me.runes) {
		return '\000'
	}
	return me.runes[me.Current.Index+1]
}

func (me *Lexer) Next() rune {
	if me.Peek() == '\n' {
		me.Current.Column = 0
		me.Current.Row++
	}
	me.Current.Index++
	if me.Current.Index >= len(me.runes) {
		me.Emit(TokenEOF, "")
	}
	return me.Peek()
}

func (me *Lexer) Previous() rune {
	if me.Current.Index > 0 {
		me.Current.Index--

		if me.Peek() == '\n' {
			me.Current.Column = 0
			me.Current.Row++
		}
	}
	return me.Peek()
}

func (me *Lexer) SkipRune() {
	_ = me.Next
}

func (me *Lexer) SkipIn(set string) bool {
	_ = me.Next
	for strings.ContainsRune(set, me.Peek()) {
		if me.Next() == '\000' {
			return false
		}
	}
	return true
}

func (me *Lexer) SkipNotIn(set string) bool {
	_ = me.Next
	for !strings.ContainsRune(set, me.Peek()) {
		if me.Next() == '\000' {
			return false
		}
	}
	return true
}

func (me *Lexer) SkipWhile(should_skip func(r rune) bool) bool {
	_ = me.Next
	for should_skip(me.Peek()) {
		if me.Next() == '\000' {
			return false
		}
	}
	return true
}

func (me *Lexer) SkipWhitespace() {
	me.SkipIn(" \t")
}

func (me *Lexer) Advance() {
	me.Last = me.Current
}

func (me *Lexer) Retry() {
	me.Current = me.Last
}

func (me *Lexer) Found(kind TokenType) {
	value := me.runes[me.Last.Index:me.Current.Index]
	svalue := string(value)
	me.Emit(kind, Value(svalue))
	me.Advance()
}

func (me *Lexer) Start() {
	if me.ReadReader() {
		rule := LexNormal
		for rule != nil {
			rule = rule(me)
		}
	}
	close(me.Output)
}

func (me *Lexer) TryLexing() {
	go me.Start()

	for token := range me.Output {
		fmt.Println("Token %s", token)
	}
}

type Parser struct {
	Lexer
}

type Environment struct {
	Parent *Environment
}

func main() {
	fmt.Println("Hello World!")
}
