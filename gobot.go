package gobot

import (
	"github.com/hybridgroup/gobot/core/robot"
	"os"
	"os/signal"
)

type Gobot struct {
	Robots robot.Robots
	trap   func(chan os.Signal)
}

func NewGobot() *Gobot {
	return &Gobot{
		trap: func(c chan os.Signal) {
			signal.Notify(c, os.Interrupt)
		},
	}
}

func (g *Gobot) Start() {
	g.Robots.Start()

	c := make(chan os.Signal, 1)
	g.trap(c)

	// waiting for interrupt coming on the channel
	_ = <-c
	g.Robots.Each(func(r *robot.Robot) {
		r.GetDevices().Halt()
		r.GetConnections().Finalize()
	})
}

func (g *Gobot) Robot(name string) *robot.Robot {
	for _, r := range g.Robots {
		if r.Name == name {
			return r
		}
	}
	return nil
}
