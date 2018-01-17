package dtlb_mean_field

import (
	"fmt"
	"github.com/mwortsma/particle_systems/probutil"
	"github.com/mwortsma/particle_systems/matutil"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"time"
)

func Realization(T int, lam float64, k int, deg int) matutil.Vec {
	X := make(matutil.Vec, T)

	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Initial conditions.
	p := distuv.Poisson{-math.Log(1 - lam), r}
	X[0] = int(math.Min(p.Rand(), float64(k-1)))

	for t := 1; t < T; t++ {
		// Serve an item if non-empty
		if X[t-1] > 0 {
			X[t] = X[t-1] - 1
		}
		// get samples for all the neighbors
		neighors, _ := getDegPlusOneNeighbors(t, deg, lam, []int{X[t-1]}, false)

		// arrival at my queue
		if r.Float64() < lam && chooseNeighbor(t, deg, lam, neighors, r) {
			X[t]++
		}

		// arrival at neighboring queues
		for i := 1; i < deg+1 && X[t] < k-1; i++  {
			if r.Float64() < lam {
				curr_neighbros := []int{X[t-1], neighors[i]}
				if chooseNeighbor(t, deg, lam, curr_neighbros, r) {
					X[t]++
				}
			}
		}
	}
	return X
}

func getDegPlusOneNeighbors(t int, deg int, lam float64, queues []int, break_val bool) ([]int, bool) {
	for len(queues) < deg + 1 {
		sample := Realization(t,lam,deg,deg)
		if break_val && sample[t-1] < queues[0] {
			return []int{}, false
		}
		queues = append(queues,sample[t-1])
	}
	return queues, true
}

// true if send to neighbor 0
func chooseNeighbor(t int, deg int, lam float64, queues []int, r *rand.Rand) bool {
	queues, b := getDegPlusOneNeighbors(t,deg,lam,queues,true)
	if !b {
		return false
	}
	// First get the min value
	min := queues[0]
	for j := 1; j < deg+1; j++ {
		if queues[j] < min {
			min = queues[j]
		}
	}
	// Select, at random, a neighbor having that value.
	min_queues := make([]int, 0)
	for j := 0; j < deg+1; j++ {
		if queues[j] == min {
			min_queues = append(min_queues, j)
		}
	}
	return min_queues[r.Intn(len(min_queues))] == 0
}

func TypicalDistr(T int, lam float64, k int, deg int, steps int) probutil.Distr {
	f := func() fmt.Stringer {
		return Realization(T, lam, k, deg)
	}
	return probutil.TypicalDistrSync(f, steps)
}