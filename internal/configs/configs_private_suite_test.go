package configs

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfigs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configs Private Suite")
}
