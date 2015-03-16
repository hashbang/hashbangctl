#!/usr/bin/env python
from distutils.core import setup

setup(
    name='hashbang-config',
    version='0.2',
    packages=['hashbang-config'],
    scripts=[ 
        'bin/hashbang',
    ],
    data_files=[('/etc',['hashbang-config.conf'])],
    author='Hashbang Team',
    author_email='team@hashbang.sh',
    license='GPL 3.0',
    description='',
    long_description=open('README.md').read(),
    install_requires=[
        'provisor',
    ],
    dependency_links = [
        'http://github.com/hashbang/provisor/tarball/master#egg=provisor',
    ]
)
