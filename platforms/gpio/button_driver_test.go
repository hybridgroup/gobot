package gpio

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Button", func() {
	var (
		t TestAdaptor
		b *ButtonDriver
	)

	BeforeEach(func() {
		b = NewButtonDriver(t, "bot", "1")
	})

	It("Must be able to readState", func() {
		Expect(b.readState()).To(Equal(1))
	})

	It("Must update on button state change to on", func() {
		b.update(1)
		Expect(b.Active).To(Equal(true))
	})

	It("Must update on button state change to off", func() {
		b.update(0)
		Expect(b.Active).To(Equal(false))
	})

	It("Must be able to Start", func() {
		Expect(b.Start()).To(Equal(true))
	})
	It("Must be able to Init", func() {
		Expect(b.Init()).To(Equal(true))
	})
	It("Must be able to Halt", func() {
		Expect(b.Halt()).To(Equal(true))
	})
})
