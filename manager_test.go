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
	"github.com/stretchr/testify/require"
)

func initTestManager() *Manager {
	log.SetOutput(&NullReadWriteCloser{})
	g := NewManager()
	g.trap = func(c chan os.Signal) {
		c <- os.Interrupt
	}
	g.AddRobot(newTestRobot("Robot1"))
	g.AddRobot(newTestRobot("Robot2"))
	g.AddRobot(newTestRobot(""))
	return g
}

func initTestManager1Robot() *Manager {
	log.SetOutput(&NullReadWriteCloser{})
	g := NewManager()
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
	require.NoError(t, n.Close())
}

func TestManagerRobot(t *testing.T) {
	g := initTestManager()
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

func TestManagerToJSON(t *testing.T) {
	g := initTestManager()
	g.AddCommand("test_function", func(params map[string]interface{}) interface{} {
		return nil
	})
	json := NewJSONManager(g)
	assert.Len(t, json.Robots, g.Robots().Len())
	assert.Len(t, json.Commands, len(g.Commands()))
}

func TestManagerStart(t *testing.T) {
	g := initTestManager()
	require.NoError(t, g.Start())
	require.NoError(t, g.Stop())
	assert.False(t, g.Running())
}

func TestManagerStartAutoRun(t *testing.T) {
	g := NewManager()
	g.AddRobot(newTestRobot("Robot99"))
	go func() { _ = g.Start() }()
	time.Sleep(10 * time.Millisecond)
	assert.True(t, g.Running())

	// stop it
	require.NoError(t, g.Stop())
	assert.False(t, g.Running())
}

func TestManagerStartDriverErrors(t *testing.T) {
	g := initTestManager1Robot()
	e := errors.New("driver start error 1")
	testDriverStart = func() error {
		return e
	}

	var want error
	want = multierror.Append(want, e)
	want = multierror.Append(want, e)
	want = multierror.Append(want, e)

	assert.Equal(t, want, g.Start())
	require.NoError(t, g.Stop())

	testDriverStart = func() error { return nil }
}

func TestManagerHaltFromRobotDriverErrors(t *testing.T) {
	g := initTestManager1Robot()
	var ec int
	testDriverHalt = func() error {
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

func TestManagerStartRobotAdaptorErrors(t *testing.T) {
	g := initTestManager1Robot()
	var ec int
	testAdaptorConnect = func() error {
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
	require.NoError(t, g.Stop())

	testAdaptorConnect = func() error { return nil }
}

func TestManagerFinalizeErrors(t *testing.T) {
	g := initTestManager1Robot()
	var ec int
	testAdaptorFinalize = func() error {
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
