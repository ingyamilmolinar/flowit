package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yamil-rivera/flowit/internal/config"
)

type command struct {
	cobra       *cobra.Command
	subcommands []command
}

// RegisterCommands registers all commands and subcommands based on the provided configuration
func RegisterCommands(workflowDefinition *config.WorkflowDefinition) error {

	workflows := workflowDefinition.Flowit.Workflows
	commands := make([]command, len(workflows))
	for i, workflow := range workflows {
		for workflowID, stages := range workflow {
			commands[i].cobra = &cobra.Command{
				Use: workflowID,
			}
			commands[i].subcommands = make([]command, len(stages))
			for j, stage := range stages {
				commands[i].subcommands[j].cobra = &cobra.Command{
					Use:  stage.ID,
					Args: cobra.ExactArgs(len(stage.Args)),
					PreRunE: func(cmd *cobra.Command, args []string) error {
						//return runtime.RunConditions(stage.Conditions, buildReplacementMap(stage.Args, args), workflowDefinition.Flowit.Config.Shell)
						return nil
					},
					RunE: func(cmd *cobra.Command, args []string) error {
						//return runtime.RunActions(stage.Actions, buildReplacementMap(stage.Args, args), workflowDefinition.Flowit.Config.Shell)
						return nil
					},
				}
			}
		}
	}

	rootCommand := &cobra.Command{
		Use:   "flowit",
		Short: "A flexible git workflow manager",
		Long:  "A flexible git workflow manager",
	}

	for _, command := range commands {
		for _, subcommand := range command.subcommands {
			command.cobra.AddCommand(subcommand.cobra)
		}
		rootCommand.AddCommand(command.cobra)
	}

	if err := rootCommand.Execute(); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
