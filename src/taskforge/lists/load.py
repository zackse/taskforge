"""Implements loading lists from entry_points

List plugins are created as python modules using entry_points. The entry_point
name should be a human readable name. It is used when loading a list such as
'sqlite' or 'mongodb'. When load is called we should get back a class which
inherits  from the abstract class 'List' in taskforge.lists.
"""

import pkg_resources


def __entry_points():
    return pkg_resources.iter_entry_points('taskforge.lists')


def get_all_lists():
    """Returns a list of Tuples of list names to class objects"""
    return [(mod.name, mod.load().List) for mod in __entry_points()]


def get_list(name):
    """Return the list implementation which corresponds to name.

    Returns None if not found.
    """
    for mod in __entry_points():
        if mod.name == name:
            return mod.load().List
    return None
