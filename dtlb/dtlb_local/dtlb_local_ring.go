package dtlb_local

import (
	"fmt"
	"github.com/mwortsma/particle_systems/graphutil"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"time"
	"sync"
)

type CondDistr []map[string]probutil.Distr

// FixedPointIterations runs the fixed point iteration algorithm.
// Returns the full distribution followed by the conditional distribution,
// the typical particle distribution and the joint and typical
// distances at each step.
func RingFixedPointIteration(
	T int,
	lam float64,
	k int,
	eps float64,
	iters int,
	steps int,
	dist probutil.Distance) (probutil.Distr, probutil.Distr, []float64, []float64) {

	joint, cond, typical := make(probutil.Distr), initCond(T), make(probutil.Distr)
	joint_dists := make([]float64, 0)
	typical_dists := make([]float64, 0)

	for iter := 0; iter < iters; iter++ {
		new_joint, new_cond, new_typical, cond_misses := ringStep(T, lam, k, steps, cond)
		joint_dist := dist(joint, new_joint)
		joint_dists = append(joint_dists, joint_dist)
		typical_dist := dist(typical, new_typical)
		typical_dists = append(typical_dists, typical_dist)
		fmt.Println(fmt.Sprintf("Iteration %d, Joint Distance: %0.4f, Typical Distance: %0.4f, Misses: %d",iter, joint_dist, typical_dist, cond_misses))
		joint, cond, typical = new_joint, new_cond, new_typical
		if joint_dist < eps {
			break
		}
	}

	return joint, typical, joint_dists, typical_dists
}

func ringStep(
	T int,
	lam float64,
	k int,
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
			X, cond_missed := ringRealization(T, lam, k, old_cond)
			// update joint, typical
			rest_mutex.Lock()
			cond_misses += cond_missed
			probutil.Update(joint, X.String(), inc)
			probutil.Update(typical, X.Col(5).String(), inc)
			rest_mutex.Unlock()
			// now update cond
			cond_mutex.Lock()
			left := X.Cols([]int{0, 1, 2, 3})
			right := X.Cols([]int{6, 5, 4, 3})
			for t := 1; t < T; t++ {
				left_key, right_key := left[:t].String(), right[:t].String()
				if _, ok := cond[t-1][left_key]; !ok {
					cond[t-1][left_key] = make(probutil.Distr)
					obvserved[t-1][left_key] = 0
				}
				if _, ok := cond[t-1][right_key]; !ok {
					cond[t-1][right_key] = make(probutil.Distr)
					obvserved[t-1][right_key] = 0
				}
				probutil.Update(cond[t-1][left_key], X.Colst([]int{4, 5}, t-1).String(), 1.0)
				probutil.Update(cond[t-1][right_key], X.Colst([]int{2, 1}, t-1).String(), 1.0)
				obvserved[t-1][left_key]++
				obvserved[t-1][right_key]++
			}
			cond_mutex.Unlock()
		}()
	}
	wg.Wait()

	// Scale conditional
	for t := 0; t < T-1; t++ {
		for key := range(cond[t]) {
			for k, v := range(cond[t][key]) {
				cond[t][key][k] = v/float64(obvserved[t][key])
			}
		}
	}

	return joint, cond, typical, cond_misses
}

func ringRealization(T int, lam float64, k int, cond CondDistr) (matutil.Mat, int) {
	// n is how many nodes we need to keep track of.
	cond_misses := 0
	n := 11
	G := graphutil.Ring(n)
	X := matutil.Create(T, n)
	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Initial conditions.
	p := distuv.Poisson{-math.Log(1 - lam), r}
	for i := 0; i < n; i++ {
		X[0][i] = int(math.Min(p.Rand(), float64(k-1)))
	}

	for t := 1; t < T; t++ {
		// Copy the state of X[t-1] to X[t]
		copy(X[t], X[t-1])
		// obtain a vector of arrivals
		arrivals := make([]bool, n)
		for i := 1; i < n-1; i++ {
			arrivals[i] = r.Float64() < lam
		}
		// if there is an arrival at 1,2 we need to sample the conditional
		key := X.ColsT([]int{5, 4, 3, 2}, t).String()
		if d, ok := cond[t-1][key]; ok {
			sample := matutil.StringToVec(probutil.Sample(d, r.Float64()))
			X[t-1][0], X[t-1][1] = sample[1], sample[0]
		} else {
			// TODO
			cond_misses++
			X[t-1][0], X[t-1][1] = int(math.Min(p.Rand(), float64(k-1))), int(math.Min(p.Rand(), float64(k-1)))
		}

		// if there is an arrival at 8,9 we need to sample the conditional
		key = X.ColsT([]int{5, 6, 7, 8}, t).String()
		if d, ok := cond[t-1][key]; ok {
			sample := matutil.StringToVec(probutil.Sample(d, r.Float64()))
			X[t-1][9], X[t-1][10] = sample[0], sample[1]
		} else {
			// TODO
			cond_misses++
			X[t-1][9], X[t-1][10] = int(math.Min(p.Rand(), float64(k-1))), int(math.Min(p.Rand(), float64(k-1)))
		}

		for i := 1; i < n-1; i++ {
			// Serve an item if non-empty
			if X[t-1][i] > 0 {
				X[t][i]--
			}
			// With probability lam there is an arrival.
			// Send to minimum neighbor.
			if arrivals[i] {
				// First get the min value
				min := X[t-1][i]
				for j := 0; j < len(G[i]); j++ {
					if X[t-1][G[i][j]] < min {
						min = X[t-1][G[i][j]]
					}
				}
				// Select, at random, a neighbor having that value.
				min_neighbors := make([]int, 0)
				if min == X[t-1][i] {
					min_neighbors = append(min_neighbors, i)
				}
				for j := 0; j < len(G[i]); j++ {
					if X[t-1][G[i][j]] == min {
						min_neighbors = append(min_neighbors, G[i][j])
					}
				}
				chosen_neighbor := min_neighbors[r.Intn(len(min_neighbors))]
				if X[t][chosen_neighbor] < k-1 {
					// Only send if below buffer.
					X[t][chosen_neighbor]++
				}
			}
		}
	}

	return X.Cols([]int{2, 3, 4, 5, 6, 7, 8}), cond_misses
}

func initCond(T int) CondDistr {
	c := make(CondDistr, T-1)
	for i := 0; i < T-1; i++ {
		c[i] = make(map[string]probutil.Distr)
	}
	return c
}
