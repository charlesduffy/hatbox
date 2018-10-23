package main

//go:generate rm -f parser/y.go
//go:generate goyacc -o parser/y.go -p expr parser/gram.y
//go:generate build/mkstructures.sh

import (
		"github.com/hatbox/parser"
		"testing"
		"log"
)

/*
func TestSQLSelect1(t *testing.T){
		parser.P.Parse(string("select foo;"))
		dd := parser.P.Mkdot()
		dd.Drawdot()
}

func TestSQLSelect2(t *testing.T){
		parser.P.Parse(string("select foo from bar;"))
		dd := parser.P.Mkdot()
		dd.Drawdot()
}

func TestSQLSelect3(t *testing.T){
		parser.P.Parse(string("select foo from bar where a < 1;"))
		dd := parser.P.Mkdot()
		dd.Drawdot()
}

func TestSQLSelect4(t *testing.T){
		parser.P.Parse(string("select foo from bar where a < 1 and b > 2;"))
		dd := parser.P.Mkdot()
		dd.Drawdot()
}

func TestSQLSelect5(t *testing.T){
		parser.P.Parse(string("select foo from bar where a < 1 and b > 2 or b = 5;"))
		dd := parser.P.Mkdot()
		dd.Drawdot()
}
*/
func TestPlannerGetRangeTable(t *testing.T) {
	      var rt parser.RangeTable
	      //parser.P.Parse(string("select foo from bar , baz , bing where A < 1 and B > 2 or C = B;"))
	      //parser.P.Parse(string("select A from tab1 where A < 1 and B = 2;"))
	      rt = parser.P.GetRangeTable()
	      log.Printf("\nrange table: %+v",rt)
}


func TestPlannerGetSelection(t *testing.T) {
	      var st parser.SelectionTable
	      //parser.P.Parse(string("select foo from bar where foo < 1;"))
	      //parser.P.Parse(string("select foo from bar , baz , bing where A < 1 and B > 2 or C = B;"))
	      //parser.P.Parse(string("select foo from bar , baz , bing where A < 1 and B = 2;"))
	      parser.P.Parse(string("select A from tab1 where A < 1 and B = 2;"))
	      st = parser.P.GetSelection()
	     log.Printf("\nselection table: %+v",st)
}

func TestPlannerGetProjection(t *testing.T) {
		var pt parser.ProjectionTable
	      //parser.P.Parse(string("select foo from bar where foo < 1;"))
	      //parser.P.Parse(string("select foo from bar , baz , bing where A < 1 and B > 2 or C = B;"))
	      //parser.P.Parse(string("select foo from bar , baz , bing where A < 1 and B = 2;"))
	      parser.P.Parse(string("select A from tab1 where A < 1 and B = 2;"))
	      pt = parser.P.GetProjection()
	      log.Printf("\nprojection table: %+v",pt)
}

func TestPlannerVisualise(t *testing.T) {
		//var pt parser.ProjectionTable
	      //parser.P.Parse(string("select foo from bar where foo < 1;"))
	      //parser.P.Parse(string("select foo from bar , baz , bing where A < 1 and B > 2 or C = B;"))
	      parser.P.Parse(string("select a from foo, bar, baz where A < 1 and A = C and C = D;"))
	      //parser.P.Parse(string("select foo from bar as b1 , baz as b2 , bing as b3 where f < 1;"))
	      //parser.P.Parse(string("select A from tab1 where A < 1 and B = 2;"))
	      dt := parser.P.Mkdot()
	      dt.Drawdot()
}
