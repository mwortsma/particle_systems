package dtcp_full

import (
	"fmt"
	"github.com/mwortsma/particle_systems/graphutil"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"time"
)

func GraphRealization(T int, p, q float64, nu float64, G graphutil.Graph) matutil.Mat {
	n := len(G)
	X := matutil.Create(T, n)

	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Initial conditions.
	for i := 0; i < n; i++ {
		if r.Float64() < nu {
			X[0][i] = 1
		}
	}

	for t := 1; t < T; t++ {
		for i := 0; i < n; i++ {
			X[t][i] = X[t-1][i]
			if X[t-1][i] == 0 {
				// get the sum of the neighbors
				sum_neighbors := 0
				for j := 0; j < len(G[i]); j++ {
					sum_neighbors += X[t-1][G[i][j]]
				}
				// transition with probability (p/deg)*sum_neighbors
				if r.Float64() < (p/float64(len(G[i])))*float64(sum_neighbors) {
					X[t][i] = 1
				}
			} else {
				// if state is 1, transition back with porbability q
				if r.Float64() < q {
					X[t][i] = 0
				}
			}
		}
	}

	return X
}

func RingRealization(T int, p, q float64, nu float64, n int) matutil.Mat {
	return GraphRealization(T, p, q, nu, graphutil.Ring(n))
}

func CompleteRealization(T int, p, q float64, nu float64, n int) matutil.Mat {
	return GraphRealization(T, p, q, nu, graphutil.Complete(n))
}

func RingTypicalDistr(T int, p, q float64, nu float64, n, steps int) probutil.Distr {
	if n < 0 {
		n = 1 + 2*T
	}
	fmt.Println("Running dtcp Full Ring n=", n)
	f := func() fmt.Stringer {
		X := RingRealization(T, p, q, nu, n)
		return X.Col(0)
	}
	return probutil.TypicalDistrSync(f, steps)
}

func CompleteTypicalDistr(T int, p, q float64, nu float64, n, steps int) probutil.Distr {
	fmt.Println("Running dtcp Full Complete n=", n)
	f := func() fmt.Stringer {
		X := CompleteRealization(T, p, q, nu, n)
		return X.Col(0)
	}
	return probutil.TypicalDistrSync(f, steps)
}
