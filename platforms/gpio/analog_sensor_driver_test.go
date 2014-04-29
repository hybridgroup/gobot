package gpio

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Analog-Sensor", func() {
	var (
		t TestAdaptor
		a *AnalogSensorDriver
	)

	BeforeEach(func() {
		a = NewAnalogSensor(t)
		a.Pin = "1"
	})

	It("Must be able to Read", func() {
		Expect(a.Read()).To(Equal(99))
	})
	It("Must be able to Start", func() {
		Expect(a.Start()).To(Equal(true))
	})
	It("Must be able to Halt", func() {
		Expect(a.Halt()).To(Equal(true))
	})
	It("Must be able to Init", func() {
		Expect(a.Init()).To(Equal(true))
	})
})
