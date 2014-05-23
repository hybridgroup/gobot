package i2c

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Wiichuck", func() {
	var (
		t TestAdaptor
		w *WiichuckDriver
	)

	BeforeEach(func() {
		w = NewWiichuckDriver(t, "bot")
	})

	PIt("Must be able to Start", func() {
		Expect(w.Start()).To(Equal(true))
	})
})
