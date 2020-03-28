package configs

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("Configs", func() {

	Describe("Reading external configuration file", func() {

		Context("Reading a valid configuration", func() {

			It("should return a populated viper struct", func() {
				viper, err := readConfig("valid", "./testdata")
				Expect(err).To(BeNil())
				Expect((*viper).GetString("flowit.version")).To(Equal("0.1"))
			})

		})

		Context("Reading an invalid configuration", func() {

			It("should return an informative error", func() {
				viper, err := readConfig("incorrect-format", "./testdata")
				Expect(err).To(Not(BeNil()))
				Expect(errors.Cause(err).Error()).To(MatchRegexp("While parsing config: yaml: line [0-9]+"))
				Expect(viper).To(BeNil())
			})

		})

		Context("Reading a non existent configuration", func() {

			It("should return an informative error", func() {
				viper, err := readConfig("non-existent", "./testdata")
				Expect(err).To(Not(BeNil()))
				Expect(errors.Cause(err).Error()).To(ContainSubstring("Not Found"))
				Expect(viper).To(BeNil())
			})

		})

		Context("Reading a configuration with no read permissions", func() {

			It("should return an informative error", func() {
				fileName := "./testdata/valid.yaml"
				os.Chmod(fileName, 0000)
				defer os.Chmod(fileName, 0644)
				viper, err := readConfig("valid", "./testdata")
				Expect(err).To(Not(BeNil()))
				Expect(errors.Cause(err).Error()).To(ContainSubstring("permission denied"))
				Expect(viper).To(BeNil())
			})

		})

	})
})
