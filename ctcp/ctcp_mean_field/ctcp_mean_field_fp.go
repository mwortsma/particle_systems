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


	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

	f := probutil.ContDistr{T: T, Dt: dt}
	for iter := 0; iter < iters; iter++ {

		evolve_function := func()([]float64, matutil.Vec) {
			return evolveSystem(T, lam, nu, f, dt, r)
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
	f probutil.ContDistr,
	dt float64,
	r *rand.Rand) ([]float64, matutil.Vec)  {
	
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
				
				v := 1-nu
				if len(f.Distr) != 0 {
					v = f.Distr[int(t/f.Dt)][1]
				}	

				if r.Float64() > 1.0-math.Exp(-lam*v*dt) {
					continue
				}
				X = append(X,1)
				break
			}

		}  else {
			// Draw an exponential random variable with rate 1.
			t += r.ExpFloat64()
			if t >= T {
				return times, X
			}
			X = append(X,0)
		}
		times = append(times, t)
	}
	return times, X
}