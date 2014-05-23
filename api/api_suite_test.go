package api

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

func TestApi(t *testing.T) {
	log.SetOutput(new(null))
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api Suite")
}
