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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/taskforge/task"
)

// File is a list which is stored in a JSON file
type File struct {
	Dir string

	MemoryList `yaml:"-" json:"-"`
}

// Init will load tasks from the JSON file if found, otherwise just init an
// empty list
func (f *File) Init() error {
	if f.MemoryList != nil {
		return nil
	}

	if f.MemoryList == nil {
		f.MemoryList = make([]task.Task, 0)
	}

	if strings.HasPrefix(f.Dir, "~") {
		f.Dir = strings.Replace(f.Dir, "~", os.Getenv("HOME"), 1)
	}

	return f.load()
}

// Add a task to the list
func (f *File) Add(t task.Task) error {
	err := f.MemoryList.Add(t)
	if err != nil {
		return err
	}

	return f.save()
}

// AddMultiple tasks to the list
func (f *File) AddMultiple(t []task.Task) error {
	err := f.MemoryList.AddMultiple(t)
	if err != nil {
		return err
	}

	return f.save()
}

// Update a task in the list
func (f *File) Update(t task.Task) error {
	err := f.MemoryList.Update(t)
	if err != nil {
		return err
	}

	return f.save()
}

// Complete a task in the list
func (f *File) Complete(id string) error {
	err := f.MemoryList.Complete(id)
	if err != nil {
		return err
	}

	return f.save()
}

// AddNote to a task in the list
func (f *File) AddNote(id string, note task.Note) error {
	err := f.MemoryList.AddNote(id, note)
	if err != nil {
		return err
	}

	return f.save()
}

func (f *File) save() error {
	if err := os.MkdirAll(f.Dir, os.ModePerm); err != nil {
		return err
	}

	stateFile := filepath.Join(f.Dir, "tasks.json")
	jsn, err := json.Marshal(f.MemoryList)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(stateFile, jsn, 0644)
}

func (f *File) load() error {
	stateFile := filepath.Join(f.Dir, "state.json")

	content, err := ioutil.ReadFile(stateFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err != nil {
		return nil
	}

	return json.Unmarshal(content, &f.MemoryList)
}
