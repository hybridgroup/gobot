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
	robots *robots
	trap   func(chan os.Signal)
	Commander
	Eventer
}

// NewGobot instantiates a new Gobot
func NewGobot() *Gobot {
	return &Gobot{
		robots: &robots{},
		trap: func(c chan os.Signal) {
			signal.Notify(c, os.Interrupt)
		},
		Commander: NewCommander(),
		Eventer:   NewEventer(),
	}
}

// Start runs the main Gobot event loop
func (g *Gobot) Start() (errs []error) {
	if rerrs := g.robots.Start(); len(rerrs) > 0 {
		for _, err := range rerrs {
			log.Println("Error:", err)
			errs = append(errs, err)
		}
	}

	c := make(chan os.Signal, 1)
	g.trap(c)
	if len(errs) > 0 {
		// there was an error during start, so we immediatly pass the interrupt
		// in order to disconnect the initialized robots, connections and devices
		c <- os.Interrupt
	}

	// waiting for interrupt coming on the channel
	_ = <-c
	g.robots.Each(func(r *Robot) {
		log.Println("Stopping Robot", r.Name, "...")
		if herrs := r.Devices().Halt(); len(herrs) > 0 {
			for _, err := range herrs {
				log.Println("Error:", err)
				errs = append(errs, err)
			}
		}
		if cerrs := r.Connections().Finalize(); len(cerrs) > 0 {
			for _, err := range cerrs {
				log.Println("Error:", err)
				errs = append(errs, err)
			}
		}
	})
	return errs
}

// Robots fetch all robots associated with this Gobot instance.
func (g *Gobot) Robots() *robots {
	return g.robots
}

// AddRobot adds a new robot to our Gobot instance.
func (g *Gobot) AddRobot(r *Robot) *Robot {
	*g.robots = append(*g.robots, r)
	return r
}

// Robot find a robot with a given name.
func (g *Gobot) Robot(name string) *Robot {
	for _, robot := range *g.Robots() {
		if robot.Name == name {
			return robot
		}
	}
	return nil
}

// ToJSON retrieves a JSON representation of this Gobot.
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
