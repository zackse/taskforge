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

use super::error::{Error, ErrorKind, Result};
use clap::{App, Arg, ArgMatches, SubCommand};
use config::Config;
use std::io;
use std::io::Write;
use tsk_lib::tasks::Task;

pub fn command<'a, 'b>() -> App<'a, 'b> {
    SubCommand::with_name("show")
        .alias("s")
        .alias("n")
        .alias("next")
        .about("Show your todos")
        .help("A blank ID will show the \"current\" or \"next\" task.")
        .arg(
            Arg::with_name("json")
                .long("json")
                .short("j")
                .help("Print the task as JSON"),
        )
        .arg(
            Arg::with_name("short")
                .long("short")
                .short("s")
                .alias("title")
                .alias("t")
                .help("Print only the task title"),
        )
        .arg(Arg::with_name("task_id").help("Task ID to show, default is current task"))
}

fn pretty_print(task: &Task, args: &ArgMatches) -> Result<()> {
    if args.is_present("json") {
        writeln!(
            &mut io::stdout(),
            "{}",
            task.to_json().map_err(|e| Error::from(e))?
        ).map_err(|e| Error::from(e))?;
        return Ok(());
    }

    if args.is_present("short") {
        writeln!(&mut io::stdout(), "{}", task.title).map_err(|e| Error::from(e))?;
        return Ok(());
    }

    let mut pretty = format!(
        "Title: {}
Project: {}
Created Date: {}
Priority: {}
Completed: {}",
        task.title,
        task.context,
        task.created_date.to_string(),
        task.priority,
        task.completed()
    );

    if let Some(completed_date) = task.completed_date {
        pretty.push_str(&format!("\nCompleted Date: {}", completed_date.to_string()));
    }

    if let Some(ref body) = task.body {
        pretty.push_str(&format!("\nBody: {}", body));
    }

    if task.notes.len() > 0 {
        pretty.push_str("\nNotes:\n");
        for note in &task.notes {
            pretty.push_str(&format!(
                "\tCreated Date: {}\n\tBody: {}\n",
                note.created_date.to_string(),
                note.body
            ));
        }
    }

    writeln!(&mut io::stdout(), "{}", pretty).map_err(|e| Error::from(e))
}

pub fn show(config: &mut Config, args: &ArgMatches) -> Result<()> {
    let task = match args.value_of("task_id") {
        Some(id) => config
            .state
            .find_by_ind(
                id.parse::<usize>().map_err(|_| {
                    Error::new(
                        ErrorKind::InvalidArg("Not a number".to_string()),
                        "Invalid argument",
                    )
                })? - 1,
            )
            .ok_or(Error::new(
                ErrorKind::InvalidArg("".to_string()),
                "Could not find a task with that ID",
            ))?,
        None => config.state.current().ok_or(Error::new(
            ErrorKind::InvalidArg("".to_string()),
            "No current task to show. Provide an index to show a completed task.",
        ))?,
    };

    pretty_print(task, args)
}
