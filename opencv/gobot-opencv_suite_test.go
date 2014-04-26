package gobotOpencv

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotOpencv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot-Opencv Suite")
}
