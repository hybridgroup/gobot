package gobotGPIO

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Servo", func() {
	var (
		adaptor TestAdaptor
		driver  *Servo
	)

	BeforeEach(func() {
		driver = NewServo(adaptor)
		driver.Pin = "1"
	})

	It("Should be able to Move", func() {
		driver.Move(100)
		Expect(driver.CurrentAngle).To(Equal(uint8(100)))
	})

	It("Should be able to move to Min", func() {
		driver.Min()
		Expect(driver.CurrentAngle).To(Equal(uint8(0)))
	})

	It("Should be able to move to Max", func() {
		driver.Max()
		Expect(driver.CurrentAngle).To(Equal(uint8(180)))
	})

	It("Should be able to move to Center", func() {
		driver.Center()
		Expect(driver.CurrentAngle).To(Equal(uint8(90)))
	})

	It("Should be able to move to init servo", func() {
		driver.InitServo()
	})

	It("Must be able to Start", func() {
		Expect(driver.Start()).To(Equal(true))
	})
	It("Must be able to Init", func() {
		Expect(driver.Init()).To(Equal(true))
	})
	It("Must be able to Halt", func() {
		Expect(driver.Halt()).To(Equal(true))
	})
})
