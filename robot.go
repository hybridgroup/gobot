package gobot

import (
	"fmt"
	"log"
)

// JSONRobot a JSON representation of a robot.
type JSONRobot struct {
	Name        string            `json:"name"`
	Commands    []string          `json:"commands"`
	Connections []*JSONConnection `json:"connections"`
	Devices     []*JSONDevice     `json:"devices"`
}

// Robot software representation of a physical board. A robot is a named
// entitity that manages multiple IO devices using a set of adaptors. Additionally
// a user can specificy custom commands to control a robot remotely.
type Robot struct {
	Name        string
	commands    map[string]func(map[string]interface{}) interface{}
	Work        func()
	connections *connections
	devices     *devices
}

type robots []*Robot

// Len counts the robots associated with this instance.
func (r *robots) Len() int {
	return len(*r)
}

// Start initialises the event loop. All robots that were added will
// be automtically started as a result of this call.
func (r *robots) Start() {
	for _, robot := range *r {
		robot.Start()
	}
}

// Each enumerates thru the robts and calls specified function
func (r *robots) Each(f func(*Robot)) {
	for _, robot := range *r {
		f(robot)
	}
}

// NewRobot constructs a new named robot. Though a robot's name will be generated,
// we recommend that user take care of naming a robot for later access.
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

// AddCommand setup a new command that we be made available via the REST api.
func (r *Robot) AddCommand(name string, f func(map[string]interface{}) interface{}) {
	r.commands[name] = f
}

// Commands lists out all available commands on this robot.
func (r *Robot) Commands() map[string]func(map[string]interface{}) interface{} {
	return r.commands
}

// Command fetch a named command on this robot.
func (r *Robot) Command(name string) func(map[string]interface{}) interface{} {
	return r.commands[name]
}

// Start a robot instance and runs it's work function if any. You should not
// need to manually start a robot if already part of a Gobot application as the
// robot will be automatically started for you.
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

// Devices retrieves all devices associated with this robot.
func (r *Robot) Devices() *devices {
	return r.devices
}

// AddDevice adds a new device on this robot.
func (r *Robot) AddDevice(d Device) Device {
	*r.devices = append(*r.Devices(), d)
	return d
}

// Device finds a device by name.
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

// Connections retrieves all connections on this robot.
func (r *Robot) Connections() *connections {
	return r.connections
}

// AddConnection add a new connection on this robot.
func (r *Robot) AddConnection(c Connection) Connection {
	*r.connections = append(*r.Connections(), c)
	return c
}

// Connection finds a connection by name.
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

// ToJSON returns a JSON representation of the master robot.
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
		jsonRobot.Connections = append(jsonRobot.Connections, r.Connection(jsonDevice.Connection).ToJSON())
		jsonRobot.Devices = append(jsonRobot.Devices, jsonDevice)
	})
	return jsonRobot
}
