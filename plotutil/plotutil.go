package plotutil

import (
	"github.com/mwortsma/particle_systems/probutil"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	//"gonum.org/v1/plot/plotutil"
	"image/color"
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

func getColors(n int) []color.RGBA {
	full := []color.RGBA{color.RGBA{R: 255, G: 0, B: 0, A:255}, color.RGBA{R: 0, G: 255, B: 155, A:155}, color.RGBA{R: 0, G: 255, B: 255, A:255}}
	return full[:n]
}

func PlotDistr(title string, distributions []probutil.Distr, labels []string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	colors := getColors(len(distributions))

	p.Title.Text = title
	p.X.Label.Text = "Path"
	p.Y.Label.Text = "Probability"

	sortedKeys := getSortedKeys(distributions[0])
	for i, v := range(distributions) {

		l, err := plotter.NewLine(getPoints(v, sortedKeys))
		if err != nil {
			panic(err)
		}
		l.LineStyle.Width = vg.Points(1)
		l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
		l.LineStyle.Color = colors[i]
		p.Add(l)
		p.Legend.Add(labels[i], l)
	}


	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 10*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
}