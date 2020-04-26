package main

import (
	"github.com/yamil-rivera/flowit/internal/commands"
	"github.com/yamil-rivera/flowit/internal/config"
	"github.com/yamil-rivera/flowit/internal/utils"
)

func main() {

	const mainCommand = "flowit"
	const configFile = "git-flow"
	const samplesDir = "samples"

	_, err := config.ProcessFlowitConfig(configFile, utils.GetProjectRootDir()+"/"+samplesDir)
	utils.OptionalExit(err)
	err = commands.RegisterCommands(mainCommand)
	utils.OptionalExit(err)
}
