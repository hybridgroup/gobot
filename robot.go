package gobot

import (
	"fmt"
	"log"
)

// JSONRobot a JSON representation of a Robot.
type JSONRobot struct {
	Name        string            `json:"name"`
	Commands    []string          `json:"commands"`
	Connections []*JSONConnection `json:"connections"`
	Devices     []*JSONDevice     `json:"devices"`
}

// NewJSONRobot returns a JSONRobot given a Robot.
func NewJSONRobot(robot *Robot) *JSONRobot {
	jsonRobot := &JSONRobot{
		Name:        robot.Name,
		Commands:    []string{},
		Connections: []*JSONConnection{},
		Devices:     []*JSONDevice{},
	}

	for command := range robot.Commands() {
		jsonRobot.Commands = append(jsonRobot.Commands, command)
	}

	robot.Devices().Each(func(device Device) {
		jsonDevice := NewJSONDevice(device)
		jsonRobot.Connections = append(jsonRobot.Connections, NewJSONConnection(robot.Connection(jsonDevice.Connection)))
		jsonRobot.Devices = append(jsonRobot.Devices, jsonDevice)
	})
	return jsonRobot
}

// Robot is a named entitity that manages a collection of connections and devices.
// It containes it's own work routine and a collection of
// custom commands to control a robot remotely via the Gobot api.
type Robot struct {
	Name        string
	Work        func()
	connections *Connections
	devices     *Devices
	Commander
	Eventer
}

// Robots is a collection of Robot
type Robots []*Robot

// Len returns the amount of Robots in the collection.
func (r *Robots) Len() int {
	return len(*r)
}

// Start calls the Start method of each Robot in the collection
func (r *Robots) Start() (errs []error) {
	for _, robot := range *r {
		if errs = robot.Start(); len(errs) > 0 {
			for i, err := range errs {
				errs[i] = fmt.Errorf("Robot %q: %v", robot.Name, err)
			}
			return
		}
	}
	return
}

// Stop calls the Stop method of each Robot in the collection
func (r *Robots) Stop() (errs []error) {
	for _, robot := range *r {
		if errs = robot.Stop(); len(errs) > 0 {
			for i, err := range errs {
				errs[i] = fmt.Errorf("Robot %q: %v", robot.Name, err)
			}
			return
		}
	}
	return
}

// Each enumerates through the Robots and calls specified callback function.
func (r *Robots) Each(f func(*Robot)) {
	for _, robot := range *r {
		f(robot)
	}
}

// NewRobot returns a new Robot given a name and optionally accepts:
//
// 	[]Connection: Connections which are automatically started and stopped with the robot
//	[]Device: Devices which are automatically started and stopped with the robot
//	func(): The work routine the robot will execute once all devices and connections have been initialized and started
// A name will be automaically generated if no name is supplied.
func NewRobot(name string, v ...interface{}) *Robot {
	if name == "" {
		name = fmt.Sprintf("%X", Rand(int(^uint(0)>>1)))
	}

	r := &Robot{
		Name:        name,
		connections: &Connections{},
		devices:     &Devices{},
		Work:        nil,
		Eventer:     NewEventer(),
		Commander:   NewCommander(),
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
		}
	}

	return r
}

// Start a Robot's Connections, Devices, and work.
func (r *Robot) Start() (errs []error) {
	log.Println("Starting Robot", r.Name, "...")
	if cerrs := r.Connections().Start(); len(cerrs) > 0 {
		errs = append(errs, cerrs...)
		return
	}
	if derrs := r.Devices().Start(); len(derrs) > 0 {
		errs = append(errs, derrs...)
		return
	}
	if r.Work != nil {
		log.Println("Starting work...")
		r.Work()
	}
	return
}

// Stop stops a Robot's connections and Devices
func (r *Robot) Stop() (errs []error) {
	log.Println("Stopping Robot", r.Name, "...")
	if heers := r.Devices().Halt(); len(heers) > 0 {
		for _, err := range heers {
			errs = append(errs, err)
		}
	}

	if ceers := r.Connections().Finalize(); len(ceers) > 0 {
		for _, err := range ceers {
			errs = append(errs, err)
		}
	}

	return errs
}

// Devices returns all devices associated with this Robot.
func (r *Robot) Devices() *Devices {
	return r.devices
}

// AddDevice adds a new Device to the robots collection of devices. Returns the
// added device.
func (r *Robot) AddDevice(d Device) Device {
	*r.devices = append(*r.Devices(), d)
	return d
}

// Device returns a device given a name. Returns nil if the Device does not exist.
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
func (r *Robot) Connections() *Connections {
	return r.connections
}

// AddConnection adds a new connection to the robots collection of connections.
// Returns the added connection.
func (r *Robot) AddConnection(c Connection) Connection {
	*r.connections = append(*r.Connections(), c)
	return c
}

// Connection returns a connection given a name. Returns nil if the Connection
// does not exist.
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
