package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/chasinglogic/tsk/task"
	"github.com/spf13/cobra"
)

var New = &cobra.Command{
	Use:   "new",
	Short: "create a new task",
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
