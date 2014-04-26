package gobotNeurosky

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobotNeurosky(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot-Neurosky Suite")
}
