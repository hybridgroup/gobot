package gobot_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGobot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot Suite")
}
