package ctcp_mean_field

import (
	//"fmt"
	"github.com/mwortsma/particle_systems/matutil"
	"github.com/mwortsma/particle_systems/probutil"
	"golang.org/x/exp/rand"
	"math"
	"time"
)

// todo distance
func MeanFieldRealization(
	T float64,
	lam float64,
	nu float64,
	dt float64,
	r *rand.Rand) ([]float64, matutil.Vec) {

	X := []int{0}

	// Initial conditions.
	if r.Float64() < nu {
		X = []int{1}
	}

	// keep track of times
	times := make([]float64, 1)
	t := 0.0

	for t < T {

		if X[len(X)-1] == 0 {

			for {
				t += dt
				if t >= T {
					return times, X
				}
				_, v := MeanFieldRealization(t-dt, lam, nu, dt, r)

				if r.Float64() > 1.0-math.Exp(-lam*float64(v[len(v)-1])*dt) {
					continue
				}
				X = append(X, 1)
				break
			}

		} else {
			// Draw an exponential random variable with rate 1.
			t += r.ExpFloat64()
			if t >= T {
				return times, X
			}
			X = append(X, 0)
		}
		times = append(times, t)
	}
	return times, X
}

func RealizationTypicalDistr(T, lam float64, nu, dt float64, steps int) probutil.ContDistr {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	f := func() ([]float64, matutil.Vec) {
		return MeanFieldRealization(T, lam, nu, dt, r)
	}
	return probutil.TypicalContDistrSync(f, dt, T, 2, steps)
}
