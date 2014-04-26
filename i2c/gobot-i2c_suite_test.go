package gobotI2C

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotGpio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot-Gpio Suite")
}
