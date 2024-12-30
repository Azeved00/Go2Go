package poisson

import (
	"errors"
	"math"
	"math/rand"
)

// PoissonProcess represents a Poisson process generator.
type PoissonProcess struct {
	Lambda float64      // Rate parameter (lambda)
	RNG    *rand.Rand   // Random number generator
}

// NewPoissonProcess creates a new PoissonProcess instance.
// Returns an error if the lambda is non-positive or RNG is nil.
func NewPoissonProcess(lambda float64, rng *rand.Rand) (*PoissonProcess, error) {
	if lambda <= 0 {
		return nil, errors.New("supplied rate parameter must be positive")
	}
	if rng == nil {
		return nil, errors.New("RNG cannot be nil")
	}
	return &PoissonProcess{
		Lambda: lambda,
		RNG:    rng,
	}, nil
}

// GetLambda returns the rate parameter (lambda).
func (p *PoissonProcess) GetLambda() float64 {
	return p.Lambda
}

// GetRNG returns the random number generator in use.
func (p *PoissonProcess) GetRNG() *rand.Rand {
	return p.RNG
}

// TimeForNextEvent generates the time until the next event.
// It uses the inverse transform sampling method for exponential distribution.
func (p *PoissonProcess) TimeForNextEvent() float64 {
	return -math.Log(1.0-p.RNG.Float64()) / p.Lambda
}

// Events returns the number of events in one unit of time.
func (p *PoissonProcess) Events() int {
	return p.EventsInInterval(1.0)
}

// EventsInInterval generates the number of events in a given time interval.
// It uses the inverse transform sampling method for Poisson distribution.
func (p *PoissonProcess) EventsInInterval(time float64) int {
	n := 0
	pFactor := math.Exp(-p.Lambda * time)
	s := pFactor
	u := p.RNG.Float64()

	for u > s {
		n++
		pFactor *= p.Lambda / float64(n)
		s += pFactor
	}
	return n
}
