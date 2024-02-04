package testutil

import (
	"sync"

	"gobot.io/x/gobot/v2"
)

var _ gobot.BLEConnector = (*bleTestClientAdaptor)(nil)

type bleTestClientAdaptor struct {
	name             string
	address          string
	mtx              sync.Mutex
	withoutResponses bool

	readCharacteristicFunc  func(string) ([]byte, error)
	writeCharacteristicFunc func(string, []byte) error
	subscribeFunc           func([]byte, error)
}

func NewBleTestAdaptor() *bleTestClientAdaptor {
	return &bleTestClientAdaptor{
		address: "01:02:03:0A:0B:0C",
		readCharacteristicFunc: func(cUUID string) ([]byte, error) {
			return []byte(cUUID), nil
		},
		writeCharacteristicFunc: func(cUUID string, data []byte) error {
			return nil
		},
	}
}

func (t *bleTestClientAdaptor) SetReadCharacteristicTestFunc(f func(cUUID string) (data []byte, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.readCharacteristicFunc = f
}

func (t *bleTestClientAdaptor) SetWriteCharacteristicTestFunc(f func(cUUID string, data []byte) error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.writeCharacteristicFunc = f
}

func (t *bleTestClientAdaptor) SendTestDataToSubscriber(data []byte, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.subscribeFunc(data, err)
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
	return t.readCharacteristicFunc(cUUID)
}

func (t *bleTestClientAdaptor) WriteCharacteristic(cUUID string, data []byte) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.writeCharacteristicFunc(cUUID, data)
}

func (t *bleTestClientAdaptor) Subscribe(cUUID string, f func([]byte, error)) error {
	t.subscribeFunc = f
	return nil
}
