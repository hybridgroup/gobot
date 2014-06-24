package gobot

type Command struct {
	Name    string
	Command func(map[string]interface{}) interface{}
}
type Commands []Command

func (c commands) Add(name string, cmd Command) {
	c[name] = cmd
}

func (c commands) Each(f func(Command)) {
	for _, command := range c {
		f(command)
	}
}
