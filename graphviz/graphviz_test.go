package graphviz

import (
	"testing"
)

func TestShow(test *testing.T) {
	g := NewDigraph("bgcolor", "pink")
	n_foo := g.AddNode("foo", "color", "red", "label", "FOO\nFOO")
	n_bar := g.AddNode("bar", "color", "green")
	_ = g.AddNode("baz", "color", "yellow")
	g.AddEdge(n_foo, n_bar, "color", "gray")
	_ = g.AddEdgeByName("foo", "baz", "color", "magenta")
	g.Dotty()
}
