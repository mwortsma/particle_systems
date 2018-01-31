package dtlb_local

import (
	"fmt"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_util"
	"github.com/mwortsma/particle_systems/graphutil"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"time"
)

func LocalRingRealization(T int, p, q float64, k int, tau int, match_cols []int, match_vals matutil.Mat, r *rand.Rand) (matutil.Mat, bool) {
	// n is how many nodes we need to keep track of.
	n := 11
	G := graphutil.Ring(n)
	X := matutil.Create(T, n)

	for i := 0; i < n; i++ {
		X[0][i] = dtlb_util.Init(p, q, k, r)
	}

	
	for i, c := range match_cols {
		if tau <= 0 || T-tau <= 0 {
			X[0][c] = match_vals[0][i]
		}
	}


	for t := 1; t < T; t++ {
		// Copy the state of X[t-1] to X[t]
		copy(X[t], X[t-1])
		// obtain a vector of arrivals
		arrivals := make([]bool, n)
		for i := 1; i < n-1; i++ {
			arrivals[i] = r.Float64() < p
		}

		if arrivals[1] || arrivals[2] {
			// if there is an arrival at 1,2 we need to sample the conditional
			var left_match_vals matutil.Mat
			if tau <= 0 || t-tau < 0 {
				left_match_vals = X.ColsT([]int{5, 4, 3, 2}, t)
			} else {
				left_match_vals = X.ColsRange([]int{5, 4, 3, 2}, t-tau, t)
			}
			left_match_cols := []int{2, 3, 4, 5}
			var Y_left matutil.Mat
			var b_left bool

			for !b_left {
				Y_left, b_left = LocalRingRealization(t, p, q, k, tau, left_match_cols, left_match_vals, r)
			}
			X[t-1][0], X[t-1][1] = Y_left[t-1][5], Y_left[t-1][4]
		}

		if arrivals[8] || arrivals[9] {
			// if there is an arrival at 1,2 we need to sample the conditional
			var right_match_vals matutil.Mat
			if tau <= 0 || t-tau < 0 {
				right_match_vals = X.ColsT([]int{5, 6, 7, 8}, t)
			} else {
				right_match_vals = X.ColsRange([]int{5, 6, 7, 8}, t-tau, t)
			}
			
			right_match_cols := []int{8, 7, 6, 5}
			var Y_right matutil.Mat
			var b_right bool

			for !b_right {
				Y_right, b_right = LocalRingRealization(t, p, q, k, tau, right_match_cols, right_match_vals, r)
			}
			X[t-1][9], X[t-1][10] = Y_right[t-1][2], Y_right[t-1][1]
		}

		for i := 1; i < n-1; i++ {
			// Serve an item if non-empty
			if X[t-1][i] > 0 && r.Float64() < q {
				X[t][i]--
			}
			// With probability lam there is an arrival.
			// Send to minimum neighbor.
			if arrivals[i] {
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

		if !X.Match(match_cols, match_vals, t, tau) {
			return [][]int{}, false
		}
	}

	return X.Cols([]int{2, 3, 4, 5, 6, 7, 8}), true
}

func LocalRingRealizationTypicalDistr(T int, lam, dt float64, k int, steps int, tau int) probutil.Distr {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	p, q := dtlb_util.GetPQ(lam, dt)
	f := func() fmt.Stringer {
		X, _ := LocalRingRealization(T, p, q, k, tau, nil, nil, r)
		return X.Col(3)
	}
	return probutil.TypicalDistrSync(f, steps)
}
