package gobot

import (
	"fmt"
	"testing"
	"time"

	"github.com/hybridgroup/gobot/gobottest"
)

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

func TestEveryWhenDone(t *testing.T) {
	i := 0
	done := Every(20*time.Millisecond, func() {
		i++
	})
	<-time.After(10 * time.Millisecond)
	done <- true
	<-time.After(50 * time.Millisecond)
	if i > 1 {
		t.Error("Test should have stopped after 20ms")
	}
}

func TestAfter(t *testing.T) {
	i := 0
	After(1*time.Millisecond, func() {
		i++
	})
	<-time.After(2 * time.Millisecond)
	gobottest.Assert(t, i, 1)
}

func TestFromScale(t *testing.T) {
	gobottest.Assert(t, FromScale(5, 0, 10), 0.5)
}

func TestToScale(t *testing.T) {
	gobottest.Assert(t, ToScale(500, 0, 10), 10.0)
	gobottest.Assert(t, ToScale(-1, 0, 10), 0.0)
	gobottest.Assert(t, ToScale(0.5, 0, 10), 5.0)
}

func TestRand(t *testing.T) {
	a := Rand(1000)
	b := Rand(1000)
	if a == b {
		t.Error(fmt.Sprintf("%v should not equal %v", a, b))
	}
}
