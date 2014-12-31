package gobot

type commander struct {
	commands map[string]func(map[string]interface{}) interface{}
}

// Commander is the interface which describes the behaviour for a Driver or Adaptor
// which exposes API commands.
type Commander interface {
	// Command returns a command given a name. Returns nil if the command is not found.
	Command(string) (command func(map[string]interface{}) interface{})
	// Commands returns a map of commands.
	Commands() (commands map[string]func(map[string]interface{}) interface{})
	// AddCommand adds a command given a name.
	AddCommand(name string, command func(map[string]interface{}) interface{})
}

// NewCommander returns a new Commander.
func NewCommander() Commander {
	return &commander{
		commands: make(map[string]func(map[string]interface{}) interface{}),
	}
}

func (c *commander) Command(name string) (command func(map[string]interface{}) interface{}) {
	command, _ = c.commands[name]
	return
}

func (c *commander) Commands() map[string]func(map[string]interface{}) interface{} {
	return c.commands
}

func (c *commander) AddCommand(name string, command func(map[string]interface{}) interface{}) {
	c.commands[name] = command
}
