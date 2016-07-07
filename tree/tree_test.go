// tree_tes
package tree

import (
	"fmt"
	_ "strings"
	"testing"
)

type StringTree struct {
	Node
	Data string
}

func NewStringTree(value string) *StringTree {
	res := &StringTree{}
	res.Data = value
	res.SetParent(nil)
	return res
}

func InitStringTree(args ...interface{}) Noder {
	str := args[0].(string)
	return NewStringTree(str)
}

func TestSibling(test *testing.T) {
	n1 := NewStringTree("s1")
	n2 := NewStringTree("s2")
	n3 := NewStringTree("s3")
	AppendSibling(n1, n2)
	AppendSibling(n1, n3)
	NewSibling(n1, InitStringTree, "s4")
	Display(n1)
}

func TestNode(test *testing.T) {
	tree := NewStringTree("root")
	s1 := "l1 c1"
	s2 := "l1 c2"
	s3 := "l1 c3"
	s4 := "l2 c1"
	l1c1 := NewChild(tree, InitStringTree, s1)

	if tree.Child() != l1c1 {
		test.Errorf("Child() %v: %v<->%v", tree, tree.Child(), l1c1)
	}

	if LastChild(tree) != l1c1 {
		test.Errorf("LastChild() %v: %v<->%v", tree, tree.Child(), l1c1)
	}

	l1c2 := NewChild(tree, InitStringTree, s2)

	if tree.Child() != l1c1 {
		test.Errorf("Child() %v: %v<->%v", tree, tree.Child(), l1c1)
	}

	if LastChild(tree) != l1c2 {
		test.Errorf("LastChild() %v: %v<->%v", tree, tree.Child(), l1c2)
	}

	l1c3 := NewChild(tree, InitStringTree, s3)

	if LastChild(tree) != l1c3 {
		test.Errorf("LastChild() %v: %v<->%v", tree, tree.Child(), l1c3)
	}

	l2c1 := NewChild(l1c2, InitStringTree, s4)

	if l1c2.Child() != l2c1 {
		test.Errorf("Child() %v: %v<->%v", l1c2, l1c2.Child(), l2c1)
	}

	if LastChild(l1c2) != l2c1 {
		test.Errorf("LastChild() %v: %v<->%v", l1c2, l1c2.Child(), l2c1)
	}

	if l1c1.(*StringTree).Data != s1 {
		test.Error("Data ")
	}

	if l1c2.(*StringTree).Data != s2 {
		test.Error("Data ")
	}

	if l1c3.(*StringTree).Data != s3 {
		test.Error("Data ")
	}

	if l2c1.(*StringTree).Data != s4 {
		test.Error("Data ")
	}

	Display(tree)

	if tree.Child() != l1c1 {
		test.Errorf("Child() %v: %v<->%v", tree, tree.Child(), l1c1)
	}

	if l1c1.After() != l1c2 {
		test.Errorf("After()  %v<->%v", tree.After(), l1c2)
	}

	if l1c2.After() != l1c3 {
		test.Error("After()")
	}

	if l1c2.Child() != l2c1 {
		test.Error("Child()")
	}

	n := Walk(tree, func(node Noder) Noder {
		fmt.Println("%v", node)
		if node.(*StringTree).Data == s4 {
			return node
		}
		return nil
	})

	if n.(*StringTree).Data != s4 {
		test.Error("Data ")
	}

	test.Logf("%v", n.(*StringTree).Data)

	test.Log("Hi tree!")
}

type Tn struct {
	Node
	data string
}

func TestNoder(test *testing.T) {
	test.Log("Hi treenoder!")
}

func TestDelete(test *testing.T) {

	tree := NewStringTree("root")
	s1 := "l1 c1"
	s2 := "l1 c2"
	s3 := "l1 c3"
	s4 := "l2 c1"
	_ = NewChild(tree, InitStringTree, s1)
	l1c2 := NewChild(tree, InitStringTree, s2)
	_ = NewChild(tree, InitStringTree, s3)
	_ = NewChild(l1c2, InitStringTree, s4)
	Remove(l1c2)
	Display(tree)

	n := Walk(tree, func(node Noder) Noder {
		fmt.Println("%v", node)
		if node.(*StringTree).Data == l1c2.(*StringTree).Data {
			return node
		}
		return nil
	})

	if n != nil {
		test.Errorf("Not deleted: %v", n)
	}

}
