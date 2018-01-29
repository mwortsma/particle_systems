package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_full"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_local"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_mean_field"
	"github.com/mwortsma/particle_systems/probutil"
	"io/ioutil"
)

func main() {

	// Discrete Time Load Balancing.

	// Generel Arguments
	// -d=x (degree of a node)
	// -n=x (number of nodes)
	// -T=x (time horizon. T > 0)
	// -lam=x (P(arrival at a node))
	// -k=x (finite buffer. no queue can have capacity k)
	// -steps=x (how many samples used in generating the empirical distribtuion)
	// -file=x (which file to save the distribution to. empty = do not save)

	// Arguments specific to fixed point iteration algorithm.
	// -distance=x (options are: L1 (used as defualt), more to come)
	// -iters=x (for the fixed point algorithm, how many iterations to run)
	// -eps=x (threshold distance between typical particle distributions)

	// Types:

	// 1. Full Simulations
	// // 1.1 Ring (full_ring)
	// // // Note: Depth optional. Defaults to 1+4*T.
	// // 1.2 Complete (dtlb -full_complete -n=?)

	// 2. Local Simulations (Using the fixed point algorithm)
	// // 1.1 Ring (dtlb -local_ring)

	// 3. Mean Field Simulations (TODO: In progress...)
	// // 3.1 Mean Field (mean_field) (dtlb -mean_field -d=12)

	// Defining the arguments.
	full_ring := flag.Bool("full_ring", false, "Full sim on the ring")
	full_complete := flag.Bool("full_complete", false, "Full sim on the complete graph")

	local_ring_fp := flag.Bool("local_ring_fp", false, "Local sim fixed point on the ring")
	local_ring_realization := flag.Bool("local_ring_realization", false, "Local sim realization on the ring")

	mean_field := flag.Bool("mean_field", false, "Mean Field simulation.")

	d := flag.Int("d", -1, "degree of a noe")
	n := flag.Int("n", -1, "number of nodes")
	T := flag.Int("T", 2, "time horizon. T>0")
	k := flag.Int("k", 5, "Finite buffer. No queue can have capacity k.")
	lam := flag.Float64("lam", 0.8, "incoming rate at each node")
	steps := flag.Int("steps", 100, "how many samples used in generating the empirical distribtuion")
	var file_str string
	flag.StringVar(&file_str, "file", "", "where to save the distribution.")

	eps := flag.Float64("epsilon", 0.001, "threshold distance between typical particle distributions")
	iters := flag.Int("iters", 4, "for the fixed point algorithm, how many iterations to run")
	var distance_str string
	flag.StringVar(&distance_str, "distance", "L1", "type of distance between distributions. Options: L1, more to come...")

	flag.Parse()

	// Obtaining the distance.
	var dist probutil.Distance
	switch distance_str {
	case "L1":
		dist = probutil.L1Distance
	default:
		fmt.Println("Distance not recognized.")
		return
	}

	fmt.Println("Discrete Time Load Balancing: ")

	var distr probutil.Distr

	switch {
	case *full_ring:
		distr = dtlb_full.RingTypicalDistr(*T, *lam, *k, *n, *steps)

	case *full_complete:
		distr = dtlb_full.CompleteTypicalDistr(*T, *lam, *k, *n, *steps)

	case *local_ring_fp:
		_, distr, _, _ = dtlb_local.RingFixedPointIteration(*T, *lam, *k, *eps, *iters, *steps, dist)

	case *local_ring_realization:
		distr = dtlb_local.LocalRingRealizationTypicalDistr(*T, *lam, *k, *steps)

	case *mean_field:
		distr = dtlb_mean_field.TypicalDistr(*T, *lam, *k, *d, *steps)
	}

	b, err := json.Marshal(distr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Writing to file ...")

	if file_str != "" {
		err = ioutil.WriteFile(file_str, b, 0777)
		if err != nil {
			panic(err)
		}
	}

}
