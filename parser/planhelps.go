package parser

import "log"
// Range Table
// A list of structs identifying relations in a query
// Just an array of TRange objects (range is a keyword)
// TODO consider some get and set func for RangeTable

type TRange struct {
	catId      int    //catalogue ID of the relation
	planId     int    //plan-specific ID, generated during planning
	physName   string //name of the relation "on disk"
	relName    string //name of the relation from the catalogue
	schemaName string //name of the schema relation is in
	aliasName  string //alias of the relation, provided in query

}

type RangeTable []TRange

// Selection Table
// A list of filter expressions, one per relation in the range table
// "physical" tables can be referred to more than once and have
// different selection predicates applied so we have an entry for
// each range table item.

type TSelection struct {
	planId int   //the plan ID of the relation this filter applies to
	selExp []pnode //an expression extracted from the WHERE clause Expr
	//which applies to just the relation referred to in
	//planId
}

type SelectionTable []TSelection

// Projection table
// list of projections for the whole query,
// that is the "select list"

type TProjection struct {
	planId int
	proj []pnode
	ord  int
}

type ProjectionTable []TProjection


//Gets the "range table" - a list of tables that we need to scan from.
//Produces a table with relation catalogue name , schema name ,
//relation name , alias , projection list
func (t pnode) GetRangeTable() RangeTable {

	var rt RangeTable
	var planid int

	planid = 0

	var f = func(l pnode, _ int) (bool, pnode) {

		if l.tag == table_ref {
			rt = append(rt, TRange{
				catId:      0,
				planId:     planid,
				physName:   l.getIdent(),
				relName:    l.getIdent(),
				schemaName: "public",
				aliasName:  l.getIdentAlias()})
			planid += 1
		}
		return false, pnode{}
	}

	t.walkPnode(f, 0)

	log.Printf("range table is: %+v", rt)
	return rt
}

func (t pnode) GetSelection() SelectionTable {

	var st SelectionTable
	var planid int

	planid = 0

	var f = func(l pnode, _ int) (bool, pnode) {

		if l.tag == where_clause {
			st = append(st, TSelection{
				planId:     planid,
				selExp:  l.tree})
			planid += 1
		}
		return false, pnode{}
	}

	t.walkPnode(f, 0)

	log.Printf("selection table is: %+v", st)
	return st

}

func (t pnode) GetProjection() ProjectionTable {

	var pt ProjectionTable
	var planid int
	var order int

	planid = 0
	order = 0

	var f = func(l pnode, _ int) (bool, pnode) {

		if l.tag == select_list {
			pt = append(pt, TProjection{
				planId:     planid,
				ord:	    order,
				proj:  l.tree})
			planid += 1
			order += 1
		}
		return false, pnode{}
	}

	t.walkPnode(f, 0)

	log.Printf("Projection table is: %+v", pt)
	return pt
}
