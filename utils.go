package gobot

import (
	"github.com/tarm/goserial"
	"io"
	"math"
	"math/rand"
	"reflect"
	"time"
)

func Every(t string, f func()) {
	dur := parseDuration(t)
	go func() {
		for {
			time.Sleep(dur)
			go f()
		}
	}()
}

func After(t string, f func()) {
	dur := parseDuration(t)
	go func() {
		time.Sleep(dur)
		f()
	}()
}

func Publish(c chan interface{}, val interface{}) {
	select {
	case c <- val:
	default:
	}
}

func On(c chan interface{}, f func(s interface{})) {
	go func() {
		for s := range c {
			f(s)
		}
	}()
}

func Rand(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

func ConnectToSerial(port string, baud int) io.ReadWriteCloser {
	c := &serial.Config{Name: port, Baud: baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}
	return s
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

func parseDuration(t string) time.Duration {
	dur, err := time.ParseDuration(t)
	if err != nil {
		panic(err)
	}
	return dur
}
