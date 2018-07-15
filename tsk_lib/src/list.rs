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

use backend::BackendError;
use query;
use task::{Note, Task};

/// List should be implemented by any collection of Tasks
pub trait List<'a> {
    /// Return a new List which has all completed task if yes_or_no is true and all
    /// uncompleted tasks if yes_or_no is false.
    fn completed(&mut self, yes_or_no: bool) -> Vec<Task>;
    /// Return a new List with only tasks in the given context
    fn with_context(&mut self, context: &str) -> Vec<Task>;
    /// Evaluate the AST and return a List of matching results
    fn search(&mut self, ast: query::ast::AST) -> Vec<Task>;
    /// Add a task to the List
    fn add(&mut self, task: Task) -> Result<(), BackendError>;
    /// Add multiple tasks to the List, should be more efficient resource
    /// utilization.
    fn add_multiple(&mut self, task: &mut Vec<Task>) -> Result<(), BackendError>;
    /// Return a vector of Tasks in this List
    fn into_vec(&mut self) -> Vec<Task>;
    /// Find a task by ID
    fn find_by_id(&mut self, id: &str) -> Option<&mut Task>;
    /// Return the current task, meaning the oldest uncompleted task in the List
    fn current(&mut self) -> Option<&mut Task>;

    /// Complete a task by id
    fn complete(&mut self, id: &str) -> Result<(), BackendError>;
    /// Update a task in the list, finding the original by the ID of the given task
    fn update(&mut self, task: Task) -> Result<(), BackendError>;
    /// Add note to a task by ID
    fn add_note(&mut self, id: &str, note: Note) -> Result<(), BackendError>;
}

impl<'a> List<'a> for Vec<Task> {
    fn completed(&mut self, yes_or_no: bool) -> Vec<Task> {
        self.iter()
            .filter(|t| !t.completed_date.is_none() && yes_or_no)
            .cloned()
            .collect()
    }

    fn with_context(&mut self, context: &str) -> Vec<Task> {
        self.iter()
            .filter(|t| t.context.as_str() == context)
            .cloned()
            .collect()
    }

    // TODO: implement
    fn search(&mut self, _ast: query::ast::AST) -> Self {
        self.clone()
    }

    fn add(&mut self, task: Task) -> Result<(), BackendError> {
        self.push(task);
        self.sort();
        Ok(())
    }

    fn add_multiple(&mut self, tasks: &mut Vec<Task>) -> Result<(), BackendError> {
        self.append(tasks);
        self.sort();
        Ok(())
    }

    fn into_vec(&mut self) -> Vec<Task> {
        self.clone()
    }

    fn find_by_id(&mut self, id: &str) -> Option<&mut Task> {
        for task in self {
            if task.id.as_str() == id {
                return Some(task);
            }
        }

        None
    }

    fn current(&mut self) -> Option<&mut Task> {
        self.sort();

        for task in self {
            if task.completed_date.is_none() {
                return Some(task);
            }
        }

        None
    }

    fn complete(&mut self, id: &str) -> Result<(), BackendError> {
        for task in self {
            if task.id.as_str() == id {
                match task.completed_date {
                    Some(_) => return Err(BackendError::from("Already completed.")),
                    None => {
                        task.completed_date = Some(Local::now());
                        return Ok(());
                    }
                }
            }
        }

        Err(BackendError::NotFound)
    }

    fn update(&mut self, new_task: Task) -> Result<(), BackendError> {
        for task in self.iter_mut() {
            if task.id == new_task.id {
                task.update(new_task.clone());
                return Ok(());
            }
        }

        Err(BackendError::NotFound)
    }

    fn add_note(&mut self, id: &str, note: Note) -> Result<(), BackendError> {
        for mut t in self.iter_mut() {
            if t.id.as_str() == id {
                t.notes.push(note.clone());
                return Ok(());
            }
        }

        Err(BackendError::NotFound)
    }
}

#[cfg(test)]
pub mod tests {
    use super::*;
    use task::Task;

    #[test]
    fn test_completed() {
        let mut task2 = Task::new("test2");
        task2.complete();

        let mut tasks = vec![Task::new("test1"), task2, Task::new("test3")];

        let completed = tasks.completed(true);

        assert_eq!(completed.len(), 1);
        assert_eq!(completed[0], tasks[1])
    }

    #[test]
    fn test_with_context() {
        let mut tasks = vec![
            Task::new("test1"),
            Task::new("test2").with_context("testing"),
            Task::new("test3"),
        ];

        let testing = tasks.with_context("testing");

        assert_eq!(testing.len(), 1);
        assert_eq!(testing[0], tasks[1])
    }

    //     #[test]
    //     fn test_search() {}

    #[test]
    fn test_current() {
        let mut task1 = Task::new("test1");
        task1.complete();

        let mut tasks = vec![
            task1,
            Task::new("test2"),
            Task::new("test3").with_priority(2.0),
        ];

        let mut task2 = tasks[2].clone();
        assert_eq!(tasks.current().unwrap(), &mut task2)
    }

    #[test]
    fn test_find_by_id() {
        let mut tasks = vec![Task::new("test1"), Task::new("test2"), Task::new("test3")];
        let id = tasks[2].id.clone();
        let task2 = tasks[2].clone();
        assert_eq!(task2, *tasks.find_by_id(&id).unwrap())
    }
}
