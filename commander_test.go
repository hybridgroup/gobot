package gobot

import "testing"

func TestCommaner(t *testing.T) {
	c := NewCommander()
	c.AddCommand("test", func(map[string]interface{}) interface{} {
		return "hi"
	})

	if _, ok := c.Commands()["test"]; !ok {
		t.Errorf("Could not add command to list of Commands")
	}

	command := c.Command("test")
	Refute(t, command, nil)

	command = c.Command("booyeah")
	Assert(t, command, (func(map[string]interface{}) interface{})(nil))
}
