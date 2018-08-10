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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
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
		return defaultConfig(), nil
	}

	var c *Config
	err = toml.Unmarshal(content, &c)

	return c, err
}

func defaultConfig() *Config {
	return &Config{
		List: map[string]interface{}{
			"name": "file",
			"dir":  defaultDir(),
		},
	}
}

func findConfigFile() string {
	for _, path := range configFilePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

var config *Config

type ServerConfig struct {
	Port int
	List map[string]interface{}
}

type Config struct {
	Server ServerConfig
	List   map[string]interface{}

	listImpl list.List
}

func (c *Config) list() (list.List, error) {
	if c.listImpl == nil {
		var err error
		name, ok := c.List["name"]
		if !ok {
			return nil, errors.New("must specify a list name")
		}

		c.listImpl, err = list.GetByName(name.(string))
		if err != nil {
			return nil, err
		}

		err = mapstructure.Decode(c.List, &c.listImpl)
		if err != nil {
			return nil, err
		}
	}

	return c.listImpl, c.listImpl.Init()
}
