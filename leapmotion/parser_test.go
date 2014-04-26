package gobotLeap

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
)

var _ = Describe("Parser", func() {
	leapAdaptor := new(LeapAdaptor)
	leap := NewLeap(leapAdaptor)

	Describe("#ParseLeapFrame", func() {
		It("Takes an array of bytes and extracts Leap Frames", func() {
			file, err := ioutil.ReadFile("./test/support/example_frame.json")
			Expect(err != nil)
			parsedFrame := leap.ParseLeapFrame(file)
			Expect(parsedFrame.Hands != nil)
			Expect(parsedFrame.Pointables != nil)
			Expect(parsedFrame.Gestures != nil)
		})

		It("Returns an empty Leap Frame if passed non-Leap bytes", func() {
			parsedFrame := leap.ParseLeapFrame([]byte{})
			Expect(parsedFrame.Timestamp == 0)
		})
	})
})
