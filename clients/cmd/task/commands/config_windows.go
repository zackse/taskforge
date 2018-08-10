package commands

import (
	"os"
	"path/filepath"
)

var configFilePaths = []string{
	"Taskforge.toml",
	filepath.Join(os.Getenv("APPDATA"), "Taskforge", "config.toml"),
	"C:\\Taskforge\\config.toml",
}

func defaultDir() string {
	return filepath.Join(os.Getenv("APPDATA"), "taskforge.d")
}
