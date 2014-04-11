package gobot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"testing"
)

func TestGobot(t *testing.T) {
	log.SetOutput(new(null))
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot Suite")
}
