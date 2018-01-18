package dtcp_local

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"sync"
	"time"
)

type CondDistr []map[string]map[int]float64

func RegTreeFixedPointIteration(
	T int,
	d int,
	p float64,
	q float64,
	nu float64,
	eps float64,
	iters int,
	steps int,
	dist probutil.Distance) (probutil.Distr, probutil.Distr, []float64, []float64) {

	fmt.Println("Running dtcp local fixed point simulation d=", d)

	joint, cond, typical := make(probutil.Distr), initCond(T), make(probutil.Distr)
	joint_dists := make([]float64, 0)
	typical_dists := make([]float64, 0)

	for iter := 0; iter < iters; iter++ {
		new_joint, new_cond, new_typical, cond_misses := ringStep(T, d, p, q, nu, steps, cond)
		joint_dist := dist(joint, new_joint)
		joint_dists = append(joint_dists, joint_dist)
		typical_dist := dist(typical, new_typical)
		typical_dists = append(typical_dists, typical_dist)
		fmt.Println(fmt.Sprintf("Iteration %d, Joint Distance: %0.4f, Typical Distance: %0.4f, Misses: %d", iter, joint_dist, typical_dist, cond_misses))
		joint, cond, typical = new_joint, new_cond, new_typical
		if joint_dist < eps {
			break
		}
	}

	return joint, typical, joint_dists, typical_dists
}

func ringStep(
	T int,
	d int,
	p float64,
	q float64,
	nu float64,
	steps int,
	old_cond CondDistr) (probutil.Distr, CondDistr, probutil.Distr, int) {

	joint, cond, typical := make(probutil.Distr), initCond(T), make(probutil.Distr)

	obvserved := make([]map[string]int, T-1)
	for t := 0; t < T-1; t++ {
		obvserved[t] = make(map[string]int)
	}

	inc := 1 / float64(steps)
	cond_misses := 0

	var cond_mutex = &sync.Mutex{}
	var rest_mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(steps)

	for step := 0; step < steps; step++ {
		go func() {
			defer wg.Done()
			X, cond_missed := treeRealization(T, d,p,q,nu, old_cond)
			// update joint, typical
			rest_mutex.Lock()
			cond_misses += cond_missed
			probutil.Update(joint, X.String(), inc)
			probutil.Update(typical, X.Col(0).String(), inc)
			rest_mutex.Unlock()
			// now update cond
			cond_mutex.Lock()
			for i := 1; i < d+1; i++ {
				key := X.Cols([]int{0,i})
				for t := 1; t < T; t++ {
					keyt := key[:t].String()
					sum_neighbors := 0
					for j := 1; j < d+1; j++ {
						if j != i {
							sum_neighbors += X[t-1][j]
						}
					}
					if _, ok := cond[t-1][keyt]; !ok {
						cond[t-1][keyt] = make(map[int]float64)
						obvserved[t-1][keyt] = 0
						for k := 0; k < d+1; k++ {
							cond[t-1][keyt][k] = 0.0
						}
					}
					obvserved[t-1][keyt]++
					cond[t-1][keyt][sum_neighbors]++
				}
			}
			cond_mutex.Unlock()
		}()
	}
	wg.Wait()

	// Scale conditional
	for t := 0; t < T-1; t++ {
		for key := range cond[t] {
			for k, v := range cond[t][key] {
				cond[t][key][k] = v / float64(obvserved[t][key])
			}
		}
	}

	return joint, cond, typical, cond_misses
}

func treeRealization(T, d int, p, q float64, nu float64, cond CondDistr) (matutil.Mat, int) {
	// n is how many nodes we need to keep track of.
	cond_misses := 0
	// X[0] stores the root. X[1:d+1] store the children.
	X := matutil.Create(T, d+1)
	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Initial conditions.
	for i := 0; i < d+1; i++ {
		if r.Float64() < nu {
			X[0][i] = 1
		}
	}

	for t := 1; t < T; t++ {
		// Copy the state of X[t-1] to X[t]
		copy(X[t], X[t-1])
		
		// update the root.
		if X[t-1][0] == 0 {
			sum_neighbors := 0
			for j := 1; j < d+1; j++ {
				sum_neighbors += X[t-1][j]
			}
			if r.Float64() < (p/float64(d))*float64(sum_neighbors) {
				X[t][0] = 1
			}

		} else {
			// transition to 0 with probability q
			if r.Float64() < q {
				X[t][0] = 0
			}
		}

		// update the rest of the neighbors.
		for i := 1; i < d+1; i++ {
			if X[t-1][i] == 0 {
				
				key := X.ColsT([]int{i,0},t).String()
				sum_neighbors := X[t-1][0]

				if distr, ok := cond[t-1][key]; ok {
					sum_neighbors += probutil.SampleInt(distr, r.Float64())
				} else {
					// TODO
					cond_misses++
					sum_neighbors += r.Intn(d)
				}

				if r.Float64() < (p/float64(d))*float64(sum_neighbors) {
					X[t][i] = 1
				}

			} else {
				// transition to 0 with probability q
				if r.Float64() < q {
					X[t][i] = 0
				}
			}
		}
	}

	return X, cond_misses
}

func initCond(T int) CondDistr {
	c := make(CondDistr, T-1)
	for i := 0; i < T-1; i++ {
		c[i] = make(map[string]map[int]float64)
	}
	return c
}
