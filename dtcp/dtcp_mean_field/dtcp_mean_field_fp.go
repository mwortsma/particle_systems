package dtcp_mean_field

import (
	"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"sync"
	"time"
)

type CondDistr []map[string]map[int]float64

func MeanFieldFixedPointIteration(
	T int,
	p float64,
	q float64,
	nu float64,
	eps float64,
	iters int,
	steps int,
	dist probutil.Distance) probutil.Distr {

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	distr, f := make(probutil.Distr), make([]float64, T)

	for iter := 0; iter < iters; iter++ {

		new_distr, new_f := meanFieldStep(T, p, q, nu, steps, r, f)
		distance := dist(new_distr, distr)
		fmt.Println(fmt.Sprintf("Iteration %d, Distance: %0.4f", iter, distance))
		distr, f = new_distr, new_f
		if distance < eps {
			break
		}
	}

	return distr
}

func meanFieldStep(
	T int,
	p float64,
	q float64,
	nu float64,
	steps int,
	r *rand.Rand,
	old_f []float64) (probutil.Distr, []float64) {

	distr, f := make(probutil.Distr), make([]float64, T)

	inc := 1.0 / float64(steps)
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(steps)
	for step := 0; step < steps; step++ {
		go func() {
			defer wg.Done()
			x := fixedPointRealization(T, p, q, nu, r, old_f)
			v := x.String()

			mutex.Lock()
			defer mutex.Unlock()
			if _, ok := distr[v]; !ok {
				distr[v] = 0.
			}
			distr[v] += inc

			for t := 0; t < T; t++ {
				f[t] += inc * float64(x[t])
			}
		}()
	}
	wg.Wait()
	return distr, f
}

func fixedPointRealization(T int, p, q float64, nu float64, r *rand.Rand, f []float64) matutil.Vec {

	X := make(matutil.Vec, T)

	// Initial conditions.
	if r.Float64() < nu {
		X[0] = 1
	}

	for t := 1; t < T; t++ {
		X[t] = X[t-1]
		if X[t-1] == 0 {
			// transition with probability ratio_neighbors * p
			if r.Float64() < f[t-1]*p {
				X[t] = 1
			}
		} else {
			// if state is 1, transition back with porbability q
			if r.Float64() < q {
				X[t] = 0
			}
		}
	}

	return X
}
