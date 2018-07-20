package backends

import (
	"fmt"

	"github.com/chasinglogic/tsk/task"
)

type Backend interface {
	task.List

	Init() error
	Save() error
	Load() error
}

func GetByName(name string) (Backend, error) {
	switch name {
	case "file":
		return File{}, nil
	default:
		return nil, fmt.Errorf("no backend found with name %s", name)
	}
}
