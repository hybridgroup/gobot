package minidrone

import (
	"sync"

	"gobot.io/x/gobot/v2/platforms/ble"
)

var _ ble.BLEConnector = (*bleTestClientAdaptor)(nil)

type bleTestClientAdaptor struct {
	name                    string
	address                 string
	mtx                     sync.Mutex
	withoutResponses        bool
	testReadCharacteristic  func(string) ([]byte, error)
	testWriteCharacteristic func(string, []byte) error
}

func (t *bleTestClientAdaptor) Connect() error            { return nil }
func (t *bleTestClientAdaptor) Reconnect() error          { return nil }
func (t *bleTestClientAdaptor) Disconnect() error         { return nil }
func (t *bleTestClientAdaptor) Finalize() error           { return nil }
func (t *bleTestClientAdaptor) Name() string              { return t.name }
func (t *bleTestClientAdaptor) SetName(n string)          { t.name = n }
func (t *bleTestClientAdaptor) Address() string           { return t.address }
func (t *bleTestClientAdaptor) WithoutResponses(use bool) { t.withoutResponses = use }

func (t *bleTestClientAdaptor) ReadCharacteristic(cUUID string) ([]byte, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testReadCharacteristic(cUUID)
}

func (t *bleTestClientAdaptor) WriteCharacteristic(cUUID string, data []byte) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testWriteCharacteristic(cUUID, data)
}

func (t *bleTestClientAdaptor) Subscribe(cUUID string, f func([]byte, error)) error {
	// TODO: implement this...
	return nil
}

func (t *bleTestClientAdaptor) TestReadCharacteristic(f func(cUUID string) (data []byte, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testReadCharacteristic = f
}

func (t *bleTestClientAdaptor) TestWriteCharacteristic(f func(cUUID string, data []byte) error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testWriteCharacteristic = f
}

func NewBleTestAdaptor() *bleTestClientAdaptor {
	return &bleTestClientAdaptor{
		address: "01:02:03:04:05:06",
		testReadCharacteristic: func(cUUID string) ([]byte, error) {
			return nil, nil
		},
		testWriteCharacteristic: func(cUUID string, data []byte) error {
			return nil
		},
	}
}
