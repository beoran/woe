// raku

/* Raku is an easy to use scripting language that can also be used easily interactively

Syntax (verified LL(1) on smlweb.cpsc.ucalgary.ca)

PROGRAM -> STATEMENTS.
STATEMENTS -> STATEMENT STATEMENTS | .
STATEMENT -> DEFINITION | EXPRESSION | BLOCK .
DEFINITION -> to WORDS BLOCK.
WORDS -> word WORDS | .
EXPRESSION -> WORD_EXPRESSION | VALUE_EXPRESSION.
WORD_EXPRESSION -> word WORD_CALLOP.
WORD_CALLOP -> WORD_OPERATION | WORD_CALL.
WORD_OPERATION -> operator PARAMETERS_NONEMPTY EOX.
WORD_CALL -> PARAMETERS EOX.
VALUE_EXPRESSION -> value VALUE_CALLOP.
VALUE_CALLOP -> VALUE_OPERATION | VALUE_CALL.
VALUE_OPERATION -> operator PARAMETERS_NONEMPTY EOX.
VALUE_CALL -> EOX.
PARAMETERS_NONEMPTY -> PARAMETER PARAMETERS.
PARAMETERS -> PARAMETERS_NONEMPTY | .
PARAMETER -> BLOCK | WORDVALUE .
BLOCK -> ob STATEMENTS cb | op STATEMENTS cp | oa STATEMENTS ca | do
STATEMENTS end.
WORDVALUE -> word | VALUE.
VALUE -> string | number | symbol.
EOX -> eol | period.

Lexer:


*/
package raku

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"unicode"

	"github.com/beoran/woe/monolog"
	"github.com/beoran/woe/tree"
)

type Value string
type TokenType int64

type Position struct {
	Index  int
	Row    int
	Column int
}

const (
	TokenEOS          TokenType = TokenType('.')
	TokenComma        TokenType = TokenType(',')
	TokenSemicolumn   TokenType = TokenType(';')
	TokenColumn       TokenType = TokenType(':')
	TokenOpenParen    TokenType = TokenType('(')
	TokenCloseParen   TokenType = TokenType(')')
	TokenOpenBrace    TokenType = TokenType('{')
	TokenCloseBrace   TokenType = TokenType('}')
	TokenOpenBracket  TokenType = TokenType('[')
	TokenCloseBracket TokenType = TokenType(']')

	TokenNone     TokenType = 0
	TokenError    TokenType = -1
	TokenWord     TokenType = -2
	TokenEOL      TokenType = -3
	TokenEOF      TokenType = -4
	TokenNumber   TokenType = -5
	TokenOperator TokenType = -6
	TokenString   TokenType = -7
	TokenKeyword  TokenType = -8
	TokenLast     TokenType = -9
)

type Token struct {
	TokenType
	Value
	Position
}

var tokenTypeNames []string = []string{
	"TokenNone", "TokenError", "TokenWord", "TokenEOL", "TokenEOF", "TokenNumber", "TokenOperator", "TokenString", "TokenKeyword",
}

var keywordList []string = []string{
	"a", "do", "end", "the", "to",
}

func (me TokenType) String() string {
	if int(me) > 0 {
		return fmt.Sprintf("Token %c", rune(me))
	} else if me > TokenLast {
		return tokenTypeNames[-int(me)]
	} else {
		return fmt.Sprintf("Unknown Token %d", int(me))
	}

}

func (me Token) String() string {
	return fmt.Sprintf("Token: %s >%s< %d %d %d.", me.TokenType, string(me.Value), me.Index, me.Row, me.Column)
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
	monolog.Error("Lex Error: %s", value)
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

func IsKeyword(word string) bool {
	i := sort.SearchStrings(keywordList, word)
	if i >= len(keywordList) {
		return false
	}
	return word == keywordList[i]
}

func LexWord(me *Lexer) LexerRule {
	me.SkipNotIn(" \t\r\n'")
	if IsKeyword(me.CurrentStringValue()) {
		me.Found(TokenKeyword)
	} else {
		me.Found(TokenWord)
	}
	return LexNormal
}

func LexNumber(me *Lexer) LexerRule {
	me.SkipNotIn(" \tBBBT\r\n")
	me.Found(TokenNumber)
	return LexNormal
}

func LexWhitespace(me *Lexer) LexerRule {
	me.SkipWhitespace()
	me.Advance()
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

func LexPunctuator(me *Lexer) LexerRule {
	me.Found(TokenType(me.Peek()))
	return LexNormal
}

func LexEOL(me *Lexer) LexerRule {
	me.SkipIn("\r\n")
	me.Found(TokenEOL)
	return LexNormal
}

func LexOperator(me *Lexer) LexerRule {
	me.SkipNotIn(" \t\r\n")
	me.Found(TokenOperator)
	return LexNormal
}

func lexEscape(me *Lexer) error {
	_ = me.Next()
	return nil
}

func LexString(me *Lexer) LexerRule {
	open := me.Peek()
	do_escape := open == '"'
	peek := me.Next()
	me.Advance()
	for ; peek != '\000'; peek = me.Next() {
		if do_escape && peek == '\\' {
			if err := lexEscape(me); err != nil {
				return LexError
			}
		} else if peek == open {
			me.Found(TokenString)
			_ = me.Next()
			me.Advance()
			return LexNormal
		}
	}
	me.Error("Unexpected EOF in string.")
	return nil
}

func LexNumberOrOperator(me *Lexer) LexerRule {
	if unicode.IsDigit(me.Next()) {
		return LexNumber
	} else {
		_ = me.Previous()
		return LexOperator
	}
}

func LexNormal(me *Lexer) LexerRule {
	peek := me.Peek()
	if peek == '#' {
		return LexComment
	} else if strings.ContainsRune(" \t", peek) {
		return LexWhitespace
	} else if strings.ContainsRune(".,;:", peek) {
		return LexPunctuator
	} else if strings.ContainsRune("\r\n", peek) {
		return LexEOL
	} else if strings.ContainsRune("+-", peek) {
		return LexNumberOrOperator
	} else if strings.ContainsRune("\"`", peek) {
		return LexString
	} else if peek == '\000' {
		me.Emit(TokenEOF, "")
		return nil
	} else if unicode.IsLetter(peek) {
		return LexWord
	} else if unicode.IsDigit(peek) {
		return LexNumber
	} else {
		return LexOperator
	}
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
	monolog.Debug("read %v %d %v\n", buffer[:n], n, err)
	if n > 0 {
		me.buffer = append(me.buffer, buffer[:n]...)
		monolog.Debug("append  %s", me.buffer)
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

func (me *Lexer) ReadReader() error {
	me.buffer = make([]byte, 0)
	more, err := me.ReadReaderOnce()
	for err == nil && more {
		more, err = me.ReadReaderOnce()
	}
	me.runes = bytes.Runes(me.buffer)

	return err
}

func (me *Lexer) Peek() rune {
	if (me.Current.Index) >= len(me.runes) {
		return '\000'
	}
	return me.runes[me.Current.Index]
}

func (me *Lexer) PeekNext() rune {
	if (me.Current.Index + 1) >= len(me.runes) {
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
	_ = me.Next()
}

func (me *Lexer) SkipIn(set string) bool {
	for strings.ContainsRune(set, me.Next()) {
		monolog.Debug("SkipIn: %s %c\n", set, me.Peek())
		if me.Peek() == '\000' {
			return false
		}
	}
	return true
}

func (me *Lexer) SkipNotIn(set string) bool {
	_ = me.Next()
	for !strings.ContainsRune(set, me.Peek()) {
		if me.Next() == '\000' {
			return false
		}
	}
	return true
}

func (me *Lexer) SkipWhile(should_skip func(r rune) bool) bool {
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

func (me *Lexer) Rewind() {
	me.Current = me.Last
}

func (me *Lexer) CurrentRuneValue() []rune {
	return me.runes[me.Last.Index:me.Current.Index]
}

func (me *Lexer) CurrentStringValue() string {
	return string(me.CurrentRuneValue())
}

func (me *Lexer) Found(kind TokenType) {
	me.Emit(kind, Value(me.CurrentStringValue()))
	me.Advance()
}

func GetFunctionName(fun interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(fun).Pointer()).Name()
}

func (me *Lexer) Start() {
	if err := me.ReadReader(); err == nil || err == io.EOF {
		rule := LexNormal
		for rule != nil {
			monolog.Debug("Lexer Rule: %s\n", GetFunctionName(rule))
			rule = rule(me)
		}
	} else {
		me.Error("Could not read in input buffer: %s", err)
	}
	close(me.Output)
}

func (me *Lexer) TryLexing() {
	go me.Start()

	for token := range me.Output {
		monolog.Info("Token %s", token)
	}
}

type AstType int

const (
	AstTypeProgram = AstType(iota)
	AstTypeStatements
	AstTypeStatement
	AstTypeDefinition
	AstTypeWords
	AstTypeExpression
	AstTypeWordExpression
	AstTypeWordCallop
	AstTypeWordOperation
	AstTypeWordCall
	AstTypeValueExpression
	AstTypeValueCallop
	AstTypeValueCall
	AstTypeParametersNonempty
	AstTypeParameters
	AstTypeParameter
	AstTypeBlock
	AstTypeWordvalue
	AstTypeValue
	AstTypeEox
	AstTypeError
)

type Ast struct {
	*tree.Node
	AstType
	*Token
}

func (me *Ast) NewChild(kind AstType, token *Token) *Ast {
	child := &Ast{}
	child.AstType = kind
	child.Token = token
	child.Node = me.Node.NewChild(child)
	return child
}

func (me *Ast) Walk(walker func(ast *Ast) *Ast) *Ast {
	node_res := me.Node.Walk(
		func(node *tree.Node) *tree.Node {
			ast_res := walker(node.Data.(*Ast))
			if ast_res == nil {
				return nil
			} else {
				return ast_res.Node
			}
		})
	return node_res.Data.(*Ast)
}

func NewAst(kind AstType) *Ast {
	ast := &Ast{}
	ast.Node = tree.New(nil, ast)
	ast.AstType = kind
	ast.Token = nil
	return ast
}

type Parser struct {
	*Ast
	*Lexer
}

func (me *Parser) ParseDefinition() {
	/*
		ParseWords()
		ParseBlock()
	*/
}

func (me *Parser) ParseProgram() {
	me.Ast = NewAst(AstTypeProgram)
	token := <-me.Lexer.Output
	switch token.TokenType {
	case TokenKeyword:
		if token.Value == "to" {
			me.ParseDefinition()
			return
		}
		fallthrough
	default:
		me.Ast.NewChild(AstTypeError, &token)
	}
}

/*
	PROGRAM -> STATEMENTS.
STATEMENTS -> STATEMENT STATEMENTS | .
STATEMENT -> DEFINITION | EXPRESSION | BLOCK .
DEFINITION -> to WORDS BLOCK.
WORDS -> word WORDS | .
EXPRESSION -> WORD_EXPRESSION | VALUE_EXPRESSION.
WORD_EXPRESSION -> word WORD_CALLOP.
WORD_CALLOP -> WORD_OPERATION | WORD_CALL.
WORD_OPERATION -> operator PARAMETERS_NONEMPTY EOX.
WORD_CALL -> PARAMETERS EOX.
VALUE_EXPRESSION -> value VALUE_CALLOP.
VALUE_CALLOP -> VALUE_OPERATION | VALUE_CALL.
VALUE_OPERATION -> operator PARAMETERS_NONEMPTY EOX.
VALUE_CALL -> EOX.
PARAMETERS_NONEMPTY -> PARAMETER PARAMETERS.
PARAMETERS -> PARAMETERS_NONEMPTY | .
PARAMETER -> BLOCK | WORDVALUE .
BLOCK -> ob STATEMENTS cb | op STATEMENTS cp | oa STATEMENTS ca | do STATEMENTS end.
WORDVALUE -> word | VALUE.
VALUE -> string | number | symbol.
EOX -> eol | period.


	AstNodeBlock = AstNodeType(iota)
)
*/

type Environment struct {
	Parent *Environment
}

func main() {
	fmt.Println("Hello World!")
}
