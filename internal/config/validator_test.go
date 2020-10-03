package config

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/yamil-rivera/flowit/internal/utils"
)

var _ = Describe("Config", func() {

	Describe("Validating a valid configuration", func() {

		Context("Validating minimum configuration", func() {

			It("should return a nil error for missing optional values", func() {

				rawConfig := validConfigJustMandatoryFields()
				err := validateWorkflowDefinition(&rawConfig)
				Expect(err).To(BeNil())

			})

		})

		Context("Validating complete configuration", func() {

			It("should return a nil error for having all optional values", func() {

				config := validConfigWithOptionalFields()
				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(BeNil())
			})

		})

		Context("Validating state machines", func() {

			It("should return a nil error when parsing '!' preceded transitions", func() {

				config := validConfigWithOptionalFields()

				config.Flowit.StateMachines[0].ID = "simple-machine"
				config.Flowit.StateMachines[0].Stages = []string{"stage-1", "stage-2", "stage-3", "stage-4"}
				config.Flowit.StateMachines[0].InitialStage = "stage-1"
				config.Flowit.StateMachines[0].FinalStages = []string{"stage-4"}
				config.Flowit.StateMachines[0].Transitions = []StateMachineTransition{
					// This means that from every stage except 'stage-4' going to every stage except 'stage-1' is allowed
					{
						From: []string{"!stage-4"},
						To:   []string{"!stage-1"},
					},
				}

				config.Flowit.Workflows[0].ID = "feature"
				config.Flowit.Workflows[0].StateMachine = "simple-machine"
				config.Flowit.Workflows[0].Stages = []Stage{
					{ID: "stage-1", Actions: []string{"action-1"}},
					{ID: "stage-2", Actions: []string{"action-2"}},
					{ID: "stage-3", Actions: []string{"action-3"}},
					{ID: "stage-4", Actions: []string{"action-4"}},
				}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(BeNil())

			})

		})

	})

	Describe("Validating an invalid configuration", func() {

		Context("Validating version", func() {

			It("should return a descriptive error for a non existent version", func() {

				config := validConfigWithOptionalFields()
				rawConfig := rawify(&config)

				rawConfig.Flowit.Version = nil

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Version:"))

			})

			It("should return a descriptive error for incorrect version", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Version = ".1"
				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Version:"))

			})

		})

		Context("Validating config", func() {

			It("should return a descriptive error for incorrect config", func() {

				config := validConfigWithOptionalFields()
				config.Flowit.Config.Shell = "/nonexistent/shell"
				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Config:"))

			})

		})

		Context("Validating variables", func() {

			It("should return a descriptive error for a non existent variable", func() {

				config := validConfigWithOptionalFields()
				rawConfig := rawify(&config)

				(*rawConfig.Flowit.Variables)["my-var"] = nil

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Variables:"))

			})

		})

		Context("Validating state machines", func() {

			It("should return a descriptive error for a non existent state-machine ID", func() {
				config := validConfigWithOptionalFields()

				rawConfig := rawify(&config)

				rawConfig.Flowit.StateMachines[0].ID = nil

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("ID:"))
			})

			It("should return a descriptive error for an invalid state-machine ID", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].ID = " "

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("ID:"))
			})

			It("should return a descriptive error for non existent state-machine stages", func() {
				config := validConfigWithOptionalFields()

				rawConfig := rawify(&config)

				rawConfig.Flowit.StateMachines[0].Stages = nil

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("Stages:"))

				config = validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Stages = []string{}

				rawConfig = rawify(&config)

				err = validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("Stages:"))

			})

			It("should return a descriptive error for an invalid state-machine stage ID", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Stages = []string{" "}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("Stages:"))
			})

			It("should return a descriptive error for a non existent state-machine initial stage", func() {
				config := validConfigWithOptionalFields()

				rawConfig := rawify(&config)

				rawConfig.Flowit.StateMachines[0].InitialStage = nil

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("InitialStage:"))
			})

			It("should return a descriptive error for an invalid state-machine initial stage ID", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].InitialStage = "initial-stage"

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("InitialStage:"))
			})

			It("should return a descriptive error for a nonexistent state-machine final stage ID", func() {
				config := validConfigWithOptionalFields()

				rawConfig := rawify(&config)

				rawConfig.Flowit.StateMachines[0].FinalStages = nil

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("FinalStages:"))

				config = validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].FinalStages = []string{}

				rawConfig = rawify(&config)

				err = validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("FinalStages:"))
			})

			It("should return a descriptive error for an invalid state-machine final stage ID", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].FinalStages = []string{"final-stage"}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("FinalStages:"))
			})

			It("should return a descriptive error for non existent state-machine transitions", func() {
				config := validConfigWithOptionalFields()
				rawConfig := rawify(&config)

				rawConfig.Flowit.StateMachines[0].Transitions = nil

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))

				config = validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Transitions = []StateMachineTransition{}

				rawConfig = rawify(&config)

				err = validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))
			})

			It("should return a descriptive error for an invalid state-machine transition stage ID", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Transitions = []StateMachineTransition{
					{
						From: []string{"from"},
						To:   []string{"to"},
					},
				}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("Transitions:"))
			})

			It("should return a descriptive error for an invalid state-machine", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Stages = []string{"stage-1", "stage-2", "stage-3"}
				config.Flowit.StateMachines[0].InitialStage = "stage-1"
				config.Flowit.StateMachines[0].FinalStages = []string{"stage-3"}
				config.Flowit.StateMachines[0].Transitions = []StateMachineTransition{
					{
						From: []string{"stage-1"},
						To:   []string{"stage-3"},
					},
				}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("Cannot reach a final node from 'stage-2' stage"))
			})

			It("should return a descriptive error for an invalid state-machine", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Stages = []string{"stage-1", "stage-2", "stage-3"}
				config.Flowit.StateMachines[0].InitialStage = "stage-1"
				config.Flowit.StateMachines[0].FinalStages = []string{"stage-3"}
				config.Flowit.StateMachines[0].Transitions = []StateMachineTransition{
					{
						From: []string{"stage-2"},
						To:   []string{"stage-1"},
					},
					{
						From: []string{"stage-1"},
						To:   []string{"stage-3"},
					},
					{
						From: []string{"stage-1"},
						To:   []string{"stage-2"},
					},
				}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("'stage-1' cannot be the destination in a transition"))
			})

			It("should return a descriptive error for an invalid state-machine", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Stages = []string{"stage-1", "stage-2", "stage-3"}
				config.Flowit.StateMachines[0].InitialStage = "stage-1"
				config.Flowit.StateMachines[0].FinalStages = []string{"stage-3"}
				config.Flowit.StateMachines[0].Transitions = []StateMachineTransition{
					{
						From: []string{"stage-1"},
						To:   []string{"stage-1", "stage-2"},
					},
					{
						From: []string{"stage-2"},
						To:   []string{"stage-3"},
					},
				}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("'stage-1' cannot be the destination in a transition"))
			})

			It("should return a descriptive error for an invalid state-machine", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Stages = []string{"stage-1", "stage-2", "stage-3"}
				config.Flowit.StateMachines[0].InitialStage = "stage-1"
				config.Flowit.StateMachines[0].FinalStages = []string{"stage-3"}
				config.Flowit.StateMachines[0].Transitions = []StateMachineTransition{
					{
						From: []string{"stage-1"},
						To:   []string{"stage-2"},
					},
					{
						From: []string{"stage-2"},
						To:   []string{"stage-3"},
					},
					{
						From: []string{"stage-3"},
						To:   []string{"stage-2"},
					},
				}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("'stage-3' cannot be the source in a transition"))
			})

			It("should return a descriptive error for an invalid state-machine", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Stages = []string{"stage-1", "stage-2", "stage-3"}
				config.Flowit.StateMachines[0].InitialStage = "stage-1"
				config.Flowit.StateMachines[0].FinalStages = []string{"stage-3"}
				config.Flowit.StateMachines[0].Transitions = []StateMachineTransition{
					{
						From: []string{"stage-1"},
						To:   []string{"stage-2", "stage-3"},
					},
					{
						From: []string{"stage-2"},
						To:   []string{"stage-3"},
					},
					{
						From: []string{"stage-3"},
						To:   []string{"stage-3"},
					},
				}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("Final stage: 'stage-3' cannot be the source in a transition"))
			})

			It("should return a descriptive error for an invalid state-machine", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Stages = []string{"stage-1", "stage-2", "stage-3", "stage-4"}
				config.Flowit.StateMachines[0].InitialStage = "stage-1"
				config.Flowit.StateMachines[0].FinalStages = []string{"stage-3"}
				config.Flowit.StateMachines[0].Transitions = []StateMachineTransition{
					{
						From: []string{"stage-1"},
						To:   []string{"stage-2", "stage-3"},
					},
					{
						From: []string{"stage-4"},
						To:   []string{"stage-3"},
					},
				}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("Cannot reach a final node from 'stage-2' stage"))
			})

			It("should return a descriptive error for an invalid state-machine", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Stages = []string{"stage-1", "stage-2", "stage-3"}
				config.Flowit.StateMachines[0].InitialStage = "stage-1"
				config.Flowit.StateMachines[0].FinalStages = []string{"stage-3"}
				config.Flowit.StateMachines[0].Transitions = []StateMachineTransition{
					{
						From: []string{"stage-1"},
						To:   []string{"stage-2", "stage-3"},
					},
				}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("Cannot reach a final node from 'stage-2' stage"))
			})

			It("should return a descriptive error for an invalid state-machine", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.StateMachines[0].Stages = []string{"stage-1", "stage-2", "stage-3", "stage-4"}
				config.Flowit.StateMachines[0].InitialStage = "stage-1"
				config.Flowit.StateMachines[0].FinalStages = []string{"stage-4"}
				config.Flowit.StateMachines[0].Transitions = []StateMachineTransition{
					{
						From: []string{"stage-1"},
						To:   []string{"stage-2", "stage-4"},
					},
					{
						From: []string{"stage-2"},
						To:   []string{"stage-3"},
					},
					{
						From: []string{"stage-3"},
						To:   []string{"stage-2"},
					},
				}

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("StateMachines:"))
				Expect(err.Error()).To(ContainSubstring("Cannot reach a final node from 'stage-2' stage"))
			})

		})

		Context("Validating workflows", func() {

			It("should return a descriptive error for a non existent workflow ID", func() {
				config := validConfigWithOptionalFields()
				rawConfig := rawify(&config)

				rawConfig.Flowit.Workflows[0].ID = nil

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Workflows:"))
				Expect(err.Error()).To(ContainSubstring("Workflow ID is nil"))
			})

			It("should return a descriptive error for an invalid workflow ID", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.Workflows[0].ID = " "

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Workflows:"))
				Expect(err.Error()).To(ContainSubstring("contains whitespaces"))
			})

			It("should return a descriptive error for an invalid state machine ID", func() {
				config := validConfigWithOptionalFields()
				config.Flowit.Workflows[0].StateMachine = " "

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Workflows:"))
				Expect(err.Error()).To(ContainSubstring("not a valid state machine"))
			})

		})

		Context("Validating stages", func() {

			It("should return a descriptive error for a non existent stage ID", func() {
				config := validConfigWithOptionalFields()
				rawConfig := rawify(&config)

				rawConfig.Flowit.Workflows[0].Stages[0].ID = nil

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Workflows:"))
				Expect(err.Error()).To(ContainSubstring("stages are missing in workflow"))
			})

			It("should return a descriptive error for a stage ID that is not defined in the state machine", func() {
				config := validConfigWithOptionalFields()

				sm := config.Flowit.Workflows[0].StateMachine
				config.Flowit.Workflows[0].Stages[0].ID = "new-stage"

				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Workflows:"))
				Expect(err.Error()).To(ContainSubstring("is not a valid " + sm + " state machine stage"))
			})

			It("should return a descriptive error for an invalid stage arg", func() {
				config := validConfigWithOptionalFields()
				firstWorkflow := config.Flowit.Workflows[0]
				firstWorkflow.ID = "feature"
				firstWorkflow.Stages[len(firstWorkflow.Stages)-1] = Stage{
					ID: "finish",
					Args: []string{
						"< my-var-without-description >",
					},
					Actions: []string{
						"action1",
						"action2",
					},
				}
				config.Flowit.Workflows[0] = firstWorkflow
				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Workflows:"))
				Expect(err.Error()).To(ContainSubstring("Invalid workflow stage argument"))

			})

			It("should return a descriptive error for a non existent stage actions", func() {
				config := validConfigWithOptionalFields()
				firstWorkflow := config.Flowit.Workflows[0]
				firstWorkflow.ID = "feature"
				firstWorkflow.Stages[len(firstWorkflow.Stages)-1] = Stage{
					ID: "finish",
				}
				config.Flowit.Workflows[0] = firstWorkflow
				rawConfig := rawify(&config)

				err := validateWorkflowDefinition(rawConfig)
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("Workflows:"))
				// TODO: Error message is: Workflows: (0: 1: cannot be blank..)
				// It doesn't say anything about 'Actions' being missing
				Expect(err.Error()).To(ContainSubstring("cannot be blank"))
			})

		})

	})

})

func validConfigJustMandatoryFields() rawWorkflowDefinition {

	var config rawWorkflowDefinition
	var flowit rawMainDefinition

	config.Flowit = &flowit

	version := "0.1"

	var stateMachine rawStateMachine
	stateMachineID := "simple-machine"
	startStageID := "start"
	finishStageID := "finish"
	stateMachineStages := []*string{
		&startStageID, &finishStageID,
	}
	stateMachineInitialStage := startStageID
	stateMachineFinalStages := []*string{
		&finishStageID,
	}
	stateMachineTransitions := []*rawStateMachineTransition{
		{
			From: []*string{
				&startStageID,
			},
			To: []*string{
				&finishStageID,
			},
		},
	}

	stateMachine.ID = &stateMachineID
	stateMachine.Stages = stateMachineStages
	stateMachine.InitialStage = &stateMachineInitialStage
	stateMachine.FinalStages = stateMachineFinalStages
	stateMachine.Transitions = stateMachineTransitions

	startStageAction1 := "start action1"
	startStageAction2 := "start action2"
	startStage := rawStage{
		ID:      &startStageID,
		Actions: []*string{&startStageAction1, &startStageAction2},
	}
	finishStageAction1 := "finish action1"
	finishStageAction2 := "finish action2"
	finishStage := rawStage{
		ID:      &finishStageID,
		Actions: []*string{&finishStageAction1, &finishStageAction2},
	}
	workflowID := "feature"
	workflowType := rawWorkflow{
		ID:           &workflowID,
		StateMachine: &stateMachineID,
		Stages: []*rawStage{
			&startStage,
			&finishStage,
		},
	}

	mainConfig := rawMainDefinition{
		Version: &version,
		StateMachines: []*rawStateMachine{
			&stateMachine,
		},
		Workflows: []*rawWorkflow{
			&workflowType,
		},
	}

	config.Flowit = &mainConfig
	return config
}

func validConfigWithOptionalFields() WorkflowDefinition {

	var flowit Flowit

	flowit.Version = "0.1"
	flowit.Config = Config{
		AbortOnFailedAction: true,
		Shell:               "/usr/bin/env bash",
	}
	flowit.Variables = map[string]interface{}{
		"var1": "value",
		"var2": 12345,
		"var3": "${env-variable}",
	}
	flowit.StateMachines = []StateMachine{
		{
			ID: "simple-machine",
			Stages: []string{
				"start", "finish",
			},
			InitialStage: "start",
			FinalStages:  []string{"finish"},
			Transitions: []StateMachineTransition{
				{
					From: []string{"start"},
					To:   []string{"finish"},
				},
			},
		},
	}
	flowit.Workflows = []Workflow{
		{
			ID:           "feature",
			StateMachine: flowit.StateMachines[0].ID,
			Stages: []Stage{
				{
					ID:   "start",
					Args: []string{"< my-var-1 | My-desc-1 >", "< my-var-2 | My-desc-2 >"},
					Conditions: []string{
						"start condition1",
					},
					Actions: []string{"start action1", "start action2"},
				},
				{
					ID:   "finish",
					Args: []string{"< my-var-1 | My-desc-1 >", "< my-var-2 | My-desc-2 >"},
					Conditions: []string{
						"finish condition1",
					},
					Actions: []string{"finish action1", "finish action2"},
				},
			},
		},
	}

	var config WorkflowDefinition
	config.Flowit = flowit
	return config
}

func rawify(config *WorkflowDefinition) *rawWorkflowDefinition {
	var rawConfig rawWorkflowDefinition
	if err := utils.DeepCopy(config, &rawConfig); err != nil {
		return nil
	}
	return &rawConfig
}
