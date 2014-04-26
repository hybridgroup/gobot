package gobotSpark

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotSpark(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot-Spark Suite")
}
