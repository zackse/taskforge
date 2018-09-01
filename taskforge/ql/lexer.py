from .tokens import Token, Type


class Lexer:
    """Scans input producing tokens"""

    def __init__(self, query):
        self.data = query
        self.pos = 0
        self.read_pos = 0
        self.ch = ''
        self._read_char()

    def __iter__(self):
        return self

    def __next__(self):
        self._skip_whitespace()

        if self.ch == '':
            raise StopIteration
        elif self.ch == '^':
            if self._peek_char() == '=':
                self._read_char()
                token = Token('!=')
            elif self._peek_char() == '^':
                self._read_char()
                token = Token('!~')
            else:
                token = Token('~')
        elif self.ch == '!':
            literal = self.ch
            self._read_char()
            literal += self.ch
            token = Token(literal)
        elif self.ch == '>' or self.ch == '<':
            literal = self.ch
            if self._peek_char() == '=':
                self._read_char()
                literal += self.ch
            token = Token(literal)
        elif self.ch == '"' or self.ch == "'":
            # skip the opening quote
            self._read_char()

            token = Token(self._quoted_string())
            if self.ch != '"' and self.ch != "'":
                token = Token(
                    'unexpected eof: no closing quote',
                    token_type=Type.UNEXPECTED)
        elif self.ch == '-':
            # skip the -
            self._read_char()
            token = Token(self._unquoted_string(), token_type=Type.STRING)
        elif self.ch.isdigit():
            # numbers can be followed by non space characters and
            # cause lexing errors when followed by a non-space
            # character and therefore do not need the additional read
            # at the bottom of this function since it would "skip" characters
            # like )
            return Token(self._number())
        elif self.ch.isalpha():
            # same as above for numbers
            return Token(self._unquoted_string())
        else:
            token = Token(self.ch)

        self._read_char()
        return token

    def next_token(self):
        """Return the next token from the input"""
        return self.__next__()

    def _read_char(self):
        """Read a character from input advancing the cursor"""
        if self.read_pos >= len(self.data):
            self.ch = ''
        else:
            self.ch = self.data[self.read_pos]

        self.pos = self.read_pos
        self.read_pos += 1

    def _peek_char(self):
        """Return the next character"""
        if self.read_pos > len(self.data):
            return ''
        else:
            return self.data[self.read_pos]

    def _read(self, valid):
        """Takes a function which takes a string and returns a boolean
        indicating if we should keep reading. Returns the full string
        which matched the valid function."""
        start = self.pos
        while valid(self.ch):
            self._read_char()

        return self.data[start:self.pos]

    def _skip_whitespace(self):
        while self.ch.isspace():
            self._read_char()

    def _unquoted_string(self):
        return self._read(lambda c: c.isalpha())

    def _number(self):
        return self._read(lambda c: c.isdigit() or c == '.')

    def _quoted_string(self):
        return self._read(lambda c: c != '"' and c != "'")
