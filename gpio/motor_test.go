package gobotGPIO

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Motor", func() {
	var (
		adaptor TestAdaptor
		driver  *Motor
	)

	BeforeEach(func() {
		driver = NewMotor(adaptor)
		driver.Pin = "1"
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
	It("Must be able to tell if IsOn", func() {
		driver.CurrentState = 1
		Expect(driver.IsOn()).To(BeTrue())
		driver.CurrentMode = "analog"
		driver.CurrentSpeed = 100
		Expect(driver.IsOn()).To(BeTrue())
	})
	It("Must be able to tell if IsOff", func() {
		Expect(driver.IsOff()).To(Equal(true))
	})
	It("Should be able to turn On", func() {
		driver.On()
		Expect(driver.CurrentState).To(Equal(uint8(1)))
		driver.CurrentMode = "analog"
		driver.CurrentSpeed = 0
		driver.On()
		Expect(driver.CurrentSpeed).To(Equal(uint8(255)))
	})
	It("Should be able to turn Off", func() {
		driver.Off()
		Expect(driver.CurrentState).To(Equal(uint8(0)))
		driver.CurrentMode = "analog"
		driver.CurrentSpeed = 100
		driver.Off()
		Expect(driver.CurrentSpeed).To(Equal(uint8(0)))
	})
	It("Should be able to Toggle", func() {
		driver.Off()
		driver.Toggle()
		Expect(driver.IsOn()).To(BeTrue())
		driver.Toggle()
		Expect(driver.IsOn()).NotTo(BeTrue())
	})
	It("Should be able to set to Min speed", func() {
		driver.Min()
	})
	It("Should be able to set to Max speed", func() {
		driver.Max()
	})
	It("Should be able to set Speed", func() {
		Expect(true)
	})
	It("Should be able to set Forward", func() {
		driver.Forward(100)
		Expect(driver.CurrentSpeed).To(Equal(uint8(100)))
		Expect(driver.CurrentDirection).To(Equal("forward"))
	})
	It("Should be able to set Backward", func() {
		driver.Backward(100)
		Expect(driver.CurrentSpeed).To(Equal(uint8(100)))
		Expect(driver.CurrentDirection).To(Equal("backward"))
	})
	It("Should be able to set Direction", func() {
		driver.Direction("none")
		driver.DirectionPin = "2"
		driver.Direction("forward")
		driver.Direction("backward")
	})

})
