package workflow_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/google/uuid"
	"github.com/yamil-rivera/flowit/internal/config"
	w "github.com/yamil-rivera/flowit/internal/workflow"
)

var _ = Describe("Workflow", func() {

	service := w.NewService()

	wd := config.Flowit{
		Version: "1",
		Config: config.Config{
			CheckpointExecution: true,
			Shell:               "bash",
		},
		Variables: map[string]interface{}{
			"variable": "value",
		},
		StateMachines: []config.StateMachine{
			{
				ID: "id",
				Stages: []string{
					"stage-1",
					"stage-2",
				},
				InitialStage: "stage-1",
				FinalStages: []string{
					"stage-2",
				},
				Transitions: []config.StateMachineTransition{
					{
						From: []string{
							"stage-1",
						},
						To: []string{
							"stage-2",
						},
					},
				},
			},
		},
	}

	Context("Creating Workflow", func() {

		It("should create a new workflow successfully", func() {
			workflow := service.CreateWorkflow("my-workflow", wd)

			id, err := uuid.Parse(workflow.ID)
			Expect(err).To(BeNil())
			Expect(workflow.Preffix).To(Equal(id.String()[:6]))
			Expect(workflow.Name).To(Equal("my-workflow"))
			Expect(workflow.IsActive).To(BeFalse())
			Expect(len(workflow.Executions)).To(Equal(0))
			Expect(workflow.LatestExecution).To(BeNil())
			Expect(workflow.State).To(Equal(wd))

			Expect(workflow.Metadata.Started).To(Equal(uint64(0)))
			Expect(workflow.Metadata.Updated).To(Equal(uint64(0)))
			Expect(workflow.Metadata.Finished).To(Equal(uint64(0)))
			Expect(workflow.Metadata.Version).To(Equal(uint64(0)))

		})

	})

	Context("Starting Workflow Execution", func() {

		It("should start a new workflow execution successfully", func() {

			wd.Variables = map[string]interface{}{
				"variable": "value",
			}
			args := []string{
				"arg-1",
			}
			workflow := service.CreateWorkflow("my-workflow", wd)

			// start first execution
			before := time.Now()
			execution := service.StartExecution(workflow, "origin", "stage-1", args)
			after := time.Now()

			// assert workflow and first execution
			workflowID := workflow.ID
			_, err := uuid.Parse(workflow.ID)
			Expect(err).To(BeNil())
			Expect(workflow.Preffix).To(Equal(workflowID[:6]))
			Expect(workflow.Name).To(Equal("my-workflow"))
			Expect(workflow.IsActive).To(BeTrue())
			Expect(workflow.State).To(Equal(wd))
			started := workflow.Metadata.Started
			Expect(workflow.Metadata.Started).To(BeNumerically(">=", uint64(before.UnixNano())))
			Expect(workflow.Metadata.Started).To(BeNumerically("<=", uint64(after.UnixNano())))
			Expect(workflow.Metadata.Updated).To(Equal(workflow.Metadata.Started))
			Expect(workflow.Metadata.Finished).To(Equal(uint64(0)))
			Expect(len(workflow.Executions)).To(Equal(1))
			Expect(workflow.LatestExecution).To(Equal(&(workflow.Executions[0])))
			Expect(workflow.LatestExecution).To(Equal(execution))

			firstExecutionID := execution.ID
			_, err = uuid.Parse(execution.ID)
			Expect(err).To(BeNil())
			Expect(execution.FromStage).To(Equal("origin"))
			Expect(execution.Stage).To(Equal("stage-1"))
			Expect(execution.Args).To(Equal(args))
			Expect(execution.Metadata.Version).To(BeEquivalentTo(0))
			Expect(execution.Metadata.Started).To(Equal(workflow.Metadata.Started))
			Expect(execution.Metadata.Finished).To(BeEquivalentTo(0))

			// start second execution
			args = []string{}
			before = time.Now()
			execution = service.StartExecution(workflow, "stage-1", "stage-2", args)
			after = time.Now()

			// assert workflow and second execution
			Expect(workflow.ID).To(Equal(workflowID))
			Expect(workflow.Preffix).To(Equal(workflowID[:6]))
			Expect(workflow.IsActive).To(BeTrue())
			Expect(workflow.Metadata.Started).To(Equal(started))
			Expect(workflow.Metadata.Updated).To(BeNumerically(">=", uint64(before.UnixNano())))
			Expect(workflow.Metadata.Updated).To(BeNumerically("<=", uint64(after.UnixNano())))
			Expect(workflow.Metadata.Finished).To(Equal(uint64(0)))
			Expect(len(workflow.Executions)).To(Equal(2))
			Expect(workflow.LatestExecution).To(Equal(&(workflow.Executions[0])))
			Expect(workflow.LatestExecution).To(Equal(execution))

			secondExecutionID := execution.ID
			_, err = uuid.Parse(execution.ID)
			Expect(secondExecutionID).ToNot(Equal(firstExecutionID))
			Expect(err).To(BeNil())
			Expect(execution.FromStage).To(Equal("stage-1"))
			Expect(execution.Stage).To(Equal("stage-2"))
			Expect(execution.Args).To(Equal(args))
			Expect(execution.Metadata.Version).To(BeEquivalentTo(0))
			Expect(execution.Metadata.Started).To(Equal(workflow.Metadata.Updated))
			Expect(execution.Metadata.Finished).To(BeEquivalentTo(0))

		})

	})

	Context("Finish a Workflow Execution", func() {

		It("should finish an active workflow execution successfully", func() {

			workflow := service.CreateWorkflow("my-workflow", wd)

			// start and finish execution
			execution := service.StartExecution(workflow, "origin", "stage-1", nil)
			before := time.Now()
			err := service.FinishExecution(workflow, execution, w.STARTED)
			after := time.Now()

			// assert
			Expect(err).To(BeNil())
			Expect(workflow.IsActive).To(BeTrue())
			Expect(workflow.Metadata.Finished).To(Equal(uint64(0)))
			Expect(execution.Metadata.Finished).To(BeNumerically(">=", uint64(before.UnixNano())))
			Expect(execution.Metadata.Finished).To(BeNumerically("<=", uint64(after.UnixNano())))

		})

		It("should finish a final workflow execution successfully", func() {

			workflow := service.CreateWorkflow("my-workflow", wd)

			// start and finish execution
			execution := service.StartExecution(workflow, "origin", "stage-1", nil)
			before := time.Now()
			err := service.FinishExecution(workflow, execution, w.FINISHED)
			after := time.Now()

			// assert
			Expect(err).To(BeNil())
			Expect(workflow.IsActive).To(BeFalse())
			Expect(workflow.Metadata.Finished).To(BeNumerically(">=", uint64(before.UnixNano())))
			Expect(workflow.Metadata.Finished).To(BeNumerically("<=", uint64(after.UnixNano())))
			Expect(execution.Metadata.Finished).To(Equal(workflow.Metadata.Finished))

		})

		It("should fail if an already terminated workflow execution is finished", func() {

			workflow := service.CreateWorkflow("my-workflow", wd)

			execution := service.StartExecution(workflow, "origin", "stage-1", nil)
			err := service.FinishExecution(workflow, execution, w.STARTED)
			Expect(err).To(BeNil())

			err = service.FinishExecution(workflow, execution, w.FINISHED)
			Expect(err).ToNot(BeNil())

		})

	})

})
