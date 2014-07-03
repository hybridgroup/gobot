package gobot

import (
	"fmt"
	"log"
)

type JSONRobot struct {
	Name        string            `json:"name"`
	Commands    []string          `json:"commands"`
	Connections []*JSONConnection `json:"connections"`
	Devices     []*JSONDevice     `json:"devices"`
}

type Robot struct {
	Name        string
	Commands    map[string]func(map[string]interface{}) interface{}
	Work        func()
	connections *connections
	devices     *devices
}

type robots struct {
	robots []*Robot
}

func (r *robots) Len() int {
	return len(r.robots)
}
func (r *robots) Add(robot *Robot) *Robot {
	r.robots = append(r.robots, robot)
	return robot
}

func (r *robots) Start() {
	for _, robot := range r.robots {
		robot.Start()
	}
}

func (r *robots) Each(f func(*Robot)) {
	for _, robot := range r.robots {
		f(robot)
	}
}

func NewRobot(name string, v ...interface{}) *Robot {
	if name == "" {
		name = fmt.Sprintf("%X", Rand(int(^uint(0)>>1)))
	}
	r := &Robot{
		Name:        name,
		Commands:    make(map[string]func(map[string]interface{}) interface{}),
		connections: &connections{},
		devices:     &devices{},
	}

	log.Println("Initializing Robot", r.Name, "...")
	if len(v) > 0 {
		if v[0] == nil {
			v[0] = connections{}
		}
		log.Println("Initializing connections...")
		for _, connection := range v[0].([]Connection) {
			c := r.Connections().Add(connection)
			log.Println("Initializing connection", c.name(), "...")
		}
	}
	if len(v) > 1 {
		if v[1] == nil {
			v[1] = devices{}
		}
		log.Println("Initializing devices...")
		for _, device := range v[1].([]Device) {
			d := r.Devices().Add(device)
			log.Println("Initializing device", d.name(), "...")
		}
	}
	if len(v) > 2 {
		if v[2] == nil {
			v[2] = func() {}
		}
		r.Work = v[2].(func())
	}
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

func (r *Robot) Devices() *devices {
	return r.devices
}

func (r *Robot) Device(name string) Device {
	if r == nil {
		return nil
	}
	for _, device := range r.devices.devices {
		if device.name() == name {
			return device
		}
	}
	return nil
}

func (r *Robot) Connections() *connections {
	return r.connections
}

func (r *Robot) Connection(name string) Connection {
	if r == nil {
		return nil
	}
	for _, connection := range r.connections.connections {
		if connection.name() == name {
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

	r.Devices().Each(func(device Device) {
		jsonDevice := device.ToJSON()
		jsonRobot.Connections = append(jsonRobot.Connections, jsonDevice.Connection)
		jsonRobot.Devices = append(jsonRobot.Devices, jsonDevice)
	})
	return jsonRobot
}
