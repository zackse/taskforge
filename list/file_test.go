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

package list

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/chasinglogic/taskforge/task"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

func TestFileList(t *testing.T) {
	for _, test := range ListTests {
		config := Config{
			"directory": fmt.Sprintf(".test.file.%s", test.Name),
		}

		l := &File{}

		require.Nil(t, mapstructure.Decode(config, &l))
		require.Nil(t, l.Init())

		t.Run(test.Name, func(t *testing.T) {
			defer os.RemoveAll(fmt.Sprintf(".test.file.%s", test.Name))
			require.Nil(t, test.Test(l))
		})
	}
}

func TestFileListPersistence(t *testing.T) {
	testDir := ".test.file.persistence"
	l := &File{
		Dir: testDir,
	}

	require.Nil(t, l.Init())
	defer os.RemoveAll(testDir)

	tasks := []task.Task{
		task.New("task 1"),
		task.New("task 2"),
		task.New("task 3"),
	}

	require.Nil(t, l.AddMultiple(tasks))
	require.Nil(t, l.Complete(tasks[1].ID))
	tasks, err := l.Slice()
	require.Nil(t, err)

	other := &File{
		Dir: testDir,
	}

	require.Nil(t, other.Init())
	slice, err := other.Slice()

	jsn1, _ := json.MarshalIndent(tasks, "", "\t")
	jsn2, _ := json.MarshalIndent(slice, "", "\t")

	require.Nil(t, err)

	for i := range slice {
		if slice[i].ID != tasks[i].ID {
			t.Errorf("Expected %s Got %s", jsn1, jsn2)
			return
		}

		if slice[i].IsCompleted() && !tasks[i].IsCompleted() {
			t.Errorf("Expected %s Got %s", jsn1, jsn2)
			return
		}
	}
}

func BenchmarkFindWhereIDAtEnd(b *testing.B) {
	tasks := MemoryList{
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
		task.New("task 1"),
	}

	idToFind := tasks[len(tasks)-1].ID
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tasks.findWhere(func(t task.Task) bool {
			return t.ID == idToFind
		})
	}
}
