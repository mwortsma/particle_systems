package main

import (
	"github.com/mwortsma/particle_systems/dtlb/dtlb_local_ring"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_mean_field"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_full_gengraph"
	"github.com/mwortsma/particle_systems/probutil"
	"github.com/mwortsma/particle_systems/plotutil"
	"fmt"
	"flag"
	"strings"
)

func main() {

	// type of compare. Options are
	// 1. dtlb-local-vs-full-ring (discrete time load balancing). Status: Complete.
	// 2. dtlb-complete-vs-mean-field. Status: TODO. Not working as expected.
	var type_string string
    flag.StringVar(&type_string, "type", "dtlb-local-vs-full-ring", "Type of compare")

	// type of distances. Options are
	// 1. L1 (L1 distance)
	var dist_string string
    flag.StringVar(&dist_string, "distance", "L1", "Type of distance")

	// General arguments
	eps := flag.Float64("epsilon", 0.001, "Epsilon")
	iters := flag.Int("iters", 4, "Number of Iterations")
	steps := flag.Int("steps", 100000, "Steps per Iteration")

	// dtlb-local-vs-full-ring arguments
	T := flag.Int("T", 4, "Steps")
	lam := flag.Float64("lambda", 0.8, "Rate")
	k := flag.Int("buffer", 5, "Buffer")

	// dtlb-complete-vs-mean-field
	n := flag.Int("n", 5, "number of particles")

	flag.Parse()

	var dist probutil.Distance
	switch dist_string {
	case "L1":
		dist = probutil.L1Distance
	default:
		fmt.Println("Distance not recognized.")
		return
	}


	var distrs []probutil.Distr
	var labels []string
	var title string

	switch type_string {
	case "dtlb-local-vs-full-ring":
		fmt.Println("dtlb-local-vs-full-ring")
		distrs = []probutil.Distr{
			func() probutil.Distr { 
				_, local,_,_ := dtlb_local_ring.FixedPointIteration(
					*T, *lam, *k, *eps, *iters, *steps, dist) 
				return local
				}(),
			dtlb_full_gengraph.RingTypicalDistr(*T, *lam, *k, *steps),
		}
		labels = []string{
			"Full",
			"Local",
		}
		title = fmt.Sprintf("DTLB Full vs Local T=%d, lambda=%0.3f, buffer=%d, steps=%0.2e",
		 *T, *lam, *k, float64(*steps))




	case "dtlb-complete-vs-mean-field":
		fmt.Println("dtlb-complete-vs-mean-field")
		deg := *n
		deg--
		distrs = []probutil.Distr{
			dtlb_mean_field.TypicalDistr(*T, *lam, *k, *n, *steps),
			dtlb_full_gengraph.CompleteTypicalDistr(*T, *lam, *k, *steps, deg),
		}
		labels = []string{
			"Mean Field",
			"Complete",
		}
		title = fmt.Sprintf("DTLB Mean Field vs Complete T=%d, lambda=%0.3f, buffer=%d, steps=%0.2e, n=%d",
		 *T, *lam, *k, float64(*steps), *n)



	default:
		fmt.Println("Type not recognized.")
		return
	}

	for i := 0; i < len(distrs); i++ {
		for j := i+1; j < len(distrs); j++ {
			d := dist(distrs[i], distrs[j])
			fmt.Println(fmt.Sprintf("Distance %s vs %s: %f", labels[i], labels[j], d))
		}
	}
	file := "plots/" + strings.Replace(title, " ", "_", -1) + ".png"
	plotutil.PlotDistr(file, title, distrs, labels)
}
