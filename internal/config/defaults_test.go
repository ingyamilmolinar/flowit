package config

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	Describe("Setting defaults on a configuration file", func() {

		It("should return defaults only for nil values", func() {
			config := rawConfig{
				CheckpointExecution: nil,
				Shell:               nil,
			}
			mainDefinition := rawMainDefinition{
				Config: &config,
			}
			workflowDefinition := rawWorkflowDefinition{
				Flowit: &mainDefinition,
			}
			setDefaults(&workflowDefinition)
			Expect(*workflowDefinition.Flowit.Config.CheckpointExecution).To(Equal(true))
			Expect(*workflowDefinition.Flowit.Config.Shell).To(Equal(os.Getenv("SHELL")))

			shell := "bash"
			config = rawConfig{
				CheckpointExecution: nil,
				Shell:               &shell,
			}
			mainDefinition = rawMainDefinition{
				Config: &config,
			}
			workflowDefinition = rawWorkflowDefinition{
				Flowit: &mainDefinition,
			}
			setDefaults(&workflowDefinition)
			Expect(*workflowDefinition.Flowit.Config.CheckpointExecution).To(Equal(true))
			Expect(*workflowDefinition.Flowit.Config.Shell).To(Equal(shell))
		})

	})

})
