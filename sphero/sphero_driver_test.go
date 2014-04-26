package gobotSphero

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SpheroDriver", func() {
	var (
		driver  *SpheroDriver
		adaptor *SpheroAdaptor
	)

	BeforeEach(func() {
		adaptor = new(SpheroAdaptor)
		adaptor.sp = sp{}
		driver = NewSphero(adaptor)
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
