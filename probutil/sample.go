package probutil

func Sample(d Distr, r float64) string {
	s := 0.
	for k, v := range d {
		if s += v; s > r {
			return k
		}
	}
	return ""
}

func SampleInt(d map[int]float64, r float64) int {
	s := 0.
	for k, v := range d {
		if s += v; s > r {
			return k
		}
	}
	return -1
}