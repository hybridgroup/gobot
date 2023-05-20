package microbit

import (
	"sync"

	"gobot.io/x/gobot/v2/platforms/ble"
)

var _ ble.BLEConnector = (*bleTestClientAdaptor)(nil)

type bleTestClientAdaptor struct {
	name             string
	address          string
	mtx              sync.Mutex
	withoutResponses bool

	testSubscribe           func([]byte, error)
	testReadCharacteristic  func(string) ([]byte, error)
	testWriteCharacteristic func(string, []byte) error
}

func (t *bleTestClientAdaptor) Connect() (err error)      { return }
func (t *bleTestClientAdaptor) Reconnect() (err error)    { return }
func (t *bleTestClientAdaptor) Disconnect() (err error)   { return }
func (t *bleTestClientAdaptor) Finalize() (err error)     { return }
func (t *bleTestClientAdaptor) Name() string              { return t.name }
func (t *bleTestClientAdaptor) SetName(n string)          { t.name = n }
func (t *bleTestClientAdaptor) Address() string           { return t.address }
func (t *bleTestClientAdaptor) WithoutResponses(use bool) { t.withoutResponses = use }

func (t *bleTestClientAdaptor) ReadCharacteristic(cUUID string) (data []byte, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testReadCharacteristic(cUUID)
}

func (t *bleTestClientAdaptor) WriteCharacteristic(cUUID string, data []byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testWriteCharacteristic(cUUID, data)
}

//nolint:revive // in tests it is be helpful to see the meaning of the parameters by name
func (t *bleTestClientAdaptor) Subscribe(cUUID string, f func([]byte, error)) (err error) {
	t.testSubscribe = f
	return
}

func (t *bleTestClientAdaptor) TestReadCharacteristic(f func(cUUID string) (data []byte, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testReadCharacteristic = f
}

func (t *bleTestClientAdaptor) TestWriteCharacteristic(f func(cUUID string, data []byte) (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testWriteCharacteristic = f
}

func (t *bleTestClientAdaptor) TestReceiveNotification(data []byte, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testSubscribe(data, err)
}

func NewBleTestAdaptor() *bleTestClientAdaptor {
	return &bleTestClientAdaptor{
		address: "01:02:03:04:05:06",
		testReadCharacteristic: func(cUUID string) (data []byte, e error) {
			return
		},
		testWriteCharacteristic: func(cUUID string, data []byte) (e error) {
			return
		},
		testSubscribe: func([]byte, error) {},
	}
}
