package gobotOpencv

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Camera", func() {
	var (
		driver *Camera
	)

	BeforeEach(func() {
		driver = NewCamera(new(Opencv))
	})

	PIt("Must be able to Start", func() {
		Expect(driver.Start()).To(Equal(true))
	})
	PIt("Must be able to Init", func() {
		Expect(driver.Init()).To(Equal(true))
	})
	PIt("Must be able to Halt", func() {
		Expect(driver.Halt()).To(Equal(true))
	})
})
