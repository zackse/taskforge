"""
Decorators and configuration file loading for the CLI
"""

import os
import sys

import toml

from ..lists.sqlite import SQLiteList

CONFIG_FILES = [
    'taskforge.toml',
    os.path.join(os.getenv('HOME'), '.taskforge.d', 'config.toml'),
    '/etc/taskforge.d/config.toml'
]

def default_config():
    """Returns a dict with the default config values"""
    return {
        'list': {
            'name': 'sqlite',
            'config': {
                'directory': '~/.taskforge.d'
            }
        },
        'server': {
            'port': 8080,
            'list': {
                'name': 'sqlite',
                'config': {
                    'directory': '~/.taskforge.d'
                }
            }
        }
    }

def load_config():
    """Loads the config file from the default locations"""
    for filename in CONFIG_FILES:
        if os.path.isfile(filename):
            with open(filename) as config_file:
                return toml.load(config_file)
    return default_config()


def config(func):
    """Load config and inject it as the keyword argument cfg"""
    def wrapper(*args, **kwargs):
        kwargs['cfg'] = load_config()
        return func(*args, **kwargs)
    return wrapper


LISTS = {
    'sqlite': SQLiteList,
    'file': SQLiteList,
}

def load_list(cfg):
    """Load the correct List implementation based on the provided config"""
    impl = LISTS.get(cfg['list']['name'])
    if impl is None:
        print('Unknown list: {}'.format(cfg['list']['name']))
        sys.exit(1)

    return impl(**cfg['list']['config'])


def inject_list(func):
    """Injects a keyword argument task_list which has the configured
    list loaded"""
    @config
    def wrapper(*args, **kwargs):
        kwargs['task_list'] = load_list(kwargs['cfg'])
        del kwargs['cfg']
        return func(*args, **kwargs)
    return wrapper
