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

var eventError = func(e *Event) (err error) {
	if e == nil {
		err = errors.New("Event does not exist")
		log.Println(err.Error())
		return
	}
	return
}

func logFailure(t *testing.T, message string) {
	_, file, line, _ := runtime.Caller(2)
	s := strings.Split(file, "/")
	t.Errorf("%v:%v: %v", s[len(s)-1], line, message)
}
func Assert(t *testing.T, a interface{}, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		logFailure(t, fmt.Sprintf("%v - \"%v\", should equal,  %v - \"%v\"",
			a, reflect.TypeOf(a), b, reflect.TypeOf(b)))
	}
}

func Refute(t *testing.T, a interface{}, b interface{}) {
	if reflect.DeepEqual(a, b) {
		logFailure(t, fmt.Sprintf("%v - \"%v\", should not equal,  %v - \"%v\"",
			a, reflect.TypeOf(a), b, reflect.TypeOf(b)))
	}
}

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
func Publish(e *Event, val interface{}) (err error) {
	if err = eventError(e); err == nil {
		e.Write(val)
	}
	return
}

// On adds `f` to callbacks that are executed on specified event
func On(e *Event, f func(s interface{})) (err error) {
	if err = eventError(e); err == nil {
		e.Callbacks = append(e.Callbacks, callback{f, false})
	}
	return
}

// Once adds `f` to callbacks that are executed on specified event
// and sets flag to be called only once
func Once(e *Event, f func(s interface{})) (err error) {
	if err = eventError(e); err == nil {
		e.Callbacks = append(e.Callbacks, callback{f, true})
	}
	return
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
