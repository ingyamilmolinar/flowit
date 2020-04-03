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
				flowitptr, err := config.ProcessFlowitConfig("valid", "./testdata")
				flowit := (*flowitptr)
				Expect(err).To(BeNil())
				Expect(flowit.Version).To(Equal("0.1"))
				Expect(flowit.Config.Shell).To(Equal("/usr/bin/env bash"))
				/* #gomnd */
				Expect(flowit.Variables["gerrit-port"]).To(Equal(29418))
				Expect(flowit.Workflow.Branches[0].ID).To(Equal("master"))
				Expect(flowit.Workflow.Stages[0]["development"][0]["actions"].([]interface{})[0].(string)).
					To(Equal("git checkout master"))
			})

		})

		Context("Processing an invalid configuration", func() {

			It("should return a descriptive error", func() {
				flowit, err := config.ProcessFlowitConfig("incorrect-types", "./testdata")
				Expect(err).To(Not(BeNil()))
				Expect(errors.Cause(err).Error()).To(MatchRegexp("[0-9]+ error\\(s\\) decoding:"))
				Expect(flowit).To((BeNil()))
			})

		})

	})
})
