package dtpp_local

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"github.com/mwortsma/particle_systems/dtmc"
	"math"
)

/*
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
*/


func Run(T,tau int, d int, beta float64, k int, n int) probutil.ContDistr {
	
	fmt.Println("Running")

	// nu uniform
	nu_f := func(v matutil.Vec) float64 {
		prob := 1.0
		for i := 0; i < len(v); i++ {
			prob *= (1.0/float64(k))
		}
		return prob
	}


	Ham := func(r int, v matutil.Vec) float64 {
		match := 0
		sum := r
		for _, val := range v {
			if val == r {
				match++
			}
			sum += val
		}
		return float64(sum - match)
	}



	Q_help := func(r int, s int, v matutil.Vec) float64 {
		prob := 1.0/float64(n*(k-1))
		ham_r := Ham(r,v)
		ham_s := Ham(s,v)
		if ham_r > ham_s {
			prob *= math.Exp(-beta * (ham_r - ham_s))
		}
		return prob
	}

	// prob transition from s to r
	Q := func(r int, s int, v matutil.Vec) float64 {
		if r != s {
			return Q_help(r,s,v)
		} else {
			prob := 1.0
			for i := 0; i < k; i++ {
				if i != s {
					prob -= Q_help(i,s,v)
				}
			}
			return prob
		}
	}

	distr := dtmc.DTMCRegtreeTDistr(T, tau, d, Q, nu_f,k)

	fmt.Println(distr)

	return distr
}
