extern crate clap;
extern crate serde;
extern crate serde_json;
extern crate taskhero;
#[macro_use]
extern crate serde_derive;

pub mod commands;
pub mod config;

use clap::App;
use commands::error::{Error, ErrorKind};
use config::Config;
use std::io;
use std::process;

fn main() {
    let matches = App::new("taskhero")
        .version("0.1.0")
        .author("Mathew Robinson <chasinglogic@gmail.com>")
        .about(
            "
Manage your tasks.

Copyright 2018 Mathew Robinson <chasinglogic@gmail.com>. All rights reserved.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.",
        )
        .subcommand(commands::new::command())
        .subcommand(commands::search::command())
        .subcommand(commands::complete::command())
        .get_matches();

    let mut config = match Config::load() {
        Ok(config) => config,
        Err(err) => {
            if err.kind() != io::ErrorKind::NotFound {
                println!("ERROR: Unable to load config: {}", err);
            }
            Config::new()
        }
    };

    let res = match matches.subcommand() {
        ("new", Some(args)) => commands::new(&mut config, args),
        ("search", Some(args)) => commands::search(&mut config, args),
        ("complete", Some(args)) => commands::complete(&mut config, args),
        (command, _) => Err(Error::new(
            ErrorKind::InvalidCommand(command.to_string()),
            "Unknown command",
        )),
    };

    if let Err(err) = res {
        if err.ignoreable() {
            process::exit(0)
        }

        println!("{}", err);
        process::exit(err.code())
    }
}
