package main

import (
	"flag"
	"fmt"
	//"github.com/mwortsma/particle_systems/ctcp/ctcp_local"
	"github.com/mwortsma/particle_systems/ctcp/ctcp_full"
	//"github.com/mwortsma/particle_systems/ctcp/ctcp_mean_field"
	"github.com/mwortsma/particle_systems/probutil"
	"encoding/json"
	"io/ioutil"
)

func main() {

	// Continuous Time Contact Process.

	// Generel Arguments
	// -d=x (degree of a node)
	// -n=x (number of nodes)
	// -T=x (time horizon. T > 0)
	// -lam=x (incoming rate at each)
	// -nu=x (P(X_0 = 1))
	// -steps=x (how many samples used in generating the empirical distribtuion)
	// -file=x (which file to save the distribution to. empty = do not save)

	// Arguments specific to fixed point iteration algorithm.
	// -distance=x (options are: L1 (used as defualt), more to come)
	// -iters=x (for the fixed point algorithm, how many iterations to run)
	// -eps=x (threshold distance between typical particle distributions)

	// Types:

	// 1. Full Simulations
	// // 1.1 Ring (full_ring)

	// Defining the arguments.
	full_ring := flag.Bool("full_ring", false, "Full sim on the ring")
	full_complete := flag.Bool("full_complete", false, "Full sim on the complete graph")

	mean_field := flag.Bool("mean_field", false, "Mean Field simulation.")

	n := flag.Int("n", -1, "number of nodes")
	T := flag.Int("T", 2, "time horizon. T>0")
	lam := flag.Float64("lam", 0.8, "incoming rate at each node")
	dt := flag.Float64("dt", 0.01, "How to discritize time.")
	nu := flag.Float64("nu", 0.5, "P(X_0 = 1)")
	steps := flag.Int("steps", 100, "how many samples used in generating the empirical distribtuion")
	var file_str string
	flag.StringVar(&file_str, "file", "", "where to save the distribution.")

	/*
	eps := flag.Float64("epsilon", 0.001, "threshold distance between typical particle distributions")
	iters := flag.Int("iters", 4, "for the fixed point algorithm, how many iterations to run")
	var distance_str string
	flag.StringVar(&distance_str, "distance", "L1", "type of distance between distributions. Options: L1, more to come...")
	*/

	flag.Parse()

	fmt.Println("Continuous Time Contact Process: ")

	var distr probutil.ContDistr

	switch {
	case *full_ring:
		distr = ctcp_full.RingTypicalDistr(*T, *lam, *nu, *dt, *n, *steps)

	case *full_complete:
		fmt.Println("Not yet implemented")

	case *mean_field:
		fmt.Println("Not yet implemented.")
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