package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var complete = &cobra.Command{
	Use:     "complete",
	Aliases: []string{"done", "d"},
	Short:   "Complete tasks by ID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("TASK_ID is a required argument")
			os.Exit(1)
		}

		backend, err := config.backend()
		if err != nil {
			fmt.Println("ERROR Unable to load backend:", err)
			os.Exit(1)
		}

		if err := backend.Complete(args[0]); err != nil {
			fmt.Println("ERROR Unable to complete task:", err)
			os.Exit(1)
		}
	},
}
