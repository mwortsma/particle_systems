package dtcp_full

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"time"
)

type node struct {
	children []*node
	parent *node
	state matutil.Vec
	is_leaf bool
	is_root bool
}

func RegTreeRealization(T,d int, p,q float64, nu float64) matutil.Vec {
	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// create tree
	var root node
	root.createNode(T,d,nu,&node{},T-1,true,r)

	for t := 1; t < T; t++ {
		// transition will be called for the whole tree recursively
		root.transition(t,d,p,q,r)
	}

	return root.state
}


func RegTreeTypicalDistr(T,d int, p,q float64, nu float64,steps int) probutil.Distr {
	f := func() fmt.Stringer {
		return RegTreeRealization(T,d,p,q,nu)
	}
	return probutil.TypicalDistrSync(f, steps)
}

// Helpers
func (n *node) createNode(
	T int,
	d int,
	nu float64,
	parent *node,
	depth int, 
	is_root bool, 
	r *rand.Rand){

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
			child.createNode(T,d,nu,n,depth-1,false,r)
			n.children[c] = &child
		}
		if n.is_root {
			var child node
			child.createNode(T,d,nu,n,depth-1,false,r)
			n.children = append(n.children, &child)
		}
	}
	// create state
	n.state = make(matutil.Vec,T)
	// initial conditions
	if r.Float64() < nu {
		n.state[0] = 1
	}
}

func (n *node) transition(t,d int, p,q float64, r *rand.Rand) {
	n.state[t] = n.state[t-1]
	if n.state[t-1] == 0 {
		// get the sum of the neighbors
		sum_neighbors := 0
		if !n.is_root {
			sum_neighbors += n.parent.state[t-1]
		}
		for _, c := range(n.children){
			sum_neighbors += c.state[t-1]
		}
		// transition with probability p/sum_neighbors
		if r.Float64() < (p/float64(d))*float64(sum_neighbors) {
			n.state[t] = 1
		}
	} else {
		// if state is 1, transition back with porbability q
		if r.Float64() < p {
			n.state[t] = 0
		}
	}
	// call transition on children
	for _, c := range(n.children){
		c.transition(t,d,p,q,r)
	}
}
