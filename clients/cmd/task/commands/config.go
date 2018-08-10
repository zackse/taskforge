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

package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/chasinglogic/taskforge/list"
	"github.com/mitchellh/mapstructure"
)

func loadConfigFile(path string) (*Config, error) {
	if _, err := os.Stat(filepath.Dir(path)); err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0644); err != nil {
			return nil, err
		}
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

	if c != nil && c.DefaultContext == "" {
		c.DefaultContext = "default"
	}

	return c, err
}

func defaultConfig() *Config {
	return &Config{
		DefaultContext: "default",
		List:           "file",
		ListConfig: map[string]interface{}{
			"dir": filepath.Join(os.Getenv("HOME"), ".tasks.d"),
		},
	}
}

func findConfigFile() string {
	possiblePaths := []string{
		".taskforge.yml",
		filepath.Join(os.Getenv("HOME"), ".tasks.d", "config.yml"),
		"/etc/taskforge/config.yml",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return filepath.Join(os.Getenv("HOME"), ".tasks.d", "config.yml")
}

var config *Config

type Config struct {
	DefaultContext string `yaml:"default_context"`
	List           string
	ListConfig     map[string]interface{} `yaml:"list_config" json:"list_config"`

	listImpl list.List
}

func (c *Config) list() (list.List, error) {
	if c.listImpl == nil {
		var err error
		c.listImpl, err = list.GetByName(c.List)
		if err != nil {
			return nil, err
		}

		err = mapstructure.Decode(c.ListConfig, &c.listImpl)
		if err != nil {
			return nil, err
		}
	}

	return c.listImpl, c.listImpl.Init()
}
