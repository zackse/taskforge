use super::error::{Error, Result};
use clap::{App, AppSettings, Arg, ArgMatches, SubCommand};
use std::io;
use std::io::Write;
use taskhero::config::Config;

pub fn command<'a, 'b>() -> App<'a, 'b> {
    SubCommand::with_name("search")
        .about("Show your todos")
        .help("A blank search will list all todos")
        .arg(
            Arg::with_name("completed")
                .long("completed")
                .short("d")
                .help("Show completed tasks in search"),
        )
        .arg(
            Arg::with_name("context")
                .short("c")
                .takes_value(true)
                .value_name("CONTEXT")
                .help("A simple and faster way to limit results to a single context"),
        )
        .arg(
            Arg::with_name("query")
                .multiple(true)
                .help("Query to search for tasks"),
        )
        .setting(AppSettings::TrailingVarArg)
}

pub fn search(config: &mut Config, args: &ArgMatches) -> Result<()> {
    let _query = match args.values_of("title") {
        Some(words) => Some(
            words
                .map(|s| s.to_string())
                .collect::<Vec<String>>()
                .join(" "),
        ),
        None => None,
    };

    let list = match args.value_of("context") {
        Some(context) => config.state.context(context),
        None => config.state.clone(),
    };

    // TODO: Add querying here.
    for (i, task) in list
        .tasks
        .iter()
        .filter(|x| !x.completed() || args.is_present("completed"))
        .enumerate()
    {
        writeln!(&mut io::stdout(), "{} | {}", i + 1, task.title)
            .map(|_| ())
            .map_err(|e| Error::from(e))?;
    }

    Ok(())
}
