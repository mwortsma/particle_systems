package probutil

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

type Distr map[string]float64

type Distance func(Distr, Distr) float64

func Sample(d Distr, r float64) string {
	s := 0.
	for k, v := range d {
		if s += v; s > r {
			return k
		}
	}
	return ""
}

func Update(d Distr, s string, f float64) {
	if _, ok := d[s]; !ok {
		d[s] = 0.
	}
	d[s] += f
}

func TypicalDistrSync(f func() fmt.Stringer, steps int) Distr {
	distr := make(map[string]float64)
	inc := 1.0 / float64(steps)
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
	return distr
}

func TypicalDistr(f func() fmt.Stringer, steps int) Distr {
	distr := make(map[string]float64)
	inc := 1.0 / float64(steps)
	for step := 0; step < steps; step++ {
		v := f().String()
		if _, ok := distr[v]; !ok {
			distr[v] = 0.
		}
		distr[v] += inc
	}
	return distr
}

func L1Distance(d1, d2 Distr) (dist float64) {
	for k, v1 := range d1 {
		v2, ok := d2[k]
		if !ok {
			dist += v1
		} else {
			dist += math.Abs(v1 - v2)
		}
	}
	for k, v := range d2 {
		if _, ok := d1[k]; !ok {
			dist += v
		}
	}
	return
}

func SharedSortedKeys(distrs []Distr) []string {
	keys_map := make(map[string]bool)
	for _, distr := range distrs {
		for k := range distr {
			keys_map[k] = true
		}
	}
	keys := make([]string, 0)
	for k := range keys_map {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
