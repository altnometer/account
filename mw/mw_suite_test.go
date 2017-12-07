package mw_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMw(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mw Suite")
}
