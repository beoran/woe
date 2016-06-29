// raku_test.go
package raku

import (
	"strings"
	"testing"
)

func HelperTryLexing(me *Lexer, test *testing.T) {
	go me.Start()
	test.Logf("Lexing started:")
	test.Logf("Lexer buffer: %v", me.buffer)

	for token := range me.Output {
		test.Logf("Token %s", token)
	}
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
