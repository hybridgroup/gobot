package gobot

import (
	"strings"
	"testing"
	"time"

	"gobot.io/x/gobot/gobottest"
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

func TestEveryWhenStopped(t *testing.T) {
	sem := make(chan bool)

	done := Every(50*time.Millisecond, func() {
		sem <- true
	})

	select {
	case <-sem:
		done.Stop()
	case <-time.After(60 * time.Millisecond):
		t.Errorf("Every was not called")
	}

	select {
	case <-time.After(60 * time.Millisecond):
	case <-sem:
		t.Error("Every should have stopped")
	}
}

func TestAfter(t *testing.T) {
	i := 0
	sem := make(chan bool)

	After(1*time.Millisecond, func() {
		i++
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("After was not called")
	}

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
	a := Rand(10000)
	b := Rand(10000)
	if a == b {
		t.Errorf("%v should not equal %v", a, b)
	}
}

func TestDefaultName(t *testing.T) {
	name := DefaultName("tester")
	gobottest.Assert(t, strings.Contains(name, "tester"), true)
}
