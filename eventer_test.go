package gobot

import (
	"testing"
)

func TestEventerAddEvent(t *testing.T) {
	e := NewEventer()
	e.AddEvent("test")

	if _, ok := e.Events()["test"]; !ok {
		t.Errorf("Could not add event to list of Event names")
	}
}

func TestEventerDeleteEvent(t *testing.T) {
	e := NewEventer()
	e.AddEvent("test1")
	e.DeleteEvent("test1")

	if _, ok := e.Events()["test1"]; ok {
		t.Errorf("Could not add delete event from list of Event names")
	}
}
