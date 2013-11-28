package gobot

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

type Robot struct {
	Connections   []interface{}
	Devices       []interface{}
	Name          string
	Commands      map[string]interface{} `json:"-"`
	RobotCommands []string
	Work          func()        `json:"-"`
	connections   []*Connection `json:"-"`
	devices       []*Device     `json:"-"`
}

func (r *Robot) Start() {
	if r.Name == "" {
		rand.Seed(time.Now().UTC().UnixNano())
		i := rand.Int()
		r.Name = fmt.Sprintf("Robot %v", i)
	}
	for k, _ := range r.Commands {
		r.RobotCommands = append(r.RobotCommands, k)
	}
	r.initConnections()
	r.initDevices()
	r.startConnections()
	r.startDevices()
	r.Work()
	select {}
}

func (r *Robot) initConnections() {
	r.connections = make([]*Connection, len(r.Connections))
	fmt.Println("Initializing connections...")
	for i := range r.Connections {
		fmt.Println("Initializing connection " + reflect.ValueOf(r.Connections[i]).Elem().FieldByName("Name").String() + "...")
		r.connections[i] = NewConnection(r.Connections[i], r)
	}
}

func (r *Robot) initDevices() {
	r.devices = make([]*Device, len(r.Devices))
	fmt.Println("Initializing devices...")
	for i := range r.Devices {
		fmt.Println("Initializing device " + reflect.ValueOf(r.Devices[i]).Elem().FieldByName("Name").String() + "...")
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

func (r *Robot) GetDevices() []*Device {
	return r.devices
}

func (r *Robot) GetDevice(name string) *Device {
	for i := range r.devices {
		if r.devices[i].Name == name {
			return r.devices[i]
		}
	}
	return nil
}
