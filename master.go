package gobot

import (
	"os"
	"os/signal"
	"sync/atomic"

	multierror "github.com/hashicorp/go-multierror"
)

// JSONMaster is a JSON representation of a Gobot Master.
type JSONMaster struct {
	Robots   []*JSONRobot `json:"robots"`
	Commands []string     `json:"commands"`
}

// NewJSONMaster returns a JSONMaster given a Gobot Master.
func NewJSONMaster(gobot *Master) *JSONMaster {
	jsonGobot := &JSONMaster{
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

// Master is the main type of your Gobot application and contains a collection of
// Robots, API commands that apply to the Master, and Events that apply to the Master.
type Master struct {
	robots  *Robots
	trap    func(chan os.Signal)
	AutoRun bool
	running atomic.Value
	Commander
	Eventer
}

// NewMaster returns a new Gobot Master
func NewMaster() *Master {
	m := &Master{
		robots: &Robots{},
		trap: func(c chan os.Signal) {
			signal.Notify(c, os.Interrupt)
		},
		AutoRun:   true,
		Commander: NewCommander(),
		Eventer:   NewEventer(),
	}
	m.running.Store(false)
	return m
}

// Start calls the Start method on each robot in its collection of robots. On
// error, call Stop to ensure that all robots are returned to a sane, stopped
// state.
func (g *Master) Start() (err error) {
	if rerr := g.robots.Start(!g.AutoRun); rerr != nil {
		err = multierror.Append(err, rerr)
		return
	}

	g.running.Store(true)

	if g.AutoRun {
		c := make(chan os.Signal, 1)
		g.trap(c)

		// waiting for interrupt coming on the channel
		<-c

		// Stop calls the Stop method on each robot in its collection of robots.
		g.Stop()
	}

	return err
}

// Stop calls the Stop method on each robot in its collection of robots.
func (g *Master) Stop() (err error) {
	if rerr := g.robots.Stop(); rerr != nil {
		err = multierror.Append(err, rerr)
	}

	g.running.Store(false)
	return
}

// Running returns if the Master is currently started or not
func (g *Master) Running() bool {
	return g.running.Load().(bool)
}

// Robots returns all robots associated with this Gobot Master.
func (g *Master) Robots() *Robots {
	return g.robots
}

// AddRobot adds a new robot to the internal collection of robots. Returns the
// added robot
func (g *Master) AddRobot(r *Robot) *Robot {
	*g.robots = append(*g.robots, r)
	return r
}

// Robot returns a robot given name. Returns nil if the Robot does not exist.
func (g *Master) Robot(name string) *Robot {
	for _, robot := range *g.Robots() {
		if robot.Name == name {
			return robot
		}
	}
	return nil
}
