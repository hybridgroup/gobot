package tello

import (
	"github.com/go-daq/crc8"
	"github.com/howeyc/crc16"
)

// CalculateCRC8 calculates the starting CRC8 byte for packet.
func CalculateCRC8(bytes []byte) byte {
	return crc8.Checksum(bytes, crc8.MakeTable(0x77))
}

// CalculateCRC16 calculates the ending CRC16 bytes for packet.
func CalculateCRC16(pkt []byte) (low, high byte) {
	i := crc16.Checksum(pkt, crc16.MakeTable(0x3692))
	low = ((byte)(i & 0xFF))
	high = ((byte)(i >> 8 & 0xFF))
	return
}
