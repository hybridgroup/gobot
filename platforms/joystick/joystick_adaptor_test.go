package joystick

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JoystickAdaptor", func() {
	var (
		j *JoystickAdaptor
	)

	BeforeEach(func() {
		j = NewJoystickAdaptor("bot")
	})

	PIt("Must be able to Finalize", func() {
		Expect(j.Finalize()).To(Equal(true))
	})
	PIt("Must be able to Connect", func() {
		Expect(j.Connect()).To(Equal(true))
	})
})
