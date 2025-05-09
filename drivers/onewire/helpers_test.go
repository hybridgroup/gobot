package onewire

import (
	"errors"
	"sync"
)

type oneWireAdaptorMock struct {
	mtx          sync.Mutex
	familyCode   byte
	serialNumber uint64
	sendCommands []string
	lastValue    int
	retErr       bool
}

func newOneWireTestAdaptor() *oneWireAdaptorMock {
	return &oneWireAdaptorMock{}
}

func (am *oneWireAdaptorMock) GetOneWireConnection(familyCode byte, serialNumber uint64) (Connection, error) {
	am.mtx.Lock()
	defer am.mtx.Unlock()

	if am.retErr {
		return nil, errors.New("GetOneWireConnection error")
	}
	am.familyCode = familyCode
	am.serialNumber = serialNumber

	return am, nil
}

// implementations of gobot.OneWireOperations
func (am *oneWireAdaptorMock) ID() string { return "" }

func (am *oneWireAdaptorMock) ReadData(command string, data []byte) error {
	am.sendCommands = append(am.sendCommands, command)

	return nil
}

func (am *oneWireAdaptorMock) WriteData(command string, data []byte) error {
	am.sendCommands = append(am.sendCommands, command)

	return nil
}

func (am *oneWireAdaptorMock) ReadInteger(command string) (int, error) {
	am.sendCommands = append(am.sendCommands, command)
	if am.retErr {
		return 0, errors.New("ReadInteger error")
	}

	return am.lastValue, nil
}

func (am *oneWireAdaptorMock) WriteInteger(command string, val int) error {
	am.sendCommands = append(am.sendCommands, command)
	if am.retErr {
		return errors.New("WriteInteger error")
	}

	return nil
}

func (am *oneWireAdaptorMock) Close() error { return nil }

// implementations of gobot.Connection, respectively gobot.Adaptor
func (am *oneWireAdaptorMock) Name() string        { return "" }
func (am *oneWireAdaptorMock) SetName(name string) {}
func (am *oneWireAdaptorMock) Connect() error      { return nil }
func (am *oneWireAdaptorMock) Finalize() error     { return nil }

type oneWireSystemDeviceMock struct {
	id          string
	lastValue   int
	lastData    []byte
	retErr      error
	lastCommand string
}

//nolint:unparam // ok here
func newOneWireTestSystemDevice(id string) *oneWireSystemDeviceMock {
	return &oneWireSystemDeviceMock{id: id}
}

func (dm *oneWireSystemDeviceMock) ID() string { return dm.id }

func (dm *oneWireSystemDeviceMock) ReadData(command string, data []byte) error {
	dm.lastCommand = command
	copy(data, dm.lastData)

	return dm.retErr
}

func (dm *oneWireSystemDeviceMock) WriteData(command string, data []byte) error {
	dm.lastCommand = command
	dm.lastData = make([]byte, len(data))
	copy(dm.lastData, data)

	return dm.retErr
}

func (dm *oneWireSystemDeviceMock) ReadInteger(command string) (int, error) {
	dm.lastCommand = command

	return dm.lastValue, dm.retErr
}

func (dm *oneWireSystemDeviceMock) WriteInteger(command string, val int) error {
	dm.lastCommand = command
	dm.lastValue = val

	return dm.retErr
}

func (dm *oneWireSystemDeviceMock) Close() error { return dm.retErr }
