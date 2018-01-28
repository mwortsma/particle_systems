package dtcp_local

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"time"
)

func LocalRegTreeRealization(T, d int, p, q float64, nu float64, match_cols []int, match_vals matutil.Mat, r *rand.Rand) (matutil.Mat, bool) {
	// n is how many nodes we need to keep track of.
	// X[0] stores the root. X[1:d+1] store the children.
	X := matutil.Create(T, d+1)

	// Initial conditions.
	for i := 0; i < d+1; i++ {
		if r.Float64() < nu {
			X[0][i] = 1
		}
	}

	for i, c := range match_cols {
		X[0][c] = match_vals[0][i]
	}

	for t := 1; t < T; t++ {
		// Copy the state of X[t-1] to X[t]
		copy(X[t], X[t-1])

		// update the root.
		if X[t-1][0] == 0 {
			sum_neighbors := 0
			for j := 1; j < d+1; j++ {
				sum_neighbors += X[t-1][j]
			}
			if r.Float64() < (p/float64(d))*float64(sum_neighbors) {
				X[t][0] = 1
			}

		} else {
			// transition to 0 with probability q
			if r.Float64() < q {
				X[t][0] = 0
			}
		}

		// update the rest of the neighbors.
		for i := 1; i < d+1; i++ {
			if X[t-1][i] == 0 {

				var Y matutil.Mat
				var b bool
				rec_match_cols := []int{i, 0}
				rec_match_vals := X.ColsT([]int{0, i}, t)

				for !b {
					Y, b = LocalRegTreeRealization(t, d, p, q, nu, rec_match_cols, rec_match_vals, r)
				}

				sum_neighbors := 0
				for j := 1; j < d+1; j++ {
					sum_neighbors += Y[t-1][j]
				}

				if r.Float64() < (p/float64(d))*float64(sum_neighbors) {
					X[t][i] = 1
				}

			} else {
				// transition to 0 with probability q
				if r.Float64() < q {
					X[t][i] = 0
				}
			}
		}

		if !X.Match(match_cols, match_vals, t) {
			return [][]int{}, false
		}
	}

	return X, true
}

func LocalRegTreeRealizationTypicalDistr(T, d int, p, q float64, nu float64, steps int) probutil.Distr {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	f := func() fmt.Stringer {
		X, _ := LocalRegTreeRealization(T, d, p, q, nu, nil, nil, r)
		return X.Col(1)
	}
	return probutil.TypicalDistr(f, steps)
}
