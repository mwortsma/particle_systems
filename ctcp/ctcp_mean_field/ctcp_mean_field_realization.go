package ctcp_mean_field

import (
	"fmt"
	"github.com/mwortsma/particle_systems/probutil"
	"github.com/mwortsma/particle_systems/matutil"
	"golang.org/x/exp/rand"
	"time"
	"math"
)

// todo distance
func MeanFieldRealization(
	T float64,
	lam float64,
	nu float64,
	dt float64) ([]float64, matutil.Vec)  {
	
	X := make(matutil.Vec, 1)

	// Ger random number to be used throughout
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	// Initial conditions.
	if r.Float64() < nu {
		X[0] = 1
	}

	// keep track of times
	times := make([]float64, 1)
	t := 0.0

	for t < T {
		if X[len(X)-1] == 0 {

			for t <= T - dt {
				_, v := MeanFieldRealization(t, lam, nu, dt)
				t += dt
				if r.Float64() < 1.0-math.Exp(-lam*float64(v[len(v)-1])*dt) {
					break
				}
			}

		}  else {
			// Draw an exponential random variable with rate 1.
			t += r.ExpFloat64()
		}

		if t >= T {
			break
		}

		times = append(times, t)
		X = append(X, 1-X[len(X)-1])
	}
	fmt.Println(times)
	fmt.Println(X)

	return times, X
}

func RealizationTypicalDistr(T, lam float64, nu, dt float64, steps int) probutil.ContDistr {
	f := func() ([]float64, matutil.Vec) {
		return MeanFieldRealization(T,lam,nu,dt)
	}
	return probutil.TypicalContDistrSync(f, dt, T, 2, steps)
}
