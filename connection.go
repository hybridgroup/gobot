package gobot

import (
	"errors"
	"log"
)

type JSONConnection struct {
	Name    string `json:"name"`
	Adaptor string `json:"adaptor"`
}

type Connection AdaptorInterface

type connections []Connection

func (c *connections) Len() int {
	return len(*c)
}

func (c *connections) Each(f func(Connection)) {
	for _, connection := range *c {
		f(connection)
	}
}

// Start() starts all the connections.
func (c *connections) Start() error {
	var err error
	log.Println("Starting connections...")
	for _, connection := range *c {
		info := "Starting connection " + connection.Name()
		if connection.Port() != "" {
			info = info + " on port " + connection.Port()
		}
		log.Println(info + "...")
		if connection.Connect() == false {
			err = errors.New("Could not start connection")
			break
		}
	}
	return err
}

// Finalize() finalizes all the connections.
func (c *connections) Finalize() {
	for _, connection := range *c {
		connection.Finalize()
	}
}
