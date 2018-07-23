package co_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCogo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cogo Suite")
}
