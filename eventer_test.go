package gobot

import (
	"testing"
	"time"

	"gobot.io/x/gobot/gobottest"
)

func TestEventerAddEvent(t *testing.T) {
	e := NewEventer()
	e.AddEvent("test")

	if _, ok := e.Events()["test"]; !ok {
		t.Errorf("Could not add event to list of Event names")
	}
	gobottest.Assert(t, e.Event("test"), "test")
}

func TestEventerDeleteEvent(t *testing.T) {
	e := NewEventer()
	e.AddEvent("test1")
	e.DeleteEvent("test1")

	if _, ok := e.Events()["test1"]; ok {
		t.Errorf("Could not add delete event from list of Event names")
	}
}

func TestEventerOn(t *testing.T) {
	e := NewEventer()
	e.AddEvent("test")

	sem := make(chan bool)
	e.On("test", func(data interface{}) {
		sem <- true
	})

	go func() {
		e.Publish("test", true)
	}()

	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("On was not called")
	}
}

func TestEventerOnce(t *testing.T) {
	e := NewEventer()
	e.AddEvent("test")

	sem := make(chan bool)
	e.Once("test", func(data interface{}) {
		sem <- true
	})

	go func() {
		e.Publish("test", true)
	}()

	select {
	case <-sem:
	case <-time.After(10 * time.Millisecond):
		t.Errorf("Once was not called")
	}

	go func() {
		e.Publish("test", true)
	}()

	select {
	case <-sem:
		t.Errorf("Once was called twice")
	case <-time.After(10 * time.Millisecond):
	}
}
