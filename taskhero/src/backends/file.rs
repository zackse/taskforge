use list::List;
use serde_json;

use std::collections::HashMap;
use std::fs;
use std::io;
use std::io::{Read, Write};
use std::path::PathBuf;

/// A simple JSON file saving backend
pub struct Backend {}

impl Backend {
    pub fn new() -> Backend {
        Backend {}
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

impl super::Backend for Backend {
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
