package dtlb_util

import (        
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
)

func GetPQ(lam,dt float64) (float64, float64) {
	p := 1.0-math.Exp(-dt*lam)
	q := 1.0-math.Exp(-dt)
	return p,q
}

func Init(p,q float64, k int, r *rand.Rand) int {
	poiss := distuv.Poisson{-math.Log(1 - (p/q)), r}
	return int(math.Min(poiss.Rand(), float64(k-1)))
}