package client

import (
	"bytes"
	"encoding/binary"
	"sync"
)

type nwFrameGenerator struct {
	seq   map[byte]byte
	hlen  int
	mutex *sync.Mutex
}

func newNetworkFrameGenerator() *nwFrameGenerator {
	nwg := nwFrameGenerator{
		seq:   make(map[byte]byte), // each frame id has it's own sequence number
		hlen:  7,                   // size of ARNETWORKAL_Frame_t header
		mutex: &sync.Mutex{},
	}
	return &nwg
}

// generate the "NetworkFrame" as bytes buffer
func (nwg *nwFrameGenerator) generate(cmd *bytes.Buffer, frameType byte, id byte) *bytes.Buffer {
	nwg.mutex.Lock()
	defer nwg.mutex.Unlock()

	// func networkFrameGenerator() func(*bytes.Buffer, byte, byte) NetworkFrame {
	//
	// ARNETWORKAL_Frame_t
	//
	// uint8  type  - frame type ARNETWORK_FRAME_TYPE
	// uint8  id    - identifier of the buffer sending the frame
	// uint8  seq   - sequence number of the frame
	// uint32 size  - size of the frame
	//

	if _, ok := nwg.seq[id]; !ok {
		nwg.seq[id] = 0
	}

	nwg.seq[id]++ // automatically rollover to 0 when > 255 is fine

	ret := &bytes.Buffer{}
	ret.WriteByte(frameType)
	ret.WriteByte(id)
	ret.WriteByte(nwg.seq[id])

	size := &bytes.Buffer{}
	//nolint:gosec // TODO: fix later
	if err := binary.Write(size, binary.LittleEndian, uint32(cmd.Len()+nwg.hlen)); err != nil {
		panic(err)
	}

	ret.Write(size.Bytes())
	ret.Write(cmd.Bytes())

	return ret
}
