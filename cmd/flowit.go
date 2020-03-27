package main

import (
	"github.com/yamil-rivera/flowit/internal/commands"
	"github.com/yamil-rivera/flowit/internal/configs"
	"github.com/yamil-rivera/flowit/internal/utils"
)

const mainCommand = "flowit"
const configFile = "git-flow"

func main() {
	_, err := configs.ProcessFlowitConfig(configFile, utils.GetRootDirectory()+"/samples/")
	utils.ExitIfErr(err)
	commands.RegisterCommands(mainCommand)
}
