package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/yamil-rivera/flowit/internal/config"
	"github.com/yamil-rivera/flowit/internal/io"
)

func main() {

	const configFile = "git-flow"
	const samplesDir = "samples"

	// TODO: We need to expose our CLI FSM workflow manager service with a clean interface
	optionalExit(config.LoadConfiguration(configFile, io.GetProjectRootDir()+"/"+samplesDir))
	// optionalExit(commands.RegisterCommands(getVersion()))
}

func getVersion() string {
	version, err := ioutil.ReadFile(io.GetProjectRootDir() + "/cmd/version")
	optionalExit(err)
	return string(version)
}

func optionalExit(err error) {
	if err != nil {
		io.Logger.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}
