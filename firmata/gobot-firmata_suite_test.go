package gobotFirmata

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotFirmata(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot-Firmata Suite")
}
