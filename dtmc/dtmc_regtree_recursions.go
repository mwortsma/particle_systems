package dtmc

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"math"
)

type Transition func(int, int, matutil.Vec) float64

func DTMCRegtreeRecursionsFull(T, tau int, d int, Q Transition, nu func(matutil.Vec) float64) []probutil.Distr {

	if tau < 0 {
		tau = math.MaxInt32
	}

	j := make([]probutil.Distr, T)
	j[0] = make(probutil.Distr)

	// j := make(probutil.Distr)
	var c probutil.Conditional

	for _, init := range matutil.BinaryStrings(d+1) {
		j[0][matutil.Mat([][]int{init}).String()] = nu(init)
	}

	fmt.Println(j)

	for t := 1; t < T; t++ {

		j[t] = make(probutil.Distr)

		c = getConditional(t-1, tau, d, j[t-1])

		j[t] = getJoint(t, tau, d, Q, j[t-1], c)

	}

	fmt.Println("Exiting")

	return j
}

func DTMCRegtreeRecursions(T, tau int, d int, Q Transition, nu func(matutil.Vec) float64) probutil.Distr {

	if tau < 0 {
		tau = math.MaxInt32
	}

	j := make(probutil.Distr)

	var c probutil.Conditional

	for _, init := range matutil.BinaryStrings(d+1) {
		j[matutil.Mat([][]int{init}).String()] = nu(init)
	}

	for t := 1; t < T; t++ {

		c = getConditional(t-1, tau, d, j)

		j = getJoint(t, tau, d, Q, j, c)

	}

	return j
}

func getJoint(t, tau int, d int, Q Transition, j probutil.Distr, c probutil.Conditional) probutil.Distr {
	fmt.Println("Obtaining joint at", t)

	jnew := make(probutil.Distr)

	conditional := make(map[string]float64)

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
				conditional[full.String()] = 0
			} else {
				lastrow := prev[len(prev)-1]
				prob := 1.0
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
				conditional[full.String()] = prob
				prob *= prob_prev
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
	fmt.Println("Obtaining Conditional at", t)

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
			lastrow := matutil.Vec(children[l-1]).String()
			// TODO: debug
			if _, ok := ct[hist_str][lastrow]; !ok {
				ct[hist_str][lastrow] = 0
			}
			if denom > 0 {
				full := matutil.Concat(history, children)
				ct[hist_str][lastrow] += jt[full.String()]/denom
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

