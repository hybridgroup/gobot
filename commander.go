package gobot

import "errors"

type commander struct {
	commands map[string]func(map[string]interface{}) interface{}
}

type Commander interface {
	Command(string) (command func(map[string]interface{}) interface{}, err error)
	Commands() (commands map[string]func(map[string]interface{}) interface{})
	AddCommand(name string, command func(map[string]interface{}) interface{})
}

func NewCommander() Commander {
	return &commander{
		commands: make(map[string]func(map[string]interface{}) interface{}),
	}
}

// Command retrieves a command by name
func (c *commander) Command(name string) (command func(map[string]interface{}) interface{}, err error) {
	command, ok := c.commands[name]
	if ok {
		return
	}
	err = errors.New("Unknown Command")
	return
}

// Commands returns a map of driver commands
func (c *commander) Commands() map[string]func(map[string]interface{}) interface{} {
	return c.commands
}

// AddCommand links specified command name to `f`
func (c *commander) AddCommand(name string, command func(map[string]interface{}) interface{}) {
	c.commands[name] = command
}
