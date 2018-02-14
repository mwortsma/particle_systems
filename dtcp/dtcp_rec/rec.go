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
				return 1-(p/float64(d))*float64(sum_neighbors)
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

func JointProb(T,tau int,d int, j []probutil.Distr,state matutil.Mat) float64{
	prob := 1.0
	t := T-1
	curr_j := j[t]
	lastrows := matutil.BinaryStrings(d+1)
	for len(state) > tau + 1 {
		denom := 0.0
		rel_state := state[len(state)-(tau+1):len(state)-1]
		for _, lastrow := range lastrows {
			full := append(rel_state, lastrow)
			denom += curr_j[full.String()]
		}
		if denom > 0 {
			prob *= curr_j[state[len(state)-(tau+1):].String()]/denom
		} else {
			return 0.0
		}
		t = t - 1
		curr_j = j[t]
		state = state[0:len(state)-1]
	}
	prob *= curr_j[state.String()]
	return prob

}

func Run(T,tau int, d int, p,q float64, nu float64) probutil.Distr {
	
	fmt.Println("Running")
	nu_f := getNu(nu,d)
	Q := getQ(p,q,d)

	j_array := dtmc.DTMCRegtreeRecursionsFull(T, tau, d, Q, nu_f)

	for _, j := range j_array {
		fmt.Println(j)
		fmt.Println("\n")
	}

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
		f[path] += JointProb(T,tau,d,j_array,state)
	}

	return f

}