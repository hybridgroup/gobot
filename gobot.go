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
	Robots   []*Robot
	Commands map[string]func(map[string]interface{}) interface{}
	trap     func(chan os.Signal)
}

func NewGobot() *Gobot {
	return &Gobot{
		Commands: make(map[string]func(map[string]interface{}) interface{}),
		trap: func(c chan os.Signal) {
			signal.Notify(c, os.Interrupt)
		},
	}
}

func (g *Gobot) AddRobot(r *Robot) *Robot {
	g.Robots = append(g.Robots, r)
	return r
}

func (g *Gobot) AddCommand(name string, f func(map[string]interface{}) interface{}) {
	g.Commands[name] = f
}

func (g *Gobot) Start() {
	Robots(g.Robots).Start()

	c := make(chan os.Signal, 1)
	g.trap(c)

	// waiting for interrupt coming on the channel
	_ = <-c
	Robots(g.Robots).Each(func(r *Robot) {
		log.Println("Stopping Robot", r.Name, "...")
		r.Devices().Halt()
		r.Connections().Finalize()
	})
}

func (g *Gobot) Robot(name string) *Robot {
	for _, r := range g.Robots {
		if r.Name == name {
			return r
		}
	}
	return nil
}

func (g *Gobot) ToJSON() *JSONGobot {
	jsonGobot := &JSONGobot{
		Robots:   []*JSONRobot{},
		Commands: []string{},
	}
	for command := range g.Commands {
		jsonGobot.Commands = append(jsonGobot.Commands, command)
	}
	for _, robot := range g.Robots {
		jsonGobot.Robots = append(jsonGobot.Robots, robot.ToJSON())
	}
	return jsonGobot
}
