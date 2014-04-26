package gobotGPIO

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Analog-Sensor", func() {
	var (
		someAdaptor TestAdaptor
		someDriver  *AnalogSensor
	)

	BeforeEach(func() {
		someDriver = NewAnalogSensor(someAdaptor)
		someDriver.Pin = "1"
	})

	It("Must be able to Read", func() {
		Expect(someDriver.Read()).To(Equal(99))
	})
	It("Must be able to Start", func() {
		Expect(someDriver.Start()).To(Equal(true))
	})
	It("Must be able to Halt", func() {
		Expect(someDriver.Halt()).To(Equal(true))
	})
	It("Must be able to Init", func() {
		Expect(someDriver.Init()).To(Equal(true))
	})
})
