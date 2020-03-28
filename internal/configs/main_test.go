package configs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/yamil-rivera/flowit/internal/configs"
)

var _ = Describe("Configs", func() {

	Describe("Processing external configuration file", func() {

		Context("Processing a correct configuration", func() {

			It("should return a populated Flowit structure", func() {
				flowitptr, err := configs.ProcessFlowitConfig("valid", "./testdata")
				flowit := (*flowitptr)
				Expect(err).To(BeNil())
				Expect(*flowit.Version).To(Equal("0.1"))
			})

		})

	})
})
