package gobot

import "testing"

func TestEventer(t *testing.T) {
	e := NewEventer()
	e.AddEvent("test")

	if _, ok := e.Events()["test"]; !ok {
		t.Errorf("Could not add event to list of Events")
	}

	event := e.Event("test")
	Refute(t, event, nil)

	event = e.Event("booyeah")
	Assert(t, event, (*Event)(nil))
}
