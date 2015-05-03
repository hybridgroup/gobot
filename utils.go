package gobot

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

var (
	// ErrUnknownEvent is the error resulting if the specified Event does not exist
	ErrUnknownEvent = errors.New("Event does not exist")
)

var eventError = func(e *Event) (err error) {
	if e == nil {
		err = ErrUnknownEvent
		log.Println(err.Error())
		return
	}
	return
}

var errFunc = func(t *testing.T, message string) {
	t.Errorf(message)
}

func logFailure(t *testing.T, message string) {
	_, file, line, _ := runtime.Caller(2)
	s := strings.Split(file, "/")
	errFunc(t, fmt.Sprintf("%v:%v: %v", s[len(s)-1], line, message))
}

// Assert checks if a and b are equal, emis a t.Errorf if they are not equal.
func Assert(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		logFailure(t, fmt.Sprintf("%v - \"%v\", should equal,  %v - \"%v\"",
			a, reflect.TypeOf(a), b, reflect.TypeOf(b)))
	}
}

// Refute checks if a and b are equal, emis a t.Errorf if they are equal.
func Refute(t *testing.T, a interface{}, b interface{}) {
	if reflect.DeepEqual(a, b) {
		logFailure(t, fmt.Sprintf("%v - \"%v\", should not equal,  %v - \"%v\"",
			a, reflect.TypeOf(a), b, reflect.TypeOf(b)))
	}
}

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

// Publish emits val to all subscribers of e. Returns ErrUnknownEvent if Event
// does not exist.
func Publish(e *Event, val interface{}) (err error) {
	if err = eventError(e); err == nil {
		e.Write(val)
	}
	return
}

// On executes f when e is Published to. Returns ErrUnknownEvent if Event
// does not exist.
func On(e *Event, f func(s interface{})) (err error) {
	if err = eventError(e); err == nil {
		e.Callbacks = append(e.Callbacks, callback{f, false})
	}
	return
}

// Once is similar to On except that it only executes f one time. Returns
//ErrUnknownEvent if Event does not exist.
func Once(e *Event, f func(s interface{})) (err error) {
	if err = eventError(e); err == nil {
		e.Callbacks = append(e.Callbacks, callback{f, true})
	}
	return
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
