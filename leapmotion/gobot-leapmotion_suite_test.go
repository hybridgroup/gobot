package gobotLeap

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotLeapmotion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot-Leapmotion Suite")
}
