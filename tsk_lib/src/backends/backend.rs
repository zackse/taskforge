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

use list::List;

use std::collections::HashMap;
use std::io;

pub trait Backend {
    /// Save the list using the given config.
    fn save(&self, config: &HashMap<String, String>, list: List) -> Result<(), io::Error>;
    /// Load the list using the given config.
    fn load(&self, config: &HashMap<String, String>) -> Result<List, io::Error>;
}
