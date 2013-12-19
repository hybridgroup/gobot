package gobot

import (
	"fmt"
	"math/rand"
	"time"
)

type Robot struct {
	Connections   []Connection
	Devices       []Device
	Name          string
	Commands      map[string]interface{} `json:"-"`
	RobotCommands []string               `json:"Commands"`
	Work          func()                 `json:"-"`
	connections   []*connection          `json:"-"`
	devices       []*device              `json:"-"`
}

func (r *Robot) Start() {
	m := GobotMaster()
	m.Robots = []Robot{*r}
	m.Start()
}

func (r *Robot) startRobot() {
	r.initName()
	r.initCommands()
	r.initConnections()
	r.initDevices()
	r.startConnections()
	r.startDevices()
	if r.Work != nil {
		r.Work()
	}
}

func (r *Robot) initName() {
	if r.Name == "" {
		rand.Seed(time.Now().UTC().UnixNano())
		i := rand.Int()
		r.Name = fmt.Sprintf("Robot %v", i)
	}
}

func (r *Robot) initCommands() {
	for k, _ := range r.Commands {
		r.RobotCommands = append(r.RobotCommands, k)
	}
}

func (r *Robot) initConnections() {
	r.connections = make([]*connection, len(r.Connections))
	fmt.Println("Initializing connections...")
	for i := range r.Connections {
		fmt.Sprintln("Initializing connection %v...", FieldByNamePtr(r.Connections[i], "Name"))
		r.connections[i] = NewConnection(r.Connections[i], r)
	}
}

func (r *Robot) initDevices() {
	r.devices = make([]*device, len(r.Devices))
	fmt.Println("Initializing devices...")
	for i := range r.Devices {
		fmt.Sprintln("Initializing device %v...", FieldByNamePtr(r.Devices[i], "Name"))
		r.devices[i] = NewDevice(r.Devices[i], r)
	}
}

func (r *Robot) startConnections() {
	fmt.Println("Starting connections...")
	for i := range r.connections {
		fmt.Println("Starting connection " + r.connections[i].Name + "...")
		r.connections[i].Connect()
	}
}

func (r *Robot) startDevices() {
	fmt.Println("Starting devices...")
	for i := range r.devices {
		fmt.Println("Starting device " + r.devices[i].Name + "...")
		r.devices[i].Start()
	}
}

func (r *Robot) GetDevices() []*device {
	return r.devices
}

func (r *Robot) GetDevice(name string) *device {
	for i := range r.devices {
		if r.devices[i].Name == name {
			return r.devices[i]
		}
	}
	return nil
}
