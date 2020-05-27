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
				err := config.LoadConfiguration("valid", "./testdata")
				Expect(err).To(BeNil())
				Expect(config.GetVersion()).To(Equal("0.1"))
				Expect(config.GetConfig().Shell).To(Equal("/usr/bin/env bash"))
				/* #gomnd */
				Expect(config.GetVariables()["gerrit-port"]).To(Equal(float64(29418)))
				Expect(config.GetBranches()[0].ID).To(Equal("master"))
				Expect(config.GetWorkflows()[0].Stages[0].Actions[0]).
					To(Equal("git checkout master"))
			})

		})

		Context("Processing an invalid configuration", func() {

			It("should return a descriptive error", func() {
				err := config.LoadConfiguration("incorrect-types", "./testdata")
				Expect(err).To(Not(BeNil()))
				Expect(errors.Cause(err).Error()).To(MatchRegexp("[0-9]+ error\\(s\\) decoding:"))
			})

		})

	})
})
