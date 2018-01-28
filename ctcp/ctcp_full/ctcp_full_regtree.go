package ctcp_full

import (
	"github.com/mwortsma/particle_systems/ctmc"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"time"
)

type node struct {
	children []*node
	parent   *node
	state    matutil.Vec
	times 	 []float64
	is_leaf  bool
	is_root  bool
	id int
	num_neighbors int
}

func RegTreeRealization(depth int, T, lam float64, d int, nu float64, r *rand.Rand) ([]float64, matutil.Vec) {
	// create tree
	var root node
	root.createNode(T, d, nu, &node{}, depth, true, r, 0)
	t := 0.0

	for {
		rates, events := GetCPTreeRatesAndEvents(&root, lam)
		if len(rates) == 0 {
			break
		}
		event_index, time_inc := ctmc.StepCTMC(rates)
		chosen_event := events[event_index]
		t += time_inc

		if t >= T {
			break
		}
		root.transition(t, chosen_event)
	}
	return root.times, root.state
}

func RegTreeTypicalDistr(depth int, T, lam float64, d int, nu float64,  dt float64, steps int) probutil.ContDistr {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	f := func() ([]float64, matutil.Vec) {
		return RegTreeRealization(depth, T, lam, d, nu, r)
	}
	return probutil.TypicalContDistrSync(f, dt, T, 2, steps)
}

// Helpers
func (n *node) createNode(
	T float64,
	d int,
	nu float64,
	parent *node,
	depth int,
	is_root bool,
	r *rand.Rand,
	id int) int {

	n.id = id
	n.num_neighbors = 0

	// set parent
	n.is_root = is_root
	if !n.is_root {
		n.parent = parent
		n.num_neighbors += 1
	}
	// create children
	if depth == 0 {
		n.is_leaf = true
	} else {
		n.children = make([]*node, d-1)
		for c := 0; c < d-1; c++ {
			var child node
			id = child.createNode(T, d, nu, n, depth-1, false, r, id+1)
			n.children[c] = &child
			n.num_neighbors += 1
		}
		if n.is_root {
			var child node
			id = child.createNode(T, d, nu, n, depth-1, false, r, id+1)
			n.children = append(n.children, &child)
			n.num_neighbors += 1
		}
	}
	// create state
	n.state = make(matutil.Vec, 1)
	n.times = make([]float64, 1)
	// initial conditions
	if r.Float64() < nu {
		n.state[0] = 1
	}

	return id + 1
}

func (n *node) transition(t float64, e ctmc.Event) bool {
	if n.id != e.Index {
		// call transition on children
		for _, c := range n.children {
			if c.transition(t, e) {
				return true
			}
		}
		return false
	}
	n.times = append(n.times, t)
	n.state = append(n.state, n.state[len(n.state)-1] + e.Inc)
	return true
}


func GetCPTreeRatesAndEvents(
	n *node, 
	lam float64) ([]float64, []ctmc.Event) {

	rates := make([]float64, 0)
	events := make([]ctmc.Event, 0)

	if n.state[len(n.state) - 1] == 1 {
		// recover
		rates = append(rates, 1)
		events = append(events, ctmc.Event{Index: n.id, Inc: -1})

		for _, c := range n.children {
			if c.state[len(c.state) - 1] == 0 {
				// infect
				rates = append(rates, lam/float64(c.num_neighbors))
				events = append(events, ctmc.Event{Index: c.id, Inc: 1})
			}
		}
	} else {
		for _, c := range n.children {
			if c.state[len(c.state) - 1] == 1 {
				// infect
				rates = append(rates, lam/float64(n.num_neighbors))
				events = append(events, ctmc.Event{Index: n.id, Inc: 1})
			}
		}
	}

	for _, c := range n.children {
		rec_rates, rec_events := GetCPTreeRatesAndEvents(c, lam)
		rates = append(rates, rec_rates...)
		events = append(events, rec_events...)
	}

	return rates, events
}

