package main

import (
	"github.com/yamil-rivera/flowit/internal/commands"
	"github.com/yamil-rivera/flowit/internal/config"
	"github.com/yamil-rivera/flowit/internal/utils"
)

func main() {

	const configFile = "git-flow"
	const samplesDir = "samples"

	workflowDefinition, err := config.ProcessWorkflowDefinition(configFile, utils.GetProjectRootDir()+"/"+samplesDir)
	utils.OptionalExit(err)
	err = commands.RegisterCommands(workflowDefinition)
	utils.OptionalExit(err)
}
