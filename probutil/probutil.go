package probutil

import (
	"fmt"
	"sync"
)

type Distr map[string]float64

func TypicalDistrSync(f func() fmt.Stringer, steps int) Distr {
	distr := make(map[string]float64)
	inc := 1.0/float64(steps)
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(steps)
	for step := 0; step < steps; step++ {
		go func() {
			defer wg.Done()
			v := f().String()
			
			mutex.Lock()
			defer mutex.Unlock()
			if _, ok := distr[v]; !ok {
				distr[v] = 0.
			}
			distr[v] += inc
		}()
	}
	wg.Wait()
	fmt.Println(len(distr))
	return distr
}

func TypicalDistr(f func() fmt.Stringer, steps int) Distr {
	distr := make(map[string]float64)
	inc := 1.0/float64(steps)
	for step := 0; step < steps; step++ {
		v := f().String()
		if _, ok := distr[v]; !ok {
			distr[v] = 0.
		}
		distr[v] += inc
	}
	return distr
}