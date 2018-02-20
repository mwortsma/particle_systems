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

func GetShared(T,tau int, d int, beta float64, k,n int) (func(matutil.Vec)float64, func(int,int,matutil.Vec)float64,  func(int,int,matutil.Vec)float64) {

		// nu uniform
	nu_f := func(v matutil.Vec) float64 {
		prob := 1.0
		for i := 0; i < len(v); i++ {
			prob *= (1.0/float64(k))
		}
		return prob
	}


	Ham := func(r, n int, v matutil.Vec) float64 {
		match := 0
		sum := r
		for _, val := range v {
			if val == r {
				match++
			}
			sum += val
		}
		return float64(sum) - float64(match)
	}



	Q_help := func(r int, s int, n int, v matutil.Vec) float64 {
		prob := 1.0/float64(n*(k-1))
		ham_r := Ham(r,n,v)
		ham_s := Ham(s,n,v)
		if ham_r > ham_s {
			prob *= math.Exp(-beta * (ham_r - ham_s))
		}
		return prob
	}

	// prob transition from s to r
	Q := func(r int, s int, v matutil.Vec) float64 {
		if r != s {
			return Q_help(r,s,n,v)
		} else {
			prob := 1.0
			for i := 0; i < k; i++ {
				if i != s {
					prob -= Q_help(i,s,n,v)
				}
			}
			return prob
		}
	}

	return nu_f, Ham, Q
}


func Run(T,tau int, d int, beta float64, k int, n int) probutil.ContDistr {
	
	fmt.Println("Running")

	nu_f, _, Q := GetShared(T,tau,d,beta,k,n)

	distr := dtmc.DTMCRegtreeTDistr(T, tau, d, Q, nu_f,k)

	fmt.Println(distr)

	return distr
}

func EndRun(T,tau int, d int, beta float64, k int, n int) probutil.Distr {
	
	fmt.Println("Running")

	nu_f, _, Q := GetShared(T,tau,d,beta,k,n)

	distr := dtmc.DTMCRegtreeEndDistr(T, tau, d, Q, nu_f,k)

	fmt.Println(distr)

	return distr
}

func FullEndRun(T,tau int, d int, beta float64, k int, n int) probutil.Distr {
	
	fmt.Println("Running on a ring")

	H := func(v matutil.Vec) float64 {
		match := 0
		sum := v[0]
		if v[0] != v[n-1] {
			match++
		}
		for i := 1; i < n; i++ {
			sum += v[i]
			if v[i-1] != v[i] {
				match++
			}
		}
		return float64(sum) - float64(match)
	}

	distr := make(probutil.Distr)

	strings := matutil.QStrings(n,k)

	for _, str := range strings {
		distr[str[0:1].String()] += math.Exp(-beta*H(str))
	}

	sum_distr := 0.0
	for _, val := range distr {
		sum_distr += val
	}

	for k, _ := range distr {
		distr[k] /= sum_distr
	}

	return distr
}

