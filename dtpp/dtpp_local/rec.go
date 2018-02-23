package dtpp_local

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"github.com/mwortsma/particle_systems/dtmc"
	"math"
	"time"
	"golang.org/x/exp/rand"
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

	return distr
}

func EndRun(T,tau int, d int, beta float64, k int, n int) probutil.Distr {
	
	fmt.Println("Running")

	nu_f, _, Q := GetShared(T,tau,d,beta,k,n)

	distr := dtmc.DTMCRegtreeEndDistr(T, tau, d, Q, nu_f,k)

	return distr
}

func FullEndRun(T,tau int, d int, beta float64, k int, n int) probutil.Distr {
	
	fmt.Println("Running on a ring")

	H := func(v matutil.Vec) float64 {
		match := 0
		sum := v[0]
		if v[0] == v[n-1] {
			match++
		}
		for i := 1; i < n; i++ {
			sum += v[i]
			if v[i-1] == v[i] {
				match++
			}
		}
		return float64(sum) - float64(match)
	}

	distr := make(probutil.Distr)
	sum_distr := 0.0
	strings := matutil.QStrings(n,k)

	for _, s := range strings {
		str := matutil.Vec([]int{s[1], s[0], s[2]}).String()
		prob_s := math.Exp(-beta*H(s))
		if _, ok := distr[str]; !ok {
			distr[str] = prob_s
		} else {
			distr[str] += prob_s
		}
		sum_distr += prob_s
	}


	for k, v := range distr {
		distr[k] = v/sum_distr
		fmt.Println(k)
		fmt.Println(distr[k])
	}

	fmt.Println(distr)
	return distr
}

func MCMC_byt(T,tau int, d int, beta float64, k int, n int, steps int) probutil.ContDistr {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	t_array := make([]float64, T)
	for i := 0; i < T; i++ {
		t_array[i] = float64(i)
	}

	f := func() ([]float64, matutil.Vec) {
		realization := MCMC_realization(T,tau,d,beta,k,n,r)
		return t_array, realization.Col(n/2)
	}
	return probutil.TypicalContDistrSync(f, 1, float64(T), 3, steps)
}

func MCMC_end(T,tau int, d int, beta float64, k int, n int, steps int) probutil.Distr {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	f := func() fmt.Stringer {
		realization := MCMC_realization(T,tau,d,beta,k,n,r)
		return matutil.Vec([]int{realization[T-1][1], realization[T-1][0], realization[T-1][2]})
	}
	return probutil.TypicalDistrSync(f, steps)
}

func MCMC_realization(T,tau int, d int, beta float64, k int, n int,r *rand.Rand) matutil.Mat {

	state := matutil.Create(T,n)

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

	for i := 0; i < n; i++ {
		random := r.Float64()
		if random < (1.0/float64(k)) {
			state[0][i] = 0
		} else if random < (2.0/float64(k)) {
			state[0][i] = 1
		} else {
			state[0][i] = 2
		}
	}


	for t := 1; t < T; t++ {
		copy(state[t], state[t-1])
		site := rand.Intn(n)
		new_value := rand.Intn(k)
		for new_value == state[t-1][site] {
			new_value = rand.Intn(k)
		}
		neighbors := make([]int, 0)
		if site > 0 {
			neighbors = append(neighbors, state[t-1][site-1])
		} else {
			neighbors = append(neighbors, state[t-1][n-1])
		}
		if site < n-1 {
			neighbors = append(neighbors, state[t-1][site+1])
		} else {
			neighbors = append(neighbors, state[t-1][0])
		}

		ham_r := Ham(new_value,n,neighbors)
		ham_s := Ham(state[t-1][site],n,neighbors)
		if ham_r < ham_s || r.Float64() < math.Exp(-beta * (ham_r - ham_s)) {
			state[t][site] = new_value
		}
	}

	return state
}

