# doc-extract

doc-extract is a tool for extracting specially tagged comments in Go
source code.  Tagged comments start with a blank line containing
`+extract`.  Both grouped line comments (`//`) and block comments (`/*
*/`) are supported.

## Installation

Go 1.16 and later:

    go install github.com/joeshaw/doc-extract@latest

Go 1.15 and earlier:

    go get -u github.com/joeshaw/doc-extract

## Usage

Simply tag comments in your Go source code with `+extract` as the first
line of your comment.  For example,

```go
package main

// +extract
// API Documentation
// =================
//
// (FIXME: add documentation)
```

Then run the `doc-extract` command, providing it with a directory of
Go source files and an output file:

    doc-extract ./example example.txt

Source files are processed in lexicographic order, _except_ that a file
named `doc.go` gets highest priority.  Comments within a file are
processed in the order they appear.

This predictable ordering allows you to add, for instance, a header to
the output file by adding it to `doc.go`.

## Example

In the `example` directory is an [example using API
Blueprint](example/README.md).

## Motivation

This tool is especially useful for extracting documentation from
comments when GoDoc isn't appropriate.  For example, documentation for
HTTP APIs.  This could be used to extract [API
Blueprint](https://apiblueprint.org/) for processing by
[Snowboard](https://github.com/bukalapak/snowboard).  Or it could be
used to pull ReStructured Text (RST) for processing by
[Sphinx](http://sphinx-doc.org/) and the
[httpdomain](https://pythonhosted.org/sphinxcontrib-httpdomain/)
extension.

This replaces my now-deprecated
[rst-extract](https://github.com/joeshaw/rst-extract) tool.

## Contributing

Issues and pull requests are welcome.  When filing a PR, please make
sure the code has been run through `gofmt` and that the tests pass.

## License

MIT
