package main

import (
	"github.com/yamil-rivera/flowit/internal/commands"
	"github.com/yamil-rivera/flowit/internal/configs"
	"github.com/yamil-rivera/flowit/internal/utils"
)

func main() {
	
	const mainCommand = "flowit"
	const configFile = "git-flow"

	_, err := configs.ProcessFlowitConfig(configFile, utils.GetRootDirectory()+"/samples/")
	utils.ExitIfErr(err)
	err = commands.RegisterCommands(mainCommand)
	utils.ExitIfErr(err)
}
