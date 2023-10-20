package gobot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEvery(t *testing.T) {
	i := 0
	begin := time.Now()
	sem := make(chan time.Time, 1)
	Every(2*time.Millisecond, func() {
		i++
		if i == 2 {
			sem <- time.Now()
		}
	})
	<-sem
	if time.Since(begin) < 4*time.Millisecond {
		t.Error("Test should have taken at least 4 milliseconds")
	}
}

func TestEveryWhenStopped(t *testing.T) {
	sem := make(chan bool)

	done := Every(100*time.Millisecond, func() {
		sem <- true
	})

	select {
	case <-sem:
		done.Stop()
	case <-time.After(190 * time.Millisecond):
		done.Stop()
		t.Errorf("Every was not called")
	}

	select {
	case <-time.After(190 * time.Millisecond):
	case <-sem:
		t.Error("Every should have stopped")
	}
}

func TestAfter(t *testing.T) {
	i := 0
	sem := make(chan bool)

	After(100*time.Millisecond, func() {
		i++
		sem <- true
	})

	select {
	case <-sem:
	case <-time.After(190 * time.Millisecond):
		t.Errorf("After was not called")
	}

	assert.Equal(t, 1, i)
}

func TestFromScale(t *testing.T) {
	assert.Equal(t, 0.5, FromScale(5, 0, 10))
}

func TestToScale(t *testing.T) {
	assert.Equal(t, 10.0, ToScale(500, 0, 10))
	assert.Equal(t, 0.0, ToScale(-1, 0, 10))
	assert.Equal(t, 5.0, ToScale(0.5, 0, 10))
}

func TestRescale(t *testing.T) {
	assert.Equal(t, 5.0, Rescale(500, 0, 1000, 0, 10))
	assert.Equal(t, 490.0, Rescale(-1.0, -1, 0, 490, 350))
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
	assert.Contains(t, name, "tester")
}
