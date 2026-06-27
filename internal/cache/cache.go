package cache

import "fmt"

const size = 30

type Value struct {
	Query      string
	Embeddings []float64
}

type Node struct {
	Val   Value
	Left  *Node
	Right *Node
}

type Queue struct {
	Head   *Node
	Tail   *Node
	Length int
}

type Hash map[string]*Node

type Cache struct {
	Queue Queue
	Hash  Hash
}

func NewCache() Cache {
	head := &Node{}
	tail := &Node{}
	head.Right = tail
	tail.Left = head
	return Cache{
		Queue: Queue{Head: head, Tail: tail},
		Hash:  Hash{},
	}
}

func (c *Cache) Check(query string, embeddings []float64) {
	node := &Node{}
	if val, ok := c.Hash[query]; ok {
		node = c.Remove(val)
	} else {
		node = &Node{Val: Value{Query: query, Embeddings: embeddings}}
	}
	c.Add(node)
	c.Hash[query] = node
}

func (c *Cache) Remove(n *Node) *Node {
	fmt.Printf("Remove %s\n", n.Val.Query)
	left := n.Left
	right := n.Right
	right.Left = left
	left.Right = right
	c.Queue.Length--
	delete(c.Hash, n.Val.Query)
	return n
}

func (c *Cache) Add(n *Node) {
	fmt.Printf("Add Node %s\n", n.Val.Query)
	temp := c.Queue.Head.Right
	c.Queue.Head.Right = n
	n.Left = c.Queue.Head
	n.Right = temp
	temp.Left = n
	c.Queue.Length++
	if c.Queue.Length > size {
		c.Remove(c.Queue.Tail.Left)
	}
}

func (c *Cache) Get(query string) ([]float64, bool) {
	if val, ok := c.Hash[query]; ok {
		return val.Val.Embeddings, true
	}
	return nil, false
}

func (c *Cache) Display() {
	node := c.Queue.Head.Right
	fmt.Printf("%d - [", c.Queue.Length)
	for i := 0; i < c.Queue.Length; i++ {
		fmt.Printf("{%s}", node.Val.Query)
		if i < c.Queue.Length-1 {
			fmt.Print("<-->")
		}
		node = node.Right
	}
	fmt.Printf("]\n")
}