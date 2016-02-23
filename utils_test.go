package gobot

import (
	"fmt"
	"testing"
	"time"
)

type nopFailer struct{ msg string }

func (tf *nopFailer) Errorf(format string, args ...interface{}) {
	tf.msg = fmt.Sprintf(format, args...)
}

func TestAssert(t *testing.T) {
	nf := &nopFailer{}

	Assert(nf, 1, 1)
	if nf.msg != "" {
		t.Errorf("Assert failed: 1 should equal 1")
	}

	Assert(nf, 1, 2)

	if nf.msg != `utils_test.go:23: 1 - "int", should equal,  2 - "int"` {
		t.Errorf("Assert failed: 1 should not equal 2")
	}
}

func TestRefute(t *testing.T) {
	nf := &nopFailer{}

	Refute(nf, 1, 2)
	if nf.msg != "" {
		t.Errorf("Refute failed: 1 should not be 2")
	}

	Refute(nf, 1, 1)
	if nf.msg != `utils_test.go:38: 1 - "int", should not equal,  1 - "int"` {
		t.Errorf("Refute failed: 1 should not be 1")
	}
}

func TestEvery(t *testing.T) {
	i := 0
	begin := time.Now().UnixNano()
	sem := make(chan int64, 1)
	Every(2*time.Millisecond, func() {
		i++
		if i == 2 {
			sem <- time.Now().UnixNano()
		}
	})
	end := <-sem
	if end-begin < 4000000 {
		t.Error("Test should have taken at least 4 milliseconds")
	}
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
	c := make(chan interface{}, 1)

	cb := callback{
		f: func(val interface{}) {
			c <- val
		},
	}

	e := &Event{Callbacks: []callback{cb}}
	Publish(e, 1)
	// this enforeces Event's callback goroutine to actually start
	<-time.After(10 * time.Millisecond)
	Publish(e, 2)
	Publish(e, 3)
	Publish(e, 4)
	i := <-c
	Assert(t, i, 1)

	var e1 = (*Event)(nil)
	Assert(t, Publish(e1, 4), ErrUnknownEvent)
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

	var e1 = (*Event)(nil)
	err := On(e1, func(data interface{}) {
		i = data.(int)
	})
	Assert(t, err, ErrUnknownEvent)
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

	var e1 = (*Event)(nil)
	err := Once(e1, func(data interface{}) {
		i = data.(int)
	})
	Assert(t, err, ErrUnknownEvent)
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
