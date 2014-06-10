package pebble

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotPebble(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot-Pebble Suite")
}
