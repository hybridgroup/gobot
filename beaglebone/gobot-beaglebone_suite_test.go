package gobotBeaglebone

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotBeaglebone(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot-Beaglebone Suite")
}
