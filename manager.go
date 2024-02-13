package gobot

import (
	"os"
	"os/signal"
	"sync/atomic"
)

// JSONManager is a JSON representation of a Gobot Manager.
type JSONManager struct {
	Robots   []*JSONRobot `json:"robots"`
	Commands []string     `json:"commands"`
}

// NewJSONManager returns a JSONManager given a Gobot Manager.
func NewJSONManager(gobot *Manager) *JSONManager {
	jsonGobot := &JSONManager{
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

// Manager is the main type of your Gobot application and contains a collection of
// Robots, API commands that apply to the Manager, and Events that apply to the Manager.
type Manager struct {
	robots  *Robots
	trap    func(chan os.Signal)
	AutoRun bool
	running atomic.Value
	Commander
	Eventer
}

// NewManager returns a new Gobot Manager
func NewManager() *Manager {
	m := &Manager{
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
func (g *Manager) Start() error {
	if err := g.robots.Start(!g.AutoRun); err != nil {
		return err
	}

	g.running.Store(true)

	if !g.AutoRun {
		return nil
	}

	c := make(chan os.Signal, 1)
	g.trap(c)

	// waiting for interrupt coming on the channel
	<-c

	// Stop calls the Stop method on each robot in its collection of robots.
	return g.Stop()
}

// Stop calls the Stop method on each robot in its collection of robots.
func (g *Manager) Stop() error {
	err := g.robots.Stop()
	g.running.Store(false)
	return err
}

// Running returns if the Manager is currently started or not
func (g *Manager) Running() bool {
	return g.running.Load().(bool) //nolint:forcetypeassert // no error return value, so there is no better way
}

// Robots returns all robots associated with this Gobot Manager.
func (g *Manager) Robots() *Robots {
	return g.robots
}

// AddRobot adds a new robot to the internal collection of robots. Returns the
// added robot
func (g *Manager) AddRobot(r *Robot) *Robot {
	*g.robots = append(*g.robots, r)
	return r
}

// Robot returns a robot given name. Returns nil if the Robot does not exist.
func (g *Manager) Robot(name string) *Robot {
	for _, robot := range *g.Robots() {
		if robot.Name == name {
			return robot
		}
	}
	return nil
}
