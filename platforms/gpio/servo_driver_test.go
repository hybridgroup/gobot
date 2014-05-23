package gpio

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Servo", func() {
	var (
		t TestAdaptor
		s *ServoDriver
	)

	BeforeEach(func() {
		s = NewServoDriver(t, "bot", "1")
	})

	It("Should be able to Move", func() {
		s.Move(100)
		Expect(s.CurrentAngle).To(Equal(uint8(100)))
	})

	It("Should be able to move to Min", func() {
		s.Min()
		Expect(s.CurrentAngle).To(Equal(uint8(0)))
	})

	It("Should be able to move to Max", func() {
		s.Max()
		Expect(s.CurrentAngle).To(Equal(uint8(180)))
	})

	It("Should be able to move to Center", func() {
		s.Center()
		Expect(s.CurrentAngle).To(Equal(uint8(90)))
	})

	It("Should be able to move to init servo", func() {
		s.InitServo()
	})

	It("Must be able to Start", func() {
		Expect(s.Start()).To(Equal(true))
	})
	It("Must be able to Init", func() {
		Expect(s.Init()).To(Equal(true))
	})
	It("Must be able to Halt", func() {
		Expect(s.Halt()).To(Equal(true))
	})
})
