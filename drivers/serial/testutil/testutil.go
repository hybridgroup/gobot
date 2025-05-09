package testutil

import "fmt"

type serialTestAdaptor struct {
	isConnected bool
	name        string

	simulateConnectErr bool
	simulateReadErr    bool
	simulateWriteErr   bool
}

func NewSerialTestAdaptor() *serialTestAdaptor {
	return &serialTestAdaptor{}
}

func (t *serialTestAdaptor) SetSimulateConnectError(val bool) {
	t.simulateConnectErr = val
}

func (t *serialTestAdaptor) SetSimulateReadError(val bool) {
	t.simulateReadErr = val
}

func (t *serialTestAdaptor) SetSimulateWriteError(val bool) {
	t.simulateWriteErr = val
}

func (t *serialTestAdaptor) IsConnected() bool {
	return t.isConnected
}

func (t *serialTestAdaptor) SerialRead(b []byte) (int, error) {
	if t.simulateReadErr {
		return 0, fmt.Errorf("read error")
	}

	return len(b), nil
}

func (t *serialTestAdaptor) SerialWrite(b []byte) (int, error) {
	if t.simulateWriteErr {
		return 0, fmt.Errorf("write error")
	}

	return len(b), nil
}

// gobot.Adaptor interfaces
func (t *serialTestAdaptor) Connect() error {
	if t.simulateConnectErr {
		return fmt.Errorf("connect error")
	}

	t.isConnected = true
	return nil
}

func (t *serialTestAdaptor) Finalize() error  { return nil }
func (t *serialTestAdaptor) Name() string     { return t.name }
func (t *serialTestAdaptor) SetName(n string) { t.name = n }
