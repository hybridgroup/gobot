package spi

type spiTestAdaptor struct{}

func (ctr *spiTestAdaptor) GetSpiConnection(busNum, chipNum, mode, bits int, maxSpeed int64) (device Connection, err error) {
	return &spiTestDevice{}, nil
}

func (ctr *spiTestAdaptor) GetSpiDefaultBus() int        { return 0 }
func (ctr *spiTestAdaptor) GetSpiDefaultChip() int       { return 0 }
func (ctr *spiTestAdaptor) GetSpiDefaultMode() int       { return 0 }
func (ctr *spiTestAdaptor) GetSpiDefaultBits() int       { return 0 }
func (ctr *spiTestAdaptor) GetSpiDefaultMaxSpeed() int64 { return 0 }

type spiTestDevice struct {
}

func (c *spiTestDevice) Tx(w, r []byte) error { return nil }
func (c *spiTestDevice) Close() error         { return nil }

func newSpiTestAdaptor() *spiTestAdaptor {
	return &spiTestAdaptor{}
}
