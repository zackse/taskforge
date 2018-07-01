use chrono::prelude::*;
use serde_json;
use std::cmp::{Ord, Ordering, PartialOrd};
use std::fmt;

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq, Eq)]
pub struct Note {
    pub created_date: DateTime<Local>,
    pub body: String,
}

impl fmt::Display for Note {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match serde_json::to_string_pretty(self) {
            Ok(jsn) => write!(f, "{}", jsn),
            Err(_) => write!(f, "{}", self.body),
        }
    }
}

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq, Eq)]
pub struct Task {
    pub title: String,
    pub context: String,
    pub created_date: DateTime<Local>,
    pub priority: i64,
    pub notes: Vec<Note>,
    pub completed_date: Option<DateTime<Local>>,
    pub body: Option<String>,
}

impl Task {
    pub fn new(title: &str) -> Task {
        Task {
            title: title.to_string(),
            context: "default".to_string(),
            created_date: Local::now(),
            priority: 0,
            completed_date: None,
            body: None,
            notes: Vec::new(),
        }
    }

    pub fn with_context(mut self, context: &str) -> Task {
        self.context = context.to_string();
        self
    }

    pub fn with_priority(mut self, priority: i64) -> Task {
        self.priority = priority;
        self
    }

    pub fn with_body(mut self, body: &str) -> Task {
        self.body = Some(body.to_string());
        self
    }

    pub fn complete(&mut self) {
        self.completed_date = Some(Local::now());
    }

    pub fn completed(&self) -> bool {
        match self.completed_date {
            Some(_) => true,
            None => false,
        }
    }

    pub fn add_note(&mut self, message: &str) {
        self.notes.push(Note {
            body: message.to_string(),
            created_date: Local::now(),
        })
    }

    pub fn to_json(&self) -> Result<String, serde_json::Error> {
        serde_json::to_string_pretty(self)
    }
}

impl fmt::Display for Task {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match serde_json::to_string_pretty(self) {
            Ok(jsn) => write!(f, "{}", jsn),
            Err(_) => write!(f, "{}", self.title),
        }
    }
}

impl PartialOrd for Task {
    fn partial_cmp(&self, other: &Task) -> Option<Ordering> {
        Some(self.cmp(other))
    }
}

impl Ord for Task {
    fn cmp(&self, other: &Task) -> Ordering {
        match self.created_date.date().cmp(&other.created_date.date()) {
            Ordering::Equal => self.priority.cmp(&other.priority).reverse(),
            order => order,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_new_task() {
        let task = Task::new("Some title");
        assert!(task.title == "Some title".to_string());
        assert!(task.completed_date.is_none());
        assert!(!task.completed());
    }

    #[test]
    fn test_complete_task() {
        let mut task = Task::new("Test");
        task.complete();
        assert!(task.completed_date.is_some());
        assert!(task.completed());
    }
}
