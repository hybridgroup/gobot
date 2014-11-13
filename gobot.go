package gobot

import (
	"log"
	"os"
	"os/signal"
)

// JSONGobot holds a JSON representation of a Gobot.
type JSONGobot struct {
	Robots   []*JSONRobot `json:"robots"`
	Commands []string     `json:"commands"`
}

// Gobot is a container composed of one or more robots
type Gobot struct {
	robots   *robots
	commands map[string]func(map[string]interface{}) interface{}
	trap     func(chan os.Signal)
}

// NewGobot returns a new Gobot
func NewGobot() *Gobot {
	return &Gobot{
		robots:   &robots{},
		commands: make(map[string]func(map[string]interface{}) interface{}),
		trap: func(c chan os.Signal) {
			signal.Notify(c, os.Interrupt)
		},
	}
}

/*
AddCommand creates a new command and adds it to the Gobot. This command
will be available via HTTP using '/commands/name'

Example:
	gbot.AddCommand( 'rollover', func( params map[string]interface{}) interface{} {
		fmt.Println( "Rolling over - Stand by...")
	})

	With the api package setup, you can now get your Gobot to rollover using: http://localhost:3000/commands/rollover
*/
func (g *Gobot) AddCommand(name string, f func(map[string]interface{}) interface{}) {
	g.commands[name] = f
}

// Commands returns all available commands on this Gobot.
func (g *Gobot) Commands() map[string]func(map[string]interface{}) interface{} {
	return g.commands
}

// Command reutnrs a command given a name.
func (g *Gobot) Command(name string) func(map[string]interface{}) interface{} {
	return g.commands[name]
}

// Start calls the Start method on each robot in it's collection of robots, and
// stops all robots on reception of a SIGINT. Start will block the execution of
// your main function until it receives the SIGINT.
func (g *Gobot) Start() {
	g.robots.Start()

	c := make(chan os.Signal, 1)
	g.trap(c)

	// waiting for interrupt coming on the channel
	_ = <-c
	g.robots.Each(func(r *Robot) {
		log.Println("Stopping Robot", r.Name, "...")
		r.Devices().Halt()
		r.Connections().Finalize()
	})
}

// Robots returns all robots associated with this Gobot instance.
func (g *Gobot) Robots() *robots {
	return g.robots
}

// AddRobot adds a new robot to the internal collection of robots. Returns the
// added robot
func (g *Gobot) AddRobot(r *Robot) *Robot {
	*g.robots = append(*g.robots, r)
	return r
}

// Robot returns a robot given name. Returns nil on no robot.
func (g *Gobot) Robot(name string) *Robot {
	for _, robot := range *g.Robots() {
		if robot.Name == name {
			return robot
		}
	}
	return nil
}

// ToJSON returns a JSON representation of this Gobot.
func (g *Gobot) ToJSON() *JSONGobot {
	jsonGobot := &JSONGobot{
		Robots:   []*JSONRobot{},
		Commands: []string{},
	}

	for command := range g.Commands() {
		jsonGobot.Commands = append(jsonGobot.Commands, command)
	}

	g.robots.Each(func(r *Robot) {
		jsonGobot.Robots = append(jsonGobot.Robots, r.ToJSON())
	})
	return jsonGobot
}
