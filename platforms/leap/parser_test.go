package leap

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Parser", func() {
	a := NewLeapMotionAdaptor()
	d := NewLeapMotionDriver(a)

	Describe("#ParseLeapFrame", func() {
		It("Takes an array of bytes and extracts Frames", func() {
			file, err := ioutil.ReadFile("./test/support/example_frame.json")
			Expect(err != nil)
			parsedFrame := d.ParseFrame(file)
			Expect(parsedFrame.Hands != nil)
			Expect(parsedFrame.Pointables != nil)
			Expect(parsedFrame.Gestures != nil)
		})

		It("Returns an empty Frame if passed non-Leap bytes", func() {
			parsedFrame := d.ParseFrame([]byte{})
			Expect(parsedFrame.Timestamp == 0)
		})
	})
})
