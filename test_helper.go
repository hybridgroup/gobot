package gobot

type testDriver struct {
	Driver
}

//func (me *testDriver) Start() bool { return true }
func (me *testDriver) Start() bool { return true }

type testAdaptor struct {
	Adaptor
}

func (me *testAdaptor) Finalize() bool   { return true }
func (me *testAdaptor) Connect() bool    { return true }
func (me *testAdaptor) Disconnect() bool { return true }
func (me *testAdaptor) Reconnect() bool  { return true }

func newTestDriver(name string) *testDriver {
	d := new(testDriver)
	d.Name = name
	return d
}
func newTestAdaptor(name string) *testAdaptor {
	a := new(testAdaptor)
	a.Name = name
	return a
}
