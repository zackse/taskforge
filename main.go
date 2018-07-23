package main

import (
	"fmt"

	"github.com/chasinglogic/tsk/commands"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unkown"
)

func init() {
	commands.Root.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tsk-%s-%s\n", commit, version)
		fmt.Println(`license`)
	},
}

func main() {
	commands.Root.Execute()
}
