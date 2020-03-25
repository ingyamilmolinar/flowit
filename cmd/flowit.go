package main

import (
	"github.com/yamil-rivera/flowit/internal/commands"
	"github.com/yamil-rivera/flowit/internal/configs"
	logging "github.com/yamil-rivera/flowit/internal/utils"
)

func main() {
	logging.InitLogger()
	configs.ReadConfig("git-flow")
	commands.RegisterCommands("flowit")
}
