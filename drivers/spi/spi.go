package spi

import (
	"golang.org/x/exp/io/spi"
	"sync"
	"golang.org/x/exp/io/spi/driver"
	"fmt"
)

type Connection struct {
	connection *spi.Device
	address int
	mutex *sync.Mutex
}

// NewConnection creates and returns a new connection to a specific
// spi device on a bus and address
func NewConnection(opener driver.Opener, address int) (connection *Connection) {
	dev, err := spi.Open(opener)
	if err != nil {
		panic(err)
	}

	return &Connection{connection: dev, address: address, mutex: &sync.Mutex{}}
}

// TODO: Read - the exp/io/spi implementation doesn't currently support reading

// Write data to a spi device
func (c *Connection) Write(data []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	fmt.Println(data)
	err := c.connection.Tx(data, nil)

	return err
}

// Close the connection
func (c *Connection) Close() error {
	err := c.connection.Close()
	return err
}