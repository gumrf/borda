package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	command := &cobra.Command{
		Use:   "bordactl",
		Short: "Control manager for Borda",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := cmd.Flags().GetBool("version")
			if err != nil {
				return err
			}
			if version {
				fmt.Println("Borda v0.0.1")
			}

			return nil
		},
	}

	command.Flags().BoolP("version", "v", false, "Print version")

	// command.AddCommand(
	// 	commands.ImportTasksCommand(),
	// 	commands.ServeCommand(),
	// )

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
