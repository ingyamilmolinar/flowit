package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	Describe("Validating configuration", func() {

		Context("Validating a valid configuration", func() {

			It("should return a nil error", func() {

				var flowit rawFlowitConfig

				version := "0.1"
				branchID := "master"
				branchName := "master"
				branchEternal := true
				branchProtected := true
				workflowBranch := rawBranch{
					ID:        &branchID,
					Name:      &branchName,
					Eternal:   &branchEternal,
					Protected: &branchProtected,
				}
				startStage := "start"
				finishStage := "finish"
				action := "action"
				actions := []*string{
					&action,
				}
				workflowStage := rawStage{
					Start:   &startStage,
					Finish:  &finishStage,
					Actions: actions,
				}
				workflow := rawWorkflow{
					Branches: []*rawBranch{
						&workflowBranch,
					},
					Stages: []*rawBranchType{
						{
							"dev": []*rawStage{
								&workflowStage,
							},
						},
					},
				}

				flowit.Version = &version
				flowit.Workflow = &workflow
				err := validateConfig(&flowit)
				Expect(err).To(BeNil())

			})

		})

	})
})
