// Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"

	"github.com/chasinglogic/taskforge/clients/cmd/task/commands"
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
		fmt.Printf("taskforge-%s-%s\n", commit, version)
		fmt.Println(`license`)
	},
}

func main() {
	commands.Root.Execute()
}
