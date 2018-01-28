package ctmc

import (
	"github.com/mwortsma/particle_systems/graphutil"
)

func GetCPRatesAndEvents(
	X []int, 
	lam float64, 
	G graphutil.Graph,
	k int) ([]float64, []Event) {

	rates := make([]float64, 0)
	events := make([]Event, 0)

	for i := range(X) {

		if X[i] == 1 {
			// recover
			rates = append(rates, 1)
			events = append(events, Event{Index: i, Inc: -1})

			for _, j := range(G[i]) {
				if X[j] == 0 {
					// infect
					rates = append(rates, lam/float64(k))
					events = append(events, Event{Index: j, Inc: 1})
				}
			}
		}
	}
	return rates, events
}