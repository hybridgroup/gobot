package gobotGPIO

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DirectPin", func() {
	var (
		adaptor TestAdaptor
		driver  *DirectPin
	)

	BeforeEach(func() {
		driver = NewDirectPin(adaptor)
		driver.Pin = "1"
	})

	It("Should be able to DigitalRead", func() {
		Expect(driver.DigitalRead()).To(Equal(1))
	})

	It("Should be able to DigitalWrite", func() {
		driver.DigitalWrite(1)
	})

	It("Should be able to AnalogRead", func() {
		Expect(driver.AnalogRead()).To(Equal(99))
	})

	It("Should be able to AnalogWrite", func() {
		driver.AnalogWrite(100)
	})

	It("Should be able to PwmWrite", func() {
		driver.PwmWrite(100)
	})

	It("Should be able to ServoWrite", func() {
		driver.ServoWrite(100)
	})

	It("Should be able to Start", func() {
		Expect(driver.Start()).To(BeTrue())
	})
	It("Should be able to Halt", func() {
		Expect(driver.Halt()).To(BeTrue())
	})
	It("Should be able to Init", func() {
		Expect(driver.Init()).To(BeTrue())
	})
})
