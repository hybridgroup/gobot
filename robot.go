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

// Robot is a named entitity that manages a collection of connections and devices.
// It containes it's own work routine and a collection of
// custom commands to control a robot remotely via the Gobot api.
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

// NewRobot returns a new Robot given a name and optionally accepts:
//
// 	[]Connection: Connections which are automatically started and stopped with the robot
//	[]Device: Devices which are automatically started and stopped with the robot
//	func(): The work routine the robot will execute once all devices and connections have been initialized and started
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

// AddCommand adds a new command to the robot's collection of commands
func (r *Robot) AddCommand(name string, f func(map[string]interface{}) interface{}) {
	r.commands[name] = f
}

// Commands returns all available commands on the robot.
func (r *Robot) Commands() map[string]func(map[string]interface{}) interface{} {
	return r.commands
}

// Command returns the command given a name.
func (r *Robot) Command(name string) func(map[string]interface{}) interface{} {
	return r.commands[name]
}

// Start starts all the robot's connections and drivers and runs it's work function.
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

// Devices returns all devices associated with this robot.
func (r *Robot) Devices() *devices {
	return r.devices
}

// AddDevice adds a new device to the robots collection of devices. Returns the
// added device.
func (r *Robot) AddDevice(d Device) Device {
	*r.devices = append(*r.Devices(), d)
	return d
}

// Device returns a device given a name. Returns nil on no device.
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

// Connections returns all connections associated with this robot.
func (r *Robot) Connections() *connections {
	return r.connections
}

// AddConnection adds a new connection to the robots collection of connections.
// Returns the added connection.
func (r *Robot) AddConnection(c Connection) Connection {
	*r.connections = append(*r.Connections(), c)
	return c
}

// Connection returns a connection given a name. Returns nil on no connection.
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

// ToJSON returns a JSON representation of the robot.
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
