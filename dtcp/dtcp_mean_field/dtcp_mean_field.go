package dtcp_mean_field

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"time"
)

func Realization(T int, p, q float64, nu float64, deg int) matutil.Vec {

	X := make(matutil.Vec, T)

	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Initial conditions.
	if r.Float64() < nu {
		X[0] = 1
	}


	for t := 1; t < T; t++ {
		X[t] = X[t-1]
		if X[t-1] == 0 {
			// get the sum of the neighbors
			sum_neighbors := 0
			for j := 0; j < deg; j++ {
				sample := Realization(t,p,q,nu,deg)
				sum_neighbors += sample[t-1]
			}
			// transition with probability (p/deg)*sum_neighbors
			if r.Float64() < (p/float64(deg))*float64(sum_neighbors) {
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

func TypicalDistr(T int, p, q float64, nu float64, deg, steps int) probutil.Distr {
	f := func() fmt.Stringer {
		return Realization(T, p, q, nu, deg)
	}
	return probutil.TypicalDistrSync(f, steps)
}
