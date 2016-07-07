// tree project tree.go
// a relativey simple recursive intrusive tree with an arbitrary amount of children on each level.
package tree

import "fmt"

/* Add a Node to the Struct you want to use these functions with,
and write an initializer maker func(...interface{}) Noder for use with
the New* functions. Everything else goes by itself.
*/
type Node struct {
	Parent_ Noder
	Child_  Noder
	Before_ Noder
	After_  Noder
}

type Noder interface {
	Child() Noder
	Parent() Noder
	Before() Noder
	After() Noder
	SetChild(Noder) Noder
	SetParent(Noder) Noder
	SetBefore(Noder) Noder
	SetAfter(Noder) Noder
}

func (me *Node) Child() Noder {
	return me.Child_
}

func (me *Node) Parent() Noder {
	return me.Parent_
}

func (me *Node) After() Noder {
	return me.After_
}

func (me *Node) Before() Noder {
	return me.Before_
}

func (me *Node) SetChild(val Noder) Noder {
	me.Child_ = val
	return me.Child_
}

func (me *Node) SetParent(val Noder) Noder {
	me.Parent_ = val
	return me.Parent_
}

func (me *Node) SetAfter(val Noder) Noder {
	me.After_ = val
	return me.After_
}

func (me *Node) SetBefore(val Noder) Noder {
	me.Before_ = val
	return me.Before_
}

func NewNoder(parent Noder, maker func(...interface{}) Noder, args ...interface{}) Noder {
	child := maker(args...)
	child.SetParent(parent)
	return child
}

func LastSibling(me Noder) Noder {
	var res Noder = me
	for res != nil && res.After() != nil {
		res = res.After()
	}
	return res
}

func LastChild(me Noder) Noder {
	return LastSibling(me.Child())
}

/* Detaches, I.E removes this node and all it's children from the parent tree. */
func Remove(me Noder) Noder {
	parent := me.Parent()
	before := me.Before()
	after := me.After()
	if before != nil {
		before.SetAfter(after)
	}
	if after != nil {
		after.SetBefore(before)
	}
	if parent != nil {
		/* Special case if me is the first child of it's parent. */
		if me == parent.Child() {
			parent.SetChild(after)
		}
	}
	me.SetParent(nil)
	return me
}

func InsertSibling(me, sibling Noder) Noder {
	after := me.After()
	me.SetAfter(sibling)
	sibling.SetBefore(me)
	sibling.SetAfter(after)
	if after != nil {
		after.SetBefore(sibling)
	}
	sibling.SetParent(me.Parent())
	return sibling
}

func AppendSibling(me, sibling Noder) Noder {
	return InsertSibling(LastSibling(me), sibling)
}

func AppendChild(me, child Noder) Noder {
	child.SetParent(me)
	if me.Child() == nil {
		me.SetChild(child)
	} else {
		AppendSibling(me.Child(), child)
	}
	return child
}

func NewSibling(me Noder, maker func(...interface{}) Noder, args ...interface{}) Noder {
	node := NewNoder(me.Parent(), maker, args...)
	return AppendSibling(me, node)
}

func NewChild(me Noder, maker func(...interface{}) Noder, args ...interface{}) Noder {
	node := NewNoder(me, maker, args...)
	return AppendChild(me, node)
}

func Walk(me Noder, walker func(me Noder) Noder) Noder {
	if found := walker(me); found != nil {
		return found
	}
	if me.Child() != nil {
		if found := Walk(me.Child(), walker); found != nil {
			return found
		}
	}
	if me.After() != nil {
		if found := Walk(me.After(), walker); found != nil {
			return found
		}
	}
	return nil
}

func Display(me Noder) {
	Walk(me, func(node Noder) Noder {
		fmt.Printf("Tree: %v\n", node)
		return nil
	})
}

/*
interface Walker {
	Walk(walker func(me *Node) *Node) *Node
}


func WalkWalker(walker Walker) {

}
*/
