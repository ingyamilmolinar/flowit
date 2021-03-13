package main

import (
	"io/ioutil"
	"os"

	"github.com/yamil-rivera/flowit/internal/command"
	"github.com/yamil-rivera/flowit/internal/config"
	"github.com/yamil-rivera/flowit/internal/fsm"
	"github.com/yamil-rivera/flowit/internal/io"
	"github.com/yamil-rivera/flowit/internal/repository"
	"github.com/yamil-rivera/flowit/internal/runtime"
	"github.com/yamil-rivera/flowit/internal/workflow"
)

func main() {

	workflowDefinition, err := config.Load(io.GetProjectRootDir() + "/samples/test.yaml")
	optionalExit(err)

	repositoryService := repository.NewService()

	workflowService := workflow.NewService()

	fsmServiceFactory := fsm.NewServiceFactory()

	runtimeService := runtime.NewService(repositoryService, fsmServiceFactory, workflowService)

	commandService := command.NewService(runtimeService, fsmServiceFactory, repositoryService, workflowDefinition)

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
		io.Printf("%v\n", err)
		// TODO: Do not show "exit status 1"
		os.Exit(1)
	}
}
