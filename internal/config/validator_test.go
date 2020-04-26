package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/yamil-rivera/flowit/internal/utils"
)

var _ = Describe("Config", func() {

	Describe("Validating a valid configuration", func() {

		It("should return a nil error for missing optional values", func() {

			rawFlowit := validConfigJustMandatoryFields()
			err := validateConfig(&rawFlowit)
			Expect(err).To(BeNil())

		})

		It("should return a nil error for having all optional values", func() {

			flowit := validConfigWithOptionalFields()
			rawFlowit := rawify(&flowit)

			err := validateConfig(rawFlowit)
			Expect(err).To(BeNil())

		})

	})

	Describe("Validating an invalid configuration", func() {

		Context("Validating version", func() {

			It("should return a descriptive error for incorrect version", func() {

				flowit := validConfigWithOptionalFields()
				flowit.Version = ".1"
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Version:"))

			})

		})

		Context("Validating config", func() {

			It("should return a descriptive error for incorrect config", func() {

				flowit := validConfigWithOptionalFields()
				flowit.Config.Shell = "/nonexistent/shell"
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Config:"))

			})

		})

		Context("Validating branches", func() {

			It("should return a descriptive error for invalid branch ID", func() {

				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Branches[0].ID = "$<variable>"
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("ID:"))

				flowit = validConfigWithOptionalFields()
				flowit.Workflow.Branches[0].ID = " "
				rawFlowit = rawify(&flowit)

				err = validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("ID:"))

			})

			It("should return a descriptive error for invalid branch name", func() {

				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Branches[0].Name = " "
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Name:"))

			})

			It("should return a descriptive error for a wrongly defined transition", func() {

				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Branches[0].ID = "master"
				flowit.Workflow.Branches[0].Eternal = true
				flowit.Workflow.Branches[1].ID = "feature"
				flowit.Workflow.Branches[1].Eternal = false
				flowit.Workflow.Branches[1].Protected = false
				// transitions should not be defined on an eternal branch
				flowit.Workflow.Branches[0].Transitions = []transition{
					{
						From: "feature",
						To: []string{
							"feature:local",
						},
					},
				}
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))

			})

			It("should return a descriptive error for an undefined transition", func() {

				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Branches[1].ID = "feature"
				flowit.Workflow.Branches[1].Eternal = false
				flowit.Workflow.Branches[1].Protected = false
				// transitions should be defined on a non eternal branch
				flowit.Workflow.Branches[1].Transitions = []transition{}
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))

			})

			It("should return a descriptive error for an incorrect transition target", func() {

				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Branches[0].ID = "master"
				flowit.Workflow.Branches[0].Eternal = true
				flowit.Workflow.Branches[0].Protected = true
				flowit.Workflow.Branches[1].ID = "feature"
				flowit.Workflow.Branches[1].Eternal = false
				flowit.Workflow.Branches[1].Transitions = []transition{
					{
						From: "master",
						To: []string{
							// Forgot :remote or :local
							"master",
						},
					},
				}
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))

			})

			It("should return a descriptive error for an incorrect transition branch ID", func() {

				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Branches[0].ID = "master"
				flowit.Workflow.Branches[0].Eternal = true
				flowit.Workflow.Branches[0].Protected = true
				flowit.Workflow.Branches[1].ID = "feature"
				flowit.Workflow.Branches[1].Eternal = false
				flowit.Workflow.Branches[1].Transitions = []transition{
					{
						From: "invalid-branch-id",
						To: []string{
							"master:local",
						},
					},
				}
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))

				flowit = validConfigWithOptionalFields()
				flowit.Workflow.Branches[0].ID = "master"
				flowit.Workflow.Branches[0].Eternal = true
				flowit.Workflow.Branches[0].Protected = true
				flowit.Workflow.Branches[1].ID = "feature"
				flowit.Workflow.Branches[1].Eternal = false
				flowit.Workflow.Branches[1].Transitions = []transition{
					{
						From: "master",
						To: []string{
							"invalid-branch-id:local",
						},
					},
				}
				rawFlowit = rawify(&flowit)

				err = validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))

			})

			// TODO: Check for duplicated repeated transitions

		})

		Context("Validating tags", func() {

			It("should return a descriptive error for an invalid tag ID", func() {

				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Tags[0].ID = " "
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Tags:"))
				Expect(err.Error()).To(ContainSubstring("ID:"))

			})

			It("should return a descriptive error for an invalid tag format", func() {

				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Tags[0].Format = " "
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Tags:"))
				Expect(err.Error()).To(ContainSubstring("Format:"))

			})

			// TODO: Check for repeated tag workflows & stages & branches
			It("should return a descriptive error for an invalid tag workflow", func() {

				flowit := validConfigWithOptionalFields()
				var existingStage string
				for k := range flowit.Workflow.Stages[0] {
					existingStage = flowit.Workflow.Stages[0][k][0]["id"].(string)
					break
				}
				flowit.Workflow.Tags[0].Stages["missing-workflow"] = []string{
					existingStage,
				}
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Tags:"))
				Expect(err.Error()).To(ContainSubstring("Stages:"))
			})

			It("should return a descriptive error for an invalid stage on a valid workflow", func() {
				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Stages[0]["my-workflow"] = []stage{
					{
						"id":   "my-stage",
						"args": "arg1 arg2",
						"conditions": []string{
							"condition-1",
						},
						"actions": []string{
							"action-1",
						},
					},
				}
				flowit.Workflow.Tags[0].Stages["my-workflow"] = []string{
					"non-existant-stage",
				}
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Tags:"))
				Expect(err.Error()).To(ContainSubstring("Stages:"))
			})

			It("should return a descriptive error for an invalid tag branch", func() {
				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Tags[0].Branches = []string{
					"non-existant-branch",
				}
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Tags:"))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
			})

		})

		Context("Validating stages", func() {

			It("should return a descriptive error for a non existant stage ID", func() {
				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Stages[0] = workflowType{
					"workflow-type": []stage{
						{
							"actions": []string{
								"action1",
							},
						},
					},
				}
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Stages:"))
				Expect(err.Error()).To(ContainSubstring("id"))
			})

			It("should return a descriptive error for a non existant stage actions", func() {
				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Stages[0] = workflowType{
					"workflow-type": []stage{
						{
							"id": "my-id",
						},
					},
				}
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Stages:"))
				Expect(err.Error()).To(ContainSubstring("actions"))
			})

			It("should return a descriptive error for a non supported stage property", func() {
				flowit := validConfigWithOptionalFields()
				flowit.Workflow.Stages[0] = workflowType{
					"workflow-type": []stage{
						{
							"id":   "my-id",
							"args": "arg1 arg2",
							"conditions": []string{
								"condition1",
							},
							"actions": []string{
								"action1",
							},
							"non-supported-property": "whatever",
						},
					},
				}
				rawFlowit := rawify(&flowit)

				err := validateConfig(rawFlowit)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Stages:"))
				Expect(err.Error()).To(ContainSubstring("non-supported-property"))
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
