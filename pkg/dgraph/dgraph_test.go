package dgraph

import (
	"fmt"
	"testing"
)

func TestBFS(t *testing.T) {
	a := &GraphNode{Val: 1}
	b := &GraphNode{Val: 2}
	c := &GraphNode{Val: 3}
	d := &GraphNode{Val: 4}
	a.AddEdge(b)
	a.AddEdge(c)
	b.AddEdge(d)

	BFS(a, func(node *GraphNode) {
		fmt.Println(node.Val)
	})
}

func BFS(root *GraphNode, f func(*GraphNode)) {
	visited := map[*GraphNode]bool{}
	q := Queue{}
	q.Enqueue(root)
	for q.IsNotEmpty() {
		node := q.Dequeue()
		visited[node] = true
		f(node)
		for _, other := range node.Edges {
			if _, ok := visited[other]; ok {
				continue
			}
			q.Enqueue(other)
		}
	}
}

type GraphNode struct {
	Val   int
	Edges []*GraphNode
}

func (gn *GraphNode) AddEdge(node *GraphNode) {
	gn.Edges = append(gn.Edges, node)
}

type Queue struct {
	arr []*GraphNode
}

func (q *Queue) Enqueue(node *GraphNode) {
	q.arr = append(q.arr, node)
}

func (q *Queue) Dequeue() *GraphNode {
	if len(q.arr) == 0 {
		return nil
	}
	node := q.arr[0]
	q.arr[0] = nil
	q.arr = q.arr[1:]
	return node
}

func (q *Queue) IsNotEmpty() bool {
	return len(q.arr) > 0
}
