package gobotJoystick

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JoystickDriver", func() {
	var (
		driver *JoystickDriver
	)

	BeforeEach(func() {
		driver = NewJoystick(new(JoystickAdaptor))
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
