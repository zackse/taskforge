//+build !windows

package commands

import (
	"os"
	"path/filepath"
)

var configFilePaths = []string{
	"Taskforge.toml",
	"taskforge.toml",
	filepath.Join(os.Getenv("HOME"), ".taskforge.d", "config.toml"),
	"/etc/taskforge/config.toml",
}

func defaultDir() string {
	return filepath.Join(os.Getenv("HOME"), "taskforge.d")
}
