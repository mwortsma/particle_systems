package plotutil

import (
	"github.com/mwortsma/particle_systems/probutil"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"sort"
)

func getPoints(distr probutil.Distr, keys []string) plotter.XYs {
	pts := make(plotter.XYs, len(keys))
	for i := 0; i < len(keys); i++ {
		pts[i].X = float64(i)
		pts[i].Y = distr[keys[i]]
	}
	return pts
}

func getSortedKeys(distr probutil.Distr) []string {
	keys := make([]string, 0, len(distr))
    for k := range distr {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    return keys
}

func PlotDistr(distr probutil.Distr) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	sortedKeys := getSortedKeys(distr)
	err = plotutil.AddLinePoints(p,
		"First", getPoints(distr, sortedKeys))
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}