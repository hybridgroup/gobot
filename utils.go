package gobot

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
	"time"
)

// Every triggers f every t time.Duration until the end of days, or when a Stop()
// is called on the Ticker that is returned by the Every function.
// It does not wait for the previous execution of f to finish before
// it fires the next f.
func Every(t time.Duration, f func()) *time.Ticker {
	ticker := time.NewTicker(t)

	go func() {
		for {
			<-ticker.C
			f()
		}
	}()

	return ticker
}

// After triggers f after t duration.
func After(t time.Duration, f func()) {
	time.AfterFunc(t, f)
}

// Rand returns a positive random int up to maximum
func Rand(maximum int) int {
	i, _ := rand.Int(rand.Reader, big.NewInt(int64(maximum)))
	return int(i.Int64())
}

// FromScale returns a converted input from minimum, maximum to 0.0...1.0.
func FromScale(input, minimum, maximum float64) float64 {
	return (input - math.Min(minimum, maximum)) / (math.Max(minimum, maximum) - math.Min(minimum, maximum))
}

// ToScale returns a converted input from 0...1 to minimum...maximum scale.
// If input is less than minimum then ToScale returns minimum.
// If input is greater than maximum then ToScale returns maximum
func ToScale(input, minimum, maximum float64) float64 {
	i := input*(math.Max(minimum, maximum)-math.Min(minimum, maximum)) + math.Min(minimum, maximum)
	switch {
	case i < math.Min(minimum, maximum):
		return math.Min(minimum, maximum)
	case i > math.Max(minimum, maximum):
		return math.Max(minimum, maximum)
	default:
		return i
	}
}

// Rescale performs a direct linear rescaling of a number from one scale to another.
func Rescale(input, fromMin, fromMax, toMin, toMax float64) float64 {
	return (input-fromMin)*(toMax-toMin)/(fromMax-fromMin) + toMin
}

// DefaultName returns a sensible random default name
// for a robot, adaptor or driver
func DefaultName(name string) string {
	return fmt.Sprintf("%s-%X", name, Rand(int(^uint(0)>>1)))
}
