package hs200

import (
	"fmt"
	"net"
	"time"
)

type Driver struct {
	cmd     []byte
	udpconn net.Conn
	tcpconn net.Conn
}

func NewDriver(tcpaddress string, udpaddress string) (*Driver, error) {
	tc, terr := net.Dial("tcp", tcpaddress)
	if terr != nil {
		return nil, terr
	}
	uc, uerr := net.Dial("udp4", udpaddress)
	if uerr != nil {
		return nil, uerr
	}
	command := []byte{
		0xff, //unknown header sent with apparently constant value
		0x04, //unknown header sent with apparently constant value
		0x7e, // 0x7f,//vertical lift up/down
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

	return &Driver{command, uc, tc}, nil
}

func (d Driver) flightloop() {
	i := 0
	for {
		time.Sleep(20 * time.Millisecond)
		if i%5 == 0 {
			d.tcpconn.Write([]byte("remote\r\n"))
		}
		d.sendUDP()
		i++
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
	go d.flightloop()
}

func (d Driver) TakeOff() {
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
	d.cmd[2] = 0x3f
	d.cmd[10] = checksum(d.cmd)
}
