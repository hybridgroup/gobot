package mavlink

/** @file
 *	@brief MAVLink comm protocol built from common.xml
 *	@see http://pixhawk.ethz.ch/software/mavlink
 */

import (
	"bytes"
	"encoding/binary"
	"io"
	"time"
)

const MAVLINK_BIG_ENDIAN = 0
const MAVLINK_LITTLE_ENDIAN = 1
const MAVLINK_STX = 254
const MAVLINK_ENDIAN = MAVLINK_LITTLE_ENDIAN
const MAVLINK_ALIGNED_FIELDS = 1
const MAVLINK_CRC_EXTRA = 1
const X25_INIT_CRC = 0xffff
const X25_VALIDATE_CRC = 0xf0b8

var sequence uint16 = 0

func generateSequence() uint8 {
	sequence = (sequence + 1) % 256
	return uint8(sequence)
}

type MAVLinkMessage interface {
	Id() uint8
	Len() uint8
	Crc() uint8
	Pack() []byte
	Decode([]byte)
}

type MAVLinkPacket struct {
	Protocol    uint8
	Length      uint8
	Sequence    uint8
	SystemID    uint8
	ComponentID uint8
	MessageID   uint8
	Data        []uint8
	Checksum    uint16
}

func ReadMAVLinkPacket(r io.Reader) (*MAVLinkPacket, error) {
	for {
		header, err := read(r, 1)
		if err != nil {
			return nil, err
		}
		if header[0] == 254 {
			length, err := read(r, 1)
			if err != nil {
				return nil, err
			} else if length[0] > 250 {
				continue
			}
			m := &MAVLinkPacket{}
			data, err := read(r, int(length[0])+7)
			if err != nil {
				return nil, err
			}
			data = append([]byte{header[0], length[0]}, data...)
			m.Decode(data)
			return m, nil
		}
	}
}

func CraftMAVLinkPacket(SystemID uint8, ComponentID uint8, Message MAVLinkMessage) *MAVLinkPacket {
	return NewMAVLinkPacket(
		0xFE,
		Message.Len(),
		generateSequence(),
		SystemID,
		ComponentID,
		Message.Id(),
		Message.Pack(),
	)
}

func NewMAVLinkPacket(Protocol uint8, Length uint8, Sequence uint8, SystemID uint8, ComponentID uint8, MessageID uint8, Data []uint8) *MAVLinkPacket {
	m := &MAVLinkPacket{
		Protocol:    Protocol,
		Length:      Length,
		Sequence:    Sequence,
		SystemID:    SystemID,
		ComponentID: ComponentID,
		MessageID:   MessageID,
		Data:        Data,
	}
	m.Checksum = crcCalculate(m)
	return m
}

func (m *MAVLinkPacket) MAVLinkMessage() (MAVLinkMessage, error) {
	return NewMAVLinkMessage(m.MessageID, m.Data)
}

func (m *MAVLinkPacket) Pack() []byte {
	data := new(bytes.Buffer)
	binary.Write(data, binary.LittleEndian, m.Protocol)
	binary.Write(data, binary.LittleEndian, m.Length)
	binary.Write(data, binary.LittleEndian, m.Sequence)
	binary.Write(data, binary.LittleEndian, m.SystemID)
	binary.Write(data, binary.LittleEndian, m.ComponentID)
	binary.Write(data, binary.LittleEndian, m.MessageID)
	data.Write(m.Data)
	binary.Write(data, binary.LittleEndian, m.Checksum)
	return data.Bytes()
}

func (m *MAVLinkPacket) Decode(buf []byte) {
	m.Protocol = buf[0]
	m.Length = buf[1]
	m.Sequence = buf[2]
	m.SystemID = buf[3]
	m.ComponentID = buf[4]
	m.MessageID = buf[5]
	m.Data = buf[6 : 6+int(m.Length)]
	checksum := buf[7+int(m.Length):]
	m.Checksum = uint16(checksum[1])<<8 | uint16(checksum[0])
}

func read(r io.Reader, length int) ([]byte, error) {
	buf := []byte{}
	for length > 0 {
		tmp := make([]byte, length)
		i, err := r.Read(tmp[:])
		if err != nil {
			return nil, err
		} else {
			length -= i
			buf = append(buf, tmp...)
			if length != i {
				<-time.After(1 * time.Millisecond)
			} else {
				break
			}
		}
	}
	return buf, nil
}

/**
 * @brief Accumulate the X.25 CRC by adding one char at a time.
 *
 * The checksum function adds the hash of one char at a time to the
 * 16 bit checksum (uint16).
 *
 * @param data to hash
 * @param crcAccum the already accumulated checksum
 **/
func crcAccumulate(data uint8, crcAccum uint16) uint16 {
	/*Accumulate one byte of data into the CRC*/
	var tmp uint8

	tmp = data ^ (uint8)(crcAccum&0xff)
	tmp ^= (tmp << 4)
	crcAccum = (uint16(crcAccum) >> 8) ^ (uint16(tmp) << 8) ^ (uint16(tmp) << 3) ^ (uint16(tmp) >> 4)
	return crcAccum
}

/**
 * @brief Initiliaze the buffer for the X.25 CRC
 *
 */
func crcInit() uint16 {
	return X25_INIT_CRC
}

/**
 * @brief Calculates the X.25 checksum on a byte buffer
 *
 * @return the checksum over the buffer bytes
 **/
func crcCalculate(m *MAVLinkPacket) uint16 {
	crc := crcInit()

	for _, v := range m.Pack()[1 : m.Length+6] {
		crc = crcAccumulate(v, crc)
	}
	message, _ := m.MAVLinkMessage()
	crc = crcAccumulate(message.Crc(), crc)
	return crc
}
