package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// RegisterCommands registers all possible subcommands based on the provided configuration
func RegisterCommands(mainCommand string) error {
	rootCmd := &cobra.Command{
		Use: mainCommand,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	err := rootCmd.Execute()
	if err != nil {
		return errors.Wrap(err, "Command execution error")
	}
	return nil
}
