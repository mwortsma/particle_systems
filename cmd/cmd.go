package main

import (
	"github.com/mwortsma/particle_systems/dtlb/dtlb_full"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_local"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_mean_field"
	//"github.com/mwortsma/particle_systems/dtlb/dtlb_local"
	"github.com/mwortsma/particle_systems/dtcp/dtcp_mean_field"
	"flag"
	"fmt"
	"github.com/mwortsma/particle_systems/dtcp/dtcp_full"
	"github.com/mwortsma/particle_systems/plotutil"
	"github.com/mwortsma/particle_systems/probutil"
	"strings"
)

func main() {

	// Type of simulation. Options are
	// 1. dtlb (discrete time load balancing)
	// 2. dtcp (discrete time contact process)
	// ad-hoc (put ad-hoc stuff here)
	var type_string string
	flag.StringVar(&type_string, "type", "dtlb", "Type of simulation.")

	// type of distances. Options are
	// 1. L1 (L1 distance)
	var dist_string string
	flag.StringVar(&dist_string, "distance", "L1", "Type of distance")

	// Choose which you would like to simulate
	full_ring := flag.Bool("full_ring", false, "Full ring simulation.")
	local_ring := flag.Bool("local_ring", false, "Local ring simulation.")

	full_complete := flag.Bool("full_complete", false, "Full complete graph simulation.")

	full_tree := flag.Bool("full_tree", false, "Full tree simulation.")
	local_tree := flag.Bool("local_tree", false, "Local tree simulation.")

	mean_field := flag.Bool("mean_field", false, "Mean Field simulation.")

	eps := flag.Float64("epsilon", 0.001, "Epsilon")
	iters := flag.Int("iters", 4, "Number of Iterations")
	steps := flag.Int("steps", 100000, "Steps per Iteration")

	T := flag.Int("T", 4, "Steps")
	lam := flag.Float64("lambda", 0.8, "Rate")
	k := flag.Int("buffer", 5, "Buffer")

	// important to default -1
	n := flag.Int("n", -1, "number of particles")
	d := flag.Int("d", -1, "degree")

	nu := flag.Float64("nu", 0.5, "P(X(0) = 1)")
	p := flag.Float64("p", 2.0/3.0, "transition 0->1 with prob (p/d)*sum(neighbors)")
	q := flag.Float64("q", 1.0/3.0, "transition 1->0 with prob q")

	flag.Parse()

	var dist probutil.Distance
	switch dist_string {
	case "L1":
		dist = probutil.L1Distance
	default:
		fmt.Println("Distance not recognized.")
		return
	}

	distrs := make([]probutil.Distr, 0)
	labels := make([]string, 0)
	title := ""

	switch type_string {

	//////////////////////////////////////////////////////////////////
	//////////////     Discrete time load balancing     //////////////
	//////////////////////////////////////////////////////////////////

	case "dtlb":

		fmt.Println("Discrete Time Load Balncing")

		if *full_ring {
			name := "Full (ring)"
			if *n > 0 {
				name = name + fmt.Sprintf(" n=%d", *n)
			}
			fmt.Println("Running ", name)
			full_distr := dtlb_full.RingTypicalDistr(*T, *lam, *k, *n, *steps)
			distrs = append(distrs, full_distr)
			labels = append(labels, name)
		}
		if *full_complete {
			name := fmt.Sprintf("Full (complete) n=%d", *n)
			fmt.Println("Running", name)
			full_distr := dtlb_full.CompleteTypicalDistr(*T, *lam, *k, *n, *steps)
			distrs = append(distrs, full_distr)
			labels = append(labels, name)
		}
		if *full_tree {
			fmt.Println("full-tree not implemented.")
		}

		if *local_ring {
			fmt.Println("Running local (ring)")
			_, local_distr, _, _ := dtlb_local.RingFixedPointIteration(
				*T, *lam, *k, *eps, *iters, *steps, dist)
			distrs = append(distrs, local_distr)
			labels = append(labels, "Local (ring)")
		}
		if *local_tree {
			fmt.Println("Local Tree not implemented.")
		}

		if *mean_field {
			if *d < 0 {
				*d = *n-1
			}
			name := fmt.Sprintf("Mean Field degree=%d", *d)
			fmt.Println("Running ", name)
			mean_field_distr := dtlb_mean_field.TypicalDistr(*T, *lam, *k, *d, *steps)
			distrs = append(distrs, mean_field_distr)
			labels = append(labels, name)
		}

		title = fmt.Sprintf("Discrete Time Load Balncing T=%d, lambda=%0.3f, buffer=%d, steps=%0.2e", *T, *lam, *k, float64(*steps))

	//////////////////////////////////////////////////////////////////
	//////////////    Discrete time contact process     //////////////
	//////////////////////////////////////////////////////////////////

	case "dtcp":

		fmt.Println("Discrete Time Contact Process")

		if *full_ring {
			name := "Full (ring)"
			if *n > 0 {
				name = name + fmt.Sprintf(" n=%d", *n)
			}
			fmt.Println("Running ", name)
			full_distr := dtcp_full.RingTypicalDistr(*T, *p, *q, *nu, *n, *steps)
			distrs = append(distrs, full_distr)
			labels = append(labels, name)
		}
		if *full_complete {
			name := fmt.Sprintf("Full (complete) n=%d", *n)
			fmt.Println("Running", name)
			full_distr := dtcp_full.CompleteTypicalDistr(*T, *p, *q, *nu, *n, *steps)
			distrs = append(distrs, full_distr)
			labels = append(labels, name)
		}
		if *full_tree {
			name := fmt.Sprintf("Full Regular Tree d=%d", *d)
			fmt.Println("Running", name)
			full_distr := dtcp_full.RegTreeTypicalDistr(*T, *d, *p, *q, *nu, *steps)
			distrs = append(distrs, full_distr)
			labels = append(labels, name)
		}

		if *mean_field {
			if *d < 0 {
				*d = *n-1
			}
			name := fmt.Sprintf("Mean Field degree=%d", *d)
			fmt.Println("Running ", name)
			mean_field_distr := dtcp_mean_field.TypicalDistr(*T,*p,*q,*nu,*d,*steps)
			distrs = append(distrs, mean_field_distr)
			labels = append(labels, name)
		}

		title = fmt.Sprintf("Discrete Time Contact Process T=%d, p=%0.3f, q=%0.3f, steps=%0.2e", *T, *p, *q, float64(*steps))

	case "ad-hoc":
		fmt.Println("Nothing ad-hoc right now.")
		return
	default:
		fmt.Println("Type not recognized.")
		return
	}

	title = title + " distances"
	for i := 0; i < len(distrs); i++ {
		for j := i + 1; j < len(distrs); j++ {
			d := dist(distrs[i], distrs[j])
			fmt.Println(fmt.Sprintf("Distance %s vs %s: %f", labels[i], labels[j], d))
			title = title + fmt.Sprintf(" %0.4f", d)
		}
	}

	file := "plots/" + strings.Replace(title, " ", "_", -1) + ".png"
	plotutil.PlotDistr(file, title, distrs, labels)
}
