// tree project tree.go
// a relativey simple recursive tree with an arbitrary amount of children on each level.
package tree

type Node struct {
	Child  *Node
	After  *Node
	Before *Node
	Parent *Node
	Data   interface{}
}

func NewEmpty() *Node {
	return &Node{}
}

func New(parent *Node, data interface{}) *Node {
	node := NewEmpty()
	node.Parent = parent
	node.Data = data
	return node
}

func (me *Node) LastSibling() *Node {
	res := me
	for res != nil && res.After != nil {
		res = res.After
	}
	return res
}

func (me *Node) InsertSibling(sibling *Node) *Node {
	after := me.After
	me.After = sibling
	sibling.Before = me
	sibling.After = after
	sibling.Parent = me.Parent
	return sibling
}

func (me *Node) AppendSibling(sibling *Node) *Node {
	return me.LastSibling().InsertSibling(sibling)
}

func (me *Node) AppendChild(child *Node) *Node {
	child.Parent = me
	if me.Child == nil {
		me.Child = child
	} else {
		me.Child.AppendSibling(child)
	}
	return child
}

func (me *Node) NewSibling(data interface{}) *Node {
	node := New(me.Parent, data)
	return me.AppendSibling(node)
}

func (me *Node) NewChild(data interface{}) *Node {
	node := New(me, data)
	return me.AppendChild(node)
}

func (me *Node) Walk(walker func(me *Node) *Node) *Node {
	node := me
	if found := walker(node); found != nil {
		return found
	}
	if me.Child != nil {
		if found := me.Child.Walk(walker); found != nil {
			return found
		}
	}
	if me.After != nil {
		if found := me.After.Walk(walker); found != nil {
			return found
		}
	}
	return nil
}
