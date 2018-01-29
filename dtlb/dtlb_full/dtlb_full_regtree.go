package dtlb_full

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_util"
	"golang.org/x/exp/rand"
	"time"
)

type node struct {
	children []*node
	parent   *node
	state    matutil.Vec
	is_leaf  bool
	is_root  bool
}

func RegTreeRealization(T, d int, lam,dt float64, k int, r *rand.Rand) matutil.Vec {

	// create tree
	var root node
	// TODO verify this depth
	p,q := dtlb_util.GetPQ(lam,dt)
	root.createNode(T, d, p,q,k,r, &node{}, 2*T-1, true)

	for t := 1; t < T; t++ {
		// transition will be called for the whole tree recursively
		root.transition(t, d, p, q, k, r)
	}

	return root.state
}

func RegTreeTypicalDistr(T, d int, p,q float64, k int, steps int) probutil.Distr {
	
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	// Ger random nummber to be used throughout
	f := func() fmt.Stringer {
		return RegTreeRealization(T, d, p,q,k,r)
	}
	return probutil.TypicalDistr(f, steps)
}

// Helpers
func (n *node) createNode(
	T int,
	d int,
	p,q float64,
	k int, 
	r *rand.Rand,
	parent *node,
	depth int,
	is_root bool) {

	// set parent
	n.is_root = is_root
	if !n.is_root {
		n.parent = parent
	}
	// create children
	if depth == 0 {
		n.is_leaf = true
	} else {
		n.children = make([]*node, d-1)
		for c := 0; c < d-1; c++ {
			var child node
			child.createNode(T, d, p,q,k,r, n, depth-1, false)
			n.children[c] = &child
		}
		if n.is_root {
			var child node
			child.createNode(T, d, p,q,k,r, n, depth-1, false)
			n.children = append(n.children, &child)
		}
	}

	// create state
	n.state = make(matutil.Vec, T)

	// Initial conditions.
	n.state[0] = dtlb_util.Init(p,q,k,r)
}

func (n *node) transition(t, d int, p,q float64, k int, r *rand.Rand) {

	n.state[t] = n.state[t-1]
	// call transition on children
	for _, c := range n.children {
		c.transition(t, d, p,q,k,r)
	}

	// serve an item with probability q
	if n.state[t-1] > 0 && r.Float64() < q {
		n.state[t]--
	}

	// incoming iter with probability p
	if r.Float64() < p {

		// First get the min value
		min := n.state[t-1]
		if !n.is_root && n.parent.state[t-1] < n.state[t-1] {
			min = n.parent.state[t-1]
		}
		for _, c := range n.children {
			if c.state[t-1] < min {
				min = c.state[t-1]
			}
		}

		// Select, at random, a neighbor having that value.
		min_neighbors := make([]*node, 0)
		if min == n.state[t-1] {
			min_neighbors = append(min_neighbors, n)
		}
		if !n.is_root && min == n.parent.state[t-1] {
			min_neighbors = append(min_neighbors, n.parent)
		}
		for _, c := range n.children {
			if c.state[t-1] == min {
				min_neighbors = append(min_neighbors, c)
			}
		}

		chosen_neighbor := min_neighbors[r.Intn(len(min_neighbors))]
		if chosen_neighbor.state[t] < k-1 {
			// Only send if below buffer.
			chosen_neighbor.state[t]++
		}
	}

}


