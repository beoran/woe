// tree_tes
package tree

import (
	_ "strings"
	"testing"
)

func TestNode(test *testing.T) {
	tree := New(nil, "root")
	s1 := "l1 c1"
	s2 := "l1 c2"
	s3 := "l1 c3"
	s4 := "l2 c1"
	l1c1 := tree.NewChild(s1)
	l1c2 := tree.NewChild(s2)
	l1c3 := tree.NewChild(s3)
	l2c1 := l1c1.NewChild(s4)

	if l1c1.Data != s1 {
		test.Error("Data ")
	}

	if l1c2.Data != s2 {
		test.Error("Data ")
	}

	if l1c3.Data != s3 {
		test.Error("Data ")
	}

	if l2c1.Data != s4 {
		test.Error("Data ")
	}

	n := tree.Walk(func(node *Node) *Node {
		if node.Data == s4 {
			return node
		}
		return nil
	})

	if n.Data != s4 {
		test.Error("Data ")
	}

	test.Logf("%v", n.Data)

	test.Log("Hi tree!")
}
