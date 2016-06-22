// raku_test.go
package raku

import (
	"strings"
	"testing"
)

func TestLexing(test *testing.T) {
	const input = `
say "hello world"

to open a door do
	set door's open to true
done
	
	
`
	lexer := OpenLexer(strings.NewReader(input))
	test.Logf("Lexer buffer: %v", lexer.buffer)
	lexer.TryLexing()
	test.Log("Hi test!")
}
