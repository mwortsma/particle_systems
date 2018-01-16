package plotutil

import (
	"github.com/mwortsma/particle_systems/probutil"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	//"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"image/color"
)

func getPoints(distr probutil.Distr, keys []string) plotter.XYs {
	pts := make(plotter.XYs, len(keys))
	for i := 0; i < len(keys); i++ {
		pts[i].X = float64(i)
		pts[i].Y = distr[keys[i]]
	}
	return pts
}

func getColors(n int) []color.RGBA {
	full := []color.RGBA{{R: 255, G: 0, B: 0, A: 255}, {R: 0, G: 0, B: 255, A: 255}, {R: 0, G: 255, B: 0, A: 255}}
	return full[:n]
}

func PlotDistr(file string, title string, distributions []probutil.Distr, labels []string) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	colors := getColors(len(distributions))

	p.Title.Text = title
	p.X.Label.Text = "Path"
	p.Y.Label.Text = "Probability"

	sortedKeys := probutil.SharedSortedKeys(distributions)
	for i, v := range distributions {

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
	p.Legend.Left = true
	p.Legend.Top = true
	p.Legend.TextStyle.Font.Size = 0.25 * vg.Inch
	p.Title.TextStyle.Font.Size = 0.25 * vg.Inch

	// Save the plot to a PNG file.
	if err := p.Save(20*vg.Inch, 10*vg.Inch, file); err != nil {
		panic(err)
	}
}
