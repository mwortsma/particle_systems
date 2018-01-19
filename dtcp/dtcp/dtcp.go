package main

import (
	"flag"
	"fmt"
	"github.com/mwortsma/particle_systems/dtcp/dtcp_local"
	"github.com/mwortsma/particle_systems/dtcp/dtcp_full"
	"github.com/mwortsma/particle_systems/dtcp/dtcp_mean_field"
	"github.com/mwortsma/particle_systems/probutil"
	"encoding/json"
	"io/ioutil"
)

func main() {

	// Discrete Time Contact Process.

	// Generel Arguments
	// -d=x (degree of a node)
	// -n=x (number of nodes)
	// -T=x (time horizon. T > 0)
	// -p=x (transition 0->1 with probability (p/d)*sum(neighbors))
	// -q=x (transntion 1->0 with probability q)
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
	// // // Note: Depth optional. Defaults to 1+2*T.
	// // 1.2 Complete (dtcp -full_complete -n=?)
	// // 1.3 Regular Tree (dtcp -full_tree -d=3 -depth=?)
	// // // Note: Depth optional. Defaults to T.
	
	// 2. Local Simulations (Using the fixed point algorithm)
	// // 1.1 Ring (dtcp -local_ring)
	// // 1.2 Regular Tree (dtcp -regular_tree -d=3)

	// 3. Mean Field Simulations (TODO: In progress...)
	// // 3.1 Mean Field (mean_field) (dtcp -mean_field -d=12)

	// Defining the arguments.
	full_ring := flag.Bool("full_ring", false, "Full sim on the ring")
	full_complete := flag.Bool("full_complete", false, "Full sim on the complete graph")
	full_tree := flag.Bool("full_tree", false, "Full sim on a regular tree")

	local_ring := flag.Bool("local_ring", false, "Local sim on the ring")
	local_tree := flag.Bool("local_tree", false, "Local sim on a regular tree")

	mean_field := flag.Bool("mean_field", false, "Mean Field simulation.")

	d := flag.Int("d", -1, "degree of a noe")
	n := flag.Int("n", -1, "number of nodes")
	T := flag.Int("T", 2, "time horizon. T>0")
	p := flag.Float64("p", 2.0/3.0, "transition 0->1 with probability (p/d)*sum(neighbors)")
	q := flag.Float64("q", 1.0/3.0, "transntion 1->0 with probability q")
	nu := flag.Float64("nu", 0.5, "P(X_0 = 1)")
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

	fmt.Println("Discrete Time Contact Process: ")

	var distr probutil.Distr

	switch {
	case *full_ring:
		distr = dtcp_full.RingTypicalDistr(*T, *p, *q, *nu, *n, *steps)

	case *full_complete:
		distr = dtcp_full.CompleteTypicalDistr(*T, *p, *q, *nu, *n, *steps)

	case *full_tree:
		distr = dtcp_full.RegTreeTypicalDistr(*T, *d, *p, *q, *nu, *steps)

	case *local_ring:
		_, distr, _, _ = dtcp_local.RegTreeFixedPointIteration(*T,2,*p,*q,*nu,*eps,*iters,*steps,dist)	

	case *local_tree:
		_, distr, _, _ = dtcp_local.RegTreeFixedPointIteration(*T,*d,*p,*q,*nu,*eps,*iters,*steps,dist)			

	case *mean_field:
		distr = dtcp_mean_field.TypicalDistr(*T, *p, *q, *nu, *d, *steps)
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