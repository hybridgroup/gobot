package gobot

import (
	"fmt"
	"log"
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
	if r.startConnections() != true {
		panic("Could not start connections")
	}
	if r.startDevices() != true {
		panic("Could not start devices")
	}
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
	log.Println("Initializing connections...")
	for i := range r.Connections {
		log.Println("Initializing connection ", FieldByNamePtr(r.Connections[i], "Name"), "...")
		r.connections[i] = NewConnection(r.Connections[i], r)
	}
}

func (r *Robot) initDevices() {
	r.devices = make([]*device, len(r.Devices))
	log.Println("Initializing devices...")
	for i := range r.Devices {
		log.Println("Initializing device ", FieldByNamePtr(r.Devices[i], "Name"), "...")
		r.devices[i] = NewDevice(r.Devices[i], r)
	}
}

func (r *Robot) startConnections() bool {
	log.Println("Starting connections...")
	success := true
	for i := range r.connections {
		log.Println("Starting connection " + r.connections[i].Name + "...")
		if r.connections[i].Connect() == false {
			success = false
			break
		}
	}
	return success
}

func (r *Robot) startDevices() bool {
	log.Println("Starting devices...")
	success := true
	for i := range r.devices {
		log.Println("Starting device " + r.devices[i].Name + "...")
		if r.devices[i].Start() == false {
			success = false
			break
		}
	}
	return success
}

func (r *Robot) finalizeConnections() {
	for i := range r.connections {
		r.connections[i].Finalize()
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
