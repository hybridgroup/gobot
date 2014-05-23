package i2c

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("BlinkM", func() {
	var (
		t TestAdaptor
		b *BlinkMDriver
	)

	BeforeEach(func() {
		b = NewBlinkMDriver(t, "bot")
	})

	It("Must be able to Start", func() {
		Expect(b.Start()).To(Equal(true))
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
