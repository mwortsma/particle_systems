package probutil

import (
	"fmt"
	"sync"
)

type Distr map[string]float64

func TypicalDistr(f func() fmt.Stringer, steps int) Distr {
	distr := make(map[string]float64)
	inc := 1.0/float64(steps)
	var mutex = &sync.Mutex{}

	for step := 0; step < steps; step++ {
		go func() {
			v := f()
			
			mutex.Lock()
			defer mutex.Unlock()
			if _, ok := distr[v.String()]; ok {
				distr[v.String()] = 0.
			}
			distr[v.String()] += inc
		}()
	}
	return distr
}