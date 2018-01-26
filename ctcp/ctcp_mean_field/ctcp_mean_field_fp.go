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
func MeanFieldFixedPointIteration(
	T float64,
	lam float64,
	nu float64,
	dt float64,
	eps float64,
	iters int,
	steps int,
	dist probutil.ContDistance) probutil.ContDistr {

	f := probutil.ContDistr{T: T, Dt: dt}
	for iter := 0; iter < iters; iter++ {

		evolve_function := func()([]float64, matutil.Vec) {
			return evolveSystem(T, lam, nu, f)
		}

		new_f := probutil.TypicalContDistrSync(evolve_function,dt,T,2,steps)

		distance := 1.0
		if iter >= 1 {
			distance = dist(f,new_f)
		}

		f = new_f

		fmt.Println(fmt.Sprintf("Iteration %d, Distance %0.5f", iter, distance))

		if distance < eps {
			break
		}

	}

	return f
}



func evolveSystem(
	T float64,
	lam float64, 
	nu float64,
	f probutil.ContDistr) ([]float64, matutil.Vec){

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

	for {
		if X[len(X)-1] == 0 {

			if len(f.Distr) == 0 {
				// TODO this is in the case of no distributin f
				infected := 1.0-nu
				t += r.ExpFloat64() / (lam * infected)
			} else {
				for {
					// prob of an arrival in this process is 1-e^{-lambda*infected*dt}
					if t >= T - f.Dt {
						return times, X
					} else if r.Float64() < 1.0-math.Exp(-lam*f.Distr[int(t/f.Dt)][1]*f.Dt) {
						t += f.Dt
						break
					}
					t += f.Dt
				}
			}
			
		} else {
			// Draw an exponential random variable with rate 1.
			t += r.ExpFloat64()
		}
		if t >= T {
			break
		}
		times = append(times, t)
		X = append(X, 1-X[len(X)-1])
	}

	return times, X
}
