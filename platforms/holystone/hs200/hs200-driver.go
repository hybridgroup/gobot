package hs200

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Driver struct {
	mutex   sync.RWMutex
	stop    chan struct{}
	cmd     []byte
	enabled bool
	udpconn net.Conn
	tcpconn net.Conn
}

func NewDriver(tcpaddress string, udpaddress string) (*Driver, error) {
//	tc, terr := net.Dial("tcp", tcpaddress)
//	if terr != nil {
//		return nil, terr
//	}
	uc, uerr := net.Dial("udp4", udpaddress)
	if uerr != nil {
		return nil, uerr
	}
	command := []byte{
		0xff, //unknown header sent with apparently constant value
		0x04, //unknown header sent with apparently constant value
		0x3f, //vertical lift up/down
		0x3f, //rotation rate left/right
		0xc0, //advance forward / backward
		0x3f, //strafe left / right
		0x90, //yaw (used as a setting to trim the yaw of the uav)
		0x10, //pitch (used as a setting to trim the pitch of the uav)
		0x10, //roll (used as a setting to trim the roll of the uav)
		0x40, //throttle
		0x00, //this is a sanity check; 255 - ((sum of flight controls from index 1 to 9) % 256)
	}
	command[10] = checksum(command)

	return &Driver{stop: make(chan struct{}), cmd: command, udpconn: uc, tcpconn: nil}, nil
}

func (d *Driver) flightLoop(stop chan struct{}) {
	udpTick := time.NewTicker(50 * time.Millisecond)
	defer udpTick.Stop()
	tcpTick := time.NewTicker(1000 * time.Millisecond)
	defer tcpTick.Stop()
	for {
		select {
		case <-udpTick.C:
			d.mutex.RLock()
			defer d.mutex.RUnlock()
			d.sendUDP()
		case <-tcpTick.C:
			//d.tcpconn.Write([]byte("1\r\n"))
		case <-stop:
			d.mutex.Lock()
			defer d.mutex.Unlock()
			d.enabled = false
			return
		}
	}
}

func checksum(c []byte) byte {
	var sum byte
	for i := 1; i < 10; i++ {
		sum += c[i]
	}
	return 255 - sum
}
func (d Driver) sendUDP() {
	d.udpconn.Write(d.cmd)
}

func (d Driver) Enable() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if !d.enabled {
		go d.flightLoop(d.stop)
		d.enabled = true
	}
}

func (d Driver) Disable() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if d.enabled {
		d.stop <- struct{}{}
	}
}

func (d Driver) TakeOff() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.cmd[2] = 0x7e
	d.cmd[10] = checksum(d.cmd)
}

func (d Driver) VerticalControl(delta int) {
	current := int(d.cmd[2])
	current += delta
	if current > 255 {
		current = 255
	}
	if current < 0 {
		current = 0
	}
	fmt.Printf("Setting to %v", byte(current))
	d.cmd[2] = byte(delta)
	d.cmd[10] = checksum(d.cmd)

}

func (d Driver) Land() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.cmd[2] = 0
	d.cmd[10] = checksum(d.cmd)
}
