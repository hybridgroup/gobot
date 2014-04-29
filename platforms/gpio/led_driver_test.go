package gpio

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Led", func() {
	var (
		t TestAdaptor
		l *LedDriver
	)

	BeforeEach(func() {
		l = NewLedDriver(t)
		l.Pin = "1"
	})

	It("Must be able to tell if IsOn", func() {
		Expect(l.IsOn()).NotTo(BeTrue())
	})

	It("Must be able to tell if IsOff", func() {
		Expect(l.IsOff()).To(BeTrue())
	})

	It("Should be able to turn On", func() {
		Expect(l.On()).To(BeTrue())
		Expect(l.IsOn()).To(BeTrue())
	})

	It("Should be able to turn Off", func() {
		Expect(l.Off()).To(BeTrue())
		Expect(l.IsOff()).To(BeTrue())
	})

	It("Should be able to Toggle", func() {
		l.Off()
		l.Toggle()
		Expect(l.IsOn()).To(BeTrue())
		l.Toggle()
		Expect(l.IsOff()).To(BeTrue())
	})

	It("Should be able to set Brightness", func() {
		l.Brightness(150)
	})

	It("Must be able to Start", func() {
		Expect(l.Start()).To(Equal(true))
	})
	It("Must be able to Init", func() {
		Expect(l.Init()).To(Equal(true))
	})
	It("Must be able to Halt", func() {
		Expect(l.Halt()).To(Equal(true))
	})
})
