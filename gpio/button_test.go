package gobotGPIO

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Button", func() {
	var (
		someAdaptor TestAdaptor
		someDriver  *Button
	)

	BeforeEach(func() {
		someDriver = NewButton(someAdaptor)
		someDriver.Pin = "1"
	})

	It("Must be able to readState", func() {
		Expect(someDriver.readState()).To(Equal(1))
	})

	It("Must update on button state change to on", func() {
		someDriver.update(1)
		Expect(someDriver.Active).To(Equal(true))
	})

	It("Must update on button state change to off", func() {
		someDriver.update(0)
		Expect(someDriver.Active).To(Equal(false))
	})

	It("Must be able to Start", func() {
		Expect(someDriver.Start()).To(Equal(true))
	})
	It("Must be able to Init", func() {
		Expect(someDriver.Init()).To(Equal(true))
	})
	It("Must be able to Halt", func() {
		Expect(someDriver.Halt()).To(Equal(true))
	})
})
