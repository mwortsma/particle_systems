package graphutil

import (
	"bytes"
	"fmt"
)

type Graph [][]int

// Ring makes a ring Graph with n nodes.
func Ring(n int) Graph {
	G := make(Graph, n)
	for i := 0; i < n; i++ {
		G[i] = []int{(i - 1 + n) % n, (i + 1) % n}
	}
	return G
}

// Complete makes a complete Graph with n nodes.
func Complete(n int) Graph {
	G := make(Graph, n)
	for i := 0; i < n; i++ {
		G[i] = make([]int, 0)
		for j := 0; j < n; j++ {
			if j != i {
				G[i] = append(G[i], j)
			}
		}
	}
	return G
}

// String gets a description of G
func (G Graph) Print() {
	var buffer bytes.Buffer
	for node, neighbors := range G {
		buffer.WriteString(fmt.Sprintf("Node %d: %v\n", node, neighbors))
	}
	fmt.Println(buffer.String())
}
