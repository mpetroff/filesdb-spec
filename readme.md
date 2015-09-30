# FilesDB Specification and Reference Implementation

A reference creation utility, written in Python, is provided in `generate.py`,
and a reference file server, written in Go, is provided in `server.go`. The
remainder of this document serves as a basic format specification.



## Specification

### Abstract

FilesDB is a specification for storing small files in
[SQLite](http://sqlite.org/) databases for immediate usage and for transfer.
FilesDB files, known as **filesets**, must implement the specification below
to ensure compatibility with devices.


## Database Specifications

Filesets are expected to be valid SQLite databases of
[version 3.0.0](http://sqlite.org/formatchng.html) or higher.
Only core SQLite features are permitted; filesets **cannot require extensions**.


### Database

Note: the schema outlined is meant to be followed as an interface.
SQLite views that produce compatible results are equally valid.
For convenience, this specification refers to tables and virtual
tables (views) as tables.


### Schema

The database is required to contain a table or view named `files`.

This table must yield exactly two columns named `filename` and
`data`. A typical create statement for the `files` table:

    CREATE TABLE files (filename text, data blob);


### Content

The `filename` column contains the path to the specified file relative to the
root of the fileset. Unix-style paths should be used, i.e. `/` is used as the
directory separator.

The `data` column contains the file's contents.


## License

The text of this specification is licensed under a
[Creative Commons Attribution 3.0 United States License](http://creativecommons.org/licenses/by/3.0/us/)
and is based on the [MBTiles 1.2 specification](https://github.com/mapbox/mbtiles-spec/blob/master/1.2/spec.md).
The reference implementation is licensed under an [MIT license](https://opensource.org/licenses/MIT).


## Author

[Matthew Petroff](https://mpetroff.net)
