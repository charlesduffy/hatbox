package main
//import "fmt"
import "log"

//Parse tree node types enumeration
//TODO: figure out a way to encapsualte this in 
//a more idiomatic fashion

type nodetype int
type datumtype int

// Node types are autogenerated from gram.y in Makefile

// Parse tree "wrapper" struct
type children []Pnode		//TODO consider if this is really necessary

type ptree struct {
	tree	children
	query	string
}

// Parse tree node
type Pnode struct {
	tag	int
	tree	[]Pnode
	val	*Expr
}

// Expression tree
type Expr struct {
	data	Datum
	left	*Expr
	right	*Expr
}

// Datum interface
// Need to add some methods here 
// but leave it blank for the time being

type Datum struct {
	value	datumval
	dtype	datumtype
}

func (p *Pnode) append_node(n Pnode) {
	p.tree = append(p.tree, n)
}

func make_identifier(i string) Datum {
	return Datum{
				value: i,
				dtype: IDENTIFIER}
}


func make_scalar_expr(d Datum, l *Expr , r *Expr) *Expr {
	return &Expr{
			data: d,
			left: l,
			right: r}
}


func Walk_ptree(t ptree) {
	for _ , p := range t.tree {
	  p.walk_pnode()
	}
}

type userfunc func(Pnode)

func (t Pnode) walk_pnode() {

	//traverse 'tree' slice left-depth first
	var p Pnode

	if (t.tree != nil) {
		for _ , p = range t.tree {

			p.walk_pnode()
		}

	}
	//print the current node
	log.Printf("%s %d %+v ",typName(t.tag), t.tag, t.val)
}

func typName(t int) string {
	return NodeYNames[t]
}

func (t Pnode) getRangeTable() RangeTable {

// 1. make a new empty RangeTable
// 2. traverse the parse tree until we get to the from_clause
// 3. iterate over the table_ref objects in the from_clause
// 4. for each table_ref object, make a TRange and Append() it 
//    to the RangeTable
// 5. When we're finished, return the RangeTable
return nil
}

func (t Pnode) getSelection() SelectionTable {

return nil
}

func (t Pnode) getProjection() ProjectionTable {

return nil
}




