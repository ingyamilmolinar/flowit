package configs

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configs", func() {

	Describe("Reading external configuration file", func() {

		Context("Reading a correct configuration", func() {

			It("should return a populated viper struct", func() {
				viper, err := readConfig("valid", "./testdata")
				Expect(err).To(BeNil())
				Expect((*viper).GetString("flowit.version")).To(Equal("0.1"))
			})

		})

	})
})
