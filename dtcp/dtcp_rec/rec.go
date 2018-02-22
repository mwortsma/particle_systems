package dtcp_rec

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"github.com/mwortsma/particle_systems/dtmc"
	"math"
)

func getQ(p,q float64,d int) func(int,int,matutil.Vec)float64 {
	return func(k int, s int, v matutil.Vec) float64 {
		if s == 1 {
			if k == 1 {
				return 1-q
			} else {
				return q
			}
		} else {
			sum_neighbors := 0
			for i := 0; i < len(v); i++ {
				sum_neighbors += v[i]
			}
			if k == 1 {
				return (p/float64(d))*float64(sum_neighbors)
			} else {
				return 1.0-(p/float64(d))*float64(sum_neighbors)
			}
		}
		return 0.0
	}
}

func getNu(nu float64,d int) func(v matutil.Vec)float64 {
	return func(v matutil.Vec) float64 {
		prob := 1.0
		for i := 0; i < len(v); i++ {
			if v[i] == 1 {
				prob *= (1-nu)
			} else {
				prob *= nu
			}
		}
		return prob
	}
}

func JointProb(T,tau int,d int, j []probutil.Distr,p []probutil.Distr, state matutil.Mat) float64{
	prob := 1.0
	t := T-1
	for len(state) > tau + 1 {
		rel_state := state[len(state)-(tau+2):]
		prob *= p[t][rel_state.String()]
		t = t - 1
		state = state[0:len(state)-1]
	}
	prob *= j[t][state.String()]
	return prob

}

func FullRun(T,tau int, d int, p,q float64, nu float64) probutil.Distr {
	
	fmt.Println("Running")
	nu_f := getNu(nu,d)
	Q := getQ(p,q,d)

	j_array, p_array := dtmc.DTMCRegtreeRecursionsFull(T, tau, d, Q, nu_f,2)

	f := make(probutil.Distr)

	if tau < 0 {
		tau = math.MaxInt32
	}
	fmt.Println("T is ", T)

	states := matutil.BinaryMats(T, d+1)
	for _, state := range states {
		path := state.Col(0).String()
		if _, ok := f[path]; !ok {
			f[path] = 0.0
		}
		f[path] += JointProb(T,tau,d,j_array,p_array,state)
	}

	return f
}

func Run(T,tau int, d int, p,q float64, nu float64) probutil.ContDistr {
	
	fmt.Println("Running")

	nu_f := func(v matutil.Vec) float64 {
		prob := 1.0
		for i := 0; i < len(v); i++ {
			if v[i] == 1 {
				prob *= (1-nu)
			} else {
				prob *= nu
			}
		}
		return prob
	}

	Q := func(k int, s int, v matutil.Vec) float64 {
		if s == 1 {
			if k == 1 {
				return 1-q
			} else {
				return q
			}
		} else {
			sum_neighbors := 0
			for i := 0; i < len(v); i++ {
				sum_neighbors += v[i]
			}
			if k == 1 {
				return (p/float64(d))*float64(sum_neighbors)
			} else {
				return 1.0-(p/float64(d))*float64(sum_neighbors)
			}
		}
		return 0.0
	}

	return dtmc.DTMCRegtreeTDistr(T, tau, d, Q, nu_f,2)
}
