package plotutil

import (
	"github.com/mwortsma/particle_systems/probutil"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
    "github.com/gonum/plot/vg/draw"
	"sort"
)

func getPoints(distr Distr, keys []string) plotter.XYs {
	pts = make(plotter.XYs, len(keys))
	for i = 0; i < len(keys); i++ {
		pts[i].X = i
		pts[i].Y = distr[keys[i]]
	}
	return pts
}

func getSortedKeys(distr Distr) {
	keys := make([]int, 0, len(distr))
    for k := range distr {
        keys = append(keys, k)
    }
    return sort.Strings(keys)
}

func PlotDistr(distr Distr) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	sortedKeys := getSortedKeys(distr)
	err = plotutil.AddLinePoints(p,
		"First", getPoints(distr, sortedKeys),
		"Second", getPoints(distr, sortedKeys),
		"Third", getPoints(distr, sortedKeys))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}