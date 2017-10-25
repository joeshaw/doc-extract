package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestSorter(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "testdata", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("Expected len(pkgs) == 1, got %d", len(pkgs))
	}

	expected := []string{"testdata.go", "doc.go", "a.go", "z.go"}
	sorted := sortedFiles(pkgs["testdata"])

	if len(sorted) != len(expected) {
		t.Fatalf("Expected %d, got %d", len(expected), len(sorted))
	}

	for i := 0; i < len(sorted); i++ {
		if sorted[i].name != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], sorted[i].name)
		}
	}
}

func TestExtractComment(t *testing.T) {
	tests := []struct {
		comment  string
		expected bool
	}{
		{"//+extract", true},
		{"// +extract", true},
		{"//           +extract\n", true},
		{"/*+extract*/", true},
		{"/* +extract */", true},
		{"/*           +extract */", true},
		{"/*\n+extract */", true},
		{"/* +extract\nfoo */", true},
		{"/* foo\n+extract */", false},
		{"// extract", false},
		{"// foo +extract", false},
	}

	for i, test := range tests {
		var c ast.Comment
		c.Text = test.comment

		cgrp := &ast.CommentGroup{List: []*ast.Comment{&c}}
		_, ok := extractComment(cgrp)
		if ok != test.expected {
			t.Fatalf("Iteration %d: Expected %t, got %t: %q", i, test.expected, ok, test.comment)
		}
	}
}

func TestParsePackage(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, "testdata", nil, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	if len(pkgs) != 1 {
		t.Fatalf("Expected len(pkgs) == 1, got %d", len(pkgs))
	}

	comments := extractPackageComments(pkgs["testdata"])

	expected := []string{
		"A comment inside testdata.go\n",
		"More text, in doc.go\n",
		"Here's a comment in a.go\n",
		"An interesting\nmulti-line\ncomment inside\nz.go\n",
	}

	if len(comments) != len(expected) {
		t.Fatalf("Expected %#v, got %#v", expected, comments)
	}

	for i := 0; i < len(comments); i++ {
		if comments[i] != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], comments[i])
		}
	}
}
