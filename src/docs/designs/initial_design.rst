Taskforge
=========

This document describes the goals of Taskforge the library, and it’s
design.

Goals
-----

-  Task management library that can be used with multiple frontends

   -  Supports querying tasks

      -  `Query Language Design <#query-language>`__
      -  `Task Design <#task-data>`__
      -  `List Design <#lists>`__

-  Supports saving and loading of tasks from multiple “services” via
   List implementations

   -  **MVP:** Only supports one list: `SQLite <#sqlite-list>`__

-  Query Language supports both simple string searches and complex field
   based queries

   -  For Example:

      -  “milk sugar” will search through all tasks for the words “milk
         sugar” in the title, body, and notes of tasks
      -  “title = WORK” will find tasks whose title is the string WORK

Design
------

Query Language
++++++++++++++

The query language for tasks will accept two “modes.”

The first mode is a simple string search. A query which takes the form:
``WORD^`` such as: ``milk and sugar`` is a simple string search. It will
be a single expression which is a String Literal of “milk and sugar” if
an interpreter finds a single String Literal expression as the root node
of an AST then it should do a fuzzy string match on the title, body, and
notes of Tasks.

The second mode is a tree of infix expressions. An infix expression as
the form ``FIELD_NAME OPERATOR VALUE`` or, in the case of logical
operators, ``EXPRESSION AND_OR EXPRESSION``. As in math and programming
languages, the ``(`` and ``)`` characters will denote “grouped”
expressions which will push them down the order of operations so they
are evaluated first. The final kind of expression is with logical
operators. A string literal inside of parens will be treated as an
``EXPRESSION`` will be interpreted as a fuzzy search as above so that
``(STRING_LITERAL) AND_OR EXPRESSION`` can be used for convenience.

There are three “literal” expressions to represent values:

-  String Literal
-  Number Literal
-  Date Literal

Dates are any string which have the format ``YYYY-MM-DD HH:MM (AM|PM)``
whether quoted or unquoted.

Important notes:

-  ``FIELD_NAME`` is lexed as a String Literal token. The parser will
   validate that it is a valid field name if it is part of an infix
   expression. Else the parser will concatenate multiple ``FIELD_NAMES`` into
   a single String Literal expression.
-  All numbers are lexed as floats, however in a query string both 5 and
   5.0 are valid.

Taskforge query language has one “prefix operator” and is not seen by the parser
or interpreter (so is not a true operator at all), and that is ``-``. The lexer
will use this during tokenization of unquoted strings to change what would
normally be a keyword into a string token. Take our above example of ``milk and
sugar``. The lexer would normally interpret this as ``STRING AND STRING``. If we
instead want this to be taken as ``STRING STRING STRING`` we must put a ``-`` in
front of and. This means the final query is ``milk -and sugar``. The ``-`` is
simply ignored by anything other than the lexer.

Valid infix operators are:

-  ``=`` equality so that ``title = foo`` means if title is equal to
   “foo”
-  ``!=`` or ``^=`` negative equality so that ``title != foo`` means
   find titles which are not equal to “foo” The ``!=`` form is preferred
   but the ``!`` character is troublesome in a shell environment so the
   ``^=`` form is provided as a convenience.
-  ``>`` and ``>=`` Greater than and Greater than or equal to so that
   ``priority > 5`` means priority is greater than ``5.0`` similarly
   ``priority >= 5`` simply includes 5.0 as a valid value.
-  ``<`` and ``<=`` Less than and Less than or equal to. The inverse of
   the above.
-  ``^`` or ``+`` A “LIKE” operator for strings, performs fuzzy matching
   instead of strict equality. The ``+`` is the preferred form however
   is inconvenient for terminal use so ``^`` is also valid.
-  ``AND`` or ``and`` both the upper and lower case forms of ``and`` are
   acceptable. These perform a logical and of two expressions.
-  ``OR`` or ``or`` both the upper and lower case forms of ``or`` are
   acceptable. These perform a logical or of two expressions.

Some example queries with literate explanations of interpreter behavior:

-  ``title = "take out the trash"``

   -  Find all tasks which have the title “take out the trash”

-  ``title ^ "take out the trash"``

   -  Find all tasks whose title contains the string “take out the
      trash”

-  ``("milk and sugar") and priority > 5``

   -  Find all tasks which have the string “milk and sugar” fuzzy
      matched on their title, body, and notes. Additionally verify that
      they have a priority greater than 5.0

-  ``milk -and sugar``

   -  Find all tasks which have the string “milk and sugar” fuzzy
      matched on their title, body, and notes.

-  ``(priority > 5 and title ^ "take out the trash") or (context = "work" and (priority >= 2 or ("my little pony")))``

   -  Find all tasks which either have a priority greater than 5.0 and a
      title containing the string “take out the trash” or which are in
      the context “work” and have a priority greater than or equal to 2
      or have the string “my little pony” in their title, body, and
      notes.

AST
^^^

The AST for the query language returned by the parser is a class which
has a single member variable ``expression``.

.. vale off

.. code:: python

   class AST:
       """Abstract syntax tree for the Taskforge query language"""

       def __init__(self, expression):
           self.expression = expression

       def __eq__(self, other):
           return self.expression == other.expression

       def __repr__(self):
           return self.expression.__repr__()

``expression`` is an Expression class object. The Expression class is as
follows:

.. code:: python

   class Expression:
       """An expression is a statement that yields a value"""

       ...implementation details

       def __init__(self, token, left=None, right=None):
           self.token = token
           self.value = None
           self.operator = None
           self.left = None
           self.right = None

           if token.token_type == Type.STRING:
               self.value = token.literal
           elif token.token_type == Type.NUMBER:
               self.value = float(token.literal)
           elif token.token_type == Type.BOOLEAN:
               self.value = bool(token.literal)
           elif token.token_type == Type.DATE:
               self.value = Expression.parse_date(token.literal)
           else:
               self.operator = token
               self.left = left
               self.right = right

       def __repr__(self):
           if self.is_infix() and self.token.token_type in [Type.AND, Type.OR]:
               return '({} {} {})'.format(
                   self.left,
                   self.operator.literal,
                   self.right)
           elif self.is_infix():
               return '({} {} {})'.format(
                   self.left.value
                   if self.left is not None
                   else self.left,
                   self.operator.literal,
                   self.right)
           elif type(self.value) is str:
               return "'{}'".format(self.value)
           else:
               return '{}'.format(self.value)

       def __eq__(self, other):
           if self.is_infix():
               return (other.is_infix() and
                       self.left == other.left and
                       self.operator == other.operator and
                       self.right == other.right)
           else:
               return (self.value == other.value and
                       self.token == other.token)

       def is_infix(self):
           """Indicates whether this expression is an infix expression"""
           return self.operator is not None

       def is_literal(self):
           """Indicates whether this expression is a literal value"""
           return self.value is not None

       def is_comparison_infix(self):
           """Indicates if this is a value comparison expression"""
           return self.is_infix() and not self.is_logical_infix()

       def is_logical_infix(self):
           """Indicates if this is a logical AND/OR expression"""
           return self.is_and_infix() or self.is_or_infix()

       def is_and_infix(self):
           """Indicates if this is a logical AND expression"""
           return (self.is_infix() and
                   self.operator.token_type == Type.AND)

       def is_or_infix(self):
           """Indicates if this is a logical OR expression"""
           return (self.is_infix() and
                   self.operator.token_type == Type.OR)

       def is_str_literal(self):
           """Indicates whether this expression is a string value"""
           return (self.is_literal() and
                   self.token.token_type == Type.STRING)

       def is_date_literal(self):
           """Indicates whether this expression is a date value"""
           return (self.is_literal() and
                   self.token.token_type == Type.DATE)

       def is_number_literal(self):
           """Indicates whether this expression is a number value"""
           return (self.is_literal() and
                   self.token.token_type == Type.NUMBER)

       def is_boolean_literal(self):
           """Indicates whether this expression is a boolean value"""
           return (self.is_literal() and
                   self.token.token_type == Type.BOOLEAN)

.. vale on

Task Data
+++++++++

The pseudo-code representation of a task is:

.. vale off

.. code:: json

   {
       id: String,
       title: String,
       context: String
       created_date: Date,
       completed_date: Date | null,
       body: String,
       priority: Float,
       notes: [Note]
   }

A Note will be represented as:

.. code:: json

   {
       id: String,
       created_date: Date,
       body: String,
   }

.. vale on

All ID’s will be hex strings of python std library uuids regardless of
list storage. This is a nice, 0 dependency, and easy to use UUID that
can be made into a string.

Task Lists
++++++++++

List will be an abstract class which all list implementations will need
to subclass, it has the following definition:

.. code:: python

   class List(ABC):
       """An abstract base class that all list implementations but derive from."""

       @abstractmethod
       def search(self, ast):
           """Evaluate the AST and return a List of matching results"""
           raise NotImplementedError

       @abstractmethod
       def add(self, task):
           """Add a task to the List"""
           raise NotImplementedError

       @abstractmethod
       def add_multiple(self, tasks):
           """Add multiple tasks to the List, should be more efficient
           resource utilization."""
           raise NotImplementedError

       @abstractmethod
       def list(self):
           """Return a python list of the Task in this List"""
           raise NotImplementedError

       @abstractmethod
       def find_by_id(self, id):
           """Find a task by id"""
           raise NotImplementedError

       @abstractmethod
       def current(self):
           """Return the current task, meaning the oldest uncompleted
           task in the List"""
           raise NotImplementedError

       @abstractmethod
       def complete(self, id):
           """Complete a task by id"""
           raise NotImplementedError

       @abstractmethod
       def update(self, task):
           """Update a task in the listist, finding the original by the
           id of the given task"""
           raise NotImplementedError

       @abstractmethod
       def add_note(self, id, note):
           """Add note to a task by id"""
           raise NotImplementedError

Additionally each list will be instantiated with the dictionary of it’s
configuration from the config file using the ``**dictionary`` syntax.
This means that a list will need to implement keyword arguments in it’s
``__init__`` constructor for all configuration items. During this it
will need to check for missing required arguments or invalid
configurations and raise a InvalidConfigError with a human readable
message. Additionally any connecting or loading of files necessary for
use will happen during object construction.

Future Work / Ideas
-------------------

Future ideas and features I will implement are as follows:

-  Additional Lists:

   -  Postgres
   -  MongoDB

-  GUI Frontends (QT is a good choice)
-  Modifier statements on queries such as ``LIMIT`` or ``ORDER BY``
-  Task custom fields
