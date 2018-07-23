package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/chasinglogic/tsk/task"
	"github.com/spf13/cobra"
)

var new = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n", "create"},
	Short:   "Create a new task",
	Run: func(cmd *cobra.Command, args []string) {
		title := strings.Join(args, " ")
		t := task.New(title)

		backend, err := config.backend()
		if err != nil {
			fmt.Println("ERROR Unable to load backend:", err)
			os.Exit(1)
		}

		err = backend.Add(t)
		if err != nil {
			fmt.Println("ERROR Unable to add task:", err)
		}

		err = backend.Save()
		if err != nil {
			fmt.Println("ERROR Unable to save to backend:", err)
		}
	},
}
