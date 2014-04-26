package gobotArdrone

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotArdrone(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot-Ardrone Suite")
}
