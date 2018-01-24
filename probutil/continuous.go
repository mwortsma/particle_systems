package probutil

import (
	"fmt"
	"sync"
	"github.com/mwortsma/particle_systems/matutil"
)

type ContDistr struct {
	Dt float64
	T float64
	K int
	Distr [][]float64
}

func TypicalContDistrSync(
	f func()([]float64, matutil.Vec), 
	dt float64, 
	T float64,
	k, steps int) ContDistr {

	length := int(float64(T)/dt)
	cdistr := make([][]float64, length)
	for i := 0; i < length; i++ {
		cdistr[i] = make([]float64, k)
	}

	inc := 1.0/float64(steps)

	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(steps)
	for step := 0; step < steps; step++ {
		go func() {
			defer wg.Done()
			times, X := f()
			fmt.Println(X)
			fmt.Println(times)
			mutex.Lock()
			defer mutex.Unlock()

			curr_index := 0
			curr_time := 0.0
			for i := 0; i < length; i++ {
				for curr_index < len(times) - 1 && 
				times[curr_index+1] < curr_time {
					curr_index += 1
				}
				fmt.Println("curr index", curr_index, "with val", X[curr_index])
				cdistr[i][X[curr_index]] += inc
				curr_time += dt
			}
		}()
	}
	wg.Wait()
	return ContDistr{dt, T, k, cdistr}
}