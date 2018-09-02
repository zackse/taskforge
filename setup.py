"""
A task management tool that integrates with 3rd party services
"""
from setuptools import find_packages, setup

with open('requirements.txt') as reqs:
    dependencies = reqs.read().split('\n')

with open('README.md') as f:
    long_description = f.read()

setup(
    name='taskforge',
    version='0.1.0',
    url='https://github.com/chasinglogic/taskforge',
    license='AGPL-3.0',
    author='Mathew Robinson',
    author_email='chasinglogic@gmail.com',
    description='A task management library and tool that integrates'
    ' with 3rd party services',
    long_description=long_description,
    long_description_content_type='text/markdown',
    packages=find_packages(exclude=['tests']),
    include_package_data=True,
    zip_safe=False,
    platforms='any',
    install_requires=dependencies,
    extras_requires={
        'cli': [
            'toml==0.9.4'
        ]
    },
    entry_points={
        'console_scripts': [
            'task = taskforge.cli:main',
        ],
    },
    classifiers=[
        # As from https://pypi.org/classifiers/
        'Development Status :: 4 - Beta',
        # 'Development Status :: 5 - Production/Stable',
        'Environment :: Console',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: GNU Affero General Public License v3',
        'Operating System :: POSIX',
        'Operating System :: MacOS',
        'Operating System :: Microsoft :: Windows',
        'Programming Language :: Python',
        'Programming Language :: Python :: 3',
        'Topic :: Software Development :: Libraries :: Python Modules',
    ]
)
