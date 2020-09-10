package runtime

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/yamil-rivera/flowit/internal/utils"
)

type Service struct {
	executor Executor
}

type Executor interface {
	Execute(command string) (string, error)
}

func NewService(executor Executor) *Service {
	return &Service{executor}
}

func NewUnixShellExecutor(shell string) Executor {
	return UnixShellExecutor{shell}
}

type UnixShellExecutor struct {
	shell string
}

func (e UnixShellExecutor) Execute(command string) (string, error) {
	shellArgs := strings.Split(e.shell, " ")
	mainCommand := shellArgs[0]
	restOfArgs := append(shellArgs[1:], "-c", command)
	cmd := exec.Command(mainCommand, restOfArgs...)
	out, err := cmd.Output()
	trimmedOut := strings.TrimSuffix(string(out), "\n")
	if err != nil {
		return trimmedOut, errors.Wrap(err, "Error executing command: "+command+" with shell: "+e.shell)
	}
	return trimmedOut, nil
}

func (s *Service) Execute(commands []string, variables map[string]interface{}) ([]string, error) {

	out, err := s.runCommands(commands, variables)
	if err != nil {
		return out, errors.WithStack(err)
	}

	return out, nil
}

func (s Service) runCommands(commands []string, variables map[string]interface{}) ([]string, error) {

	if len(commands) == 0 {
		return nil, nil
	}

	var outs []string
	for _, command := range commands {
		parsedCommand, err := utils.EvaluateVariablesInExpression(command, variables)
		if err != nil {
			return outs, errors.Wrap(err, "Error evaluating variables in command: "+command)
		}
		out, err := s.executor.Execute(parsedCommand)
		outs = append(outs, out)
		if err != nil {
			return outs, errors.Wrap(err, "Error executing command: "+parsedCommand)
		}
	}

	return outs, nil
}
