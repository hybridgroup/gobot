package gobot

import (
	"testing"

	"github.com/hybridgroup/gobot/gobottest"
)

func TestEventer(t *testing.T) {
	e := NewEventer()
	e.AddEvent("test")

	if _, ok := e.Events()["test"]; !ok {
		t.Errorf("Could not add event to list of Events")
	}

	event := e.Event("test")
	gobottest.Refute(t, event, nil)

	event = e.Event("booyeah")
	gobottest.Assert(t, event, (*Event)(nil))
}
