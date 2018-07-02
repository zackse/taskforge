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
use clap::{App, AppSettings, Arg, ArgMatches, SubCommand};
use config::Config;
use taskhero::tasks::Task;

pub fn command<'a, 'b>() -> App<'a, 'b> {
    SubCommand::with_name("complete")
        .about("Complete tasks")
        .arg(Arg::with_name("task").multiple(true).help(
            "Title or ID of task to complete, if not provided the current task will be completed",
        ))
        .setting(AppSettings::TrailingVarArg)
}

pub fn complete(config: &mut Config, args: &ArgMatches) -> Result<()> {
    {
        let task: &mut Task = match args.values_of("task") {
            Some(mut words) => {
                if let Ok(id) = words.nth(0).unwrap().parse::<usize>() {
                    if id == 0 {
                        return Err(Error::new(
                            ErrorKind::InvalidArg("0".to_string()),
                            "Cannot use 0 as a task ID",
                        ));
                    }

                    config.state.find_by_ind(id - 1)
                } else {
                    match config.state.find_by_title(
                        &words
                            .map(|s| s.to_string())
                            .collect::<Vec<String>>()
                            .join(" "),
                    ) {
                        Some(task) => task,
                        None => {
                            return Err(Error::new(
                                ErrorKind::InvalidArg("".to_string()),
                                "Unable to find task with that title",
                            ));
                        }
                    }
                }
            }
            None => config.state.current(),
        };

        task.complete();
    }
    config.save().map_err(|e| Error::from(e))
}
