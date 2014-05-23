package joystick

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JoystickDriver", func() {
	var (
		d *JoystickDriver
	)

	BeforeEach(func() {
		d = NewJoystickDriver(NewJoystickAdaptor("bot"), "bot", "/dev/null")
	})

	PIt("Must be able to Start", func() {
		Expect(d.Start()).To(Equal(true))
	})
	PIt("Must be able to Init", func() {
		Expect(d.Init()).To(Equal(true))
	})
	PIt("Must be able to Halt", func() {
		Expect(d.Halt()).To(Equal(true))
	})
})
