package testmocks

import (
	"strings"

	"github.com/yamil-rivera/flowit/internal/runtime"
	w "github.com/yamil-rivera/flowit/internal/workflow"
)

type mockService struct {
	workflows map[string][]w.Workflow
}

func NewRepositoryMock() runtime.RepositoryService {
	return mockService{
		workflows: make(map[string][]w.Workflow),
	}
}

func (s mockService) PutWorkflow(workflow w.Workflow) error {
	s.workflows[workflow.Name] = append(s.workflows[workflow.Name], workflow)
	return nil
}

func (s mockService) GetWorkflow(workflowName, workflowID string) (w.OptionalWorkflow, error) {
	for wn, wfs := range s.workflows {
		if wn != workflowName {
			continue
		}
		for _, wf := range wfs {
			if wf.ID == workflowID {
				return w.NewWorkflowOptional(wf), nil
			}
		}
	}
	return w.OptionalWorkflow{}, nil
}

func (s mockService) GetWorkflows(workflowName string, count int, excludeInactive bool) ([]w.Workflow, error) {
	return s.workflows[workflowName], nil
}

func (s mockService) GetAllWorkflows(excludeInactive bool) ([]w.Workflow, error) {
	var workflows []w.Workflow
	for _, ws := range s.workflows {
		for _, wf := range ws {
			workflows = append(workflows, wf)
		}
	}
	return workflows, nil
}

func (s mockService) GetWorkflowFromPreffix(workflowName, workflowIDPreffix string) (w.OptionalWorkflow, error) {
	for wn, wfs := range s.workflows {
		if wn != workflowName {
			continue
		}
		for _, wf := range wfs {
			if strings.Contains(wf.ID, workflowIDPreffix) {
				return w.NewWorkflowOptional(wf), nil
			}
		}
	}
	return w.OptionalWorkflow{}, nil
}
