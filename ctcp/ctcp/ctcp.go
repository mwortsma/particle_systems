package main

import (
	"flag"
	"fmt"
	//"github.com/mwortsma/particle_systems/ctcp/ctcp_local"
	"github.com/mwortsma/particle_systems/ctcp/ctcp_full"
	"github.com/mwortsma/particle_systems/ctcp/ctcp_mean_field"
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

	// Types: TODO

	// Defining the arguments.
	full_ring := flag.Bool("full_ring", false, "Full sim on the ring")
	full_complete := flag.Bool("full_complete", false, "Full sim on the complete graph")
	full_regtree := flag.Bool("full_regtree", false, "Full sim on a regular tree")


	mean_field_fp := flag.Bool("mean_field_fp", false, "Mean Field fp simulation.")
	mean_field_realization := flag.Bool("mean_field_realization", false, "Mean Field realization.")


	n := flag.Int("n", -1, "number of nodes")
	T := flag.Float64("T", 2, "time horizon. T>0")
	lam := flag.Float64("lam", 0.8, "incoming rate at each node")
	dt := flag.Float64("dt", 0.25, "How to discritize time.")
	nu := flag.Float64("nu", 0.5, "P(X_0 = 1)")

	depth := flag.Int("depth", 3, "Depth of tree")
	d := flag.Int("d", 3, "Degree of tree")
	steps := flag.Int("steps", 100, "how many samples used in generating the empirical distribtuion")
	var file_str string
	flag.StringVar(&file_str, "file", "", "where to save the distribution.")

	
	eps := flag.Float64("epsilon", 0.001, "threshold distance between typical particle distributions")
	iters := flag.Int("iters", 4, "for the fixed point algorithm, how many iterations to run")

	flag.Parse()

	fmt.Println("Continuous Time Contact Process: ")

	var distr probutil.ContDistr

	switch {
	case *full_ring:
		distr = ctcp_full.RingTypicalDistr(*T, *lam, *nu, *dt, *n, *steps)

	case *full_complete:
		distr = ctcp_full.CompleteTypicalDistr(*T, *lam, *nu, *dt, *n, *steps)

	case *full_regtree:
		distr = ctcp_full.RegTreeTypicalDistr(*depth, *T, *lam, *d, *nu , *dt, *steps)

	case *mean_field_fp:
		distr = ctcp_mean_field.MeanFieldFixedPointIteration(
			*T,*lam,*nu,*dt,*eps,*iters,*steps,probutil.ContL1Distance)

	case *mean_field_realization:
		distr = ctcp_mean_field.RealizationTypicalDistr(*T, *lam, *nu, *dt, *steps)

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