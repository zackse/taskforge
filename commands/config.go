package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/chasinglogic/tsk/backends"
	"github.com/mitchellh/mapstructure"
)

func loadConfigFile(path string) (*Config, error) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	} else if err != nil {
		c := defaultConfig()

		yml, err := yaml.Marshal(c)
		if err != nil {
			return c, err
		}

		return c, ioutil.WriteFile(path, yml, 0644)
	}

	var c *Config
	err = yaml.Unmarshal(content, &c)
	return c, err
}

func defaultConfig() *Config {
	return &Config{
		Backend: "file",
		BackendConfig: map[string]interface{}{
			"dir": filepath.Join(os.Getenv("HOME"), ".tasks.d"),
		},
	}
}

func findConfigFile() string {
	possiblePaths := []string{
		".tsk.yml",
		filepath.Join(os.Getenv("HOME"), ".tasks.d", "tsk.yml"),
		"/etc/tsk/tsk.yml",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return filepath.Join(os.Getenv("HOME"), ".tasks.d", "tsk.yml")
}

var config *Config

type Config struct {
	Backend       string
	BackendConfig map[string]interface{} `yaml:"backend_config" json:"backend_config"`

	backendImpl backends.Backend
}

func (c *Config) backend() (backends.Backend, error) {
	if c.backendImpl == nil {
		fmt.Println("loading backend")

		var err error
		c.backendImpl, err = backends.GetByName(c.Backend)
		if err != nil {
			return nil, err
		}

		err = mapstructure.Decode(c.BackendConfig, &c.backendImpl)
		if err != nil {
			return nil, err
		}

		if c.backendImpl == nil {
			return nil, fmt.Errorf("backend %s didn't initialize but returned no error", c.Backend)
		}

		err = c.backendImpl.Init()
		if err != nil {
			return nil, err
		}
	}

	fmt.Println(c.backendImpl)

	return c.backendImpl, c.backendImpl.Load()
}
