package gobot

import (
	"fmt"
	"log"
	"reflect"
)

// JSONConnection is a JSON representation of a Connection.
type JSONConnection struct {
	Name    string `json:"name"`
	Adaptor string `json:"adaptor"`
}

// NewJSONConnection returns a JSONConnection given a Connection.
func NewJSONConnection(connection Connection) *JSONConnection {
	return &JSONConnection{
		Name:    connection.Name(),
		Adaptor: reflect.TypeOf(connection).String(),
	}
}

// A Connection is an instance of an Adaptor
type Connection Adaptor

// Connections represents a collection of Connection
type Connections []Connection

// Len returns connections length
func (c *Connections) Len() int {
	return len(*c)
}

// Each enumerates through the Connections and calls specified callback function.
func (c *Connections) Each(f func(Connection)) {
	for _, connection := range *c {
		f(connection)
	}
}

// Start calls Connect on each Connection in c
func (c *Connections) Start() (errs []error) {
	log.Println("Starting connections...")
	for _, connection := range *c {
		info := "Starting connection " + connection.Name()

		if porter, ok := connection.(Porter); ok {
			info = info + " on port " + porter.Port()
		}

		log.Println(info + "...")

		if errs = connection.Connect(); len(errs) > 0 {
			for i, err := range errs {
				errs[i] = fmt.Errorf("Connection %q: %v", connection.Name(), err)
			}
			return
		}
	}
	return
}

// Finalize calls Finalize on each Connection in c
func (c *Connections) Finalize() (errs []error) {
	for _, connection := range *c {
		if cerrs := connection.Finalize(); cerrs != nil {
			for i, err := range cerrs {
				cerrs[i] = fmt.Errorf("Connection %q: %v", connection.Name(), err)
			}
			errs = append(errs, cerrs...)
		}
	}
	return errs
}
