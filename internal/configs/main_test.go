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
				Expect(*flowit.Config.Shell).To(Equal("/usr/bin/env bash"))
				Expect((*flowit.Variables)["gerrit-port"]).To(Equal(29418))
				Expect(*flowit.Workflow.Branches[0].ID).To(Equal("master"))
				Expect((*flowit.Workflow.Stages[1])["actions"].([]interface{})[0].(string)).To(Equal("git checkout master"))
			})

		})

	})
})
