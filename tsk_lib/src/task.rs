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

use chrono::prelude::*;
use md5::Digest;
use md5::Md5;
use serde_json;

use std::cmp::{Ord, Ordering, PartialOrd};
use std::fmt;
use std::str;

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq)]
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

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq)]
pub struct Task {
    pub id: String,
    pub title: String,
    pub context: String,
    pub created_date: DateTime<Local>,
    pub priority: f64,
    pub notes: Vec<Note>,
    pub completed_date: Option<DateTime<Local>>,
    pub body: Option<String>,
}

impl Task {
    pub fn new(title: &str) -> Task {
        let mut id = title.to_string();
        let created_date = Local::now();

        id.push_str(":");
        id.push_str(&format!("{}", created_date));

        let mut hasher = Md5::default();
        hasher.input(id.as_bytes());

        Task {
            // TODO: Remove this unwrap
            id: str::from_utf8(hasher.result().as_slice())
                .unwrap()
                .to_string(),
            title: title.to_string(),
            context: "default".to_string(),
            created_date: created_date,
            priority: 0.0,
            completed_date: None,
            body: None,
            notes: Vec::new(),
        }
    }

    pub fn with_context(mut self, context: &str) -> Task {
        self.context = context.to_string();
        self
    }

    pub fn with_priority(mut self, priority: f64) -> Task {
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

    pub fn update(&mut self, mut task: Task) {
        self.title = task.title;
        self.context = task.context;
        self.priority = task.priority;
        self.completed_date = task.completed_date;
        self.body = task.body;
        self.notes.append(&mut task.notes);
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

impl Eq for Task {}

impl Ord for Task {
    fn cmp(&self, other: &Task) -> Ordering {
        match self.partial_cmp(other) {
            Some(order) => order,
            None => Ordering::Less,
        }
    }
}

impl PartialOrd for Task {
    fn partial_cmp(&self, other: &Task) -> Option<Ordering> {
        match self.priority.partial_cmp(&other.priority) {
            Some(Ordering::Equal) => Some(self.created_date.date().cmp(&other.created_date.date())),
            Some(order) => Some(order.reverse()),
            None => None,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::prelude::*;
    use chrono::Duration;

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

    #[test]
    fn test_simple_task_ordering() {
        let list = List::new(vec![
            Task::new("test1").with_priority(1),
            Task::new("test3").with_priority(3),
            Task::new("test2").with_priority(2),
            Task::new("test0"),
        ]);

        let mut priority = 3;
        for task in list {
            assert_eq!(task.priority, priority);
            priority = priority - 1;
        }
    }

    #[test]
    fn test_multi_day_task_ordering() {
        let mut yesterday = Task::new("test2").with_priority(2);
        yesterday.created_date = Local::now() - Duration::days(1);

        let tasks = vec![
            Task::new("test1").with_priority(1),
            yesterday,
            Task::new("test3").with_priority(3),
            Task::new("test0").with_priority(2),
        ];

        let list = List::new(tasks.clone());

        let mut iter = list.into_iter();
        assert_eq!(iter.next().unwrap(), tasks[1]);
        assert_eq!(iter.next().unwrap(), tasks[2]);
        assert_eq!(iter.next().unwrap(), tasks[0]);
        assert_eq!(iter.next().unwrap(), tasks[3]);
    }
}
