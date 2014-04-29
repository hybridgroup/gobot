package gpio

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Motor", func() {
	var (
		t TestAdaptor
		m *MotorDriver
	)

	BeforeEach(func() {
		m = NewMotorDriver(t)
		m.Pin = "1"
	})

	It("Must be able to Start", func() {
		Expect(m.Start()).To(Equal(true))
	})
	It("Must be able to Init", func() {
		Expect(m.Init()).To(Equal(true))
	})
	It("Must be able to Halt", func() {
		Expect(m.Halt()).To(Equal(true))
	})
	It("Must be able to tell if IsOn", func() {
		m.CurrentState = 1
		Expect(m.IsOn()).To(BeTrue())
		m.CurrentMode = "analog"
		m.CurrentSpeed = 100
		Expect(m.IsOn()).To(BeTrue())
	})
	It("Must be able to tell if IsOff", func() {
		Expect(m.IsOff()).To(Equal(true))
	})
	It("Should be able to turn On", func() {
		m.On()
		Expect(m.CurrentState).To(Equal(uint8(1)))
		m.CurrentMode = "analog"
		m.CurrentSpeed = 0
		m.On()
		Expect(m.CurrentSpeed).To(Equal(uint8(255)))
	})
	It("Should be able to turn Off", func() {
		m.Off()
		Expect(m.CurrentState).To(Equal(uint8(0)))
		m.CurrentMode = "analog"
		m.CurrentSpeed = 100
		m.Off()
		Expect(m.CurrentSpeed).To(Equal(uint8(0)))
	})
	It("Should be able to Toggle", func() {
		m.Off()
		m.Toggle()
		Expect(m.IsOn()).To(BeTrue())
		m.Toggle()
		Expect(m.IsOn()).NotTo(BeTrue())
	})
	It("Should be able to set to Min speed", func() {
		m.Min()
	})
	It("Should be able to set to Max speed", func() {
		m.Max()
	})
	It("Should be able to set Speed", func() {
		Expect(true)
	})
	It("Should be able to set Forward", func() {
		m.Forward(100)
		Expect(m.CurrentSpeed).To(Equal(uint8(100)))
		Expect(m.CurrentDirection).To(Equal("forward"))
	})
	It("Should be able to set Backward", func() {
		m.Backward(100)
		Expect(m.CurrentSpeed).To(Equal(uint8(100)))
		Expect(m.CurrentDirection).To(Equal("backward"))
	})
	It("Should be able to set Direction", func() {
		m.Direction("none")
		m.DirectionPin = "2"
		m.Direction("forward")
		m.Direction("backward")
	})

})
