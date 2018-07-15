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

use std::collections::HashMap;
use std::env;
use std::fs;
use std::io;
use std::io::{Read, Write};
use std::path::PathBuf;

use backend::{Backend, BackendError};
use list::List;
use query::ast::AST;
use task::{Note, Task};

/// A simple JSON file saving backend
pub struct FileBackend {
    tasks: Vec<Task>,
}

impl FileBackend {
    pub fn new() -> FileBackend {
        FileBackend { tasks: Vec::new() }
    }

    pub fn default_config() -> HashMap<String, String> {
        let mut config = HashMap::new();

        config.insert(
            "dir".to_string(),
            match env::var("TASK_DIR") {
                Ok(task_dir) => task_dir,
                Err(_) => match env::var("HOME") {
                    Ok(mut home) => {
                        home.push_str(".tasks.d");
                        home
                    }
                    Err(_) => ".tasks.d".to_string(),
                },
            },
        );

        config
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

impl<'a> List<'a> for FileBackend {
    fn completed(&mut self, yes_or_no: bool) -> Vec<Task> {
        self.tasks.completed(yes_or_no)
    }

    fn with_context(&mut self, context: &str) -> Vec<Task> {
        self.tasks.with_context(context)
    }

    fn search(&mut self, ast: AST) -> Vec<Task> {
        self.tasks.search(ast)
    }

    fn add(&mut self, task: Task) -> Result<(), BackendError> {
        self.tasks.add(task)
    }

    fn add_multiple(&mut self, tasks: &mut Vec<Task>) -> Result<(), BackendError> {
        self.tasks.add_multiple(tasks)
    }

    fn into_vec(&mut self) -> Vec<Task> {
        self.tasks.into_vec()
    }

    fn find_by_id(&mut self, id: &str) -> Option<&mut Task> {
        self.tasks.find_by_id(id)
    }

    fn current(&mut self) -> Option<&mut Task> {
        self.tasks.current()
    }

    fn complete(&mut self, id: &str) -> Result<(), BackendError> {
        self.tasks.complete(id)
    }

    fn update(&mut self, task: Task) -> Result<(), BackendError> {
        self.tasks.update(task)
    }

    fn add_note(&mut self, id: &str, note: Note) -> Result<(), BackendError> {
        self.tasks.add_note(id, note)
    }
}

impl<'a> From<Vec<Task>> for FileBackend {
    fn from(tasks: Vec<Task>) -> FileBackend {
        FileBackend {
            tasks: tasks.clone(),
        }
    }
}

impl<'a> Backend<'a> for FileBackend {
    fn save(&self, config: &HashMap<String, String>) -> Result<(), BackendError> {
        let file = FileBackend::file(config)?;
        let list = serde_json::to_string_pretty(&self.tasks)?;
        fs::File::create(file)?
            .write_all(list.as_bytes())
            .map(|_| ())
            .map_err(|e| BackendError::from(e))
    }

    fn load(&mut self, config: &HashMap<String, String>) -> Result<(), BackendError> {
        let mut file = fs::File::open(FileBackend::file(config)?)?;
        let mut contents = String::new();
        file.read_to_string(&mut contents)?;
        let tasks = serde_json::from_str(&contents)
            .map_err(|e| io::Error::new(io::ErrorKind::Other, format!("{}", e)))?;
        self.tasks = tasks;
        Ok(())
    }
}
