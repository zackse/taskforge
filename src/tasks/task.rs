use chrono::prelude::*;
use serde_json;
use std::fmt;
use std::fmt::Display;

#[derive(Debug, Serialize, Deserialize)]
pub struct Note {
    pub created_date: DateTime<Local>,
    pub body: String,
}

impl Display for Note {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match serde_json::to_string_pretty(self) {
            Ok(jsn) => write!(f, "{}", jsn),
            Err(_) => write!(f, "{}", self.body),
        }
    }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Task {
    pub title: String,
    pub context: String,
    pub created_date: DateTime<Local>,
    pub notes: Vec<Note>,
    pub completed_date: Option<DateTime<Local>>,
    pub priority: Option<i64>,
    pub body: Option<String>,
}

impl Task {
    pub fn new(title: String) -> Task {
        Task {
            title: title,
            context: "default".to_string(),
            created_date: Local::now(),
            completed_date: None,
            priority: None,
            body: None,
            notes: Vec::new(),
        }
    }

    pub fn in_context(mut self, context: String) -> Task {
        self.context = context;
        self
    }

    pub fn with_priority(mut self, priority: i64) -> Task {
        self.priority = Some(priority);
        self
    }

    pub fn with_body(mut self, body: String) -> Task {
        self.body = Some(body);
        self
    }

    pub fn complete(&mut self) {
        self.completed_date = Some(Local::now());
    }

    pub fn add_note(mut self, message: String) {
        self.notes.push(Note {
            body: message,
            created_date: Local::now(),
        })
    }
}

impl Display for Task {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match serde_json::to_string_pretty(self) {
            Ok(jsn) => write!(f, "{}", jsn),
            Err(_) => write!(f, "{}", self.title),
        }
    }
}
