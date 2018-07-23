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


package task

import (
	"fmt"
	"testing"
	"time"
)

func compareLists(list1, list2 MemoryList) bool {
	if len(list1) != len(list2) {
		return false
	}

	for i := range list1 {
		if list1[i].ID != list2[i].ID {
			return false
		}
	}

	return true
}

func TestMemoryList_Completed(t *testing.T) {
	type args struct {
		completed bool
	}

	fixture := MemoryList{
		Task{Title: "completed 1", CompletedDate: time.Now()},
		Task{Title: "not completed"},
		Task{Title: "completed 2", CompletedDate: time.Now()},
	}

	tests := []struct {
		name string
		ml   MemoryList
		args args
		want []Task
	}{
		{
			name: "find completed",
			ml:   fixture,
			args: args{true},
			want: []Task{
				fixture[0],
				fixture[2],
			},
		},
		{
			name: "find incomplete",
			ml:   fixture,
			args: args{false},
			want: []Task{
				fixture[1],
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ml.Completed(tt.args.completed)

			if !compareLists(got, tt.want) {
				t.Errorf("MemoryList.Completed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryList_Context(t *testing.T) {
	otherContext := New("other context")
	otherContext.Context = "other"

	fixture := MemoryList{
		New("default context"),
		otherContext,
		New("default context"),
	}

	type args struct {
		context string
	}

	tests := []struct {
		name string
		ml   MemoryList
		args args
		want []Task
	}{
		{
			name: "other context",
			ml:   fixture,
			args: args{"other"},
			want: []Task{
				otherContext,
			},
		},
		{
			name: "default context",
			ml:   fixture,
			args: args{"default"},
			want: []Task{
				fixture[0],
				fixture[2],
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ml.Context(tt.args.context)
			if len(got) != len(tt.want) {
				t.Errorf("MemoryList.Completed() = %v, want %v", got, tt.want)
				return
			}

			for i := range got {
				if got[i].ID != tt.want[i].ID {
					t.Errorf("MemoryList.Completed() = %v, want %v", got, tt.want)
					return
				}
			}
		})
	}
}

func TestMemoryList_Add(t *testing.T) {
	list := MemoryList{
		New("task 1"),
		New("task 2"),
	}

	task3 := New("task 3")

	err := list.Add(task3)
	if err != nil {
		t.Error(err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 tasks got %d: %v", len(list), list)
		return
	}

	if list[2].ID != task3.ID {
		t.Errorf("Expected %v Got %v", task3, list[2])
	}
}

func TestMemoryList_AddMultiple(t *testing.T) {
	list := MemoryList{
		New("task 1"),
		New("task 2"),
	}

	task3 := New("task 3")
	task4 := New("task 4")

	other := []Task{
		task3,
		task4,
	}

	err := list.AddMultiple(other)
	if err != nil {
		t.Error(err)
	}

	if len(list) != 4 {
		t.Errorf("Expected 4 tasks got %d: %v", len(list), list)
		return
	}

	if list[2].ID != task3.ID {
		t.Errorf("Expected %v Got %v", task3, list[2])
	}

	if list[3].ID != task4.ID {
		t.Errorf("Expected %v Got %v", task4, list[3])
	}
}

func TestMemoryList_FindByID(t *testing.T) {
	task2 := New("task 2")
	list := MemoryList{
		New("task 1"),
		task2,
		New("task 3"),
	}

	found, _ := list.FindByID(task2.ID)
	if task2.ID != found.ID || task2.Title != found.Title {
		t.Errorf("MemoryList.FindByID() = %v, want %v", found, task2)
	}
}

func TestMemoryList_Current(t *testing.T) {
	completed := New("completed task")
	completed.Complete()
	current := New("current task")

	firstUnComplete := MemoryList{
		current,
		completed,
	}

	lastUncomplete := MemoryList{
		completed,
		current,
	}

	tests := []struct {
		name      string
		ml        MemoryList
		want      Task
		shouldErr bool
	}{
		{
			name: "all completed",
			ml: MemoryList{
				completed,
				completed,
			},
			shouldErr: true,
		},
		{
			name: "first uncompleted",
			ml:   firstUnComplete,
			want: current,
		},
		{
			name: "last uncompleted",
			ml:   lastUncomplete,
			want: current,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ml.Current()
			if err != nil && !tt.shouldErr {
				t.Errorf("MemoryList.Current() = got error: %s when expected %v", err, tt.want)
				return
			} else if err != nil && tt.shouldErr {
				return
			} else if err == nil && tt.shouldErr {
				t.Errorf("MemoryList.Current() = got no error when expected one got: %v", got)
				return
			}

			if tt.want.ID != got.ID {
				t.Errorf("MemoryList.Current() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemoryList_Complete(t *testing.T) {
	fixture := MemoryList{
		New("task 1"),
		New("task 2"),
		New("task to complete"),
		New("task 4"),
	}

	id := fixture[2].ID

	err := fixture.Complete(id)
	if err != nil {
		t.Errorf("Error completing: %s", err)
		return
	}

	fmt.Println(fixture)

	if fixture[2].CompletedDate.IsZero() {
		t.Errorf("Expected task to complete to have a completed date got: %v", fixture)
	}
}

func TestMemoryList_Update(t *testing.T) {
	toUpdate := New("task to update")

	list := MemoryList{
		New("task 1"),
		toUpdate,
		New("task 3"),
	}

	other := toUpdate
	other.Title = "task updated"

	err := list.Update(other)
	if err != nil {
		t.Errorf("Got an error updating: %s", err)
		return
	}

	if list[1].Title != "task updated" {
		t.Errorf("Expected title \"task updated\" got: %v", list[1])
	}
}

func TestMemoryList_AddNote(t *testing.T) {
	task := New("add a note to me")

	list := MemoryList{
		New("task 1"),
		task,
		New("task 3"),
	}

	note := NewNote("this is a note")

	err := list.AddNote(task.ID, note)
	if err != nil {
		t.Errorf("Got an err when expected none %s: %v", err, list)
		return
	}

	if len(list[1].Notes) != 1 {
		t.Errorf("Expected one note got: %v", list[1])
	}
}
