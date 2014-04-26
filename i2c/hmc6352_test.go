package gobotI2C

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HMC6352", func() {
	var (
		someAdaptor TestAdaptor
		someDriver  *HMC6352
	)

	BeforeEach(func() {
		someDriver = NewHMC6352(someAdaptor)
		someDriver.Interval = "1s"
	})

	It("Must be able to Start", func() {
		Expect(someDriver.Start()).To(Equal(true))
	})
})
