package ctmc

import (
	"testing"
)

func TestCTMC(t *testing.T) {
	rates := []float64{0.5, 0.0001, 1, 0.5}
	pops := make([]int, 4)
	s := 0.0
	iters := 10000
	for i := 0; i < iters; i++ {
		e, t := StepCTMC(rates)
		s += t
		pops[e] += 1
	}
	t.Log(pops)
	t.Log(float64(s)/float64(iters))
}
