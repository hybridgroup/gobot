package i2c

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HMC6352", func() {
	var (
		t TestAdaptor
		h *HMC6352Driver
	)

	BeforeEach(func() {
		h = NewHMC6352Driver(t)
	})

	It("Must be able to Start", func() {
		Expect(h.Start()).To(Equal(true))
	})
})
