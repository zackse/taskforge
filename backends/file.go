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

	task.MemoryList `yaml:"-" json:"-"`
}

func (f *File) Init() error {
	if f.MemoryList != nil {
		return nil
	}

	f.MemoryList = make([]task.Task, 0)
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
