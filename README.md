# taskforge, a task management cli

Task management tool that integrates well with other services.

## Installation

### Install from Release

For now the only way to install taskforge is via pip:

```
$ pip install taskforge[cli]
```

## Usage

```text
usage: task [-h]
            {next,n,add,new,a,query,q,s,search,list,todo,complete,done,d} ...

positional arguments:
  {next,n,add,new,a,query,q,s,search,list,todo,complete,done,d}
    next (n)            Print the next task in the list
    add (new, a)        Add a new task to the list
    query (q, s, search, list)
                        Query tasks in the list.
    todo                Print incomplete tasks in the list
    complete (done, d)  Complete tasks in the list. If no ID given will
                        complete the current task.

optional arguments:
  -h, --help            show this help message and exit
```

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. :fire: Submit a pull request :D :fire:

All pull requests should go to the develop branch not master. Thanks!

See the [DESIGN.md](https://github.com/chasinglogic/taskforge/blob/master/DESIGN.md)
for more info. Not everything laid out there is implemented yet.

## License

This code is distributed under the GNU General Public License

```text
    Copyright (C) 2018 Mathew Robinson

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
```
