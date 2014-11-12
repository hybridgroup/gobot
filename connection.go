package gobot

import (
	"errors"
	"fmt"
	"log"
)

// JSONConnection holds a JSON representation of a connection.
type JSONConnection struct {
	Name    string `json:"name"`
	Adaptor string `json:"adaptor"`
}

type Connection AdaptorInterface

type connections []Connection

// Len returns connections length
func (c *connections) Len() int {
	return len(*c)
}

// Each calls function for each connection
func (c *connections) Each(f func(Connection)) {
	for _, connection := range *c {
		f(connection)
	}
}

// Start initializes all the connections.
func (c *connections) Start() error {
	var err error
	log.Println("Starting connections...")
	for _, connection := range *c {
		info := "Starting connection " + connection.Name()
		if connection.Port() != "" {
			info = info + " on port " + connection.Port()
		}
		log.Println(info + "...")
		err = connection.Connect()
		if err != nil {
			err = errors.New(fmt.Sprintf("Could not start connection: %v", err))
			break
		}
	}
	return err
}

// Finalize finishes all the connections.
func (c *connections) Finalize() {
	for _, connection := range *c {
		connection.Finalize()
	}
}
