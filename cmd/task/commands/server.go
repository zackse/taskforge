package commands

import (
	"fmt"
	"os"

	"github.com/chasinglogic/taskforge/server"
	"github.com/spf13/cobra"
)

func init() {
	serverCmd.AddCommand(serverRun)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage a task server",
}

var serverRun = &cobra.Command{
	Use:   "run",
	Short: "Run a task server",
	Run: func(cmd *cobra.Command, args []string) {
		l, err := config.list()
		if err != nil {
			fmt.Println("ERROR Unable to load list:", err)
			os.Exit(1)
		}

		s := server.New(l)
		s.Listen()
	},
}
