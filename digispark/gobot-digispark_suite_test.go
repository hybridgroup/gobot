package gobotDigispark

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotDigispark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot-Digispark Suite")
}
