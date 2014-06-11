package gobot

import (
	"math"
	"math/rand"
	"reflect"
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
	e.Callbacks = append(e.Callbacks, f)
}

func Rand(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

func Call(thing interface{}, method string, params ...interface{}) []reflect.Value {
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return reflect.ValueOf(thing).MethodByName(method).Call(in)
}

func FieldByName(thing interface{}, field string) reflect.Value {
	return reflect.ValueOf(thing).FieldByName(field)
}
func FieldByNamePtr(thing interface{}, field string) reflect.Value {
	return reflect.ValueOf(thing).Elem().FieldByName(field)
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
