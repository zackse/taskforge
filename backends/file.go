package backends

import (
	"os"
	"path/filepath"

	"github.com/chasinglogic/tsk/task"
)

type File struct {
	Dir string

	*task.MemoryList `yaml:"-" json:"-"`
}

func (f File) Init() error {
	if f.Dir != "" {
		return nil
	}

	taskDir := os.Getenv("TASK_DIR")
	if taskDir == "" {
		f.Dir = filepath.Join(os.Getenv("HOME"), ".tasks.d")
		return nil
	}

	f.Dir = taskDir
	return nil
}

func (f File) Save() error {
	return nil
}

func (f File) Load() error {
	return nil
}
