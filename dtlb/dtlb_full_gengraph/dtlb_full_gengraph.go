package dtlb_full_gengraph

import (
	"github.com/mwortsma/particle_systems/graphutil"
	"github.com/mwortsma/particle_systems/matutil"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"time"
	//"fmt"
	"math"
)

func Realization(T int, lam float64, k int, G graphutil.Graph) matutil.Mat {
	n := len(G)
	X := matutil.Create(T, n)

	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Initial conditions.
	p := distuv.Poisson{-math.Log(1 - lam), r}
	for i := 0; i < n; i++ {
		X[0][i] = int(math.Min(p.Rand(), float64(k-1)))
	}

	for t := 1; t < T; t++ {
		// Copy the state of X[t-1] to X[t]
		copy(X[t], X[t-1])
		for i := 0; i < n; i++ {
			// Serve an item if non-empty
			if X[t-1][i] > 0 {
				X[t][i]--
			}
			// With probability lam there is an arrival.
			// Send to minimum neighbor.
			if r.Float64() < lam {
				// First get the min value
				min := X[t-1][G[i][0]]
				for j := 1; j < len(G[i]); j++ {
					if X[t-1][G[i][j]] < min {
						min = X[t-1][G[i][j]]
					}
				}
				// Select, at random, a neighbor having that value.
				min_neighbors := make([]int, 0)
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

func RingRealization(T int, lam float64, k int, n int) matutil.Mat {
	return Realization(T, lam, k, graphutil.Ring(n))
}

func CompleteRealization(T int, lam float64, k int, n int) matutil.Mat {
	return Realization(T, lam, k, graphutil.Complete(n))
}

func Run() {
	T, lam, k, n := 100, 0.95, 6, 100
	RingRealization(T, lam, k, n)
	//matutil.PrintMat(RingRealization(T, lam, k, n))
}
