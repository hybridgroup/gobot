package gobot

type commander struct {
	commands map[string]func(map[string]interface{}) interface{}
}

type Commander interface {
	Command(string) (command func(map[string]interface{}) interface{})
	Commands() (commands map[string]func(map[string]interface{}) interface{})
	AddCommand(name string, command func(map[string]interface{}) interface{})
}

func NewCommander() Commander {
	return &commander{
		commands: make(map[string]func(map[string]interface{}) interface{}),
	}
}

// Command retrieves a command by name
func (c *commander) Command(name string) (command func(map[string]interface{}) interface{}) {
	command, _ = c.commands[name]
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
