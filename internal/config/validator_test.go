package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/yamil-rivera/flowit/internal/utils"
)

var _ = Describe("Config", func() {

	Describe("Validating a valid configuration", func() {

		It("should return a nil error for missing optional values", func() {

			rawConfig := validConfigJustMandatoryFields()
			err := validateConfig(&rawConfig)
			Expect(err).To(BeNil())

		})

		It("should return a nil error for having all optional values", func() {

			config := validConfigWithOptionalFields()
			rawConfig := rawify(&config)

			err := validateConfig(rawConfig)
			Expect(err).To(BeNil())

		})

	})

	Describe("Validating an invalid configuration", func() {

		Context("Validating version", func() {

			It("should return a descriptive error for incorrect version", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Version = ".1"
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Version:"))

			})

		})

		Context("Validating config", func() {

			It("should return a descriptive error for incorrect config", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Config.Shell = "/nonexistent/shell"
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Config:"))

			})

		})

		Context("Validating branches", func() {

			It("should return a descriptive error for invalid branch ID", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Branches[0].ID = "$<variable>"
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("ID:"))

				config = validConfigWithOptionalFields()
				config.Flowit.Branches[0].ID = " "
				rawConfig = rawify(&config)

				err = validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("ID:"))

			})

			It("should return a descriptive error for invalid branch name", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Branches[0].Name = " "
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Name:"))

			})

			It("should return a descriptive error for a wrongly defined transition", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Branches[0].ID = "master"
				config.Flowit.Branches[0].Eternal = true
				config.Flowit.Branches[1].ID = "feature"
				config.Flowit.Branches[1].Eternal = false
				config.Flowit.Branches[1].Protected = false
				// transitions should not be defined on an eternal branch
				config.Flowit.Branches[0].Transitions = []transition{
					{
						From: "feature",
						To: []string{
							"feature:local",
						},
					},
				}
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))

			})

			It("should return a descriptive error for an undefined transition", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Branches[1].ID = "feature"
				config.Flowit.Branches[1].Eternal = false
				config.Flowit.Branches[1].Protected = false
				// transitions should be defined on a non eternal branch
				config.Flowit.Branches[1].Transitions = []transition{}
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))

			})

			It("should return a descriptive error for an incorrect transition target", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Branches[0].ID = "master"
				config.Flowit.Branches[0].Eternal = true
				config.Flowit.Branches[0].Protected = true
				config.Flowit.Branches[1].ID = "feature"
				config.Flowit.Branches[1].Eternal = false
				config.Flowit.Branches[1].Transitions = []transition{
					{
						From: "master",
						To: []string{
							// Forgot :remote or :local
							"master",
						},
					},
				}
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))

			})

			It("should return a descriptive error for an incorrect transition branch ID", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Branches[0].ID = "master"
				config.Flowit.Branches[0].Eternal = true
				config.Flowit.Branches[0].Protected = true
				config.Flowit.Branches[1].ID = "feature"
				config.Flowit.Branches[1].Eternal = false
				config.Flowit.Branches[1].Transitions = []transition{
					{
						From: "invalid-branch-id",
						To: []string{
							"master:local",
						},
					},
				}
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))

				config = validConfigWithOptionalFields()
				config.Flowit.Branches[0].ID = "master"
				config.Flowit.Branches[0].Eternal = true
				config.Flowit.Branches[0].Protected = true
				config.Flowit.Branches[1].ID = "feature"
				config.Flowit.Branches[1].Eternal = false
				config.Flowit.Branches[1].Transitions = []transition{
					{
						From: "master",
						To: []string{
							"invalid-branch-id:local",
						},
					},
				}
				rawConfig = rawify(&config)

				err = validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))

			})

			// TODO: Check for duplicated repeated transitions

		})

		Context("Validating tags", func() {

			It("should return a descriptive error for an invalid tag ID", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Tags[0].ID = " "
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Tags:"))
				Expect(err.Error()).To(ContainSubstring("ID:"))

			})

			It("should return a descriptive error for an invalid tag format", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Tags[0].Format = " "
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Tags:"))
				Expect(err.Error()).To(ContainSubstring("Format:"))

			})

			// TODO: Check for repeated tag workflows & stages & branches
			It("should return a descriptive error for an invalid tag workflow", func() {

				config := validConfigWithOptionalFields()
				var existingStage string
				for k := range config.Flowit.Workflows[0] {
					existingStage = config.Flowit.Workflows[0][k][0].ID
					break
				}
				config.Flowit.Tags[0].Stages["missing-workflow"] = []string{
					existingStage,
				}
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Tags:"))
				Expect(err.Error()).To(ContainSubstring("Stages:"))
			})

			It("should return a descriptive error for an invalid stage on a valid workflow", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.Workflows[0]["my-workflow"] = []stage{
					{
						ID:   "my-stage",
						Args: []string{"arg1", "arg2"},
						Conditions: []string{
							"condition-1",
						},
						Actions: []string{
							"action-1",
						},
					},
				}
				config.Flowit.Tags[0].Stages["my-workflow"] = []string{
					"non-existent-stage",
				}
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Tags:"))
				Expect(err.Error()).To(ContainSubstring("Stages:"))
			})

			It("should return a descriptive error for an invalid tag branch", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.Tags[0].Branches = []string{
					"non-existent-branch",
				}
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Tags:"))
				Expect(err.Error()).To(ContainSubstring("Branches:"))
			})

		})

		Context("Validating stages", func() {

			It("should return a descriptive error for a non existent stage ID", func() {
				config := validConfigWithOptionalFields()
				firstWorkflow := config.Flowit.Workflows[0]
				config.Flowit.Workflows = []workflow{
					firstWorkflow,
					{
						"workflow": []stage{
							{
								Actions: []string{
									"action1",
								},
							},
						},
					},
				}
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Workflows:"))
				Expect(err.Error()).To(ContainSubstring("cannot be blank"))
			})

			It("should return a descriptive error for a non existent stage actions", func() {
				config := validConfigWithOptionalFields()
				firstWorkflow := config.Flowit.Workflows[0]
				config.Flowit.Workflows = []workflow{
					firstWorkflow,
					{
						"workflow": []stage{
							{
								ID: "my-id",
							},
						},
					},
				}
				rawConfig := rawify(&config)

				err := validateConfig(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Workflows:"))
				Expect(err.Error()).To(ContainSubstring("cannot be blank"))
			})

		})

	})

})

func rawify(config *WorkflowDefinition) *rawWorkflowDefinition {
	var rawConfig rawWorkflowDefinition
	if err := utils.DeepCopy(config, &rawConfig); err != nil {
		return nil
	}
	return &rawConfig
}
