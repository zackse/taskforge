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

use serde_json;
use std::env;
use std::fs;
use std::io;
use std::io::prelude::*;
use std::path::PathBuf;
use taskhero::tasks::list::List;

#[derive(Debug, Deserialize, Serialize)]
pub struct Config {
    pub state: List,
}

impl Config {
    fn dir() -> Result<PathBuf, io::Error> {
        let cfg_dir = match env::var("TASKHERO_DIR") {
            Ok(task_dir) => PathBuf::from(task_dir),
            Err(_) => match env::home_dir() {
                Some(mut home) => {
                    home.push(".taskhero");
                    home
                }
                None => PathBuf::from(".taskhero"),
            },
        };

        if !cfg_dir.exists() {
            fs::create_dir_all(&cfg_dir)?;
        }

        Ok(cfg_dir)
    }

    fn state_file() -> Result<PathBuf, io::Error> {
        let mut dir = Config::dir()?;
        dir.push("state.json");
        Ok(dir)
    }

    pub fn new() -> Config {
        Config {
            state: List::new(Vec::new()),
        }
    }

    pub fn save(&self) -> Result<(), io::Error> {
        let list = serde_json::to_string_pretty(&self.state)?;
        let file = Config::state_file()?;
        fs::File::create(file)?.write_all(list.as_bytes())
    }

    pub fn load() -> Result<Config, io::Error> {
        let mut config = Config::new();
        let mut contents = String::new();
        Config::state_file()
            .and_then(|path| fs::File::open(path))
            .and_then(|mut f| f.read_to_string(&mut contents))?;

        config.state = serde_json::from_str(&contents).map_err(|e| serde_json::Error::from(e))?;
        Ok(config)
    }
}
