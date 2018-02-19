package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/mwortsma/particle_systems/dtpp/dtpp_local"

	"github.com/mwortsma/particle_systems/probutil"
	"io/ioutil"
)

func main() {

	rec := flag.Bool("rec", false, "rec")

	d := flag.Int("d", 2, "degree of a noe")
	n := flag.Int("n", 10, "number of nodes")
	k := flag.Int("k", 3, "states")
	T := flag.Int("T", 2, "time horizon. T>0")
	beta := flag.Float64("beta", 1.5, "temp inverse")
	tau := flag.Int("tau", -1, "how many steps to look back")

	//steps := flag.Int("steps", 100, "how many samples used in generating the empirical distribtuion")
	var file_str string
	flag.StringVar(&file_str, "file", "", "where to save the distribution.")

	flag.Parse()


	var distr probutil.GenDistr

	switch {


	case *rec:
		distr = dtpp_local.Run(*T,*tau, *d, *beta,*k, *n)
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
