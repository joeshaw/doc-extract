# API Blueprint example

This directory contains a Go source file ([example.go](example.go))
which has [API Blueprint](https://apiblueprint.org) definitions for a
RESTish HTTP API in its comments.

Using doc-extract, we can extract these comments into an API Blueprint
file:

    doc-extract . example.apib

And then we'll used a Node.js-based tool called
[Aglio](https://github.com/danielgtaylor/aglio) running in Docker to
render that to a nice HTML file.

    docker run --rm -t -v $(pwd):/docs humangeo/aglio -i /docs/example.apib -o /docs/example.html

[Viola!](https://htmlviewer.github.io/?https://raw.githubusercontent.com/joeshaw/doc-extract/master/example/example.html)
