use super::tasks::list::List;
use serde_json;
use std::env;
use std::fs;
use std::io;
use std::io::prelude::*;
use std::path::PathBuf;

#[derive(Debug, Deserialize, Serialize)]
pub struct Config {
    pub state: List,
}

impl Config {
    fn dir() -> Result<PathBuf, io::Error> {
        let cfg_dir = match env::home_dir() {
            Some(mut home) => {
                home.push(".taskhero");
                home
            }
            None => PathBuf::from(".taskhero"),
        };

        if !cfg_dir.exists() {
            fs::create_dir_all(&cfg_dir)?;
        }

        Ok(cfg_dir)
    }

    fn file() -> Result<PathBuf, io::Error> {
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
        let config = serde_json::to_string_pretty(self)?;
        let file = Config::file()?;
        fs::File::create(file)?.write_all(config.as_bytes())
    }

    pub fn load() -> Result<Config, io::Error> {
        let mut contents = String::new();
        Config::file()
            .and_then(|path| fs::File::open(path))
            .and_then(|mut f| f.read_to_string(&mut contents))
            .and_then(|_| serde_json::from_str(&contents).map_err(|e| io::Error::from(e)))
    }
}
