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

use toml;

use std::collections::HashMap;
use std::env;
use std::fs;
use std::io;
use std::io::prelude::*;
use std::path::PathBuf;
use tsk_lib::backends;
use tsk_lib::list::List;

#[derive(Debug, Deserialize, Serialize)]
pub struct BackendConfig {
    name: String,
    config: HashMap<String, String>,
}

impl BackendConfig {
    fn new() -> BackendConfig {
        BackendConfig {
            name: "file".to_string(),
            config: HashMap::new(),
        }
    }

    fn backend(&self) -> Box<backends::Backend> {
        Box::new(match self.name.as_ref() {
            "file" => backends::file::Backend::new(),
            _ => backends::file::Backend::new(),
        })
    }

    fn init_default(&mut self) {
        if let None = self.config.get("dir") {
            self.config.insert(
                "dir".to_string(),
                Config::dir().unwrap().to_string_lossy().to_string(),
            );
        }
    }

    fn init(&mut self) {
        match self.name.as_ref() {
            "file" => self.init_default(),
            _ => self.init_default(),
        }
    }

    fn load(&self) -> Result<List, io::Error> {
        self.backend().load(&self.config)
    }

    fn save(&self, list: List) -> Result<(), io::Error> {
        self.backend().save(&self.config, list)
    }
}

#[derive(Debug, Deserialize, Serialize)]
pub struct Config {
    backend: BackendConfig,

    #[serde(skip_deserializing, skip_serializing)]
    pub state: List,
}

impl Config {
    fn dir() -> Result<PathBuf, io::Error> {
        let cfg_dir = match env::var("TASKHERO_DIR") {
            Ok(task_dir) => PathBuf::from(task_dir),
            Err(_) => match env::home_dir() {
                Some(mut home) => {
                    home.push(".tsk_lib");
                    home
                }
                None => PathBuf::from(".tsk_lib"),
            },
        };

        if !cfg_dir.exists() {
            fs::create_dir_all(&cfg_dir)?;
        }

        Ok(cfg_dir)
    }

    fn file() -> Result<PathBuf, io::Error> {
        let mut dir = Config::dir()?;
        dir.push("config.toml");
        Ok(dir)
    }

    pub fn new() -> Config {
        Config {
            backend: BackendConfig::new(),
            state: List::new(Vec::new()),
        }
    }

    pub fn save(&self) -> Result<(), io::Error> {
        self.backend.save(self.state.clone())?;

        let file = Config::file()?;
        let conf = toml::to_string(self)
            .map_err(|e| io::Error::new(io::ErrorKind::Other, format!("{}", e)))?;
        fs::File::create(file)?.write_all(conf.as_bytes())
    }

    pub fn load() -> Result<Config, io::Error> {
        let mut config = Config::new();
        let mut contents = String::new();
        let file = fs::File::open(Config::file()?);

        match file {
            Ok(mut f) => {
                f.read_to_string(&mut contents)?;
                config = toml::from_str(&contents)
                    .map_err(|e| io::Error::new(io::ErrorKind::Other, format!("{}", e)))?;
            }
            Err(ref e) if e.kind() == io::ErrorKind::NotFound => (),
            Err(e) => return Err(e),
        }

        config.backend.init();
        config.state = config.backend.load()?;

        Ok(config)
    }
}
