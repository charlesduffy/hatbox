package main
//import "fmt"
//import "log"

// Datum name. Usually something like "table"."column"
// but can be aliased at various stages of execution. 

// This is pretty speculative at this stage
// TODO refactor datum name structure
// TODO investigate integrating with Datum type iteself


// Tuple struct used to represent a database tuple
// internally just about everywhere, including as the 
// basic item transmitted 
type Tuple struct {
	data	[]Datum	
	dtypes	[]datumtype
//	dnames  []datumname
	width	int
}


// Range Table
// A list of structs identifying relations in a query
// Just an array of TRange objects (range is a keyword)
// TODO consider some get and set func for RangeTable


type TRange struct {

	catId		int		//catalogue ID of the relation
	planId		int		//plan-specific ID, generated during planning
	physName	string		//name of the relation "on disk" 
	relName		string		//name of the relation from the catalogue
	schemaName	string		//name of the relation from the catalogue
	aliasName	string		//alias of the relation, provided in query

}

type RangeTable []TRange

// Selection Table
// A list of filter expressions, one per relation in the range table
// "physical" tables can be referred to more than once and have 
// different selection predicates applied so we have an entry for 
// each range table item. 

type TSelection struct {

	planId		int		//the plan ID of the relation this filter applies to
	selExp		Expr		//an expression extracted from the WHERE clause Expr
					//which applies to just the relation referred to in 
					//planId
}

type SelectionTable []TSelection

// Projection table
// list of projections for the whole query,
// that is the "select list"

type TProjection struct {
	proj	[]Expr
	ord	int
}

type ProjectionTable []TProjection

