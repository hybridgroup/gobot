package gobot

import (
	"errors"
	"log"
	"reflect"
)

type Connection interface {
	Connect() bool
	Finalize() bool
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
	c := new(connection)
	s := reflect.ValueOf(adaptor).Type().String()
	c.Type = s[1:len(s)]
	c.Name = FieldByNamePtr(adaptor, "Name").String()
	c.Port = FieldByNamePtr(adaptor, "Port").String()
	c.Params = make(map[string]interface{})
	keys := FieldByNamePtr(adaptor, "Params").MapKeys()
	for k := range keys {
		c.Params[keys[k].String()] = FieldByNamePtr(adaptor, "Params").MapIndex(keys[k])
	}
	c.Robot = r
	c.Adaptor = adaptor
	return c
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
