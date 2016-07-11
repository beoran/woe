// raku

/* Raku is an easy to use scripting language that can also be used easily interactively

Desrired syntax (verified LL(1) on smlweb.cpsc.ucalgary.ca)

PROGRAM -> STATEMENTS.
STATEMENTS -> STATEMENT STATEMENTS | .
STATEMENT -> EXPRESSION EOX  | DEFINITION | BLOCK | EOX .
DEFINITION -> define WORDS BLOCK.
WORDS -> word WORDS | .
EXPRESSION -> WORDVALUE MODIFIERS.
MODIFIERS -> MODIFIER MODIFIERS | .
OPERATION ->  operator MODIFIER .
MODIFIER -> OPERATION | WORDVALUE | PARENTHESIS | BLOCK.
PARENTHESIS -> '(' EXPRESSION ')' | ot EXPRESSION ct.
BLOCK -> oe STATEMENTS ce | do STATEMENTS end .
WORDVALUE -> word | VALUE | a | the.
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
	"strings"
	"unicode"

	"github.com/beoran/woe/graphviz"
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
	TokenPeriod       TokenType = TokenType('.')
	TokenComma        TokenType = TokenType(',')
	TokenSemicolon    TokenType = TokenType(';')
	TokenColon        TokenType = TokenType(':')
	TokenOpenParen    TokenType = TokenType('(')
	TokenCloseParen   TokenType = TokenType(')')
	TokenOpenBrace    TokenType = TokenType('{')
	TokenCloseBrace   TokenType = TokenType('}')
	TokenOpenBracket  TokenType = TokenType('[')
	TokenCloseBracket TokenType = TokenType(']')

	TokenNone         TokenType = 0
	TokenError        TokenType = -1
	TokenWord         TokenType = -2
	TokenEOL          TokenType = -3
	TokenEOF          TokenType = -4
	TokenNumber       TokenType = -5
	TokenOperator     TokenType = -6
	TokenString       TokenType = -7
	TokenSymbol       TokenType = -8
	TokenFirstKeyword TokenType = -9
	TokenKeywordA     TokenType = -10
	TokenKeywordDo    TokenType = -11
	TokenKeywordEnd   TokenType = -12
	TokenKeywordThe   TokenType = -13
	TokenKeywordDef   TokenType = -14
	TokenLastKeyword  TokenType = -15
	TokenLast         TokenType = -15
)

type Token struct {
	TokenType
	Value
	Position
}

var tokenTypeMap map[TokenType]string = map[TokenType]string{
	TokenNone:       "TokenNone",
	TokenError:      "TokenError",
	TokenWord:       "TokenWord",
	TokenEOL:        "TokenEOL",
	TokenEOF:        "TokenEOF",
	TokenNumber:     "TokenNumber",
	TokenOperator:   "TokenOperator",
	TokenString:     "TokenString",
	TokenSymbol:     "TokenSymbol",
	TokenKeywordA:   "TokenKeywordA",
	TokenKeywordDo:  "TokenKeywordDo",
	TokenKeywordEnd: "TokenKeywordEnd",
	TokenKeywordThe: "TokenKeywordThe",
	TokenKeywordDef: "TokenKeywordDef",
}

var keywordMap map[string]TokenType = map[string]TokenType{
	"a":      TokenKeywordA,
	"an":     TokenKeywordA,
	"do":     TokenKeywordDo,
	"def":    TokenKeywordDef,
	"define": TokenKeywordDef,
	"end":    TokenKeywordEnd,
	"the":    TokenKeywordThe,
}

var sigilMap map[string]TokenType = map[string]TokenType{
	"[": TokenOpenBracket,
	"{": TokenOpenBrace,
	"(": TokenOpenParen,
	"]": TokenCloseBracket,
	"}": TokenCloseBrace,
	")": TokenCloseParen,
}

func (me TokenType) String() string {
	name, found := tokenTypeMap[me]
	if found {
		return name
	} else {
		if (me > 0) && (me < 256) {
			return fmt.Sprintf("TokenChar<%c>", byte(me))
		}
		return fmt.Sprintf("Unknown Token %d", int(me))
	}
}

func (me Token) String() string {
	return fmt.Sprintf("Token: %s >%s< %d %d %d.", me.TokenType, string(me.Value), me.Index, me.Row, me.Column)
}

type TokenChannel chan *Token

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
	tok := &Token{t, v, me.Current}
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

/* Returns whether or not a keyword was found, and if so, the TokenType
of the keyword.*/
func LookupKeyword(word string) (bool, TokenType) {
	kind, found := keywordMap[word]
	return found, kind
}

/* Returns whether or not a special operator or sigil was found, and if so,
returns the TokenTyp of the sigil.*/
func LookupSigil(sigil string) (bool, TokenType) {
	fmt.Printf("LookupSigil: %s\n", sigil)
	kind, found := sigilMap[sigil]
	return found, kind
}

func LexSigil(me *Lexer) LexerRule {
	me.Found(TokenType(me.Peek()))
	_ = me.Next()
	me.Advance()
	return LexNormal
}

func LexWord(me *Lexer) LexerRule {
	me.SkipNotIn(" \t\r\n'({[]})")

	iskw, kind := LookupKeyword(me.CurrentStringValue())
	if iskw {
		me.Found(kind)
	} else {
		me.Found(TokenWord)
	}
	return LexNormal
}

func LexSymbol(me *Lexer) LexerRule {
	me.SkipNotIn(" \t\r\n'({[]})")
	me.Found(TokenSymbol)
	return LexNormal
}

func LexNumber(me *Lexer) LexerRule {
	me.SkipNotIn(" \t\r\n'({[]})")
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
	me.Found(TokenType(me.Next()))
	me.Advance()
	return LexNormal
}

func LexEOL(me *Lexer) LexerRule {
	me.SkipIn("\r\n")
	me.Found(TokenEOL)
	return LexNormal
}

func LexOperator(me *Lexer) LexerRule {
	me.SkipNotIn(" \t\r\n({[]})")
	issig, kind := LookupSigil(me.CurrentStringValue())
	if issig {
		me.Found(kind)
	} else {
		me.Found(TokenOperator)
	}
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
	} else if strings.ContainsRune("([{}])", peek) {
		return LexSigil
	} else if strings.ContainsRune("$", peek) {
		return LexSymbol
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
		//me.Emit(TokenEOF, "")
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
	AstTypeOperation
	AstTypeOperations
	AstTypeCallArgs
	AstTypeValueExpression
	AstTypeValueCallop
	AstTypeParametersNonempty
	AstTypeParameters
	AstTypeParameter
	AstTypeBlock
	AstTypeWordValue
	AstTypeWord
	AstTypeValue
	AstTypeEox
	AstTypeOperator
	AstTypeParenthesis
	AstTypeModifier
	AstTypeError
)

var astTypeMap map[AstType]string = map[AstType]string{
	AstTypeProgram:            "AstTypeProgram",
	AstTypeStatements:         "AstTypeStatements",
	AstTypeStatement:          "AstTypeStatement:",
	AstTypeDefinition:         "AstTypeDefinition",
	AstTypeWords:              "AstTypeWords",
	AstTypeExpression:         "AstTypeExpression",
	AstTypeWordExpression:     "AstTypeWordExpression",
	AstTypeWordCallop:         "AstTypeWordCallop",
	AstTypeOperation:          "AstTypeOperation",
	AstTypeOperations:         "AstTypeOperations",
	AstTypeCallArgs:           "AstTypeCallArgs",
	AstTypeValueExpression:    "AstTypeValueExpression",
	AstTypeValueCallop:        "AstTypeValueCallop",
	AstTypeParametersNonempty: "AstTypeParametersNonempty",
	AstTypeParameters:         "AstTypeParameters",
	AstTypeParameter:          "AstTypeParameter",
	AstTypeBlock:              "AstTypeBlock",
	AstTypeWordValue:          "AstTypeWordValue",
	AstTypeWord:               "AstTypeWord",
	AstTypeValue:              "AstTypeValue",
	AstTypeEox:                "AstTypeEox",
	AstTypeOperator:           "AstTypeOperator",
	AstTypeParenthesis:        "AstTypeParenthesis",
	AstTypeModifier:           "AstTypeModifier",
	AstTypeError:              "AstTypeError",
}

func (me AstType) String() string {
	name, found := astTypeMap[me]
	if found {
		return name
	} else {
		return fmt.Sprintf("Unknown AstType %d", int(me))
	}
}

type Ast struct {
	tree.Node
	AstType
	*Token
}

func (me *Ast) NewChild(kind AstType, token *Token) *Ast {
	child := &Ast{}
	child.AstType = kind
	child.Token = token
	tree.AppendChild(me, child)
	return child
}

func (me *Ast) Walk(walker func(ast *Ast) *Ast) *Ast {
	node_res := tree.Walk(me,
		func(node tree.Noder) tree.Noder {
			ast_res := walker(node.(*Ast))
			if ast_res == nil {
				return nil
			} else {
				return ast_res
			}
		})
	if node_res != nil {
		return node_res.(*Ast)
	} else {
		return nil
	}
}

func (me *Ast) Remove() {
	_ = tree.Remove(me)
}

func NewAst(kind AstType) *Ast {
	ast := &Ast{}
	ast.AstType = kind
	ast.Token = nil
	return ast
}

type ParseAction func(parser *Parser) bool

type RuleType int

const (
	RuleTypeNone = RuleType(iota)
	RuleTypeAlternate
	RuleTypeSequence
)

type Rule struct {
	tree.Node
	Name string
	RuleType
	ParseAction
}

func NewRule(name string, ruty RuleType) *Rule {
	res := &Rule{}
	res.RuleType = ruty
	res.Name = name
	return res
}

func (me *Rule) NewChild(action ParseAction) *Rule {
	child := NewRule("foo", RuleTypeNone)
	tree.AppendChild(me, child)
	return child
}

func (me *Rule) Walk(walker func(rule *Rule) *Rule) *Rule {
	node_res := tree.Walk(me,
		func(node tree.Noder) tree.Noder {
			rule_res := walker(node.(*Rule))
			if rule_res == nil {
				return nil
			} else {
				return rule_res
			}
		})
	return node_res.(*Rule)
}

type Parser struct {
	*Ast
	*Lexer
	now       *Ast
	lookahead *Token
}

func (me *Parser) SetupRules() {

}

func (me *Parser) Expect(types ...TokenType) bool {
	monolog.Debug("Expecting: ", types, " from ", me.now.AstType, " have ", me.LookaheadType(), " \n")
	for _, t := range types {
		if me.LookaheadType() == t {
			monolog.Debug("Found: ", t, "\n")
			return true
		}
	}
	monolog.Debug("Not found.\n")
	return false
}

type Parsable interface {
	isParsable()
}

func (me TokenType) isParsable() {
}

func (me ParseAction) isParsable() {
}

/* Advance the lexer but only of there is no lookahead token already available in me.lookahead.
 */
func (me *Parser) Advance() *Token {
	if me.lookahead == nil {
		me.lookahead = <-me.Lexer.Output
	}
	return me.lookahead
}

func (me *Parser) DropLookahead() {
	me.lookahead = nil
}

func (me *Parser) Lookahead() *Token {
	return me.lookahead
}

func (me *Parser) LookaheadType() TokenType {
	if me.lookahead == nil {
		return TokenError
	}
	return me.Lookahead().TokenType
}

func (me *Parser) Consume(atyp AstType, types ...TokenType) bool {
	me.Advance()
	res := me.Expect(types...)
	if res {
		me.NewAstChild(atyp)
		me.DropLookahead()
	}
	return res
}

func (me *Parser) ConsumeWithoutAst(types ...TokenType) bool {
	me.Advance()
	res := me.Expect(types...)
	if res {
		me.DropLookahead()
	}
	return res
}

/*
func (me * Parser) OneOf(restype AstType, options ...Parsable) bool {
	res := false
	k, v := range options {
		switch option := v.Type {
			case TokenType: res := Consume(restype, option)
			case ParseAction: res := option(me)
		}
	}
	return res
}
*/

func (me *Parser) ParseEOX() bool {
	return me.ConsumeWithoutAst(TokenEOL, TokenPeriod)
}

func (me *Parser) ParseValue() bool {
	return me.Consume(AstTypeValue, TokenString, TokenNumber, TokenSymbol)
}

func (me *Parser) ParseWord() bool {
	return me.Consume(AstTypeWord, TokenWord, TokenKeywordA, TokenKeywordThe)
}

func (me *Parser) ParseWordValue() bool {
	me.NewAstChildDescend(AstTypeWordValue)
	res := me.ParseValue() || me.ParseWord()
	me.AstAscend(res)
	return res
}

func (me *Parser) ParseParameter() bool {
	me.NewAstChildDescend(AstTypeParameter)
	res := me.ParseWordValue() || me.ParseBlock()
	me.AstAscend(res)
	return res
}

func (me *Parser) ParseParametersNonempty() bool {
	res := false
	for me.ParseParameter() {
		res = true
	}
	return res
}

func (me *Parser) ParseParameters() bool {
	me.NewAstChildDescend(AstTypeParameters)
	_ = me.ParseParametersNonempty()
	me.AstAscend(true)
	return true
}

func (me *Parser) ParseCallArgs() bool {
	me.NewAstChildDescend(AstTypeCallArgs)
	res := me.ParseParameters() && me.ParseEOX()
	me.AstAscend(res)
	return res
}

func (me *Parser) ParseOperator() bool {
	return me.Consume(AstTypeOperator, TokenOperator)
}

/*
func (me *Parser) ParseOperation() bool {
	me.NewAstChildDescend(AstTypeOperation)
	res := me.ParseOperator() && me.ParseParameter()
	me.AstAscend(res)
	return res
}
*/

func (me *Parser) ParseOperations() bool {
	me.NewAstChildDescend(AstTypeOperations)
	res := me.ParseOperation()
	for me.ParseOperation() {
	}
	me.AstAscend(res)
	return res
}

func (me *Parser) ParseWordCallOp() bool {
	me.NewAstChildDescend(AstTypeWordCallop)
	res := me.ParseCallArgs() || me.ParseOperations()
	me.AstAscend(res)
	return res
}

func (me *Parser) ParseWordExpression() bool {
	me.NewAstChildDescend(AstTypeWordExpression)
	res := me.ParseWord() && me.ParseWordCallOp()
	me.AstAscend(res)
	return res
}

func (me *Parser) ParseValueCallOp() bool {
	me.NewAstChildDescend(AstTypeValueCallop)
	res := me.ParseCallArgs() || me.ParseOperations()
	me.AstAscend(res)
	return res
}

func (me *Parser) ParseValueExpression() bool {
	me.NewAstChildDescend(AstTypeValueExpression)
	res := me.ParseValue() && me.ParseValueCallOp()
	me.AstAscend(res)
	return false
}

func (me *Parser) NewAstChild(tyty AstType) *Ast {
	return me.now.NewChild(tyty, me.lookahead)
}

func (me *Parser) NewAstChildDescend(tyty AstType) {
	node := me.NewAstChild(tyty)
	me.now = node
}

func (me *Parser) AstAscend(keep bool) {
	if me.now.Parent() != nil {
		now := me.now
		me.now = now.Parent().(*Ast)
		if !keep {
			now.Remove()
		}
	}
}

func (me TokenType) BlockCloseForOpen() (TokenType, bool) {
	switch me {
	case TokenOpenBrace:
		return TokenCloseBrace, true
	case TokenOpenParen:
		return TokenCloseParen, true
	default:
		return TokenError, false
	}

}

func (me TokenType) ParenthesisCloseForOpen() (TokenType, bool) {
	switch me {
	case TokenOpenBracket:
		return TokenCloseBracket, true
	case TokenOpenParen:
		return TokenCloseParen, true
	default:
		return TokenError, false
	}

}

func (me *Parser) ParseBlock() bool {
	me.Advance()
	open := me.LookaheadType()
	done, ok := open.BlockCloseForOpen()
	if !ok {
		/* Not an opening of a block, so no block found. */
		return false
	}
	me.DropLookahead()
	me.NewAstChildDescend(AstTypeBlock)
	res := me.ParseStatements()
	me.AstAscend(res)
	if res {
		me.Advance()
		if me.LookaheadType() != done {
			return me.ParseError()
		}
		me.DropLookahead()
	}
	return res
}

func (me *Parser) ParseParenthesis() bool {
	me.Advance()
	open := me.LookaheadType()
	done, ok := open.ParenthesisCloseForOpen()
	if !ok {
		/* Not an opening of a parenthesis, so no parenthesis found. */
		return false
	}
	me.DropLookahead()
	me.NewAstChildDescend(AstTypeParenthesis)
	res := me.ParseExpression()
	me.AstAscend(res)
	if res {
		me.Advance()
		if me.LookaheadType() != done {
			return me.ParseError()
		}
		me.DropLookahead()
	}
	return res
}

func (me *Parser) ParseWords() bool {
	me.NewAstChildDescend(AstTypeWords)
	res := me.ParseWord()
	for me.ParseWord() {
	}
	me.AstAscend(res)
	return res
}

func (me *Parser) ParseDefinition() bool {
	me.Advance()
	res := me.Consume(AstTypeDefinition, TokenKeywordDef)
	if !res {
		return false
	}
	res = res && me.ParseWords()
	if !res {
		_ = me.ParseError()
	}
	res = res && me.ParseBlock()
	if !res {
		_ = me.ParseError()
	}
	me.AstAscend(true)
	return res
}

func (me *Parser) ParseOperation() bool {
	me.NewAstChildDescend(AstTypeOperation)
	res := me.ParseOperator() && me.ParseModifier()
	me.AstAscend(res)
	return res
}

func (me *Parser) ParseModifier() bool {
	me.NewAstChildDescend(AstTypeModifier)
	res := me.ParseOperation() || me.ParseWordValue() ||
		me.ParseParenthesis() || me.ParseBlock()
	me.AstAscend(res)
	return res
}

func (me *Parser) ParseModifiers() bool {
	for me.ParseModifier() {
	}
	return true
}

func (me *Parser) ParseError() bool {
	me.now.NewChild(AstTypeError, me.lookahead)
	fmt.Printf("Parse error: at %s\n", me.lookahead)
	return false
}

func (me *Parser) ParseExpression() bool {
	return me.ParseWordValue() && me.ParseModifiers()
}

func (me *Parser) ParseStatement() bool {
	me.NewAstChildDescend(AstTypeStatement)
	/* First case is for an empty expression/statement. */
	res := me.ParseEOX() ||
		me.ParseDefinition() ||
		(me.ParseExpression() && me.ParseEOX()) ||
		me.ParseBlock()

	me.AstAscend(res)
	return res
}

func (me *Parser) ParseEOF() bool {
	return me.Consume(AstTypeEox, TokenEOF)
}

func (me *Parser) ParseStatements() bool {
	me.NewAstChildDescend(AstTypeStatements)
	res := me.ParseStatement()

	for me.ParseStatement() {
	}

	me.AstAscend(res)
	return res
}

func (me *Parser) ParseProgram() bool {
	return me.ParseStatements() && me.ParseEOF()
}

func NewParserForLexer(lexer *Lexer) *Parser {
	me := &Parser{}
	me.Ast = NewAst(AstTypeProgram)
	me.now = me.Ast
	me.Lexer = lexer
	me.Ast.Token = &Token{}
	go me.Lexer.Start()
	return me
}

func NewParserForText(text string) *Parser {
	lexer := OpenLexer(strings.NewReader(text))
	return NewParserForLexer(lexer)
}

func (me *Ast) DotID() string {
	return fmt.Sprintf("ast_%p", me)
}

func (me *Ast) Dotty() {
	g := graphviz.NewDigraph("rankdir", "LR")
	me.Walk(func(ast *Ast) *Ast {
		label := ast.AstType.String()
		if ast.Token != nil {
			label = label + "\n" + ast.Token.String()
		}
		g.AddNode(ast.DotID(), "label", label)
		if ast.Parent() != nil {
			g.AddEdgeByName(ast.Parent().(*Ast).DotID(), ast.DotID())
		}
		return nil
	})
	g.Dotty()
}

/*

PROGRAM -> STATEMENTS.
STATEMENTS -> STATEMENT STATEMENTS | .
STATEMENT -> EXPRESSION EOX  | DEFINITION | BLOCK | EOX .
DEFINITION -> define WORDS BLOCK.
WORDS -> word WORDS | .
EXPRESSION -> WORDVALUE MODIFIERS.
MODIFIERS -> MODIFIER MODIFIERS | .
OPERATION ->  operator MODIFIER .
MODIFIER -> OPERATION | WORDVALUE | PARENTHESIS | BLOCK.
PARENTHESIS -> '(' EXPRESSION ')' | ot EXPRESSION ct.
BLOCK -> oe STATEMENTS ce | do STATEMENTS end .
WORDVALUE -> word | VALUE | a | the.
VALUE -> string | number | symbol.
EOX -> eol | period.

	AstNodeBlock = AstNodeType(iota)
)
*/

type DefineType int

const (
		DefineTypeNone = DefineType(iota),
		DefineTypeGo,
		DefineTypeUser,
		DefineTypeVar,
)

type Value interface {
	
}

type DefinePattern struct {
	Parts []string
}

type GoDefineFunc func(runtime Runtime, args ... Value) Value;

type UserDefine struct {
	DefinePattern
	* Ast	
}

type GoDefine struct {
	DefinePattern
	* GoDefineFunc
}


type Define struct {
	DefineType
	Ast * definition
}

type Environment struct {
	Parent *Environment	
}


type Runtime struct {
	Environment
}




func main() {
	fmt.Println("Hello World!")
}
