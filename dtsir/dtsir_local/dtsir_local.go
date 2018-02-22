package dtsir_local

import (
	"math"
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"github.com/mwortsma/particle_systems/dtmc"
)

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

func FullRun(T,tau int, d int, p,q float64, nu []float64) probutil.Distr {
	
	nu_f := func(v matutil.Vec) float64 {
		prob := 1.0
		for i := 0; i < len(v); i++ {
			prob *= nu[v[i]]
		}
		return prob
	}

	Q := func(k int, s int, v matutil.Vec) float64 {
		if s == 0 {
			sum_neighbors := 0
			for i := 0; i < len(v); i++ {
				if v[i] == 1 {
					sum_neighbors += v[i]
				}
			}
			if k == 1 {
				return (p/float64(d))*float64(sum_neighbors)
			} else if k == 0 {
				return 1-(p/float64(d))*float64(sum_neighbors)
			}
		} else if s == 1 {
			if k == 1 {
				return 1-q
			} else if k == 2 {
				return q
			}
		} else if s == 2 {
			if k == 2 {
				return 1.0
			}
		}
		return 0.0
	}

	j_array, p_array := dtmc.DTMCRegtreeRecursionsFull(T, tau, d, Q, nu_f,3)

	f := make(probutil.Distr)

	if tau < 0 {
		tau = math.MaxInt32
	}
	fmt.Println("T is ", T)

	states := matutil.QMats(T, d+1,3)
	for _, state := range states {
		path := state.Col(0).String()
		if _, ok := f[path]; !ok {
			f[path] = 0.0
		}
		f[path] += JointProb(T,tau,d,j_array,p_array,state)
	}

	fmt.Println(f)

	return f
}



func Run(T,tau int, d int, p,q float64, nu []float64) probutil.ContDistr {

	nu_f := func(v matutil.Vec) float64 {
		prob := 1.0
		for i := 0; i < len(v); i++ {
			prob *= nu[v[i]]
		}
		return prob
	}

	Q := func(k int, s int, v matutil.Vec) float64 {
		if s == 0 {
			sum_neighbors := 0
			for i := 0; i < len(v); i++ {
				if v[i] == 1 {
					sum_neighbors += v[i]
				}
			}
			if k == 1 {
				return (p/float64(d))*float64(sum_neighbors)
			} else if k == 0 {
				return 1-(p/float64(d))*float64(sum_neighbors)
			}
		} else if s == 1 {
			if k == 1 {
				return 1-q
			} else if k == 2 {
				return q
			}
		} else if s == 2 {
			if k == 2 {
				return 1.0
			}
		}
		return 0.0
	}

	return dtmc.DTMCRegtreeTDistr(T, tau, d, Q, nu_f,3)
}
