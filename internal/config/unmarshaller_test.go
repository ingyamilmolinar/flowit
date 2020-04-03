package config

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var _ = Describe("Config", func() {

	Describe("Unmarshalling external configuration file", func() {

		Context("Unmarshalling a valid configuration", func() {

			It("should return a populated viper struct", func() {
				viper := viper.New()
				viper.SetConfigFile("./testdata/valid.yaml")
				if err := viper.ReadInConfig(); err != nil {
					Fail(fmt.Sprintf("Error reading config %+v", err))
				}

				flowit, err := unmarshallConfig(viper)
				Expect(err).To(BeNil())
				Expect(*flowit.Version).To(Equal("0.1"))
			})

			It("should set nil on missing sections", func() {
				viper := viper.New()
				viper.SetConfigFile("./testdata/missing-sections.yaml")
				if err := viper.ReadInConfig(); err != nil {
					Fail(fmt.Sprintf("Error reading config %+v", err))
				}
				flowit, err := unmarshallConfig(viper)
				Expect(err).To(BeNil())
				Expect(flowit.Config).To(BeNil())
				Expect(flowit.Variables).To(BeNil())
			})

		})

		Context("Unmarshalling an invalid configuration", func() {

			It("should return an informative error for incorrect types", func() {
				viper := viper.New()
				viper.SetConfigFile("./testdata/incorrect-types.yaml")
				if err := viper.ReadInConfig(); err != nil {
					Fail(fmt.Sprintf("Error reading config %+v", err))
				}
				flowit, err := unmarshallConfig(viper)
				Expect(err).To(Not(BeNil()))
				Expect(errors.Cause(err).Error()).To(ContainSubstring("abort-on-failed-action"))
				Expect(errors.Cause(err).Error()).To(ContainSubstring("Config.Shell"))
				Expect(errors.Cause(err).Error()).To(ContainSubstring("Workflow.Branches"))
				Expect(flowit).To(BeNil())
			})

		})

	})
})
