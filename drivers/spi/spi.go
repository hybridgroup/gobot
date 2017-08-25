package spi

import (
	"fmt"
	"sync"

	"golang.org/x/exp/io/spi"
	"golang.org/x/exp/io/spi/driver"
)

type Connection struct {
	connection *spi.Device
	address    int
	mutex      *sync.Mutex
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
func (c *Connection) Read32(address byte, msg byte) int32 {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	// read build two arrays 8 bytes wide, 4 for addressing, and 4 for the response
	w := make([]byte, 8)
	w[0] = address
	w[1] = msg
	r := make([]byte, len(w))
	err := c.connection.Tx(w, r)
	if err != nil {
		panic(err)
	}
	return int32(int(r[4])<<24 | int(r[5])<<16 | int(r[6])<<8 | int(r[7]))
}

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
