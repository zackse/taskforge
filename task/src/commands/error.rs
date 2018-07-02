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

use std::fmt;
use std::io;
use std::result;

#[derive(Debug, Clone)]
pub enum ErrorKind {
    IO(io::ErrorKind),
    InvalidArg(String),
    InvalidCommand(String),
}

#[derive(Debug, Clone)]
pub struct Error {
    kind: ErrorKind,
    message: String,
}

impl Error {
    pub fn new(kind: ErrorKind, msg: &str) -> Error {
        Error {
            kind: kind,
            message: msg.to_string(),
        }
    }

    pub fn kind(&self) -> ErrorKind {
        self.kind.clone()
    }

    pub fn ignoreable(&self) -> bool {
        match self.kind {
            ErrorKind::IO(io::ErrorKind::BrokenPipe) => true,
            _ => false,
        }
    }

    pub fn code(&self) -> i32 {
        match self.kind {
            ErrorKind::IO(_) => 126,
            ErrorKind::InvalidArg(_) => 128,
            ErrorKind::InvalidCommand(_) => 127,
        }
    }
}

impl From<io::Error> for Error {
    fn from(err: io::Error) -> Error {
        Error {
            kind: ErrorKind::IO(err.kind()),
            message: format!("{}", err),
        }
    }
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self.kind.clone() {
            ErrorKind::IO(_) => write!(f, "{}", self.message),
            ErrorKind::InvalidArg(arg) => write!(f, "{}: {}", self.message, arg),
            ErrorKind::InvalidCommand(cmd) => write!(f, "{}: {}", self.message, cmd),
        }
    }
}

pub type Result<T> = result::Result<T, Error>;
