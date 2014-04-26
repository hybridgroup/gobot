package gobotJoystick

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotJoystick(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot-Joystick Suite")
}
