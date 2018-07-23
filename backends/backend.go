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
	"fmt"

	"github.com/chasinglogic/tsk/task"
)

// Backend is a task.List that supports saving and loading from
// a data source.
type Backend interface {
	task.List

	Init() error
	Save() error
	Load() error
}

// Config is a convenience type used to represent the unstructured config that
// will be stored by the application. All backends should accept being decoded
// using mapstructure from this type. Additionally, the Init() function should
// return an error if a required argument was not provided.
type Config map[string]interface{}

// GetByName returns the appropriate Backend implementation by
// a human redable string name. It returns a user friendly error
// string if not found.
func GetByName(name string) (Backend, error) {
	switch name {
	case "file":
		return &File{}, nil
	default:
		return nil, fmt.Errorf("no backend found with name %s", name)
	}
}
