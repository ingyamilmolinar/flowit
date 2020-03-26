package main

import (
	"github.com/yamil-rivera/flowit/internal/commands"
	"github.com/yamil-rivera/flowit/internal/configs"
	"github.com/yamil-rivera/flowit/internal/utils"
)

const mainCommand = "flowit"
const configFile = "git-flow"

func main() {
	_, err := configs.ParseConfig(configFile, utils.GetRootDirectory()+"/samples/")
	exitIfErr(err)
	commands.RegisterCommands(mainCommand)
}

func exitIfErr(err error) {
	if err != nil {
		logger := utils.GetLogger()
		logger.Error(err)
		panic(err)
	}
}
