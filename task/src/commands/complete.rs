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
use clap::{App, AppSettings, Arg, ArgMatches, SubCommand, Values};
use config::Config;
use std::process;
use taskhero::tasks::Task;

pub fn command<'a, 'b>() -> App<'a, 'b> {
    SubCommand::with_name("complete")
        .alias("c")
        .about("Complete tasks")
        .arg(Arg::with_name("task").multiple(true).help(
            "Title or ID of task to complete, if not provided the current task will be completed",
        ))
        .setting(AppSettings::TrailingVarArg)
}

fn get_task_from_input<'a>(
    config: &'a mut Config,
    words: &mut Values,
) -> Result<Option<&'a mut Task>> {
    if let Ok(id) = words.nth(0).unwrap().parse::<usize>() {
        if id == 0 {
            return Err(Error::new(
                ErrorKind::InvalidArg("0".to_string()),
                "Cannot use 0 as a task ID",
            ));
        }

        return Ok(config.state.find_by_ind(id - 1));
    }

    let task = config.state.find_by_title(&words
        .map(|s| s.to_string())
        .collect::<Vec<String>>()
        .join(" "));

    if task.is_none() {
        return Err(Error::new(
            ErrorKind::InvalidArg("".to_string()),
            "Unable to find task with that title",
        ));
    }

    Ok(task)
}

pub fn complete(config: &mut Config, args: &ArgMatches) -> Result<()> {
    // Use this inner block to release the mutable borrow of config before we save it
    {
        let mut task = match args.values_of("task") {
            Some(ref mut words) => get_task_from_input(config, words)?,
            None => config.state.current(),
        };

        match task {
            Some(ref mut task) => task.complete(),
            None => {
                println!("No task with that id found");
                process::exit(0);
            }
        }
    }

    config.save().map_err(|e| Error::from(e))
}
