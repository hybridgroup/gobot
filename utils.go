package gobot

import (
	"encoding/json"
	"math/rand"
	"net"
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

func On(cs chan interface{}) interface{} {
	for s := range cs {
		return s
	}
	return nil
}

func ConnectTo(port string) net.Conn {
	tcpPort, err := net.Dial("tcp", port)
	if err != nil {
		panic(err)
	}
	return tcpPort
}

func Call(thing interface{}, method string, params ...interface{}) (result []reflect.Value, err error) {
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = reflect.ValueOf(thing).MethodByName(method).Call(in)
	return
}

func toJson(obj interface{}) string {
	b, _ := json.Marshal(obj)
	return string(b)
}
