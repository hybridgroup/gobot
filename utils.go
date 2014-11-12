package gobot

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

// Every triggers f every t time until the end of days. It does not wait for the
// previous execution of f to finish before it fires the next f.
func Every(t time.Duration, f func()) {
	c := time.Tick(t)

	go func() {
		for {
			<-c
			go f()
		}
	}()
}

// After triggers f after t duration.
func After(t time.Duration, f func()) {
	time.AfterFunc(t, f)
}

// Publish emits val to all subscribers of e.
func Publish(e *Event, val interface{}) {
	e.Write(val)
}

// On executes f when e is Published to.
func On(e *Event, f func(s interface{})) {
	e.Callbacks = append(e.Callbacks, callback{f, false})
}

// Once is similar to On except that it only executes f one time.
func Once(e *Event, f func(s interface{})) {
	e.Callbacks = append(e.Callbacks, callback{f, true})
}

// Rand returns a positive random int up to max
func Rand(max int) int {
	i, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(i.Int64())
}

// FromScale returns a converted input from min, max to 0.0...1.0.
func FromScale(input, min, max float64) float64 {
	return (input - math.Min(min, max)) / (math.Max(min, max) - math.Min(min, max))
}

// ToScale returns a converted input from 0...1 to min...max scale.
// If input is less than min then ToScale returns min.
// If input is greater than max then ToScale returns max
func ToScale(input, min, max float64) float64 {
	i := input*(math.Max(min, max)-math.Min(min, max)) + math.Min(min, max)
	if i < math.Min(min, max) {
		return math.Min(min, max)
	} else if i > math.Max(min, max) {
		return math.Max(min, max)
	} else {
		return i
	}
}
