# tsk, a task management cli

It's like [taskwarrior](https://taskwarrior.org) but flexible in different ways.

## Installation

### Install from Release

1. Navigate to [the Releases Page](https://github.com/chasinglogic/dfm/releases)
2. Find the tar ball for your platform / architecture. For example, on 64 bit
   Mac OSX, the archive is named `dfm_{version}_darwin_amd64.tar.gz`
3. Extract the tar ball
4. Put the dfm binary in your `$PATH`

### Install from Source

Simply run go get:

```bash
$ go get github.com/chasinglogic/dfm
```

If your `$GOPATH/bin` is in your `$PATH` then you now have dfm installed.

## Usage

```text
Manage your tasks

Usage:
  tsk [command]

Available Commands:
  complete    Complete tasks by ID
  edit        Edit a task as YAML
  help        Help about any command
  new         Create a new task
  next        Show the current task
  query       Search and list tasks
  version     print version information

Flags:
  -h, --help   help for tsk

Use "tsk [command] --help" for more information about a command.
```

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. :fire: Submit a pull request :D :fire:

All pull requests should go to the develop branch not master. Thanks!

See the [DESIGN.md](https://github.com/chasinglogic/tsk/blob/master/DESIGN.md)
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