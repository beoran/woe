// raku_test.go
package raku

import (
	"strings"
	"testing"

	_ "github.com/beoran/woe/monolog"
	"github.com/beoran/woe/tree"
)

func HelperTryLexing(me *Lexer, test *testing.T) {
	go me.Start()
	test.Logf("Lexing started:")
	test.Logf("Lexer buffer: %v", me.buffer)

	for token := range me.Output {
		test.Logf("Token %s", token)
	}
}

func Assert(test *testing.T, ok bool, text string) bool {
	if !ok {
		test.Error(text)
	}
	return ok
}

func TestLexing(test *testing.T) {
	const input = `
say "hello \"world\\"

to open a door do
	set door's open to true
end

to increment variable by value do
	variable = variable + value 
end
`
	lexer := OpenLexer(strings.NewReader(input))
	HelperTryLexing(lexer, test)
	test.Log("Hi test!")
}

func TestLexing2(test *testing.T) {
	const input = `say`
	lexer := OpenLexer(strings.NewReader(input))
	HelperTryLexing(lexer, test)
	test.Log("Hi test!")
}

func TestLexing3(test *testing.T) {
	const input = `$sym`
	lexer := OpenLexer(strings.NewReader(input))
	HelperTryLexing(lexer, test)
	test.Log("Hi test!")
}

func TestParseValue(test *testing.T) {
	const input = `"hello \"world\\"`
	parser := NewParserForText(input)
	Assert(test, parser.ParseValue(), "Could not parse value")
	tree.Display(parser.Ast)
}

func TestParseValue2(test *testing.T) {
	const input = `2.1`
	parser := NewParserForText(input)
	Assert(test, parser.ParseValue(), "Could not parse value")
	tree.Display(parser.Ast)
}

func TestParseValue3(test *testing.T) {
	const input = `$sym`
	parser := NewParserForText(input)
	Assert(test, parser.ParseValue(), "Could not parse value")
	tree.Display(parser.Ast)
}

func TestParseEox(test *testing.T) {
	const input = `
`
	parser := NewParserForText(input)
	Assert(test, parser.ParseEOX(), "Could not parse EOX")
	tree.Display(parser.Ast)
}

func TestParseEox2(test *testing.T) {
	const input = `.
`
	parser := NewParserForText(input)
	Assert(test, parser.ParseEOX(), "Could not parse EOX")
	tree.Display(parser.Ast)
}

func TestParseWord(test *testing.T) {
	const input = `say`
	parser := NewParserForText(input)
	Assert(test, parser.ParseWord(), "Could not parse word")
	tree.Display(parser.Ast)
}

func TestParseWordExpression(test *testing.T) {
	const input = `say "hello world" three times
	`
	parser := NewParserForText(input)
	Assert(test, parser.ParseWordExpression(), "Could not parse word expression")
	tree.Display(parser.Ast)
}

func TestParseWordExpression2(test *testing.T) {
	const input = `val + 10 * z
	`
	parser := NewParserForText(input)
	Assert(test, parser.ParseWordExpression(), "Could not parse word expression with operators")
	tree.Display(parser.Ast)
}

func TestParseStatements(test *testing.T) {
	const input = `val + 10 * z. open door.
	`
	parser := NewParserForText(input)
	Assert(test, parser.ParseStatements(), "Could not parse statements with only a parse word expression with operators")
	tree.Display(parser.Ast)
}

func TestParseProgram(test *testing.T) {
	const input = `val + 10 * z. open door.
	`
	parser := NewParserForText(input)
	Assert(test, parser.ParseProgram(), "Could not parse program.")
	tree.Display(parser.Ast)
}

func TestParseProgram2(test *testing.T) {
	const input = `to greet someone [
say "hello" someone
]

greet bob

if mp < cost do
	say "Not enough mana!"
end else do
	say "Zap!"
end

`
	parser := NewParserForText(input)
	Assert(test, parser.ParseProgram(), "Could not parse program.")
	tree.Display(parser.Ast)
}

func TestParseblock(test *testing.T) {
	// monolog.Setup("raku_test.log", true, false)
	const input = `[
say "hello"
say "world"
]
`
	parser := NewParserForText(input)
	Assert(test, parser.ParseBlock(), "Could not parse block.")
	tree.Display(parser.Ast)
	parser.Ast.Dotty()
}
