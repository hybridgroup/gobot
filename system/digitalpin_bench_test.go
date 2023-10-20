package system

import (
	"testing"
)

func BenchmarkDigitalRead(b *testing.B) {
	a := NewAccesser()
	mockPaths := []string{
		"/sys/class/gpio/export",
		"/sys/class/gpio/unexport",
		"/sys/class/gpio/gpio10/value",
		"/sys/class/gpio/gpio10/direction",
	}
	a.UseMockFilesystem(mockPaths)

	pin := a.NewDigitalPin("", 10)
	_ = pin.Write(1)

	for i := 0; i < b.N; i++ {
		_, _ = pin.Read()
	}
}
