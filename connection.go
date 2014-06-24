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

type connections struct {
	connections []AdaptorInterface
}

func (c *connections) Len() int {
	return len(c.connections)
}

func (c *connections) Add(a AdaptorInterface) AdaptorInterface {
	c.connections = append(c.connections, a)
	return a
}

func (c *connections) Each(f func(AdaptorInterface)) {
	for _, connection := range c.connections {
		f(connection)
	}
}

// Start() starts all the connections.
func (c connections) Start() error {
	var err error
	log.Println("Starting connections...")
	for _, connection := range c.connections {
		log.Println("Starting connection " + connection.name() + "...")
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
