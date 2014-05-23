package opencv

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Window", func() {
	var (
		w *WindowDriver
	)

	BeforeEach(func() {
		w = NewWindowDriver("bot")
	})

	PIt("Must be able to Start", func() {
		Expect(w.Start()).To(Equal(true))
	})
	PIt("Must be able to Init", func() {
		Expect(w.Init()).To(Equal(true))
	})
	PIt("Must be able to Halt", func() {
		Expect(w.Halt()).To(Equal(true))
	})
})
