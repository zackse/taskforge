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

use std::io;

use list::List;
use task::Task;

pub enum BackendError {
    NotFound,
    Serialization(String),
    Network(io::Error),
    IO(io::Error),
    Other(String),
}

impl From<io::Error> for BackendError {
    fn from(err: io::Error) -> BackendError {
        match err.kind() {
            io::ErrorKind::ConnectionRefused => BackendError::Network(err),
            io::ErrorKind::ConnectionAborted => BackendError::Network(err),
            io::ErrorKind::NotConnected => BackendError::Network(err),
            io::ErrorKind::AddrInUse => BackendError::Network(err),
            io::ErrorKind::AddrNotAvailable => BackendError::Network(err),
            _ => BackendError::IO(err),
        }
    }
}

impl<'a> From<&'a str> for BackendError {
    fn from(s: &'a str) -> BackendError {
        BackendError::Other(s.to_string())
    }
}

impl From<String> for BackendError {
    fn from(s: String) -> BackendError {
        BackendError::from(s.as_ref())
    }
}

/// Backend is implemented by all structs that know how to save and load lists.
/// Ideally all backends should also implement List
pub trait Backend<'a>: List<'a> + From<Vec<&'a Task>> {
    /// Save the list owned by self.
    fn save(&self) -> Result<(), BackendError>;
    /// Load the list owned by self.
    fn load(&self) -> Result<(), BackendError>;
}
