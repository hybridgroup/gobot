package gobot

import (
	"log"
	"os"
	"os/signal"
)

type JSONGobot struct {
	Robots   []*JSONRobot `json:"robots"`
	Commands []string     `json:"commands"`
}

type Gobot struct {
	robots   *robots
	commands map[string]func(map[string]interface{}) interface{}
	trap     func(chan os.Signal)
}

func NewGobot() *Gobot {
	return &Gobot{
		robots:   &robots{},
		commands: make(map[string]func(map[string]interface{}) interface{}),
		trap: func(c chan os.Signal) {
			signal.Notify(c, os.Interrupt)
		},
	}
}

func (g *Gobot) Commands() map[string]func(map[string]interface{}) interface{} {
	return g.commands
}

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

func (g *Gobot) Robots() *robots {
	return g.robots
}

func (g *Gobot) Robot(name string) *Robot {
	for _, robot := range g.Robots().robots {
		if robot.Name == name {
			return robot
		}
	}
	return nil
}

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
