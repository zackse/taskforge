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
	"time"

	"github.com/chasinglogic/taskforge/list"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/spf13/cobra"
)

var (
	makeTimesLocal = false
	fixIds         = false
	duplicateIds   = false
)

func init() {
	doctor.Flags().BoolVarP(&makeTimesLocal, "local-times", "l", false,
		"make all times current local time")
	doctor.Flags().BoolVarP(&fixIds, "fix-ids", "i", false,
		"verify all IDs are ObjectIds")
	doctor.Flags().BoolVarP(&duplicateIds, "duplicate-ids", "d", false,
		"generate new IDs for objects with duplicates")
}

var doctor = &cobra.Command{
	Use:     "doctor",
	Aliases: []string{},
	Short:   "Fix task metadata. This can be used to cleanup common issues.",
	Run: func(cmd *cobra.Command, args []string) {
		l, err := config.list()
		if err != nil {
			fmt.Println("ERROR Unable to load list:", err)
			os.Exit(1)
		}

		tasks, err := l.Slice()
		if err != nil {
			fmt.Println("ERROR Unable to retrieve tasks:", err)
			os.Exit(1)
		}

		duplicate := map[string]struct{}{}
		changed := false
		for i := range tasks {
			if tasks[i].CreatedDate.IsZero() {
				fmt.Println("Found invalid created date.")
				changed = true

				if tasks[i].IsComplete() {
					fmt.Println("Task is completed setting created to completed date.")
					tasks[i].CreatedDate = tasks[i].CompletedDate
				} else {
					fmt.Println("Task is incomplete, setting created date to today.")
					tasks[i].CreatedDate = time.Now()
				}
			}

			if makeTimesLocal {
				if tasks[i].CreatedDate.Location().String() != time.Local.String() {
					tasks[i].CreatedDate = tasks[i].CreatedDate.Local()
					changed = true
				}

				if tasks[i].CompletedDate.Location().String() != time.Local.String() {
					tasks[i].CompletedDate = tasks[i].CompletedDate.Local()
					changed = true
				}
			}

			if fixIds {
				_, err := objectid.FromHex(tasks[i].ID)
				if err != nil {
					fmt.Println("Task has an invalid object id, regenerating.")
					tasks[i].ID = objectid.New().Hex()
					changed = true
				}
			}

			if duplicateIds {
				if _, ok := duplicate[tasks[i].ID]; ok {
					fmt.Println("Task has duplicate id, regenerating.")
					tasks[i].ID = objectid.New().Hex()
					changed = true
				} else {
					duplicate[tasks[i].ID] = struct{}{}
				}
			}
		}

		if !changed {
			fmt.Println("Nothing to change!")
			os.Exit(0)
		}

		switch l.(type) {
		case *list.File:
			fileList := l.(*list.File)
			fileList.MemoryList = list.MemoryList(tasks)
			fileList.Save()
		}
	},
}
