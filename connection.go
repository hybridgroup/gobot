package gobot

import (
	"errors"
	"log"
)

type JSONConnection struct {
	Name    string `json:"name"`
	Port    string `json:"port"`
	Adaptor string `json:"adaptor"`
}

type Connection AdaptorInterface

type connections struct {
	connections []Connection
}

func (c *connections) Len() int {
	return len(c.connections)
}

func (c *connections) Add(a Connection) Connection {
	c.connections = append(c.connections, a)
	return a
}

func (c *connections) Each(f func(Connection)) {
	for _, connection := range c.connections {
		f(connection)
	}
}

// Start() starts all the connections.
func (c connections) Start() error {
	var err error
	log.Println("Starting connections...")
	for _, connection := range c.connections {
		info := "Starting connection " + connection.Name()
		if connection.Port() != "" {
			info = info + " on Port " + connection.Port()
		}
		log.Println(info + "...")
		if connection.Connect() == false {
			err = errors.New("Could not start connection")
			break
		}
	}
	return err
}

// Filanize() finalizes all the connections.
func (c connections) Finalize() {
	for _, connection := range c.connections {
		connection.Finalize()
	}
}
