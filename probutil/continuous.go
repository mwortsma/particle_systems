package probutil

import (
	// "fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"math"
	"sync"
)

type ContDistr struct {
	Dt    float64
	T     float64
	K     int
	Distr [][]float64
}
type ContDistance func(ContDistr, ContDistr) float64

func TypicalContDistrSync(
	f func() ([]float64, matutil.Vec),
	dt float64,
	T float64,
	k int,
	steps int) ContDistr {

	length := int(float64(T) / dt)
	cdistr := make([][]float64, length)
	for i := 0; i < length; i++ {
		cdistr[i] = make([]float64, k)
	}

	inc := 1.0 / float64(steps)

	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(steps)
	for step := 0; step < steps; step++ {
		go func() {
			defer wg.Done()
			times, X := f()
			mutex.Lock()
			defer mutex.Unlock()

			curr_index := 0
			curr_time := 0.0
			for i := 0; i < length; i++ {
				for curr_index < len(times)-1 &&
					times[curr_index+1] <= curr_time {
					curr_index += 1
				}
				cdistr[i][X[curr_index]] += inc
				curr_time += dt
			}

			// fmt.Println("got", times, X, "out", cdistr)
		}()
	}
	wg.Wait()
	return ContDistr{dt, T, k, cdistr}
}

func TypicalContDistr(
	f func() ([]float64, matutil.Vec),
	dt float64,
	T float64,
	k int,
	steps int) ContDistr {

	length := int(float64(T) / dt)
	cdistr := make([][]float64, length)
	for i := 0; i < length; i++ {
		cdistr[i] = make([]float64, k)
	}

	inc := 1.0 / float64(steps)

	for step := 0; step < steps; step++ {

		times, X := f()

		curr_index := 0
		curr_time := 0.0
		for i := 0; i < length; i++ {
			for curr_index < len(times)-1 &&
				times[curr_index+1] <= curr_time {
				curr_index += 1
			}
			cdistr[i][X[curr_index]] += inc
			curr_time += dt
		}

	}

	return ContDistr{dt, T, k, cdistr}
}

func ContL1Distance(d1, d2 ContDistr) (dist float64) {
	for i := 0; i < int(float64(d1.T)/d1.Dt); i++ {
		for j := 0; j < d1.K-1; j++ {
			dist += math.Abs(d1.Distr[i][j] - d2.Distr[i][j])
		}
	}
	return
}
