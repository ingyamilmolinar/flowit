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
				cs, err := config.Load("./testdata/valid.yaml")
				Expect(err).To(BeNil())
				Expect(cs.Flowit.Version).To(Equal("0.1"))
				Expect(cs.Flowit.Config.Shell).To(Equal("/usr/bin/env bash"))
				/* #gomnd */
				Expect(cs.Flowit.Variables["gerrit-port"]).To(Equal(float64(29418)))
				Expect(cs.Flowit.Workflows[0].Stages[0].Actions[0]).
					To(Equal("git checkout master"))
			})

		})

		Context("Processing an invalid configuration", func() {

			It("should return a descriptive error", func() {
				_, err := config.Load("./testdata/incorrect-types.yaml")
				Expect(err).To(Not(BeNil()))
				Expect(errors.Cause(err).Error()).To(MatchRegexp("[0-9]+ error\\(s\\) decoding:"))
			})

		})

	})
})
