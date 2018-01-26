package ctmc

import (
	"golang.org/x/exp/rand"
	"time"
)

type Event struct{
	Index int
	Inc int
}

// StepCTMC takes in an array of rates where a rate corresponds to event.
// StepCTCM then evolves the CTMC and returns
// 1. the time increment and
// 2. the index of the chosen event.
// Works by drawing an exponential random variable with rate = sum(rates)
// And choosing an event with probability proportional to the rate
func StepCTMC(rates []float64) (event int, time_inc float64) {
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	// get the sum of the rates
	sum_rates := 0.0
	for _, rate := range(rates) {
		sum_rates += rate
	}
	// get a scaled version of the rates which sums to one
	scaled_rates := make([]float64, len(rates))
	for i, rate := range(rates) {
		scaled_rates[i] = rate/sum_rates
	}
	// draw a time ~ exp(sum_rates)
	time_inc = r.ExpFloat64() / sum_rates
	// choose an event with probability proportional to the rate
	rand := r.Float64()
	cumsum := 0.0
	for i, rate := range(scaled_rates){
		cumsum += rate
		if cumsum > rand {
			event = i
			break 
		}
	}
	return event, time_inc
}