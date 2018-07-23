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


package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	Root.AddCommand(new)
	Root.AddCommand(query)
	Root.AddCommand(complete)
	Root.AddCommand(edit)
	Root.AddCommand(next)
}

// Root is the root CLI command
var Root = &cobra.Command{
	Use:   "tsk",
	Short: "Manage your tasks",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		file := findConfigFile()

		var err error
		config, err = loadConfigFile(file)
		if err != nil {
			fmt.Println("ERROR Unable to load config file:", err)
			os.Exit(1)
		}
	},
}
