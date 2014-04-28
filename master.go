package gobot

import (
	"os"
	"os/signal"
	"runtime"
)

type Master struct {
	Robots []*Robot
	NumCPU int
	Api    *api
	trap   func(chan os.Signal)
}

// used to be GobotMaster()
func NewMaster() *Master {
	return &Master{
		NumCPU: runtime.NumCPU(),
		trap: func(c chan os.Signal) {
			signal.Notify(c, os.Interrupt)
		},
	}
}

func (m *Master) Start() {
	// this changes the amount of cores used by the program
	// to match the amount of CPUs set on master.
	runtime.GOMAXPROCS(m.NumCPU)

	if m.Api != nil {
		m.Api.start()
	}

	for _, r := range m.Robots {
		r.startRobot()
	}

	var c = make(chan os.Signal, 1)
	m.trap(c)

	// waiting on something coming on the channel
	_ = <-c
	for _, r := range m.Robots {
		r.GetDevices().Halt()
		r.finalizeConnections()
	}

}

func (m *Master) FindRobot(name string) *Robot {
	for _, robot := range m.Robots {
		if robot.Name == name {
			return robot
		}
	}
	return nil
}
