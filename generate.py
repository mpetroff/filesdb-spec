#/usr/bin/env python

'''
Creates a SQLite file database.
Copyright (c) 2015 Matthew Petroff

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
'''

from __future__ import print_function

import argparse
import os
import sqlite3
import sys

if sys.version > '3':
    buffer = memoryview

# Parse terminal arguments
parser = argparse.ArgumentParser(description='Generate a SQLite file database.',
                                 formatter_class=argparse.ArgumentDefaultsHelpFormatter)
parser.add_argument('inputDir', metavar='INPUT', help='source directory')
parser.add_argument('-o', '--output', dest='output', default=argparse.SUPPRESS,
                    help='output file (default: ./INPUT.filesdb)')
args = parser.parse_args()

# Determine output file name
if 'output' in args:
    db_name = args.output
else:
    split = os.path.split(args.inputDir)
    if split[-1] == '':
        db_name = os.path.split(args.inputDir)[-2] + '.filesdb'
    else:
        db_name = os.path.split(args.inputDir)[-1] + '.filesdb'

print(db_name)
# Check if database exists
if os.path.isfile(db_name):
    print('Database already exists! Exiting...')
    sys.exit(1)

# Open database
conn = sqlite3.connect(db_name)
cursor = conn.cursor()

# Create table
cursor.execute('create table files (filename text, data blob)')

# Insert files
for directory in os.walk(args.inputDir):
    if len(directory[2]) > 0:
        dir_name = directory[0][len(args.inputDir) + len(os.sep) - 1:]
        for file_name in directory[2]:
            data_name = os.path.join(dir_name, file_name)
            file_name = os.path.join(directory[0], file_name)
            with open(file_name, 'rb') as in_file:
                cursor.execute('insert into files values (?, ?)', (data_name, buffer(in_file.read())))

# Save and close database
conn.commit()
conn.close()
