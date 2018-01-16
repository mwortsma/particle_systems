package main

import (
	"github.com/mwortsma/particle_systems/dtlb/dtlb_local_ring"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_full_gengraph"
	"github.com/mwortsma/particle_systems/probutil"
	"github.com/mwortsma/particle_systems/plotutil"
	"fmt"
	"flag"
	"strings"
)

func main() {
	T := flag.Int("T", 4, "Steps")
	lam := flag.Float64("lambda", 0.8, "Rate")
	k := flag.Int("buffer", 5, "Buffer")
	eps := flag.Float64("epsilon", 0.001, "Epsilon")
	iters := flag.Int("iters", 4, "Number of Iterations")
	steps := flag.Int("steps", 100000, "Steps per Iteration")
	flag.Parse()
	_,_,local,_,_:= dtlb_local_ring.FixedPointIteration(*T, *lam, *k, *eps, *iters, *steps, probutil.L1Distance)
	full :=  dtlb_full_gengraph.RingTypicalDistr(*T, *lam, *k, *steps)
	distance := probutil.L1Distance(local, full)
	title := fmt.Sprintf("DTLB Full vs Local T=%d, lambda=%0.3f, buffer=%d, steps=%0.2e, L1 distance=%0.04f", *T, *lam, *k, float64(*steps), distance)
	file := strings.Replace(title, " ", "_", -1) + ".png"
	plotutil.PlotDistr(file, title, []probutil.Distr{local, full}, []string{"Local", "Full"})
	fmt.Println(distance)
}
