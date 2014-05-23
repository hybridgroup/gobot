package gobot

import (
	"log"
	"os"
	"os/signal"
)

type Gobot struct {
	Robots []*Robot
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
