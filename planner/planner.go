package planner

import "github.com/davecgh/go-spew/spew"
import "github.com/hatbox/parser"


// we need a Plan interface object

type Planable interface {
	//ok things we want in a Plan object are:
	//a way to get the Range table
	//
	GetRangeTable() parser.RangeTable
}

func PlanQuery (p Planable) {
// .1 get the range table

	rt := p.GetRangeTable()
// .2 get the select list table
	spew.Dump(rt)
// .3 get the predicate list

// .4 make scan nodes
}

