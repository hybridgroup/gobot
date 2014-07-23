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
	// start a go routine to not bloc the function
	go func() {
		for {
			// wait for the ticker to tell us to run
			<-c
			// run the passed function in another go routine
			// so we don't slow down the loop.
			go f()
		}
	}()
}

// After triggers the passed function after `t` duration.
func After(t time.Duration, f func()) {
	time.AfterFunc(t, f)
}

func Publish(e *Event, val interface{}) {
	e.Write(val)
}

func On(e *Event, f func(s interface{})) {
	e.Callbacks = append(e.Callbacks, callback{f, false})
}

func Once(e *Event, f func(s interface{})) {
	e.Callbacks = append(e.Callbacks, callback{f, true})
}

func Rand(max int) int {
	i, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(i.Int64())
}

func FromScale(input, min, max float64) float64 {
	return (input - math.Min(min, max)) / (math.Max(min, max) - math.Min(min, max))
}

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
