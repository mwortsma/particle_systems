package dtlb_full

import (
	"fmt"
	"github.com/mwortsma/particle_systems/graphutil"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_util"
	"golang.org/x/exp/rand"
	"time"
)

func GraphRealization(T int, p,q float64, k int, G graphutil.Graph, r *rand.Rand) matutil.Mat {
	n := len(G)
	X := matutil.Create(T, n)

	// Initial conditions.
	for i := 0; i < n; i++ {
		X[0][i] = dtlb_util.Init(p,q,k,r)
	}

	for t := 1; t < T; t++ {
		// Copy the state of X[t-1] to X[t]
		copy(X[t], X[t-1])
		for i := 0; i < n; i++ {
			// Serve an item if non-empty
			if X[t-1][i] > 0 && r.Float64() < q {
				X[t][i]--
			}
			// With probability lam there is an arrival.
			// Send to minimum neighbor.
			if r.Float64() < p {
				// First get the min value
				min := X[t-1][i]
				for j := 0; j < len(G[i]); j++ {
					if X[t-1][G[i][j]] < min {
						min = X[t-1][G[i][j]]
					}
				}
				// Select, at random, a neighbor having that value.
				min_neighbors := make([]int, 0)
				if min == X[t-1][i] {
					min_neighbors = append(min_neighbors, i)
				}
				for j := 0; j < len(G[i]); j++ {
					if X[t-1][G[i][j]] == min {
						min_neighbors = append(min_neighbors, G[i][j])
					}
				}
				chosen_neighbor := min_neighbors[r.Intn(len(min_neighbors))]
				if X[t][chosen_neighbor] < k-1 {
					// Only send if below buffer.
					X[t][chosen_neighbor]++
				}
			}
		}
	}

	return X
}

func RingRealization(T int, p,q float64, k int, n int, r *rand.Rand) matutil.Mat {
	return GraphRealization(T, p,q, k, graphutil.Ring(n), r)
}

func CompleteRealization(T int, p,q float64, k int, n int, r *rand.Rand) matutil.Mat {
	return GraphRealization(T, p,q, k, graphutil.Complete(n), r)
}

func RingTypicalDistr(T int, lam, dt float64, k int, n, steps int) probutil.Distr {
	fmt.Println("Running dtlb full ring T =", T, "n = ", n)
	if n < 0 {
		n = 1 + 4*T
	}
	p,q := dtlb_util.GetPQ(lam,dt)
	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	f := func() fmt.Stringer {
		X := RingRealization(T, p,q, k, n,r)
		return X.Col(0)
	}
	return probutil.TypicalDistrSync(f, steps)
}

func CompleteTypicalDistr(T int, lam, dt float64, k int, n, steps int) probutil.Distr {
	fmt.Println("Running dtlb full Complete T =", T, "n = ", n)
	p,q := dtlb_util.GetPQ(lam,dt)
	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	f := func() fmt.Stringer {
		X := CompleteRealization(T, p,q, k, n,r)
		return X.Col(0)
	}
	return probutil.TypicalDistrSync(f, steps)
}
