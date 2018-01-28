package dtcp_mean_field

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"time"
)

func Realization(T int, p, q float64, nu float64, r *rand.Rand) matutil.Vec {

	X := make(matutil.Vec, T)

	// Initial conditions.
	if r.Float64() < nu {
		X[0] = 1
	}

	for t := 1; t < T; t++ {
		X[t] = X[t-1]
		if X[t-1] == 0 {
			// get the sum of the neighbors
			v := Realization(t, p, q, nu, r)
			// transition with probability ratio_neighbors * p
			if r.Float64() < float64(v[t-1])*p {
				X[t] = 1
			}
		} else {
			// if state is 1, transition back with porbability q
			if r.Float64() < q {
				X[t] = 0
			}
		}
	}

	return X
}

func RealizationTypicalDistr(T int, p, q float64, nu float64, steps int) probutil.Distr {
	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	fmt.Println("Running dtcp mean field")
	f := func() fmt.Stringer {
		return Realization(T, p, q, nu, r)
	}
	return probutil.TypicalDistrSync(f, steps)
}
