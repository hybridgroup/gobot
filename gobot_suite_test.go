package gobot

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"log"
	"testing"
)

type null struct{}

func (null) Write(p []byte) (int, error) {
	return len(p), nil
}

func TestGobot(t *testing.T) {
	log.SetOutput(new(null))
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gobot Suite")
}
