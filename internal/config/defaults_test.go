package config

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {

	Describe("Setting defaults on a configuration file", func() {

		It("should return defaults only for nil values", func() {
			strict := true
			config := rawConfig{
				AbortOnFailedAction: nil,
				Strict:              &strict,
				Shell:               nil,
			}
			mainDefinition := rawMainDefinition{
				Config: &config,
			}
			workflowDefinition := rawWorkflowDefinition{
				Flowit: &mainDefinition,
			}
			setDefaults(&workflowDefinition)
			Expect(*workflowDefinition.Flowit.Config.AbortOnFailedAction).To(Equal(true))
			Expect(*workflowDefinition.Flowit.Config.Strict).To(Equal(strict))
			Expect(*workflowDefinition.Flowit.Config.Shell).To(Equal(os.Getenv("SHELL")))

			shell := "bash"
			config = rawConfig{
				AbortOnFailedAction: nil,
				Strict:              nil,
				Shell:               &shell,
			}
			mainDefinition = rawMainDefinition{
				Config: &config,
			}
			workflowDefinition = rawWorkflowDefinition{
				Flowit: &mainDefinition,
			}
			setDefaults(&workflowDefinition)
			Expect(*workflowDefinition.Flowit.Config.AbortOnFailedAction).To(Equal(true))
			Expect(*workflowDefinition.Flowit.Config.Strict).To(Equal(false))
			Expect(*workflowDefinition.Flowit.Config.Shell).To(Equal(shell))
		})

		It("should set tag stages to all possible stages", func() {
			tag := rawTag{}
			tags := []*rawTag{
				&tag,
			}
			stage1ID := "stage1 ID"
			stage1 := rawStage{
				ID: &stage1ID,
			}
			stage2ID := "stage2 ID"
			stage2 := rawStage{
				ID: &stage2ID,
			}
			workflowID1 := "workflow1"
			workflowID2 := "workflow2"
			workflows := []*rawWorkflow{
				{
					ID: &workflowID1,
					Stages: []*rawStage{
						&stage1,
						&stage2,
					},
				},
				{
					ID: &workflowID2,
					Stages: []*rawStage{
						&stage2,
					},
				},
			}
			mainDefinition := rawMainDefinition{
				Tags:      tags,
				Workflows: workflows,
			}
			workflowDefinition := rawWorkflowDefinition{
				Flowit: &mainDefinition,
			}
			expectedStages := rawStages{
				"workflow1": []*string{
					&stage1ID,
					&stage2ID,
				},
				"workflow2": []*string{
					&stage2ID,
				},
			}
			setDefaults(&workflowDefinition)
			Expect(*workflowDefinition.Flowit.Tags[0].Stages).To(Equal(expectedStages))
		})

		It("should set tag stages to all possible stages", func() {
			tag := rawTag{}
			tags := []*rawTag{
				&tag,
			}
			branch1ID := "branch 1"
			branch2ID := "branch 2"
			branches := []*rawBranch{
				{
					ID: &branch1ID,
				},
				{
					ID: &branch2ID,
				},
			}
			mainDefinition := rawMainDefinition{
				Tags:     tags,
				Branches: branches,
			}
			workflowDefinition := rawWorkflowDefinition{
				Flowit: &mainDefinition,
			}
			expectedBranches := []*string{
				&branch1ID,
				&branch2ID,
			}
			setDefaults(&workflowDefinition)
			Expect(workflowDefinition.Flowit.Tags[0].Branches).To(Equal(expectedBranches))
		})

	})

})
