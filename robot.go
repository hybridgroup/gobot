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
	commands    map[string]func(map[string]interface{}) interface{}
	Work        func()
	connections *connections
	devices     *devices
}

type robots []*Robot

func (r *robots) Len() int {
	return len(*r)
}

func (r *robots) Start() {
	for _, robot := range *r {
		robot.Start()
	}
}

func (r *robots) Each(f func(*Robot)) {
	for _, robot := range *r {
		f(robot)
	}
}

func NewRobot(name string, v ...interface{}) *Robot {
	if name == "" {
		name = fmt.Sprintf("%X", Rand(int(^uint(0)>>1)))
	}

	r := &Robot{
		Name:        name,
		commands:    make(map[string]func(map[string]interface{}) interface{}),
		connections: &connections{},
		devices:     &devices{},
		Work:        nil,
	}

	log.Println("Initializing Robot", r.Name, "...")

	for i := range v {
		switch v[i].(type) {
		case []Connection:
			log.Println("Initializing connections...")
			for _, connection := range v[i].([]Connection) {
				c := r.AddConnection(connection)
				log.Println("Initializing connection", c.Name(), "...")
			}
		case []Device:
			log.Println("Initializing devices...")
			for _, device := range v[i].([]Device) {
				d := r.AddDevice(device)
				log.Println("Initializing device", d.Name(), "...")
			}
		case func():
			r.Work = v[i].(func())
		default:
			fmt.Println("Unknown argument passed to NewRobot")
		}
	}

	return r
}

func (r *Robot) AddCommand(name string, f func(map[string]interface{}) interface{}) {
	r.commands[name] = f
}

func (r *Robot) Commands() map[string]func(map[string]interface{}) interface{} {
	return r.commands
}

func (r *Robot) Command(name string) func(map[string]interface{}) interface{} {
	return r.commands[name]
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

func (r *Robot) AddDevice(d Device) Device {
	*r.devices = append(*r.Devices(), d)
	return d
}

func (r *Robot) Device(name string) Device {
	if r == nil {
		return nil
	}
	for _, device := range *r.devices {
		if device.Name() == name {
			return device
		}
	}
	return nil
}

func (r *Robot) Connections() *connections {
	return r.connections
}

func (r *Robot) AddConnection(c Connection) Connection {
	*r.connections = append(*r.Connections(), c)
	return c
}

func (r *Robot) Connection(name string) Connection {
	if r == nil {
		return nil
	}
	for _, connection := range *r.connections {
		if connection.Name() == name {
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

	for command := range r.Commands() {
		jsonRobot.Commands = append(jsonRobot.Commands, command)
	}

	r.Devices().Each(func(device Device) {
		jsonDevice := device.ToJSON()
		jsonRobot.Connections = append(jsonRobot.Connections, jsonDevice.Connection)
		jsonRobot.Devices = append(jsonRobot.Devices, jsonDevice)
	})
	return jsonRobot
}
