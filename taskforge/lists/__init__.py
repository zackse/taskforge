from abc import ABC, abstractmethod


class InvalidConfigError(Exception):
    """Indicates an invalid configuration was provided to the List"""
    pass


class NotFoundError(Exception):
    """Indicates a task with the given id does not exist"""

    def __init__(self, id):
        self.id

    def __repr__(self):
        return 'no task with id {} exists'.format(self.id)


class List(ABC):
    """An abstract base class that all list implementations but derive
    from."""

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
