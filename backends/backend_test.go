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


package backends

import (
	"reflect"
	"testing"

	"github.com/chasinglogic/tsk/ql/lexer"
	"github.com/chasinglogic/tsk/ql/parser"
	"github.com/chasinglogic/tsk/task"
	"github.com/mitchellh/mapstructure"
)

type backendTest struct {
	config  Config
	backend Backend
	empty   Backend
	setup   func(Backend) error
	cleanup func(Backend) error
}

func compareTasks(task1, task2 task.Task) bool {
	return task1.Title == task2.Title && task1.ID == task2.ID
}

func compareTaskSlice(tasks, other []task.Task) bool {
	for i := range tasks {
		if !compareTasks(tasks[i], other[i]) {
			return false
		}
	}

	return true
}

func runBackendTest(t *testing.T, test backendTest) {
	backend := test.backend
	empty := test.empty
	if test.setup != nil {
		test.setup(backend)
	}

	err := mapstructure.Decode(test.config, &backend)
	if err != nil {
		t.Errorf("Error decoding config: %s", err)
		return
	}

	err = mapstructure.Decode(test.config, &empty)
	if err != nil {
		t.Errorf("Error decoding config to empty: %s", err)
		return
	}

	if err := backend.Init(); err != nil {
		t.Errorf("Error running backend init: %s", err)
		return
	}

	if err := empty.Init(); err != nil {
		t.Errorf("Error running empty init: %s", err)
		return
	}

	if test.cleanup != nil {
		defer test.cleanup(backend)
	}

	list := []task.Task{
		task.New("task 1"),
		task.New("task 2"),
		task.New("task 3"),
	}

	if err := backend.AddMultiple(list); err != nil {
		t.Errorf("Error adding multiple: %s", err)
		return
	}

	task4 := task.New("task 4")
	if err := backend.Add(task4); err != nil {
		t.Errorf("Error adding one: %s", err)
		return
	}

	expected := []task.Task{
		list[0],
		list[1],
		list[2],
		task4,
	}
	if slice := backend.Slice(); !reflect.DeepEqual(backend.Slice(), expected) {
		t.Errorf("Slice failed, expected: %v got :%v", expected, slice)
		return
	}

	if err := backend.Save(); err != nil {
		t.Errorf("Error saving: %s", err)
		return
	}

	if reflect.DeepEqual(empty.Slice(), expected) {
		t.Errorf("Unexpected state, empty backend should not match loaded backend.")
		return
	}

	if err := empty.Load(); err != nil {
		t.Errorf("Unable to load state into empty backend: %s", err)
		return
	}

	if slice := empty.Slice(); reflect.DeepEqual(slice, expected) {
		t.Errorf("Loaded state does not match, expected %v got %v", expected, slice)
		return
	}

	taskInOtherContext := task.New("other context task")
	taskInOtherContext.Context = "other"
	if err := backend.Add(taskInOtherContext); err != nil {
		t.Errorf("Error adding one: %s", err)
		return
	}

	expected = []task.Task{
		taskInOtherContext,
	}
	if slice := backend.Context("other"); !reflect.DeepEqual(slice, expected) {
		t.Errorf("Context failed, expected: %v got :%v", expected, slice)
		return
	}

	if task, err := backend.FindByID(taskInOtherContext.ID); err != nil || !compareTasks(task, taskInOtherContext) {
		t.Errorf("FindByID failed, expected %v got %v", taskInOtherContext, task)
		t.Errorf("FindByID err: %s", err)
		return
	}

	toUpdate := task.New("task to update")
	backend.Add(toUpdate)

	toUpdate.Title = "task updated"
	if err := backend.Update(toUpdate); err != nil {
		t.Errorf("Error updating: %s", err)
		return
	}

	updated, err := backend.FindByID(toUpdate.ID)
	if err != nil {
		t.Errorf("Error finding updated by id: %s", err)
		return
	}

	if updated.Title != "task updated" {
		t.Errorf("Task was not updated, expected %v got %v", toUpdate, updated)
		return
	}

	current, err := backend.Current()
	if err != nil {
		t.Errorf("Error getting current task: %s", err)
		return
	}

	if err := backend.Complete(current.ID); err != nil {
		t.Errorf("Error completing task: %s", err)
		return
	}

	newCurrent, err := backend.Current()
	if err != nil {
		t.Errorf("Error getting current task: %s", err)
		return
	}

	if compareTasks(newCurrent, current) {
		t.Errorf("Expected a new current task but got the same.")
		return
	}

	expected = []task.Task{current}
	if completedTasks := backend.Completed(true); !compareTaskSlice(completedTasks, expected) {
		t.Errorf("Expected only completed tasks %v got %v", expected, completedTasks)
		return
	}

	if err := backend.AddNote(current.ID, task.NewNote("this is a note")); err != nil {
		t.Errorf("Error adding note: %s", err)
		return
	}

	notedTask, err := backend.FindByID(current.ID)
	if err != nil {
		t.Errorf("Error retrieving by id: %s", err)
		return
	}

	if len(notedTask.Notes) != 1 {
		t.Errorf("Expected 1 note got: %v", notedTask)
		return
	}

	queryTests := []struct {
		query    string
		expected []string
	}{
		{
			query: "context = \"other\"",
			expected: []string{
				"other context task",
			},
		},
		{
			query: "title = \"task 4\" or title = \"task 1\" or title = \"task updated\"",
			expected: []string{
				"task 1",
				"task 4",
				"task updated",
			},
		},
		{
			query: "(title = \"task 1\" and context = \"default\") or (context = \"other\")",
			expected: []string{
				"task 1",
				"other context task",
			},
		},
		{
			query:    "context = \"work\"",
			expected: []string{},
		},
		{
			query: "task",
			expected: []string{
				"task 1",
				"task 2",
				"task 3",
				"task 4",
				"task updated",
				"other context task",
			},
		},
		{
			query: "priority = 0",
			expected: []string{
				"task 1",
				"task 2",
				"task 3",
				"task 4",
				"task updated",
				"other context task",
			},
		},
		{
			query:    "priority > 1",
			expected: []string{},
		},
	}

	for _, test := range queryTests {
		t.Run(test.query, func(t *testing.T) {
			p := parser.New(lexer.New(test.query))
			ast := p.Parse()

			if err := p.Error(); err != nil {
				t.Errorf("error parsing query: %s", err)
				return
			}

			list, err := backend.Search(ast)
			if err != nil {
				t.Errorf("error searching backend: %s", err)
				return
			}

			if len(list) != len(test.expected) {
				t.Errorf("Expected %v Got %v", test.expected, list)
			}

		expectedTitles:
			for x := range test.expected {
				for y := range list {
					if list[y].Title == test.expected[x] {
						continue expectedTitles
					}
				}

				t.Errorf("Expected %v Got %v", test.expected, list)
			}
		})
	}
}
