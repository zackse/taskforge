# tsk

This document describes the goals of tsk, and it's design.

# Scope

## Goals

### MVP

- Task management library that can be used with multiple frontends
  - Supports querying tasks
    - [Query Language Design](#query-language)
    - [Task Design](#task-data)
    - [List Design](#task-lists)
    - [Backend Design](#backends)
- Supports saving and loading of tasks from multiple backends
  - **MVP:** Only supports one backend: [LocalFile](#local-backend)
- Query Language supports both simple string searches and complex field based queries
  - For Example:
    - "milk and sugar" will search through all tasks for the words "milk and sugar" in the title, body, and notes of tasks
    - "title = WORK" will find tasks whose title is the string WORK
- Task management CLI as the first usable frontend
  - Supports use of different backends via a configuration file
  - Follows Unix best practices
  - Uses a subcommand interface the following commands will be supported:
    - `new` [design](#new-subcommand)
    - `note` [design](#note-subcommand)
    - `complete` [design](#complete-subcommand)
    - `query` [design](#query-subcommand)
    - `next` [design](#next-subcommand)
    - `edit` [design](#edit-subcommand)

### Future Work / Ideas

Future ideas and features I will implement are as follows:

- Additional Backends:
  - SQLite
  - Postgres
  - MongoDB
  - S3 / Object Storage
- REST API Server built on tsk\_lib
- GUI Frontends (QT is a good choice)
- Modifier statements on queries such as `LIMIT` or `ORDER BY`
- Configurable canned queries
- Task custom fields

# Design

## tsk\_lib

`tsk_lib` is the shared library which holds all of the "business logic" for tsk.
Anything which does not relate to the presentation of the data is contained
within tsk\_lib.

### Query Language

The query language for tasks will accept two "modes".

The first mode is a simple string search. A query which takes the form: `WORD^`
such as: `milk and sugar` is a simple string search. It will be a single
expression which is a String Literal of "milk and sugar" if an interpreter finds
a single String Literal expression as the root node of an AST then it should do
a fuzzy string match on the title, body, and notes of Tasks.

The second mode is a tree of infix expressions. An infix expression as the form
`FIELD_NAME OPERATOR VALUE` or, in the case of logical operators, `EXPRESSION
AND_OR EXPRESSION`. As in math and programming languages, the `(` and `)`
characters will denote "grouped" expressions which will push them down the order
of operations so they are evaluated first. The final kind of expression is with
logical operators. A string literal inside of parens will be treated as an
`EXPRESSION` will be interpreted as a fuzzy search as above so that
`(STRING_LITERAL) AND_OR EXPRESSION` can be used for convenience.

There are three "literal" expressions to represent values:

- String Literal
- Number Literal
- Date Literal

Dates are any string which have the format `YYYY-MM-DD HH:MM (AM|PM)` whether
quoted or unquoted.

Important notes:

- `FIELD_NAME` is lexed as a String Literal token. The parser will validate that it
  is a valid field name if it is part of an infix expression. Else the
  parser will concat multiple `FIELD_NAMES` into a single String Literal
  expression.  
- All numbers are lexed as floats, however in a query string both 5 and 5.0
  are valid.

There is only one "prefix operator" and is not seen by the parser or interpreter
(so is not a true operator at all), and that is `-`. This is used during
unquoted string queries to indicate that a keyword should be interpreted
literally. Take our above example of `milk and sugar`. The lexer would normally
interpret this as `STRING AND STRING`. If we instead want this to be taken as
`STRING STRING STRING` we must put a `-` in front of and. This means the final
query is `milk -and sugar`. The `-` is simply ignored by anything other than the
lexer.

Valid infix operators are:

- `=` equality so that `title = foo` means if title is equal to "foo"
- `!=` or `^=` negative equality so that `title != foo` means find titles which are not equal to "foo"
  The `!=` form is preferred but the `!` character is troublesome in a
  shell environment so the `^=` form is provided as a convenience.
- `>` and `>=` Greater than and Greater than or equal to so that `priority >
  5` means priority is greater than `5.0` similarly `priority >= 5` simply
  includes 5.0 as a valid value.
- `<` and `<=` Less than and Less than or equal to. The inverse of the above.
- `^` or `~` A "LIKE" operator for strings, performs fuzzy matching instead
  of strict equality. The `~` is the preferred form however is inconvenient
  for terminal use so `^` is also valid.
- `AND` or `and` both the upper and lower case forms of `and` are acceptable.
  These perform a logical and of two expressions.
- `OR` or `or` both the upper and lower case forms of `or` are acceptable.
  These perform a logical or of two expressions.

Some example queries with literate explanations of interpreter behavior:

- `title = "take out the trash"`
  - Find all tasks which have the title "take out the trash"
- `title ^ "take out the trash"`
  - Find all tasks whose title contains the string "take out the trash"
- `("milk and sugar") and priority > 5`
  - Find all tasks which have the string "milk and sugar" fuzzy matched on
    their title, body, and notes. Additionally verify that they have a
    priority greater than 5.0
- `milk -and sugar`
  - Find all tasks which have the string "milk and sugar" fuzzy matched on
    their title, body, and notes.
- `(priority > 5 and title ^ "take out the trash") or (context = "work" and (priority >= 2 or ("my little pony")))`
  - Find all tasks which either have a priority greater than 5.0 and a title
    containing the string "take out the trash" or which are in the context
    "work" and have a priority greater than or equal to 2 or have the string
    "my little pony" in their title, body, and notes.

#### AST

The AST for the query language returned by the parser will simply be a struct
which contains a single Expression. It will implement `Iterator` so that the
nodes of the AST can easily be traversed and each List will not need to
implement it's own AST traversal code.

### Task Data

The pseudo-code representation of a task is:

```json
{
    id: String,
    title: String,
    context: String
    created_date: Date,
    completed_date: Option<Date>,
    body: Option<String>,
    priority: Float,
    notes: Vec<Note>,
}
```

The ID of a task will be it's title and created\_date joined using a `:`
character and hashed using MD5. More pseudo-code:

```text
md5(title + ":" + created_date.to_string())
```

A Note will be represented as:

```json
{
    created_date: Date,
    body: String,
}
```

Task will implement the following traits and methods:

- `impl From<&'a str>`
- `impl From<String>`
- `impl std::fmt::Display`
- `impl Debug`
- `impl Clone`
- `new(title: &str) -> Task`
  - `with_x` where x is the remaining task fields not in new
- `impl PartialOrd`
- `impl Serialize` from serde
- `impl Deserialize` from serde

### Task Lists

List will be a trait that will be implemented by all concrete Backends.
Additionally List will be implemented for `Vec<Task>` and `Vec<&Task>`.

It has the following definition:

```rust
// Return a new List which has all completed task if yes_or_no is true and all
// uncompleted tasks if yes_or_no is false.
fn completed(&mut self, yes_or_no: bool) -> Box<List>
// Return a new List with only tasks in the given context
fn with_context(&mut self) -> Box<List>
// Evaluate the AST and return a List of matching results
fn search(&mut self, ast: tsk_lib::query::ast::AST) -> Box<List>
// Add a task to the List
fn add(&mut self, task: tsk_lib::task::Task) -> Result<(), BackendError>
// Add multiple tasks to the List, should be more efficient resource
// utilization.
fn add_multiple(&mut self, task: &Vec<Task>) -> Result<(), BackendError>
// Return a vector of Tasks in this List
fn into_vec(&mut self) -> Vec<Task>
// Find a task by ID
fn find_by_id(&mut self, id: &str) -> Option<Task>
// Return the current task, meaning the oldest uncompleted task in the List
fn current(&mut self) -> Option<Task>

// Complete a task by id
fn complete(&mut self, id: &str) -> Result<(), BackendError>
// Update a task in the list, finding the original by the ID of the given task
fn update(&mut self, task: Task) -> Result<(), BackendError>
// Add note to a task by ID
fn add_note(&mut self, id: &str, note: Note) -> Result<(), BackendError>
```

### Backends

Backend is a trait which is implemented by all of the concrete Backends. It has
the following pseudo-code definition:

```rust
fn save(&self) -> Result<(), BackendError>
fn save_list(&self, Box<List>) -> Result<(), BackendError>
fn load(&self) -> Result<(), BackendError>
fn load_list(&self, Box<List>) -> Result<(), BackendError>
```

BackendError will be an enum that has the following variants:

```text
NotFound
Serialization(String)
Network(std::io::Error)
IO(std::io::Error)
Other(String)
```

It will implement the following:

- `impl From<std::io::Error>`
- `impl From<String>`
- `impl From<T>` where T is all of the serde\_x libraries errors used in tsk

The design of BackendError is subject to change.

## tsk (CLI)

### Configuration File

tsk will be configured with YAML.

TODO: Write this section

### New Subcommand

The new subcommand will accept the following flags:

- `--body TEXT_BODY` populates the task body
- `--priority PRIORITY_NUMBER` populates the task priority
- `--context CONTEXT` populates the task context
- `--from-file PATH_TO_FILE` loads a task/s from a yaml file or csv file

It takes VarArgs and concats them into the title of a new task. So that:

```bash
tsk new write a design doc
```

Will create a task with the title "write a design doc". Flags as described above
can be passed to populate other fields of the Task. Otherwise the flag fields
will get the defaults described below:

- body: None
- priority: 0.0
- context: "default"

The VarArgs are ignored if `--from-file` is provided. If the file is a .csv then
new will assume it is a CSV with the following format:

```csv
title,body,context,priority
record_title,record_body,record_context,record_priority
```

The order of the columns is not important. Only the title, priority, and context
columns are required. Values can be omitted for the optional comments for any
record which does require them.

### Note Subcommand

The note subcommand takes no flags and one argument: the ID of the task to add
a note to. So that:

```bash
tsk note TASK_ID
```

Opens up your `$EDITOR` and allows you to input text that will then be used as
the body of a note which is attached to the task.

### Complete Subcommand

The complete subcommand takes no flags and one argument: the ID of the task to
add a note to. So that:

```bash
tsk complete TASK_ID
```

Will complete the task indicated by `TASK_ID`.

### Query Subcommand

The query subcommand takes the following flags:

- `--completed` a convenience flag to show completed results
- `--csv` print results as a CSV
- `--raw` print no decoration on task table (i.e. remove the "|" and "-"
  characters)
- `--id-only` print only matching task ID's

It takes VarArgs and concatenates them into a query using the [Query
Language](#query-language) parser. It then prints each task in a table using the
following format:

```text
--------------------------------------------
| ID      | Created Date      | Title      |
--------------------------------------------
| TASK_ID | TASK_CREATED_DATE | TASK_TITLE |
--------------------------------------------
```

If raw is given:

```text
ID      Created Date      Title
TASK_ID TASK_CREATED_DATE TASK_TITLE
```

If ID is given only a newline separated list of TASK\IDs are printed with no
headers.

### Next Subcommand

The next subcommand takes the following flags:  

- `--title-only` print only the task title
- `--id-only` print only the task id

But it takes no arguments. It returns the item currently at the "top" of the
list (sorted by oldest date and highest priority). It prints it like so:

```text
TASK_ID TASK_CREATED_DATE TASK_TITLE
```

If title or id only flags are given then only that field is printed.

### Edit Subcommand

The edit subcommand takes one argument: the task ID. It opens the indicated task
in `$EDITOR` as a yaml file and includes all fields from the task. Upon saving
and exiting the file will be read, parsed, and the task will be updated with
that info.
