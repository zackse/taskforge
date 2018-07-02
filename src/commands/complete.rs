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
