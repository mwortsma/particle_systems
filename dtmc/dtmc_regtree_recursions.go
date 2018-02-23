package dtmc

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"math"
)

type Transition func(int, int, matutil.Vec) float64

func DTMCRegtreeEndDistr(T, tau int, d int, Q Transition, nu func(matutil.Vec) float64, k int) probutil.Distr {

	j := DTMCRegtreeRecursions(T, tau, d, Q, nu, k) 

	f := make(probutil.Distr)

	for k, prob := range(j) {
		mat := matutil.StringToMat(k)
		v := matutil.Vec(mat[len(mat)-1])
		str := v.String()
		if _, ok := f[str]; !ok {
			f[str] = prob
		} else {
			f[str] += prob
		}
	}

	return f
}

func DTMCRegtreeTDistr(T, tau int, d int, Q Transition, nu func(matutil.Vec) float64, k int) probutil.ContDistr {

	js, _ := DTMCRegtreeRecursionsFull(T, tau, d, Q, nu, k) 

	f := make([][]float64, T)
	t := 0
	for _, j := range js {
		f[t] = make([]float64, k)
		for k, prob := range j {
			mat := matutil.StringToMat(k)
			v := matutil.Vec(mat.Col(0))
			f[t][v[len(v)-1]] += prob
		}
		t += 1;
	}

	return probutil.ContDistr{1, float64(T), k, f}
}


func DTMCRegtreeRecursionsFull(T, tau int, d int, Q Transition, nu func(matutil.Vec) float64, k int) ([]probutil.Distr, []probutil.Distr) {

	if tau < 0 {
		tau = math.MaxInt32
	}

	j := make([]probutil.Distr, T)
	p := make([]probutil.Distr, T)
	j[0] = make(probutil.Distr)

	// j := make(probutil.Distr)
	var c probutil.Conditional

	for _, init := range matutil.QStrings(d+1,k) {
		j[0][matutil.Mat([][]int{init}).String()] = nu(init)
	}

	for t := 1; t < T; t++ {

		//fmt.Println(j[t-1])
		//fmt.Println("\n")

		j[t] = make(probutil.Distr)
		p[t] = make(probutil.Distr)

		c = getConditional(t-1, tau, d, j[t-1], k)

		j[t], p[t] = getJoint(t, tau, d, Q, j[t-1], c, k)

	}

	//fmt.Println(j[T-1])

	fmt.Println("Exiting")

	return j, p
}

func DTMCRegtreeRecursions(T, tau int, d int, Q Transition, nu func(matutil.Vec) float64, k int) probutil.Distr {

	if tau < 0 {
		tau = math.MaxInt32
	}

	j := make(probutil.Distr)

	var c probutil.Conditional

	for _, init := range matutil.QStrings(d+1,k) {
		j[matutil.Mat([][]int{init}).String()] = nu(init)
	}

	for t := 1; t < T; t++ {

		c = getConditional(t-1, tau, d, j, k)

		j, _ = getJoint(t, tau, d, Q, j, c, k)

	}

	return j
}

func getJoint(t, tau int, d int, Q Transition, j probutil.Distr, c probutil.Conditional, k int) (probutil.Distr, probutil.Distr) {
	fmt.Println("Obtaining joint at", t)

	jnew := make(probutil.Distr)

	p := make(probutil.Distr)

	l := Min(tau,t) + 1
	r := Min(tau,t-1) + 1

	prev_vals := matutil.QMats(r,d+1,k)
	new_vals := matutil.QStrings(d+1,k)
	other_children_vals := matutil.QStrings(d-1,k)

	for _, prev := range prev_vals {
		prob_prev := j[prev.String()]
		for _, new_val := range new_vals {
			full := append(prev, new_val)
			trimmed := full[len(full)-l:]
			trimmed_str := trimmed.String()
			lastrow := prev[len(prev)-1]
			prob := 1.0
			prob *= Q(new_val[0], lastrow[0], lastrow[1:])
			for i := 1; i < d+1; i++ {
				sum_prob := 0.0
				for _, other_children := range other_children_vals {
					hist := prev.Cols([]int{i,0}).String()
					sum_prob += c[hist][other_children.String()]*
						Q(new_val[i], lastrow[i], append(other_children, lastrow[0]))
				}
				prob *= sum_prob
			}
			p[full.String()] = prob
			prob *= prob_prev
			if _, ok := jnew[trimmed_str]; !ok {
				jnew[trimmed_str] = prob
			} else {
				jnew[trimmed_str] += prob	
			}
		}
	}

	return jnew, p
}

func getConditional(t, tau int, d int, jt probutil.Distr, k int) probutil.Conditional {
	fmt.Println("Obtaining Conditional at", t)

	ct := make(probutil.Conditional)

	l := Min(tau,t) + 1

	history_vals := matutil.QMats(l,2,k)
	children_vals := matutil.QMats(l,d-1,k)


	for _, history := range history_vals {
		hist_str := history.String()
		ct[hist_str] = make(probutil.Distr)
		denom := 0.0
		for _, children := range children_vals {
			full := matutil.Concat(history, children)
			denom += jt[full.String()]
		}
		// Important
		if denom == 0 {
			continue
		}
		for _, children := range children_vals {
			lastrow := matutil.Vec(children[l-1]).String()
			// TODO: debug
			if _, ok := ct[hist_str][lastrow]; !ok {
				ct[hist_str][lastrow] = 0
			}
			full := matutil.Concat(history, children)
			ct[hist_str][lastrow] += jt[full.String()]/denom
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

