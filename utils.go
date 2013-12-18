package gobot

import (
	"encoding/json"
	"github.com/tarm/goserial"
	"io"
	"math/rand"
	"net"
	"reflect"
	"regexp"
	"time"
)

type Port io.ReadWriteCloser

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

func parseDuration(t string) time.Duration {
	dur, err := time.ParseDuration(t)
	if err != nil {
		panic(err)
	}
	return dur
}

func Rand(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max)
}

func On(cs chan interface{}, f func(s interface{})) {
	go func() {
		for s := range cs {
			f(s)
		}
	}()
}

func ConnectToTcp(port string) io.ReadWriteCloser {
	tcpPort, err := net.Dial("tcp", port)
	if err != nil {
		panic(err)
	}
	return tcpPort
}

func ConnectToSerial(port string, baud int) io.ReadWriteCloser {
	c := &serial.Config{Name: port, Baud: baud}
	s, err := serial.OpenPort(c)
	if err != nil {
		panic(err)
	}
	return s
}

func IsUrl(url string) bool {
	rp := regexp.MustCompile("([^A-Za-z0-9]+).([^A-Za-z0-9]+).([^A-Za-z0-9]+)")
	return rp.MatchString(url)
}

func Call(thing interface{}, method string, params ...interface{}) []reflect.Value {
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return reflect.ValueOf(thing).MethodByName(method).Call(in)
}

func toJson(obj interface{}) string {
	b, _ := json.Marshal(obj)
	return string(b)
}
