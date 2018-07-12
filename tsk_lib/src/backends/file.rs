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

use list::List;
use serde_json;

use std::collections::HashMap;
use std::fs;
use std::io;
use std::io::{Read, Write};
use std::path::PathBuf;

use backend::Backend;

/// A simple JSON file saving backend
pub struct FileBackend {
    config: HashMap<String, String>,
    tasks: Vec<Task>,
}

impl FileBackend {
    pub fn new() -> FileBackend {
        let mut config = HashMap::new();

        config.set("dir", match env::var("TASK_DIR") {
            Ok(task_dir) => task_dir,
            Err(_) => match env::var("HOME") {
                Some(mut home) => {
                    home.push(".tasks.d");
                    home
                }
                None => ".tasks.d".to_string(),
            },
        });


        FileBackend{
            config: config,
            tasks: Vec::new(),
        }
    }

    pub fn with_config(mut self, config: HashMap<String, String>) -> FileBackend {
        self.config = config;
        self
    }

    pub fn file(config: &HashMap<String, String>) -> Result<PathBuf, io::Error> {
        match config.get("dir") {
            Some(dir) => {
                let mut path = PathBuf::from(dir);
                path.push("state.json");
                Ok(path)
            }
            None => Err(io::Error::new(
                io::ErrorKind::NotFound,
                "dir not found in backend config",
            )),
        }
    }
}

impl Backend for FileBackend {
    fn save(&self, config: &HashMap<String, String>, list: List) -> Result<(), io::Error> {
        let file = Backend::file(config)?;
        let list = serde_json::to_string_pretty(&list)?;
        fs::File::create(file)?.write_all(list.as_bytes())
    }

    fn load(&self, config: &HashMap<String, String>) -> Result<List, io::Error> {
        let mut file = fs::File::open(Backend::file(config)?)?;
        let mut contents = String::new();
        file.read_to_string(&mut contents)?;
        serde_json::from_str(&contents)
            .map_err(|e| io::Error::new(io::ErrorKind::Other, format!("{}", e)))
    }
}

impl 
