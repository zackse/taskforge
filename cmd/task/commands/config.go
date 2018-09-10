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

	"github.com/BurntSushi/toml"
	"github.com/chasinglogic/taskforge/list"
	"github.com/mitchellh/mapstructure"
)

var config *Config

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
		List: ListConfig{
			Name: "file",
			Config: map[string]interface{}{
				"directory": defaultDir(),
			},
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

func loadList(name string, cfg map[string]interface{}) (list.List, error) {
	var err error
	var listImpl list.List

	listImpl, err = list.GetByName(name)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(cfg, &listImpl)
	if err != nil {
		return nil, err
	}

	return listImpl, nil
}

type ListConfig struct {
	Name   string
	Config map[string]interface{}
}

type ServerConfig struct {
	Port int
	Addr string
	List ListConfig
}

func (sc *ServerConfig) list() (list.List, error) {
	l, err := loadList(sc.List.Name, sc.List.Config)
	if err != nil {
		return nil, err
	}

	return l, l.Init()
}

type Config struct {
	Server ServerConfig
	List   ListConfig

	listImpl list.List
}

func (c *Config) list() (list.List, error) {
	if c.listImpl == nil {
		var err error
		c.listImpl, err = loadList(c.List.Name, c.List.Config)
		if err != nil {
			return nil, err
		}
	}

	return c.listImpl, c.listImpl.Init()
}
