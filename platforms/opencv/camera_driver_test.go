package opencv

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Camera", func() {
	var (
		c *CameraDriver
	)

	BeforeEach(func() {
		c = NewCameraDriver("bot", 0)
	})

	PIt("Must be able to Start", func() {
		Expect(c.Start()).To(Equal(true))
	})
	PIt("Must be able to Init", func() {
		Expect(c.Init()).To(Equal(true))
	})
	PIt("Must be able to Halt", func() {
		Expect(c.Halt()).To(Equal(true))
	})
})
