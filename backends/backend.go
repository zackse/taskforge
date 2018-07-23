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
