package gobotGPIO

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Led", func() {
	var (
		adaptor TestAdaptor
		driver  *Led
	)

	BeforeEach(func() {
		driver = NewLed(adaptor)
		driver.Pin = "1"
	})

	It("Must be able to tell if IsOn", func() {
		Expect(driver.IsOn()).NotTo(BeTrue())
	})

	It("Must be able to tell if IsOff", func() {
		Expect(driver.IsOff()).To(BeTrue())
	})

	It("Should be able to turn On", func() {
		Expect(driver.On()).To(BeTrue())
		Expect(driver.IsOn()).To(BeTrue())
	})

	It("Should be able to turn Off", func() {
		Expect(driver.Off()).To(BeTrue())
		Expect(driver.IsOff()).To(BeTrue())
	})

	It("Should be able to Toggle", func() {
		driver.Off()
		driver.Toggle()
		Expect(driver.IsOn()).To(BeTrue())
		driver.Toggle()
		Expect(driver.IsOff()).To(BeTrue())
	})

	It("Should be able to set Brightness", func() {
		driver.Brightness(150)
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
