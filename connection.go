package gobot

import (
	"errors"
	"log"
	"reflect"
)

type Connection interface {
	Connect() bool
	Finalize() bool
	port() string
	name() string
	params() map[string]interface{}
}

type JSONConnection struct {
	Name    string `json:"name"`
	Port    string `json:"port"`
	Adaptor string `json:"adaptor"`
}

type connection struct {
	Name    string                 `json:"-"`
	Type    string                 `json:"-"`
	Adaptor AdaptorInterface       `json:"-"`
	Port    string                 `json:"-"`
	Robot   *Robot                 `json:"-"`
	Params  map[string]interface{} `json:"-"`
}

type connections []*connection

// Start() starts all the connections.
func (c connections) Start() error {
	var err error
	log.Println("Starting connections...")
	for _, connection := range c {
		log.Println("Starting connection " + connection.Name + "...")
		if connection.Connect() == false {
			err = errors.New("Could not start connection")
			break
		}
	}
	return err
}

// Filanize() finalizes all the connections.
func (c connections) Finalize() {
	for _, connection := range c {
		connection.Finalize()
	}
}

func NewConnection(adaptor AdaptorInterface, r *Robot) *connection {
	t := reflect.ValueOf(adaptor).Type().String()
	return &connection{
		Type:    t[1:len(t)],
		Name:    adaptor.name(),
		Port:    adaptor.port(),
		Params:  adaptor.params(),
		Robot:   r,
		Adaptor: adaptor,
	}
}

func (c *connection) Connect() bool {
	log.Println("Connecting to " + c.Name + " on port " + c.Port + "...")
	return c.Adaptor.Connect()
}

func (c *connection) Finalize() bool {
	log.Println("Finalizing " + c.Name + "...")
	return c.Adaptor.Finalize()
}

func (c *connection) ToJSON() *JSONConnection {
	return &JSONConnection{
		Name:    c.Name,
		Port:    c.Port,
		Adaptor: c.Type,
	}
}

func (c *connection) port() string {
	return c.Port
}

func (c *connection) name() string {
	return c.Name
}

func (c *connection) params() map[string]interface{} {
	return c.Adaptor.params()
}
