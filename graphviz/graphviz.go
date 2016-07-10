// graphviz
package graphviz

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

var replacer *strings.Replacer

func init() {
	replacer = strings.NewReplacer("\n", "\\n", "\r", "\\r", "\t", "\\t")
}

type Attributes map[string]string

func NewAttributes(attributes ...string) Attributes {
	me := make(Attributes)
	for i := 1; i < len(attributes); i += 2 {
		key := attributes[i-1]
		value := replacer.Replace(attributes[i])
		me[key] = value
	}
	return me
}

type Node struct {
	Attributes
	ID string
}

func NewNode(id string, attributes ...string) *Node {
	me := &Node{}
	me.ID = id
	me.Attributes = NewAttributes(attributes...)
	return me
}

func (me Attributes) WriteTo(out io.Writer) {
	comma := false
	if len(me) > 0 {
		fmt.Fprintf(out, "[")
		for k, v := range me {
			if comma {
				fmt.Fprintf(out, ",")
			}
			fmt.Fprintf(out, "%s=\"%s\"", k, v)
			comma = true
		}
		fmt.Fprintf(out, "]")
	}
}

func (me Attributes) WriteForGraphTo(out io.Writer) {
	if len(me) > 0 {
		for k, v := range me {
			fmt.Fprintf(out, "%s=\"%s\";\n", k, v)
		}
	}
}

func (me *Node) WriteTo(out io.Writer) {
	fmt.Fprintf(out, "%s", me.ID)
	me.Attributes.WriteTo(out)
	fmt.Fprintf(out, ";\n")
}

type Edge struct {
	Attributes
	From *Node
	To   *Node
}

func NewEdge(from, to *Node, attributes ...string) *Edge {
	me := &Edge{}
	me.From = from
	me.To = to
	me.Attributes = NewAttributes(attributes...)
	return me
}

func (me *Edge) WriteTo(out io.Writer) {
	if (me.From != nil) && (me.To != nil) {
		fmt.Fprintf(out, "%s -> %s ", me.From.ID, me.To.ID)
		me.Attributes.WriteTo(out)
		fmt.Fprintf(out, ";\n")
	}
}

type Digraph struct {
	Attributes
	nodes []*Node
	edges []*Edge
}

func NewDigraph(attributes ...string) *Digraph {
	me := &Digraph{}
	me.Attributes = NewAttributes(attributes...)
	return me
}

func (me *Digraph) AddNode(id string, attributes ...string) *Node {
	node := NewNode(id, attributes...)
	me.nodes = append(me.nodes, node)
	return node
}

func (me *Digraph) AddEdge(from, to *Node, attributes ...string) *Edge {
	edge := NewEdge(from, to, attributes...)
	me.edges = append(me.edges, edge)
	return edge
}

func (me *Digraph) FindNode(id string) *Node {
	/* XXX stupid linear search for now... */
	for _, node := range me.nodes {
		if node.ID == id {
			return node
		}
	}
	return nil
}

func (me *Digraph) AddEdgeByName(from, to string, attributes ...string) *Edge {
	node_from := me.FindNode(from)
	node_to := me.FindNode(to)
	return me.AddEdge(node_from, node_to, attributes...)
}

func (me *Digraph) WriteTo(out io.Writer) {
	fmt.Fprintf(out, "digraph {\n")
	me.Attributes.WriteForGraphTo(out)
	for _, node := range me.nodes {
		node.WriteTo(out)
	}
	for _, edge := range me.edges {
		edge.WriteTo(out)
	}
	fmt.Fprintf(out, "\n}\n")
}

func (me *Digraph) Dotty() error {
	file, err := ioutil.TempFile("", "woe_gv_")
	if file == nil {
		return err
	}

	me.WriteTo(file)
	name := file.Name()
	file.Close()

	cmd := exec.Command("dotty", name)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil {
		return err
	}
	// cmd.Wait()
	return nil
}
