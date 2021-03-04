// The doc-extract utility extracts text from Go source comments
// tagged with the special token "+extract" in the first line.
//
// Usage: doc-extract <source dir> <output file>
//
// "source dir" is a directory containing Go source files ending
// in .go, and "output file" is the file to write.
//

// Source files are processed in lexicographic order, except that a file
// named doc.go is always processed first.  Comments within a file are
// extracted in the order they appear.  This predictable order allows
// you to add, for instance, a header to the output file by adding it to
// doc.go.
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path"
	"sort"
	"strings"
)

type file struct {
	file *ast.File
	name string
}

func sortedFiles(pkg *ast.Package) []file {
	files := make([]file, 0, len(pkg.Files))

	for name, f := range pkg.Files {
		files = append(files, file{file: f, name: path.Base(name)})
	}

	// Sort files passed in according to this heuristic:
	// 1. file with name == package name
	// 2. file with name "doc.go"
	// 3. lexicographic order
	sort.Slice(files, func(i, j int) bool {
		ni, nj := files[i].name, files[j].name

		if ni == "doc.go" {
			return true
		}

		return ni < nj
	})

	return files
}

func extractComment(cgrp *ast.CommentGroup) (string, bool) {
	s := cgrp.Text()
	parts := strings.SplitN(s, "\n", 2)
	if strings.TrimSpace(parts[0]) == "+extract" {
		return parts[1], true
	}
	return "", false
}

func extractPackageComments(pkg *ast.Package) []string {
	files := sortedFiles(pkg)

	var comments []string
	for _, f := range files {
		for _, c := range f.file.Comments {
			s, ok := extractComment(c)
			if ok {
				comments = append(comments, s)
			}
		}
	}

	return comments
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: %s <source dir> <output file>\n", os.Args[0])
		return
	}

	srcDir := os.Args[1]
	outFile := os.Args[2]

	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, srcDir, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing files in %s: %s\n", srcDir, err)
		os.Exit(1)
	}

	var out io.Writer
	if outFile == "-" {
		out = os.Stdout
	} else {
		f, err := os.Create(outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file %s: %s\n", outFile, err)
			os.Exit(1)
		}
		defer f.Close()

		out = f
	}

	for _, pkg := range pkgs {
		comments := extractPackageComments(pkg)
		for _, c := range comments {
			fmt.Fprintln(out, c)
		}
	}
}
