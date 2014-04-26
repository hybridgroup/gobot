package gobotI2C

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BlinkM", func() {
	var (
		someAdaptor TestAdaptor
		someDriver  *BlinkM
	)

	BeforeEach(func() {
		someDriver = NewBlinkM(someAdaptor)
	})

	It("Must be able to Start", func() {
		Expect(someDriver.Start()).To(Equal(true))
	})

	PIt("Should be able to set Rgb", func() {
		Expect(true)
	})

	PIt("Should be able to Fade", func() {
		Expect(true)
	})

	PIt("Should be able to get FirmwareVersion", func() {
		Expect(true)
	})

	PIt("Should be able to set Color", func() {
		Expect(true)
	})
})
