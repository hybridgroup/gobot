package gobot

import (
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func initTestMaster() *Master {
	log.SetOutput(&NullReadWriteCloser{})
	g := NewMaster()
	g.trap = func(c chan os.Signal) {
		c <- os.Interrupt
	}
	g.AddRobot(newTestRobot("Robot1"))
	g.AddRobot(newTestRobot("Robot2"))
	g.AddRobot(newTestRobot(""))
	return g
}

func initTestMaster1Robot() *Master {
	log.SetOutput(&NullReadWriteCloser{})
	g := NewMaster()
	g.trap = func(c chan os.Signal) {
		c <- os.Interrupt
	}
	g.AddRobot(newTestRobot("Robot99"))

	return g
}

func TestNullReadWriteCloser(t *testing.T) {
	n := &NullReadWriteCloser{}
	i, _ := n.Write([]byte{1, 2, 3})
	assert.Equal(t, 3, i)
	i, _ = n.Read(make([]byte, 10))
	assert.Equal(t, 10, i)
	assert.NoError(t, n.Close())
}

func TestMasterRobot(t *testing.T) {
	g := initTestMaster()
	assert.Equal(t, "Robot1", g.Robot("Robot1").Name)
	assert.Equal(t, (*Robot)(nil), g.Robot("Robot4"))
	assert.Equal(t, (Device)(nil), g.Robot("Robot4").Device("Device1"))
	assert.Equal(t, (Connection)(nil), g.Robot("Robot4").Connection("Connection1"))
	assert.Equal(t, (Device)(nil), g.Robot("Robot1").Device("Device4"))
	assert.Equal(t, "Device1", g.Robot("Robot1").Device("Device1").Name())
	assert.Equal(t, 3, g.Robot("Robot1").Devices().Len())
	assert.Equal(t, (Connection)(nil), g.Robot("Robot1").Connection("Connection4"))
	assert.Equal(t, 3, g.Robot("Robot1").Connections().Len())
}

func TestMasterToJSON(t *testing.T) {
	g := initTestMaster()
	g.AddCommand("test_function", func(params map[string]interface{}) interface{} {
		return nil
	})
	json := NewJSONMaster(g)
	assert.Equal(t, g.Robots().Len(), len(json.Robots))
	assert.Equal(t, len(g.Commands()), len(json.Commands))
}

func TestMasterStart(t *testing.T) {
	g := initTestMaster()
	assert.NoError(t, g.Start())
	assert.NoError(t, g.Stop())
	assert.False(t, g.Running())
}

func TestMasterStartAutoRun(t *testing.T) {
	g := NewMaster()
	g.AddRobot(newTestRobot("Robot99"))
	go func() { _ = g.Start() }()
	time.Sleep(10 * time.Millisecond)
	assert.True(t, g.Running())

	// stop it
	assert.NoError(t, g.Stop())
	assert.False(t, g.Running())
}

func TestMasterStartDriverErrors(t *testing.T) {
	g := initTestMaster1Robot()
	e := errors.New("driver start error 1")
	testDriverStart = func() (err error) {
		return e
	}

	var want error
	want = multierror.Append(want, e)
	want = multierror.Append(want, e)
	want = multierror.Append(want, e)

	assert.Equal(t, want, g.Start())
	assert.NoError(t, g.Stop())

	testDriverStart = func() (err error) { return }
}

func TestMasterHaltFromRobotDriverErrors(t *testing.T) {
	g := initTestMaster1Robot()
	var ec int
	testDriverHalt = func() (err error) {
		ec++
		return fmt.Errorf("driver halt error %d", ec)
	}
	defer func() { testDriverHalt = func() error { return nil } }()

	var want error
	for i := 1; i <= 3; i++ {
		e := fmt.Errorf("driver halt error %d", i)
		want = multierror.Append(want, e)
	}

	assert.Equal(t, want, g.Start())
}

func TestMasterStartRobotAdaptorErrors(t *testing.T) {
	g := initTestMaster1Robot()
	var ec int
	testAdaptorConnect = func() (err error) {
		ec++
		return fmt.Errorf("adaptor start error %d", ec)
	}
	defer func() { testAdaptorConnect = func() error { return nil } }()

	var want error
	for i := 1; i <= 3; i++ {
		e := fmt.Errorf("adaptor start error %d", i)
		want = multierror.Append(want, e)
	}

	assert.Equal(t, want, g.Start())
	assert.NoError(t, g.Stop())

	testAdaptorConnect = func() (err error) { return }
}

func TestMasterFinalizeErrors(t *testing.T) {
	g := initTestMaster1Robot()
	var ec int
	testAdaptorFinalize = func() (err error) {
		ec++
		return fmt.Errorf("adaptor finalize error %d", ec)
	}
	defer func() { testAdaptorFinalize = func() error { return nil } }()

	var want error
	for i := 1; i <= 3; i++ {
		e := fmt.Errorf("adaptor finalize error %d", i)
		want = multierror.Append(want, e)
	}

	assert.Equal(t, want, g.Start())
}
