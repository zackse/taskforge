#!/usr/bin/env python3
"""usage: python build_docs.py

Build the docs including README for taskforge.
"""

import glob
import os
import re
import pkgutil
from datetime import datetime

import jinja2

import taskforge.cli as clipkg


def get_template_files():
    """Find all jinja2 template files"""
    return glob.iglob(
        '{}/**/*.j2'.format(
            os.path.join(os.path.dirname(__file__), 'templates')),
        recursive=True)


def get_module_docs(package):
    """Return a list of all docstrings found for modules in package recursively"""
    prefix = package.__name__ + "."
    docs = [package.__doc__]
    for importer, modname, _ispkg in pkgutil.walk_packages(
            package.__path__, prefix):
        docs.append(importer.find_module(modname).load_module(modname).__doc__)
    return docs


CMD_RGX = re.compile(r"usage: (task ([a-z]{1,})?)")


def get_command_name(usage_string):
    """Get the command name from usage string."""
    matches = CMD_RGX.match(usage_string)
    return matches.group(1)


def build_context():
    """Build the jinja2 context used for rendering templates"""
    # some hard coded values
    context = {
        'package_name':
        'taskforge-cli',
        'version':
        '0.2.0',
        'current_year':
        datetime.now().strftime("%Y"),
        'contributing_doc_link':
        'https://github.com/chasinglogic/taskforge/blob/master/CONTRIBUTING.md',
        'designs_link':
        'https://github.com/chasinglogic/taskforge/blob/master/docs/designs',
    }

    usage_strings = [
        doc for doc in get_module_docs(clipkg) if doc.startswith('usage:')
    ]

    context['usage_strings'] = [{
        'command': get_command_name(usage_string),
        'usage': usage_string,
    } for usage_string in usage_strings]

    return context


def main():
    templates = get_template_files()
    context = build_context()
    for template in templates:
        with open(template) as tmpl:
            rendered = jinja2.Template(tmpl.read()).render(**context)

        targetfile = template.replace('scripts/templates/', '').replace(
            '.j2', '')
        with open(targetfile, 'w') as target:
            target.write(rendered)


if __name__ == '__main__':
    main()
