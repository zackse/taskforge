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
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/chasinglogic/tsk/task"
)

type File struct {
	Dir string

	task.MemoryList `yaml:"-" json:"-"`
}

func (f *File) Init() error {
	if f.MemoryList != nil {
		return nil
	}

	f.MemoryList = make([]task.Task, 0)

	if strings.HasPrefix(f.Dir, "~") {
		f.Dir = strings.Replace(f.Dir, "~", os.Getenv("HOME"), 1)
	}

	return nil
}

func (f *File) Save() error {
	if err := os.MkdirAll(f.Dir, os.ModePerm); err != nil {
		return err
	}

	stateFile := filepath.Join(f.Dir, "state.json")
	jsn, err := json.Marshal(f.MemoryList)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(stateFile, jsn, 0644)
}

func (f *File) Load() error {
	stateFile := filepath.Join(f.Dir, "state.json")

	content, err := ioutil.ReadFile(stateFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err != nil {
		return nil
	}

	return json.Unmarshal(content, &f.MemoryList)
}
