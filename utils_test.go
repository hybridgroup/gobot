package gobot

import (
	"fmt"
	"testing"
	"time"
)

func TestEvery(t *testing.T) {
	i := 0
	Every(2*time.Millisecond, func() {
		i++
	})
	<-time.After(5 * time.Millisecond)
	Assert(t, i, 2)
}

func TestAfter(t *testing.T) {
	i := 0
	After(1*time.Millisecond, func() {
		i++
	})
	<-time.After(2 * time.Millisecond)
	Assert(t, i, 1)
}

func TestPublish(t *testing.T) {
	e := &Event{Chan: make(chan interface{}, 1)}
	Publish(e, 1)
	Publish(e, 2)
	Publish(e, 3)
	Publish(e, 4)
	i := <-e.Chan
	Assert(t, i, 1)
}

func TestOn(t *testing.T) {
	var i int
	e := NewEvent()
	On(e, func(data interface{}) {
		i = data.(int)
	})
	Publish(e, 10)
	<-time.After(1 * time.Millisecond)
	Assert(t, i, 10)
}
func TestOnce(t *testing.T) {
	i := 0
	e := NewEvent()
	Once(e, func(data interface{}) {
		i += data.(int)
	})
	On(e, func(data interface{}) {
		i += data.(int)
	})
	Publish(e, 10)
	<-time.After(1 * time.Millisecond)
	Publish(e, 10)
	<-time.After(1 * time.Millisecond)
	Assert(t, i, 30)
}

func TestFromScale(t *testing.T) {
	Assert(t, FromScale(5, 0, 10), 0.5)
}

func TestToScale(t *testing.T) {
	Assert(t, ToScale(500, 0, 10), 10.0)
	Assert(t, ToScale(-1, 0, 10), 0.0)
	Assert(t, ToScale(0.5, 0, 10), 5.0)
}

func TestRand(t *testing.T) {
	a := Rand(1000)
	b := Rand(1000)
	if a == b {
		t.Error(fmt.Sprintf("%v should not equal %v", a, b))
	}
}
