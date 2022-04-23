package commands

import (
	"borda/internal/app"

	"github.com/spf13/cobra"
)

func ServeCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "serve",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			app.Run()
		},
	}

	return command
}
