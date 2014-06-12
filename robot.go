package gobot

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type JSONRobot struct {
	Name        string            `json:"name"`
	Commands    []string          `json:"commands"`
	Connections []*JSONConnection `json:"connections"`
	Devices     []*JSONDevice     `json:"devices"`
}

type Robot struct {
	Name        string                                              `json:"-"`
	Commands    map[string]func(map[string]interface{}) interface{} `json:"-"`
	Work        func()                                              `json:"-"`
	connections connections                                         `json:"-"`
	devices     devices                                             `json:"-"`
}

type Robots []*Robot

func (r Robots) Start() {
	for _, robot := range r {
		robot.Start()
	}
}

func (r Robots) Each(f func(*Robot)) {
	for _, robot := range r {
		f(robot)
	}
}

func NewRobot(name string, c []Connection, d []Device, work func()) *Robot {
	r := &Robot{
		Name:     name,
		Work:     work,
		Commands: make(map[string]func(map[string]interface{}) interface{}),
	}
	r.initName()
	log.Println("Initializing Robot", r.Name, "...")
	r.initConnections(c)
	r.initDevices(d)
	return r
}

func (r *Robot) AddCommand(name string, f func(map[string]interface{}) interface{}) {
	r.Commands[name] = f
}

func (r *Robot) Start() {
	log.Println("Starting Robot", r.Name, "...")
	if err := r.Connections().Start(); err != nil {
		panic("Could not start connections")
	}
	if err := r.Devices().Start(); err != nil {
		panic("Could not start devices")
	}
	if r.Work != nil {
		log.Println("Starting work...")
		r.Work()
	}
}

func (r *Robot) initName() {
	if r.Name == "" {
		rand.Seed(time.Now().UTC().UnixNano())
		i := rand.Int()
		r.Name = fmt.Sprintf("Robot%v", i)
	}
}

func (r *Robot) initConnections(c []Connection) {
	r.connections = make(connections, len(c))
	log.Println("Initializing connections...")
	for i, connection := range c {
		log.Println("Initializing connection", FieldByNamePtr(connection, "Name"), "...")
		r.connections[i] = NewConnection(connection, r)
	}
}

func (r *Robot) initDevices(d []Device) {
	r.devices = make([]*device, len(d))
	log.Println("Initializing devices...")
	for i, device := range d {
		log.Println("Initializing device", FieldByNamePtr(device, "Name"), "...")
		r.devices[i] = NewDevice(device, r)
	}
}

func (r *Robot) Devices() devices {
	return devices(r.devices)
}

func (r *Robot) Device(name string) *device {
	if r == nil {
		return nil
	}
	for _, device := range r.devices {
		if device.Name == name {
			return device
		}
	}
	return nil
}

func (r *Robot) Connections() connections {
	return connections(r.connections)
}

func (r *Robot) Connection(name string) *connection {
	if r == nil {
		return nil
	}
	for _, connection := range r.connections {
		if connection.Name == name {
			return connection
		}
	}
	return nil
}

func (r *Robot) ToJSON() *JSONRobot {
	jsonRobot := &JSONRobot{
		Name:        r.Name,
		Commands:    []string{},
		Connections: []*JSONConnection{},
		Devices:     []*JSONDevice{},
	}
	for command := range r.Commands {
		jsonRobot.Commands = append(jsonRobot.Commands, command)
	}
	for _, device := range r.Devices() {
		jsonDevice := device.ToJSON()
		jsonRobot.Connections = append(jsonRobot.Connections, jsonDevice.Connection)
		jsonRobot.Devices = append(jsonRobot.Devices, jsonDevice)
	}
	return jsonRobot
}
