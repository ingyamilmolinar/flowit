package main

import (
	"github.com/yamil-rivera/flowit/internal/commands"
	"github.com/yamil-rivera/flowit/internal/configs"
	"github.com/yamil-rivera/flowit/internal/utils"
)

func main() {
	utils.InitLogger()
	config := configs.ReadConfig("git-flow", utils.GetRootDirectory()+"/samples/")
	configs.ValidateConfig(config)
	commands.RegisterCommands("flowit")
}
