package main

//go:generate rm -f parser/y.go
//go:generate goyacc -o parser/y.go -p expr parser/gram.y
//go:generate build/mkstructures.sh

import (
		"github.com/hatbox/parser"
		"testing"
)


func TestSQLSelect1(t *testing.T){
		parser.P.Parse(string("select foo from bar;"))
		dd := parser.P.Mkdot()
		dd.Drawdot()
}

func TestPlannerGetRangeTable(t *testing.T) {
	      parser.P.Parse(string("select foo from bar;"))
	      parser.P.GetRangeTable()
}

func TestPlannerGetSelection(t *testing.T) {
	      parser.P.Parse(string("select foo from bar;"))
	      parser.P.GetSelection()
}

func TestPlannerGetProjection(t *testing.T) {
	      parser.P.Parse(string("select foo from bar;"))
	      parser.P.GetProjection()
}

