package gobot

import (
	"fmt"
)

type testStruct struct {
	i int
	f float64
}

func (t *testStruct) Hello(name string, message string) string {
	return fmt.Sprintf("Hello %v! %v", name, message)
}

type null struct{}

func (null) Write(p []byte) (int, error) {
	return len(p), nil
}

type testDriver struct {
	Driver
	Adaptor *testAdaptor
}

func (t *testDriver) Init() bool  { return true }
func (t *testDriver) Start() bool { return true }
func (t *testDriver) Halt() bool  { return true }

type testAdaptor struct {
	Adaptor
}

func (t *testAdaptor) Finalize() bool { return true }
func (t *testAdaptor) Connect() bool  { return true }

func newTestDriver(name string, adaptor *testAdaptor) *testDriver {
	t := &testDriver{
		Driver: Driver{
			Commands: make(map[string]func(map[string]interface{}) interface{}),
			Name:     name,
		},
		Adaptor: adaptor,
	}

	t.Driver.AddCommand("TestDriverCommand", func(params map[string]interface{}) interface{} {
		name := params["name"].(string)
		return fmt.Sprintf("hello %v", name)
	})

	t.Driver.AddCommand("DriverCommand", func(params map[string]interface{}) interface{} {
		name := params["name"].(string)
		return fmt.Sprintf("hello %v", name)
	})

	return t
}

func newTestAdaptor(name string) *testAdaptor {
	return &testAdaptor{
		Adaptor: Adaptor{
			Name: name,
			Params: map[string]interface{}{
				"param1": "1",
				"param2": 2,
			},
		},
	}
}

func NewTestRobot(name string) *Robot {
	return newTestRobot(name)
}
func newTestRobot(name string) *Robot {
	adaptor1 := newTestAdaptor("Connection 1")
	adaptor2 := newTestAdaptor("Connection 2")
	adaptor3 := newTestAdaptor("Connection 3")
	driver1 := newTestDriver("Device 1", adaptor1)
	driver2 := newTestDriver("Device 2", adaptor2)
	driver3 := newTestDriver("Device 3", adaptor3)
	work := func() {}
	r := NewRobot(name, []Connection{adaptor1, adaptor2, adaptor3}, []Device{driver1, driver2, driver3}, work)
	r.AddCommand("robotTestFunction", func(params map[string]interface{}) interface{} {
		message := params["message"].(string)
		robot := params["robot"].(string)
		fmt.Println(params)
		return fmt.Sprintf("hey %v, %v", robot, message)
	})
	return r
}

func newTestStruct() *testStruct {
	s := new(testStruct)
	s.i = 10
	s.f = 0.2
	return s
}
