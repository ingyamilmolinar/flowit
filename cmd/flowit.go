package main

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/command"
	"github.com/yamil-rivera/flowit/internal/config"
	"github.com/yamil-rivera/flowit/internal/fsm"
	"github.com/yamil-rivera/flowit/internal/io"
	"github.com/yamil-rivera/flowit/internal/repository"
	"github.com/yamil-rivera/flowit/internal/runtime"
	"github.com/yamil-rivera/flowit/internal/utils"
	"github.com/yamil-rivera/flowit/internal/workflow"
)

func main() {

	workflowDefinition, err := config.Load(io.GetProjectRootDir() + "/samples/test.yaml")
	optionalExit(err)

	stateMachines, err := buildFSMs(workflowDefinition.Flowit.StateMachines, workflowDefinition.Flowit.Workflows)
	optionalExit(err)

	fsmService := fsm.NewService(stateMachines)

	repositoryService := repository.NewService()

	workflowService := workflow.NewService()

	shell := workflowDefinition.Flowit.Config.Shell
	runtimeService := runtime.NewService(runtime.NewUnixShellExecutor(shell))

	commandService := command.NewService(runtimeService, fsmService, repositoryService, workflowService, workflowDefinition)

	version, err := cliVersion()
	optionalExit(err)

	optionalExit(commandService.RegisterCommands(version))
	optionalExit(commandService.Execute())
}

func cliVersion() (string, error) {
	version, err := ioutil.ReadFile(io.GetProjectRootDir() + "/cmd/version")
	return string(version), err
}

func optionalExit(err error) {
	if err != nil {
		io.Logger.Errorf("%+v", err)
		os.Exit(1)
	}
}

func buildFSMs(stateMachines []config.StateMachine, workflows []config.Workflow) ([]fsm.StateMachine, error) {

	fsmWorkflows := make([]fsm.StateMachine, len(workflows))
	for i, workflow := range workflows {
		fsmWorkflows[i].ID = workflow.ID

		var stateMachine fsm.StateMachine
		for _, sm := range stateMachines {

			if workflow.StateMachine == sm.ID {
				stateMachine.States = sm.Stages
				stateMachine.InitialState = sm.InitialStage
				stateMachine.FinalStates = sm.FinalStages
				fsmTransitions, err := buildTransitions(sm.Transitions)
				if err != nil {
					return nil, errors.WithStack(err)
				}
				stateMachine.Transitions = fsmTransitions
			}

		}

		fsmWorkflows[i].States = stateMachine.States
		fsmWorkflows[i].InitialState = stateMachine.InitialState
		fsmWorkflows[i].FinalStates = stateMachine.FinalStates
		fsmWorkflows[i].Transitions = stateMachine.Transitions

	}
	return fsmWorkflows, nil
}

func buildTransitions(configTransitions []config.StateMachineTransition) ([]fsm.StateMachineTransition, error) {

	var fsmTransitions []fsm.StateMachineTransition
	if err := utils.DeepCopy(configTransitions, &fsmTransitions); err != nil {
		return nil, errors.WithStack(err)
	}
	return fsmTransitions, nil

}
