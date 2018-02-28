package dtcp_mean_field

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"time"
)

func getExpectation(T int, p,q float64, nu float64) []float64 {
	f := make([]float64, T)
	f[0] = nu
	for i := 0; i < T; i++ {
		f_new := make([]float64, T)
		f_new[0] = nu
		for t := 1; t < i+1; t++ {
			f_new[t] = f_new[t-1]*(1-q) + (1-f_new[t-1])*(p/2.0)*f[t-1]
		}
		f = f_new
	}
	return f
}


func MeanFieldRecursionRealization(T int, p,q float64, nu float64, r *rand.Rand, f []float64) matutil.Vec {

	X := make(matutil.Vec, T)

	// Initial conditions.
	if r.Float64() < nu {
		X[0] = 1
	}

	for t := 1; t < T; t++ {
		X[t] = X[t-1]
		if X[t-1] == 0 {
			// transition with probability ratio_neighbors * p
			if r.Float64() < f[t-1]*p {
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

func RecursionTypicalDistrByt(T int, p, q float64, nu float64, steps int) probutil.ContDistr {
	k := 2
	g := getExpectation(T,p,q,nu)
	fmt.Println(g)
	f := make([][]float64, T)
	for t := 0; t < T; t++ {
		f[t] = make([]float64, k)
		f[t][0] = 1-g[t]
		f[t][1] = g[t]
	}

	return probutil.ContDistr{1, float64(T), k, f}
}


func RecursionTypicalDistr(T int, p, q float64, nu float64, steps int) probutil.Distr {
	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	g := getExpectation(T,p,q,nu)
	fmt.Println("Running dtcp mean field")
	f := func() fmt.Stringer {
		return MeanFieldRecursionRealization(T, p, q, nu, r, g)
	}
	return probutil.TypicalDistrSync(f, steps)
}