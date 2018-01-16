package main

import (
	"fmt"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_full_gengraph"
	"github.com/mwortsma/particle_systems/plotutil"
)

func main() {
	distr := dtlb_full_gengraph.RingTypicalDistr(4, 0.8, 10, 100000)
	plotutil.PlotDistr(distr)
	fmt.Println(len(distr))
}
