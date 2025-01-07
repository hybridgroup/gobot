package system

import (
	"fmt"
	"path"
)

type onewireDeviceSysfs struct {
	id        string
	sysfsPath string
	sfa       *sysfsFileAccess
}

func newOneWireDeviceSysfs(sfa *sysfsFileAccess, id string) *onewireDeviceSysfs {
	p := &onewireDeviceSysfs{
		id:        id,
		sysfsPath: path.Join("/sys/bus/w1/devices", id),
		sfa:       sfa,
	}
	return p
}

// ID returns the device id in the form "family code"-"serial number". Implements gobot.OneWireSystemDevicer.
func (o *onewireDeviceSysfs) ID() string {
	return o.id
}

// ReadData reads from the sysfs path specified by the command. Implements gobot.OneWireSystemDevicer.
func (o *onewireDeviceSysfs) ReadData(command string, data []byte) error {
	p := path.Join(o.sysfsPath, command)
	buf, err := o.sfa.read(p)
	if err != nil {
		return err
	}
	copy(data, buf)

	if len(buf) < len(data) {
		return fmt.Errorf("count of read bytes (%d) is smaller than expected (%d)", len(buf), len(data))
	}

	return nil
}

// WriteData writes to the path specified by the command. Implements gobot.OneWireSystemDevicer.
func (o *onewireDeviceSysfs) WriteData(command string, data []byte) error {
	p := path.Join(o.sysfsPath, command)
	return o.sfa.write(p, data)
}

// ReadInteger reads an integer value from the device. Implements gobot.OneWireSystemDevicer.
func (o *onewireDeviceSysfs) ReadInteger(command string) (int, error) {
	p := path.Join(o.sysfsPath, command)
	return o.sfa.readInteger(p)
}

// WriteInteger writes an integer value to the device. Implements gobot.OneWireSystemDevicer.
func (o *onewireDeviceSysfs) WriteInteger(command string, val int) error {
	p := path.Join(o.sysfsPath, command)
	return o.sfa.writeInteger(p, val)
}

// Close the 1-wire connection. Implements gobot.OneWireSystemDevicer.
func (o *onewireDeviceSysfs) Close() error {
	// currently nothing to do here - the file descriptors will be closed immediately after read/write
	return nil
}
