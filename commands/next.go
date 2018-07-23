package commands

import (
	"fmt"
	"os"

	"github.com/chasinglogic/tsk/task"
	"github.com/spf13/cobra"
)

var next = &cobra.Command{
	Use:     "next",
	Aliases: []string{"current"},
	Short:   "Show the current task",
	Run: func(cmd *cobra.Command, args []string) {
		backend, err := config.backend()
		if err != nil {
			fmt.Println("ERROR Unable to load backend:", err)
			os.Exit(1)
		}

		current, err := backend.Current()
		if err != nil && err == task.ErrNotFound {
			fmt.Println("No uncompleted tasks found!")
			return
		} else if err != nil {
			fmt.Println("ERROR unable to get current task:", err)
			return
		}

		fmt.Printf("%s %s\n", current.ID, current.Title)
	},
}
