package gobot

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventerAddEvent(t *testing.T) {
	e := NewEventer()
	e.AddEvent("test")

	if _, ok := e.Events()["test"]; !ok {
		require.Fail(t, "Could not add event to list of Event names")
	}
	assert.Equal(t, "test", e.Event("test"))
	assert.Equal(t, "", e.Event("unknown"))
}

func TestEventerDeleteEvent(t *testing.T) {
	e := NewEventer()
	e.AddEvent("test1")
	e.DeleteEvent("test1")

	if _, ok := e.Events()["test1"]; ok {
		require.Fail(t, "Could not add delete event from list of Event names")
	}
}

func TestEventerOn(t *testing.T) {
	e := NewEventer()

	sem := make(chan bool)
	_ = e.On("test", func(data interface{}) {
		sem <- true
	})

	// wait some time to ensure the eventer go routine is working
	time.Sleep(10 * time.Millisecond)

	e.Publish("test", true)

	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		require.Fail(t, "On was not called")
	}
}

func TestEventerOnce(t *testing.T) {
	e := NewEventer()

	sem := make(chan bool)
	_ = e.Once("test", func(data interface{}) {
		sem <- true
	})

	// wait some time to ensure the eventer go routine is working
	time.Sleep(10 * time.Millisecond)

	e.Publish("test", true)

	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		require.Fail(t, "Once was not called")
	}

	e.Publish("test", true)

	select {
	case <-sem:
		require.Fail(t, "Once was called twice")
	case <-time.After(10 * time.Millisecond):
	}
}
