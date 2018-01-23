package ctcp_full

import (
	//"fmt"
	"github.com/mwortsma/particle_systems/ctmc"
	"github.com/mwortsma/particle_systems/graphutil"
	"github.com/mwortsma/particle_systems/matutil"
	// "github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"time"
)

func GraphRealization(
	T int, 
	lam float64,
	nu float64, 
	G graphutil.Graph,
	k int) ([]float64, matutil.Mat) {

	n := len(G)
	X := matutil.Create(T, n)

	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// keep track of times
	times := make([]float64, T)
	time := 0.0

	// Initial conditions.
	for i := 0; i < n; i++ {
		if r.Float64() < nu {
			X[0][i] = 1
		}
	}

	for t := 1; t < T; t++ {
		rates, events := getRatesAndEvents(X[t-1], lam, G, k)
		event_index, time_inc := ctmc.StepCTMC(rates)
		chosen_event := events[event_index]
		// Copy the state of X[t-1] to X[t]
		copy(X[t], X[t-1])
		X[t][chosen_event.Index] += chosen_event.Inc
		time += time_inc
		times[t] = time
	}

	return times, X
}

func RingRealization(T int, lam float64, nu float64, n int) ([]float64, matutil.Mat) {
	return GraphRealization(T, lam, nu, graphutil.Ring(n), 2)
}

func CompleteRealization(T int, lam float64, nu float64, n int) ([]float64, matutil.Mat) {
	return GraphRealization(T, lam, nu, graphutil.Complete(n), n-1)
}

func getRatesAndEvents(
	X []int, 
	lam float64, 
	G graphutil.Graph,
	k int) ([]float64, []ctmc.Event) {

	rates := make([]float64, 0)
	events := make([]ctmc.Event, 0)

	for i := range(X) {

		if X[i] == 1 {
			// recover
			rates = append(rates, 1)
			events = append(events, ctmc.Event{Index: i, Inc: -1})

			for j := range(G[i]) {
				if X[j] == 0 {
					// infect
					rates = append(rates, lam/float64(k))
					events = append(events, ctmc.Event{Index: j, Inc: 1})
				}
			}
		}
	}
	return rates, events
}
