package dtcp_local

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
)


type CondExp []map[string]float64
type JointDistr []map[string]float64

func binomProb(x int, succ_prob float64) float64 {
	if x == 1 {
		return succ_prob
	} else {
		return 1.0-succ_prob
	}
	return 0.0
}

func transitionProb(prev int, curr int, sum float64, d int, p float64, q float64) float64 {
	if prev == 0 {
		return binomProb(curr, (p/float64(d))*sum)
	} else {
		return binomProb(curr, 1.0-q)
	}
	return 0.0
}

func getDistributions(T,d int, p,q float64, nu float64, tau int) (JointDistr, CondExp) {

	f := make(JointDistr, T)
	c := make(CondExp, T)

	fmt.Println("running.")

	// initialize maps
	for i := 0; i < T; i++ {
		f[i] = make(map[string]float64)
		c[i] = make(map[string]float64)
	}

	// l is the local neighborhood size
	l := d+1

	// initial conditions
	for _, s := range matutil.BinaryStrings(l) {
		prob := 1.0
		for _, elt := range s {
			prob *= binomProb(elt, nu)
		}
		f[0][matutil.Mat([][]int{s}).String()] = prob
	}

	for t := 1; t < T; t++ {
		// get c[t-1] using f[t-1]
		hsitorys := matutil.BinaryMats(t,2)
		paths := matutil.BinaryMats(t,d-1)
		for _, history := range hsitorys {

			hist_str := history.String()
			sum_prob := 0.0
			c[t-1][hist_str] = 0.0
			for _, path := range paths {
				full := matutil.Concat(history, path)
				prob := f[t-1][full.String()]
				sum_prob += prob
				sum := 0
				for j := 0; j < d-1; j++ {
					sum += path[t-1][j]
				}
				c[t-1][hist_str] += (float64(sum) * prob)
			}
			c[t-1][hist_str] /= sum_prob
		}

		// get f[t] using c[t-1]
		newstates := matutil.BinaryMats(t+1, d+1)
		for _, newstate := range newstates {
			str := newstate.String()
			hist := newstate.AllColsT(t)
			prev := newstate[t-1]
			curr := newstate[t]
			f[t][str] = f[t-1][hist.String()]

			sum := 0.0
			for j := 1; j < d+1; j++ {
				sum += float64(prev[j])
			}
			f[t][str] *= transitionProb(prev[0], curr[0], sum, d, p, q)
			for j := 1; j < d+1; j++ {
				sum = c[t-1][hist.Cols([]int{j,0}).String()]
				f[t][str] *= transitionProb(prev[j], curr[j],float64(prev[0]) + sum,d, p,q)
			}
		}
	}
	
	fmt.Println(f[0])
	fmt.Println("\n\n\n\n")
	fmt.Println(c[0])
	fmt.Println("\n\n\n\n")
	fmt.Println(f[1])
	
	return f, c
} 

func RecursionTypicalDistr(T,d int, p,q float64, nu float64, tau int) probutil.Distr {
	joint, _ := getDistributions(T,d, p,q, nu, tau)
	f := make(map[string]float64)
	states := matutil.BinaryMats(T, d+1)
	for _, state := range states {
		path := state.Col(0).String()
		if _, ok := f[path]; !ok {
			f[path] = 0.0
		}
		f[path] += joint[T-1][state.String()]
	}

	return f
}