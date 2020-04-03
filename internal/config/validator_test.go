package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/yamil-rivera/flowit/internal/utils"
)

var _ = Describe("Config", func() {

	Describe("Validating configuration", func() {

		Context("Validating a valid configuration", func() {

			It("should return a nil error", func() {

				var flowit FlowitConfig

				flowit.Version = "0.1"
				flowit.Workflow.Branches = []branch{
					{
						ID:        "master",
						Name:      "master",
						Eternal:   true,
						Protected: true,
					},
				}
				flowit.Workflow.Stages = []workflowType{
					{
						"dev": []stage{
							{
								"id":      "start",
								"actions": []string{"action1", "action2"},
							},
						},
					},
				}

				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(BeNil())

			})

		})

	})
})

func rawify(config *FlowitConfig) *rawFlowitConfig {
	var rawConfig rawFlowitConfig
	if err := utils.DeepCopy(config, &rawConfig); err != nil {
		return nil
	}
	return &rawConfig
}
