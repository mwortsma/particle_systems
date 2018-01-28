package ctcp_full

import (
	"fmt"
	"github.com/mwortsma/particle_systems/ctmc"
	"github.com/mwortsma/particle_systems/graphutil"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"time"
)

func GraphRealization(
	T float64, 
	lam float64,
	nu float64, 
	G graphutil.Graph,
	k int,
	r *rand.Rand) ([]float64, matutil.Mat) {

	n := len(G)
	X := make([][]int, 1)
	X[0] = make([]int, n)



	// keep track of times
	times := make([]float64, 1)
	t := 0.0

	// Initial conditions.
	for i := 0; i < n; i++ {
		if r.Float64() < nu {
			X[0][i] = 1
		}
	}

	for  {
		rates, events := ctmc.GetCPRatesAndEvents(X[len(X)-1], lam, G, k)
		if len(rates) == 0 {
			break
		}
		event_index, time_inc := ctmc.StepCTMC(rates)
		chosen_event := events[event_index]
		t += time_inc
		if t >= T {
			break
		}
		times = append(times, t)
		newstate := make([]int, n)
		copy(newstate, X[len(X)-1])
		newstate[chosen_event.Index] += chosen_event.Inc
		X = append(X, newstate)
	}

	return times, X
}

func RingRealization(T float64, lam float64, nu float64, n int, r *rand.Rand) ([]float64, matutil.Mat) {
	return GraphRealization(T, lam, nu, graphutil.Ring(n), 2, r)
}

func CompleteRealization(T float64, lam float64, nu float64, n int, r *rand.Rand) ([]float64, matutil.Mat) {
	return GraphRealization(T, lam, nu, graphutil.Complete(n), n-1, r)
}

func RingTypicalDistr(T float64, lam float64, nu, dt float64, n, steps int) probutil.ContDistr {
	fmt.Println("Running ctcp Full Ring n =", n)
	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	f := func() ([]float64, matutil.Vec) {
		times, X := RingRealization(T, lam, nu, n,r)
		return times, X.Col(0)
	}
	return probutil.TypicalContDistrSync(f, dt, T, 2, steps)
}

func CompleteTypicalDistr(T float64, lam float64, nu, dt float64, n, steps int) probutil.ContDistr {
	fmt.Println("Running ctcp Full Complete n =", n)
	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	f := func() ([]float64, matutil.Vec) {
		times, X := CompleteRealization(T, lam, nu, n,r)
		return times, X.Col(0)
	}
	return probutil.TypicalContDistrSync(f, dt, T, 2, steps)
}

