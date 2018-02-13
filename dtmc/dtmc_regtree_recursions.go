package dtmc

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"math"
)

type Transition func(int, int, matutil.Vec) float64

func DTMCRegtreeRecursions(T, tau int, d int, Q Transition, nu func(matutil.Vec) float64) (probutil.Distr, probutil.Conditional) {

	if tau < 1 {
		tau = math.MaxInt32
	}

	j := make(probutil.Distr)
	var c probutil.Conditional

	for _, init := range matutil.BinaryStrings(d+1) {
		j[init.String()] = nu(init)
	}

	for t := 1; t < T; t++ {
		c = getConditional(t, tau, d, j)
		j = getJoint(t, tau, d, Q, j, c)
	}


	return j,c
}

func getJoint(t, tau int, d int, Q Transition, j probutil.Distr, c probutil.Conditional) probutil.Distr {
	fmt.Println("Obtaining conditional at", t)

	jnew := make(probutil.Distr)

	l := Min(tau,t) + 1
	r := Min(tau,t-1) + 1

	prev_vals := matutil.BinaryMats(r,d+1)
	new_vals := matutil.BinaryStrings(d+1)
	other_children_vals := matutil.BinaryStrings(d-1)

	for _, prev := range prev_vals {
		prob_prev := j[prev.String()]
		for _, new_val := range new_vals {
			full := append(prev, new_val)
			trimmed := full[len(full)-l:]
			trimmed_str := trimmed.String()
			if prob_prev == 0 {
				jnew[trimmed.String()] = 0
			} else {
				lastrow := prev[len(prev)-1]
				prob := prob_prev
				prob *= Q(new_val[0], lastrow[0], lastrow[1:])
				for i := 1; i < d+1; i++ {
					sum_prob := 0.0
					for _, other_children := range other_children_vals {
						hist := prev.Cols([]int{i,0})
						sum_prob += c[hist.String()][other_children.String()]*
							Q(new_val[i], lastrow[i], append(other_children, lastrow[0]))
					}
					prob *= sum_prob
				}
				if _, ok := jnew[trimmed_str]; !ok {
					jnew[trimmed_str] = prob
				} else {
					jnew[trimmed_str] += prob
				}
			}
		}
	}

	return jnew
}

func getConditional(t, tau int, d int, jt probutil.Distr) probutil.Conditional {

	ct := make(probutil.Conditional)

	l := Min(tau,t) + 1

	history_vals := matutil.BinaryMats(l,2)
	children_vals := matutil.BinaryMats(l,d-1)


	for _, history := range history_vals {
		hist_str := history.String()
		ct[hist_str] = make(probutil.Distr)
		denom := 0.0
		for _, children := range children_vals {
			full := matutil.Concat(history, children)
			denom += jt[full.String()]
		}
		for _, children := range children_vals {
			lastrow := matutil.Vec(children[l-1])
			if denom == 0 {
				ct[hist_str][lastrow.String()] = 0
			} else {
				full := matutil.Concat(history, children)
				ct[hist_str][lastrow.String()] = jt[full.String()]/denom
			}
		}
	}
	
	return ct
}


// Minimum
func Min(a,b int) int {
	if a < b {
		return a
	} else {
		return b
	}
	return -1
}

