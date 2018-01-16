package main

import (
	"fmt"
	"github.com/mwortsma/particle_systems/dtlb/dtlb_full_gengraph"
	"github.com/mwortsma/particle_systems/plotutil"
	"github.com/mwortsma/particle_systems/probutil"
)

func main() {
	distr := dtlb_full_gengraph.RingTypicalDistr(4, 0.8, 10, 100000)
	distr2 := make(map[string]float64)
	for k := range(distr) {
		distr2[k] = distr[k] + 0.1
	}
	plotutil.PlotDistr("test", []probutil.Distr{distr,distr2}, []string{"first", "second"})
	fmt.Println(len(distr))
}
