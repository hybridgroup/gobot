package gobot

import (
	"fmt"
	"testing"
	"time"
)

func TestDriver(t *testing.T) {
	a := NewTestAdaptor("testAdaptor")
	d := NewDriver("",
		"testDriver",
		a,
		"1",
		5*time.Second,
	)

	Refute(t, d.Name(), "")
	Assert(t, d.Type(), "testDriver")
	Assert(t, d.Interval(), 5*time.Second)
	Assert(t, d.Pin(), "1")
	Assert(t, d.Adaptor(), a)

	d.SetPin("10")
	Assert(t, d.Pin(), "10")

	d.SetName("myDriver")
	Assert(t, d.Name(), "myDriver")

	d.SetInterval(100 * time.Second)
	Assert(t, d.Interval(), 100*time.Second)

	Assert(t, len(d.Commands()), 0)
	d.AddCommand("cmd1", func(params map[string]interface{}) interface{} {
		return fmt.Sprintf("hello from %v", params["name"])
	})
	Assert(t, len(d.Commands()), 1)
	Assert(t,
		d.Command("cmd1")(map[string]interface{}{"name": d.Name()}).(string),
		"hello from "+d.Name(),
	)

	Assert(t, len(d.Events()), 0)
	d.AddEvent("event1")
	Assert(t, len(d.Events()), 1)
	Refute(t, d.Event("event1"), nil)

	defer func() {
		r := recover()
		if r != nil {
			Assert(t, "Unknown Driver Event: event2", r)
		} else {
			t.Errorf("Did not return Unknown Event error")
		}
	}()
	d.Event("event2")
}
