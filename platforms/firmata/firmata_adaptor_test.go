package firmata

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("FirmataAdaptor", func() {
	var (
		adaptor *FirmataAdaptor
	)

	BeforeEach(func() {
		adaptor = NewFirmataAdaptor("tom", "/dev/null")
		adaptor.connect = func(f *FirmataAdaptor) {
			f.Board = newBoard(sp{})
			f.Board.Events = append(f.Board.Events, event{Name: "firmware_query"})
			f.Board.Events = append(f.Board.Events, event{Name: "capability_query"})
			f.Board.Events = append(f.Board.Events, event{Name: "analog_mapping_query"})
		}
		adaptor.Connect()
	})

	It("Must be able to Finalize", func() {
		Expect(adaptor.Finalize()).To(Equal(true))
	})
	It("Must be able to Disconnect", func() {
		Expect(adaptor.Disconnect()).To(Equal(true))
	})
	It("Must be able to InitServo", func() {
		adaptor.InitServo()
	})
	It("Must be able to ServoWrite", func() {
		adaptor.ServoWrite("1", 50)
	})
	It("Must be able to PwmWrite", func() {
		adaptor.PwmWrite("1", 50)
	})
	It("Must be able to DigitalWrite", func() {
		adaptor.DigitalWrite("1", 1)
	})
	It("DigitalRead should return -1 on no data", func() {
		Expect(adaptor.DigitalRead("1")).To(Equal(-1))
	})
	It("DigitalRead should return data", func() {
		pinNumber := "1"
		adaptor.Board.Events = append(adaptor.Board.Events, event{Name: fmt.Sprintf("digital_read_%v", pinNumber), Data: []byte{0x01}})
		Expect(adaptor.DigitalRead(pinNumber)).To(Equal(0x01))
	})
	It("AnalogRead should return -1 on no data", func() {
		Expect(adaptor.AnalogRead("1")).To(Equal(-1))
	})
	It("AnalogRead should return data", func() {
		pinNumber := "1"
		value := 133
		adaptor.Board.Events = append(adaptor.Board.Events, event{Name: fmt.Sprintf("analog_read_%v", pinNumber), Data: []byte{byte(value >> 24), byte(value >> 16), byte(value >> 8), byte(value & 0xff)}})
		Expect(adaptor.AnalogRead(pinNumber)).To(Equal(133))
	})
	It("Must be able to I2cStart", func() {
		adaptor.I2cStart(0x00)
	})
	It("I2cRead should return [] on no data", func() {
		Expect(adaptor.I2cRead(1)).To(Equal(make([]byte, 0)))
	})
	It("I2cRead should return data", func() {
		i := []byte{100}
		i2cReply := map[string][]byte{}
		i2cReply["data"] = i
		adaptor.Board.Events = append(adaptor.Board.Events, event{Name: "i2c_reply", I2cReply: i2cReply})
		Expect(adaptor.I2cRead(1)).To(Equal(i))
	})
	It("Must be able to I2cWrite", func() {
		adaptor.I2cWrite([]byte{0x00, 0x01})
	})
})
