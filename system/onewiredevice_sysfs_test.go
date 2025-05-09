package system

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_newOneWireDeviceSysfs(t *testing.T) {
	// arrange
	m := &MockFilesystem{}
	sfa := sysfsFileAccess{fs: m, readBufLen: 2}
	const id = "0815"
	// act
	d := newOneWireDeviceSysfs(&sfa, id)
	// assert
	assert.Equal(t, "/sys/bus/w1/devices/"+id, d.sysfsPath)
	assert.Equal(t, &sfa, d.sfa)
	assert.Equal(t, id, d.ID())
}

func TestOneWireDeviceReadData(t *testing.T) {
	// arrange
	const (
		id             = "0816"
		command        = "getValue"
		path           = "/sys/bus/w1/devices/" + id + "/" + command
		content        = "234"
		countBytesRead = 3
	)
	fs := newMockFilesystem([]string{path})
	sfa := sysfsFileAccess{fs: fs, readBufLen: countBytesRead}
	d := newOneWireDeviceSysfs(&sfa, id)
	fs.Files[path].Contents = content
	data := []byte{1, 1, 1}
	// act
	err := d.ReadData(command, data)
	// assert
	require.NoError(t, err)
	assert.Equal(t, []byte(content), data)
}

func TestOneWireDeviceWriteData(t *testing.T) {
	// arrange
	const (
		id      = "0817"
		command = "putValue"
		path    = "/sys/bus/w1/devices/" + id + "/" + command
	)
	fs := newMockFilesystem([]string{path})
	sfa := sysfsFileAccess{fs: fs}
	d := newOneWireDeviceSysfs(&sfa, id)
	fs.Files[path].Contents = "old content"
	data := []byte{1, 2, 3}
	// act
	err := d.WriteData(command, data)
	// assert
	require.NoError(t, err)
	assert.Equal(t, data, []byte(fs.Files[path].Contents))
}

func TestOneWireDeviceReadInteger(t *testing.T) {
	// arrange
	const (
		id             = "0818"
		command        = "getIntegerValue"
		path           = "/sys/bus/w1/devices/" + id + "/" + command
		content        = "23456"
		countBytesRead = 4
		want           = 2345 // limited by readBufLen
	)
	fs := newMockFilesystem([]string{path})
	sfa := sysfsFileAccess{fs: fs, readBufLen: countBytesRead}
	d := newOneWireDeviceSysfs(&sfa, id)
	fs.Files[path].Contents = content
	// act
	got, err := d.ReadInteger(command)
	// assert
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestOneWireDeviceWriteInteger(t *testing.T) {
	// arrange
	const (
		id      = "0818"
		command = "getIntegerValue"
		path    = "/sys/bus/w1/devices/" + id + "/" + command
		write   = 155
	)
	fs := newMockFilesystem([]string{path})
	sfa := sysfsFileAccess{fs: fs}
	d := newOneWireDeviceSysfs(&sfa, id)
	fs.Files[path].Contents = "old content"
	// act
	err := d.WriteInteger(command, write)
	// assert
	require.NoError(t, err)
	assert.Equal(t, strconv.Itoa(write), fs.Files[path].Contents)
}
