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

use super::task::Task;
use serde_json;

use std::fmt;
use std::fmt::Display;
use std::iter;
use std::slice;

#[derive(Debug, Serialize, Deserialize, Clone, PartialEq)]
pub struct List {
    pub tasks: Vec<Task>,
}

impl IntoIterator for List {
    type Item = Task;
    type IntoIter = ::std::vec::IntoIter<Task>;

    fn into_iter(self) -> Self::IntoIter {
        self.tasks.into_iter()
    }
}

impl Display for List {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match serde_json::to_string_pretty(self) {
            Ok(jsn) => write!(f, "{}", jsn),
            Err(_) => write!(f, "ERROR: Unable to serialize list"),
        }
    }
}

impl List {
    pub fn new(tasks: Vec<Task>) -> List {
        let mut list = List { tasks: tasks };

        list.tasks.sort();
        list
    }

    /// Return a reference to the "current" task.
    pub fn current<'a>(&'a mut self) -> Option<&'a mut Task> {
        let mut ind = 0;
        for (task_id, task) in self.enumerate() {
            if !task.completed() {
                ind = task_id;
                break;
            }
        }

        self.find_by_ind(ind)
    }

    /// Add a task to the List, will sort after adding.
    pub fn add(&mut self, task: Task) {
        self.tasks.push(task);
        self.tasks.sort();
    }

    /// Add multiple tasks to the List, this is more efficient than calling add multiple times
    /// since only one sort is performed. It will empty the given vector.
    pub fn add_multiple(&mut self, tasks: &mut Vec<Task>) {
        self.tasks.append(tasks);
        self.tasks.sort();
    }

    /// Return a reference to the task at the given ID / index.
    pub fn find_by_ind<'a>(&'a mut self, id: usize) -> Option<&'a mut Task> {
        if self.tasks.len() < id {
            return None;
        }

        let t: &'a mut Task = &mut self.tasks[id];
        Some(t)
    }

    /// Return a reference to the first task with the given title.
    pub fn find_by_title<'a>(&'a mut self, title: &str) -> Option<&'a mut Task> {
        let mut ind = None;
        for (i, task) in self.enumerate() {
            if task.title == title {
                ind = Some(i);
                break;
            }
        }

        match ind {
            Some(i) => self.find_by_ind(i),
            None => None,
        }
    }

    /// Return an enumerated iterator over the tasks in this list.
    pub fn enumerate(&self) -> iter::Enumerate<slice::Iter<Task>> {
        self.tasks.iter().enumerate()
    }

    pub fn to_json(&self) -> Result<String, serde_json::Error> {
        serde_json::to_string_pretty(self)
    }

    pub fn context(&self, context: &str) -> List {
        List::new(
            self.tasks
                .iter()
                .filter(|x| &x.context == context)
                .map(|x| x.clone())
                .collect(),
        )
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use chrono::prelude::*;
    use chrono::Duration;

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
            Task::new("test0"),
        ];

        let list = List::new(tasks.clone());

        let mut iter = list.into_iter();
        assert_eq!(iter.next().unwrap(), tasks[1]);
        assert_eq!(iter.next().unwrap(), tasks[2]);
        assert_eq!(iter.next().unwrap(), tasks[0]);
        assert_eq!(iter.next().unwrap(), tasks[3]);
    }
}
