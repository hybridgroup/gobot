package gobot

import (
	"log"
	"os"
	"os/signal"
)

// JSONGobot is a JSON representation of a Gobot.
type JSONGobot struct {
	Robots   []*JSONRobot `json:"robots"`
	Commands []string     `json:"commands"`
}

// NewJSONGobot returns a JSONGobt given a Gobot.
func NewJSONGobot(gobot *Gobot) *JSONGobot {
	jsonGobot := &JSONGobot{
		Robots:   []*JSONRobot{},
		Commands: []string{},
	}

	for command := range gobot.Commands() {
		jsonGobot.Commands = append(jsonGobot.Commands, command)
	}

	gobot.robots.Each(func(r *Robot) {
		jsonGobot.Robots = append(jsonGobot.Robots, NewJSONRobot(r))
	})
	return jsonGobot
}

// Gobot is the main type of your Gobot application and contains a collection of
// Robots, API commands and Events.
type Gobot struct {
	robots   *Robots
	trap     func(chan os.Signal)
	AutoStop bool
	Commander
	Eventer
}

// NewGobot returns a new Gobot
func NewGobot() *Gobot {
	return &Gobot{
		robots: &Robots{},
		trap: func(c chan os.Signal) {
			signal.Notify(c, os.Interrupt)
		},
		AutoStop:  true,
		Commander: NewCommander(),
		Eventer:   NewEventer(),
	}
}

// Start calls the Start method on each robot in its collection of robots. On
// error, call Stop to ensure that all robots are returned to a sane, stopped
// state.
func (g *Gobot) Start() (errs []error) {
	if rerrs := g.robots.Start(); len(rerrs) > 0 {
		for _, err := range rerrs {
			log.Println("Error:", err)
			errs = append(errs, err)
		}
	}

	if g.AutoStop {
		c := make(chan os.Signal, 1)
		g.trap(c)
		if len(errs) > 0 {
			// there was an error during start, so we immediatly pass the interrupt
			// in order to disconnect the initialized robots, connections and devices
			c <- os.Interrupt
		}

		// waiting for interrupt coming on the channel
		_ = <-c

		// Stop calls the Stop method on each robot in its collection of robots.
		g.Stop()
	}

	return errs
}

// Stop calls the Stop method on each robot in its collection of robots.
func (g *Gobot) Stop() (errs []error) {
	if rerrs := g.robots.Stop(); len(rerrs) > 0 {
		for _, err := range rerrs {
			log.Println("Error:", err)
			errs = append(errs, err)
		}
	}

	return errs
}

// Robots returns all robots associated with this Gobot.
func (g *Gobot) Robots() *Robots {
	return g.robots
}

// AddRobot adds a new robot to the internal collection of robots. Returns the
// added robot
func (g *Gobot) AddRobot(r *Robot) *Robot {
	*g.robots = append(*g.robots, r)
	return r
}

// Robot returns a robot given name. Returns nil if the Robot does not exist.
func (g *Gobot) Robot(name string) *Robot {
	for _, robot := range *g.Robots() {
		if robot.Name == name {
			return robot
		}
	}
	return nil
}
