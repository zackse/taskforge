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
