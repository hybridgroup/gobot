package gobot

import (
	"crypto/rand"
	"math"
	"math/big"
	"time"
)

// Every triggers f every `t` time until the end of days.
func Every(t time.Duration, f func()) {
	c := time.Tick(t)

	go func() {
		for {
			<-c
			go f()
		}
	}()
}

// After triggers the passed function after `t` duration.
func After(t time.Duration, f func()) {
	time.AfterFunc(t, f)
}

// Publish emits an event by writting value
func Publish(e *Event, val interface{}) {
	e.Write(val)
}

// On adds `f` to callbacks that are executed on specified event
func On(e *Event, f func(s interface{})) {
	e.Callbacks = append(e.Callbacks, callback{f, false})
}

// Once adds `f` to callbacks that are executed on specified event
// and sets flag to be called only once
func Once(e *Event, f func(s interface{})) {
	e.Callbacks = append(e.Callbacks, callback{f, true})
}

// Rand generates random int lower than max
func Rand(max int) int {
	i, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(i.Int64())
}

// FromScale creates a scale using min and max values
// to be used in combination with ToScale
func FromScale(input, min, max float64) float64 {
	return (input - math.Min(min, max)) / (math.Max(min, max) - math.Min(min, max))
}

// ToScale is used with FromScale to return input converted to new scale
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
