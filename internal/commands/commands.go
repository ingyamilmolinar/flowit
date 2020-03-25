package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RegisterCommands registers all possible subcommands based on the provided configuration
func RegisterCommands(mainCommand string) {
	rootCmd := &cobra.Command{
		Use: mainCommand,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
