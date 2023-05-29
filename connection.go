package gobot

import (
	"log"
	"reflect"

	multierror "github.com/hashicorp/go-multierror"
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
func (c *Connections) Start() (err error) {
	log.Println("Starting connections...")
	for _, connection := range *c {
		info := "Starting connection " + connection.Name()

		if porter, ok := connection.(Porter); ok {
			info = info + " on port " + porter.Port()
		}

		log.Println(info + "...")

		if cerr := connection.Connect(); cerr != nil {
			err = multierror.Append(err, cerr)
		}
	}
	return err
}

// Finalize calls Finalize on each Connection in c
func (c *Connections) Finalize() (err error) {
	for _, connection := range *c {
		if cerr := connection.Finalize(); cerr != nil {
			err = multierror.Append(err, cerr)
		}
	}
	return err
}
