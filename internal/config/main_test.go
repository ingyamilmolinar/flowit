package config_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"

	"github.com/yamil-rivera/flowit/internal/config"
)

var _ = Describe("Config", func() {

	Describe("Processing external configuration file", func() {

		Context("Processing a valid configuration", func() {

			It("should return a populated Flowit structure", func() {
				workflowDefinitionPtr, err := config.ProcessFlowitConfig("valid", "./testdata")
				Expect(err).To(BeNil())
				workflowDefinition := (*workflowDefinitionPtr)
				Expect(workflowDefinition.Flowit.Version).To(Equal("0.1"))
				Expect(workflowDefinition.Flowit.Config.Shell).To(Equal("/usr/bin/env bash"))
				/* #gomnd */
				Expect(workflowDefinition.Flowit.Variables["gerrit-port"]).To(Equal(float64(29418)))
				Expect(workflowDefinition.Flowit.Branches[0].ID).To(Equal("master"))
				Expect(workflowDefinition.Flowit.Workflows[0]["development"][0].Actions[0]).
					To(Equal("git checkout master"))
			})

		})

		Context("Processing an invalid configuration", func() {

			It("should return a descriptive error", func() {
				workflowDefinition, err := config.ProcessFlowitConfig("incorrect-types", "./testdata")
				Expect(err).To(Not(BeNil()))
				Expect(errors.Cause(err).Error()).To(MatchRegexp("[0-9]+ error\\(s\\) decoding:"))
				Expect(workflowDefinition).To((BeNil()))
			})

		})

	})
})
