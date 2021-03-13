package runtime_test

import (
	"errors"

	"github.com/yamil-rivera/flowit/internal/config"
	"github.com/yamil-rivera/flowit/internal/utils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/yamil-rivera/flowit/internal/fsm"
	r "github.com/yamil-rivera/flowit/internal/runtime"
	"github.com/yamil-rivera/flowit/internal/testmocks"
	"github.com/yamil-rivera/flowit/internal/workflow"
)

type mockExecutor struct{}
type mockWriter struct {
	captures []string
}

func (e mockExecutor) Config(shell string) {
	// We don't do anything
}

func (e mockExecutor) Execute(command string) (string, error) {
	if command == "FAIL" {
		return command, errors.New("Command failed")
	}
	return command, nil
}

func (w *mockWriter) Write(s string) error {
	w.captures = append(w.captures, s)
	return nil
}

var _ = Describe("Runtime", func() {

	createWorkflowDefinition := func() config.Flowit {
		return config.Flowit{
			Version: "1",
			Config: config.Config{
				CheckpointExecution: true,
			},
			Variables: map[string]interface{}{},
			StateMachines: []config.StateMachine{
				{
					ID: "simple-machine",
					Stages: []string{
						"start",
					},
					InitialStage: "start",
					FinalStages: []string{
						"finish",
					},
					Transitions: []config.StateMachineTransition{
						{
							From: []string{
								"start",
							},
							To: []string{
								"finish",
							},
						},
					},
				},
			},
			Workflows: []config.Workflow{
				{
					ID:           "feature",
					StateMachine: "simple-machine",
					Stages: []config.Stage{
						{
							ID: "start",
							Args: []string{
								"< arg-1 | test >",
								"< arg-2 | test >",
							},
							Conditions: []string{
								"COND1",
								"COND2: $<arg-1>",
							},
							Actions: []string{
								"ACTION1",
								"ACTION2: $<arg-2>",
							},
						},
					},
				},
			},
		}
	}

	fsf := fsm.NewServiceFactory()
	ws := workflow.NewService()

	Context("Executing stages", func() {

		It("should execute successfully for a new workflow", func() {
			rs := testmocks.NewRepositoryMock()
			service := r.NewService(rs, fsf, ws)

			args := []string{
				"1",
				"2",
			}
			workflowName := "feature"
			stageID := "start"

			wd := createWorkflowDefinition()
			writer := &mockWriter{}
			err := service.Run(utils.OptionalString{}, args, workflowName, stageID, wd, mockExecutor{}, writer)
			Expect(err).ToNot(HaveOccurred())
			Expect(writer.captures).To(ContainElements([]string{
				"COND1",
				"COND2: 1",
				"ACTION1",
				"ACTION2: 2",
			}))
		})

		It("should execute successfully for an existing workflow", func() {
			rs := testmocks.NewRepositoryMock()
			w := workflow.Workflow{
				ID:      "12345",
				Preffix: "12345",
				Name:    "feature",
				State:   createWorkflowDefinition(),
			}
			err := rs.PutWorkflow(w)
			Expect(err).ToNot(HaveOccurred())
			service := r.NewService(rs, fsf, ws)

			wd := config.Flowit{}
			args := []string{
				"1",
				"2",
			}
			workflowName := "feature"
			stageID := "start"

			writer := &mockWriter{}
			// TODO: Consider changing service.Run() to accept either a workflowID or a workflowDefinition
			err = service.Run(utils.NewStringOptional(w.ID), args, workflowName, stageID, wd, mockExecutor{}, writer)
			Expect(err).ToNot(HaveOccurred())
			Expect(writer.captures).To(ContainElements([]string{
				"COND1",
				"COND2: 1",
				"ACTION1",
				"ACTION2: 2",
			}))
		})

		It("should fail to run stage with wrong arguments", func() {
			rs := testmocks.NewRepositoryMock()
			service := r.NewService(rs, fsf, ws)

			args := []string{
				"1",
			}
			workflowName := "feature"
			stageID := "start"

			wd := createWorkflowDefinition()
			writer := &mockWriter{}
			err := service.Run(utils.OptionalString{}, args, workflowName, stageID, wd, mockExecutor{}, writer)
			Expect(err).To(HaveOccurred())
		})

		It("should not execute actions if one condition fails", func() {
			rs := testmocks.NewRepositoryMock()
			service := r.NewService(rs, fsf, ws)

			args := []string{
				"1",
				"2",
			}
			workflowName := "feature"
			stageID := "start"

			wd := createWorkflowDefinition()
			wd.Workflows[0].Stages[0].Conditions = []string{
				"COND1",
				"COND2: $<arg-1>",
				"FAIL",
			}
			writer := &mockWriter{}
			err := service.Run(utils.OptionalString{}, args, workflowName, stageID, wd, mockExecutor{}, writer)
			Expect(err).To(HaveOccurred())
			Expect(writer.captures).To(ContainElements([]string{
				"COND1",
				"COND2: 1",
			}))
			Expect(writer.captures).ToNot(ContainElements([]string{
				"ACTION1",
				"ACTION2: 2",
			}))
		})

		It("should save checkpoint with failed action", func() {
			rs := testmocks.NewRepositoryMock()
			service := r.NewService(rs, fsf, ws)

			args := []string{
				"1",
				"2",
			}
			workflowName := "feature"
			stageID := "start"

			wd := createWorkflowDefinition()
			wd.Config.CheckpointExecution = true
			wd.Workflows[0].Stages[0].Actions = []string{
				"ACTION1",
				"ACTION2: $<arg-2>",
				"FAIL",
			}
			writer := &mockWriter{}
			err := service.Run(utils.OptionalString{}, args, workflowName, stageID, wd, mockExecutor{}, writer)
			Expect(err).To(HaveOccurred())
			Expect(writer.captures).To(ContainElements([]string{
				"COND1",
				"COND2: 1",
				"ACTION1",
				"ACTION2: 2",
				"FAIL",
			}))

			workflows, err := rs.GetWorkflows("feature", 1, true)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(workflows)).To(Equal(1))

			writer = &mockWriter{}
			err = service.Run(utils.NewStringOptional(workflows[0].Preffix), args, workflowName, stageID, wd, mockExecutor{}, writer)
			Expect(err).To(HaveOccurred())
			Expect(writer.captures).To(ContainElements([]string{
				"COND1",
				"COND2: 1",
				"FAIL",
			}))
			Expect(writer.captures).ToNot(ContainElements([]string{
				"ACTION1",
				"ACTION2: 2",
			}))
		})

		It("should fail to resume a failed checkpoint stage if given different arguments", func() {
			rs := testmocks.NewRepositoryMock()
			service := r.NewService(rs, fsf, ws)

			args := []string{
				"1",
				"2",
			}
			workflowName := "feature"
			stageID := "start"

			wd := createWorkflowDefinition()
			wd.Config.CheckpointExecution = true
			wd.Workflows[0].Stages[0].Actions = []string{
				"ACTION1",
				"ACTION2: $<arg-2>",
				"FAIL",
			}
			writer := &mockWriter{}
			err := service.Run(utils.OptionalString{}, args, workflowName, stageID, wd, mockExecutor{}, writer)
			Expect(err).To(HaveOccurred())

			workflows, err := rs.GetWorkflows("feature", 1, true)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(workflows)).To(Equal(1))

			args = []string{
				"1",
			}
			writer = &mockWriter{}
			err = service.Run(utils.NewStringOptional(workflows[0].Preffix), args, workflowName, stageID, wd, mockExecutor{}, writer)
			Expect(err).To(HaveOccurred())
			Expect(writer.captures).ToNot(ContainElement("COND1"))
		})

		It("should fail to execute an incorrect stage", func() {
			rs := testmocks.NewRepositoryMock()
			service := r.NewService(rs, fsf, ws)

			args := []string{}
			workflowName := "feature"
			stageID := "finish"

			wd := createWorkflowDefinition()
			writer := &mockWriter{}
			err := service.Run(utils.OptionalString{}, args, workflowName, stageID, wd, mockExecutor{}, writer)
			Expect(err).To(HaveOccurred())
			Expect(writer.captures).ToNot(ContainElement("COND1"))
		})

	})

})
