package repository_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/yamil-rivera/flowit/internal/models"
	r "github.com/yamil-rivera/flowit/internal/repository"
)

var _ = Describe("Repository", func() {

	execution := models.Execution{
		ID:    "2",
		State: "state",
		Metadata: models.ExecutionMetadata{
			Version:  0xABABABAB,
			Started:  0xBCBCBCBC,
			Finished: 0xCDCDCDCD,
		},
	}
	workflow := models.Workflow{
		ID:           "1",
		Name:         "workflow",
		DefinitionID: "definition",
		IsActive:     true,
		Executions: []models.Execution{
			execution,
		},
		LatestExecution: &execution,
		Variables: map[string]string{
			"my-var": "my-val",
		},
		Metadata: models.WorkflowMetadata{
			Version:  0xDEDEDEDE,
			Started:  0xEFEFEFEF,
			Updated:  0xABABABAB,
			Finished: 0xBCBCBCBC,
		},
	}

	Context("Storing workflows", func() {

		It("should successfully save and retrieve a populated workflow", func() {

			defer r.DeleteDB()

			err := r.PutWorkflow(workflow)
			Expect(err).To(BeNil())
			optionalWorkflow, err := r.GetWorkflow("definition", "1")
			Expect(err).To(BeNil())
			savedWorkflow, err := optionalWorkflow.Get()
			Expect(err).To(BeNil())
			Expect(savedWorkflow).To(Equal(workflow))

		})

		It("should successfully overwrite a workflow", func() {

			defer r.DeleteDB()

			err := r.PutWorkflow(workflow)
			Expect(err).To(BeNil())
			expectedWorkflow := workflow
			expectedWorkflow.Name = "other workflow"
			err = r.PutWorkflow(expectedWorkflow)
			Expect(err).To(BeNil())
			overwrittenWorkflowOption, err := r.GetWorkflow("definition", "1")
			Expect(err).To(BeNil())
			overwrittenWorkflow, err := overwrittenWorkflowOption.Get()
			Expect(err).To(BeNil())
			Expect(overwrittenWorkflow).To(Equal(expectedWorkflow))

		})

	})

	Context("Retrieving workflows", func() {

		It("should successfully retrieve a workflow", func() {

			defer r.DeleteDB()

			err := r.PutWorkflow(workflow)
			Expect(err).To(BeNil())

			workflow2 := workflow
			workflow2.ID = "2"
			workflow2.Name = "workflow 2"

			err = r.PutWorkflow(workflow2)
			Expect(err).To(BeNil())

			firstWorkflowOptional, err := r.GetWorkflow("definition", "1")
			Expect(err).To(BeNil())
			firstWorkflow, err := firstWorkflowOptional.Get()
			Expect(err).To(BeNil())
			Expect(firstWorkflow).To(Equal(workflow))

		})

		It("should return an empty optional when workflow does not exist", func() {

			defer r.DeleteDB()

			firstWorkflowOptional, err := r.GetWorkflow("definition", "1")
			Expect(err).To(BeNil())
			_, err = firstWorkflowOptional.Get()
			Expect(err).To(Not(BeNil()))

			err = r.PutWorkflow(workflow)
			Expect(err).To(BeNil())

			firstWorkflowOptional, err = r.GetWorkflow("Definition", "1")
			Expect(err).To(BeNil())
			_, err = firstWorkflowOptional.Get()
			Expect(err).To(Not(BeNil()))

		})

		It("should successfully retrieve a workflow from prefix", func() {

			defer r.DeleteDB()

			workflow1 := workflow
			workflow1.ID = "100"
			workflow1.Name = "workflow 1"

			workflow2 := workflow
			workflow2.ID = "200"
			workflow2.Name = "workflow 2"

			workflow3 := workflow
			workflow3.ID = "300"
			workflow3.Name = "workflow 3"

			err := r.PutWorkflow(workflow1)
			Expect(err).To(BeNil())
			err = r.PutWorkflow(workflow2)
			Expect(err).To(BeNil())
			err = r.PutWorkflow(workflow3)
			Expect(err).To(BeNil())

			workflowOptional, err := r.GetWorkflowFromPreffix("definition", "1")
			Expect(err).To(BeNil())
			workflowWithPrefix, err := workflowOptional.Get()
			Expect(err).To(BeNil())
			Expect(workflowWithPrefix).To(Equal(workflow1))

			workflowOptional, err = r.GetWorkflowFromPreffix("definition", "2")
			Expect(err).To(BeNil())
			workflowWithPrefix, err = workflowOptional.Get()
			Expect(err).To(BeNil())
			Expect(workflowWithPrefix).To(Equal(workflow2))

			workflowOptional, err = r.GetWorkflowFromPreffix("definition", "3")
			Expect(err).To(BeNil())
			workflowWithPrefix, err = workflowOptional.Get()
			Expect(err).To(BeNil())
			Expect(workflowWithPrefix).To(Equal(workflow3))

		})

		It("should return an empty optional when a workflow does not start with prefix", func() {

			defer r.DeleteDB()

			workflow1 := workflow
			workflow1.ID = "01"
			workflow1.Name = "workflow 1"

			err := r.PutWorkflow(workflow1)
			Expect(err).To(BeNil())

			workflowOptional, err := r.GetWorkflowFromPreffix("definition", "1")
			Expect(err).To(BeNil())
			_, err = workflowOptional.Get()
			Expect(err).To(Not(BeNil()))

		})

		It("should successfully retrieve a list of n workflows", func() {

			defer r.DeleteDB()

			workflow1 := workflow
			workflow1.ID = "1"
			workflow1.Name = "workflow 1"

			workflow2 := workflow
			workflow2.ID = "2"
			workflow2.Name = "workflow 2"

			workflow3 := workflow
			workflow3.ID = "3"
			workflow3.Name = "workflow 3"

			err := r.PutWorkflow(workflow1)
			Expect(err).To(BeNil())
			err = r.PutWorkflow(workflow2)
			Expect(err).To(BeNil())
			err = r.PutWorkflow(workflow3)
			Expect(err).To(BeNil())

			workflows, err := r.GetWorkflows("definition", 0, false)
			Expect(err).To(BeNil())
			Expect(len(workflows)).To(Equal(3))
			Expect(workflows[0]).To(Equal(workflow1))
			Expect(workflows[1]).To(Equal(workflow2))
			Expect(workflows[2]).To(Equal(workflow3))

			workflows, err = r.GetWorkflows("definition", 2, false)
			Expect(err).To(BeNil())
			Expect(len(workflows)).To(Equal(2))
			Expect(workflows[0]).To(Equal(workflow1))
			Expect(workflows[1]).To(Equal(workflow2))

			workflows, err = r.GetWorkflows("definition", 4, false)
			Expect(len(workflows)).To(Equal(3))
			Expect(workflows[0]).To(Equal(workflow1))
			Expect(workflows[1]).To(Equal(workflow2))
			Expect(workflows[2]).To(Equal(workflow3))
			Expect(err).To(Not(BeNil()))

		})

		It("should successfully retrieve a list of active workflows", func() {

			defer r.DeleteDB()

			workflow1 := workflow
			workflow1.ID = "1"
			workflow1.Name = "workflow 1"
			workflow1.IsActive = true

			workflow2 := workflow
			workflow2.ID = "2"
			workflow2.Name = "workflow 2"
			workflow2.IsActive = false

			workflow3 := workflow
			workflow3.ID = "3"
			workflow3.Name = "workflow 3"
			workflow3.IsActive = true

			err := r.PutWorkflow(workflow1)
			Expect(err).To(BeNil())
			err = r.PutWorkflow(workflow2)
			Expect(err).To(BeNil())
			err = r.PutWorkflow(workflow3)
			Expect(err).To(BeNil())

			workflows, err := r.GetWorkflows("definition", 0, true)
			Expect(err).To(BeNil())
			Expect(len(workflows)).To(Equal(2))
			Expect(workflows[0]).To(Equal(workflow1))
			Expect(workflows[1]).To(Equal(workflow3))

			workflows, err = r.GetWorkflows("definition", 1, true)
			Expect(err).To(BeNil())
			Expect(len(workflows)).To(Equal(1))
			Expect(workflows[0]).To(Equal(workflow1))

			workflows, err = r.GetWorkflows("definition", 3, true)
			Expect(len(workflows)).To(Equal(2))
			Expect(workflows[0]).To(Equal(workflow1))
			Expect(workflows[1]).To(Equal(workflow3))
			Expect(err).To(Not(BeNil()))

		})

	})

	Context("Deleting workflows", func() {

		It("should successfully delete a workflow", func() {

			defer r.DeleteDB()

			err := r.PutWorkflow(workflow)
			Expect(err).To(BeNil())
			err = r.DeleteWorkflow("definition", "1")
			Expect(err).To(BeNil())
			deletedWorkflowOptional, err := r.GetWorkflow("definition", "1")
			Expect(err).To(BeNil())
			_, err = deletedWorkflowOptional.Get()
			Expect(err).To(Not(BeNil()))

		})

		It("should return an error if workflow does not exist", func() {

			defer r.DeleteDB()

			err := r.DeleteWorkflow("definition", "1")
			Expect(err).To(Not(BeNil()))

		})

	})

	Context("Deleting the DB", func() {

		It("should successfully wipe out the DB", func() {

			err := r.PutWorkflow(workflow)
			Expect(err).To(BeNil())

			workflows, err := r.GetWorkflows("definition", 0, false)
			Expect(err).To(BeNil())
			Expect(len(workflows)).To(Equal(1))

			r.DeleteDB()

			workflows, err = r.GetWorkflows("definition", 0, false)
			Expect(err).To(BeNil())
			Expect(len(workflows)).To(Equal(0))

		})

	})

})