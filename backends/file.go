package backends

import (
	"encoding/json"
	"io/ioutil"
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
	stateFile := filepath.Join(f.Dir, "state.json")

	jsn, err := json.Marshal(f.MemoryList)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(stateFile, jsn, os.ModePerm)
}

func (f File) Load() error {
	stateFile := filepath.Join(f.Dir, "state.json")

	content, err := ioutil.ReadFile(stateFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(content, &f.MemoryList)
}
