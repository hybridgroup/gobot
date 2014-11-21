package gobot

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

// JSONConnection holds a JSON representation of a connection.
type JSONConnection struct {
	Name    string `json:"name"`
	Adaptor string `json:"adaptor"`
}

// ToJSON returns a json representation of an adaptor
func NewJSONConnection(connection Connection) *JSONConnection {
	return &JSONConnection{
		Name:    connection.Name(),
		Adaptor: reflect.TypeOf(connection).String(),
	}
}

type Connection Adaptor

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
func (c *connections) Start() (errs []error) {
	log.Println("Starting connections...")
	for _, connection := range *c {
		info := "Starting connection " + connection.Name()
		if connection.Port() != "" {
			info = info + " on port " + connection.Port()
		}
		log.Println(info + "...")
		if errs = connection.Connect(); len(errs) > 0 {
			for i, err := range errs {
				errs[i] = errors.New(fmt.Sprintf("Connection %q: %v", connection.Name(), err))
			}
			return
		}
	}
	return
}

// Finalize finishes all the connections.
func (c *connections) Finalize() (errs []error) {
	for _, connection := range *c {
		if cerrs := connection.Finalize(); cerrs != nil {
			for i, err := range cerrs {
				cerrs[i] = errors.New(fmt.Sprintf("Connection %q: %v", connection.Name(), err))
			}
			errs = append(errs, cerrs...)
		}
	}
	return errs
}
