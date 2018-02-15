package parser

import "log"
//import "github.com/davecgh/go-spew/spew"

//Parse tree node types enumeration
//TODO: figure out a way to encapsualte this in
//a more idiomatic fashion

type nodetype int
type datumtype int

// Node types are autogenerated from gram.y in Makefile

// Pnode 'attribute' map keys
const (
	att_alias = iota // Node alias (IDENTIFIER)
	att_order        // Order in a list (int)
	att_subq         // Is a subquery (bool)
)


// Parse tree node
type pnode struct {
	tag  int
	tree []pnode
	attr map[int]interface{}
	dat  Datum
}


// Datum interface
// Need to add some methods here
// but leave it blank for the time being

type Datum struct {
	value datumval
	dtype datumtype
}

func (p *pnode) appendNode(n pnode) {
	p.tree = append(p.tree, n)
}

func makeNode(n int) pnode {
	return pnode{tag: n}
}

func (p *pnode) addAttr(att int, val interface{}) {

	if p.attr == nil {
		p.attr = make(map[int]interface{})
	}

	p.attr[att] = val
}

func (p *pnode) addDatum(d datumval, t datumtype) {
	p.dat = Datum{
		value: d,
		dtype: t}
}

func (p *pnode) addDatum0(d Datum) {
	p.dat = d
}

func makeOperScalarExpr(d datumtype, l pnode, r pnode) pnode {

	n := pnode{tag: scalar_expr,
		dat: Datum{
			value: nil,
			dtype: d}}
	n.appendNode(l)
	n.appendNode(r)
	return n
}

func makeScalarExpr(d Datum, l pnode, r pnode) pnode {

	n := pnode{tag: scalar_expr}
	n.dat = d
	n.appendNode(l)
	n.appendNode(r)
	return n
}

type PUserFunc func(pnode, int) (bool, pnode)

func (t pnode) walkPnode(fn PUserFunc, depth int) (bool, pnode) {

	//traverse 'tree' slice left-depth first

	var p pnode

	ret, q := fn(t, depth)
	if ret == true {
		return true, q
	}

	if t.tree != nil {
		for _, p = range t.tree {

			ret, q = p.walkPnode(fn, depth+1)

			if ret == true {
				return true, q
			}
		}
	}

	return false, q
}

func typName(t int) string {
	return NodeYNames[t]
}

// For a Pnode that holds a single identifier, return it as string
func (t pnode) getIdent() string {

	if t.dat.dtype == IDENTIFIER {
		return t.dat.value.(string)
	} else {
		return ""
	}
}

// For a Pnode that has an alias, return it as a string

func (t pnode) getIdentAlias() string {

	if alias, ok := t.attr[att_alias].(string); ok {
		return alias
	} else {
		return ""
	}
}



//Gets the list of tables that we need to scan from.
//Produces a table with relation catalogue name , schema name ,
//relation name , alias , projection list
//func (pt ParseTree) getRangeTable() RangeTable {
//	t := pt.tree
//	return t[0].getRangeTable()
//}//


func (t pnode) getRangeTable() RangeTable {

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

func (t pnode) getSelection() SelectionTable {

	return nil
}

func (t pnode) getProjection() ProjectionTable {

	return nil
}

// Walk parse tree for debugging purposes
func (t pnode) walkParseTree() {

	var f = func(l pnode, _ int) (bool, pnode) {
		log.Print("calling func\n")
		//log.Printf("fn: current pnode: %s %d %+v ",typName(l.tag), l.tag, l.val)
		return false, pnode{}
	}
	_, a := t.walkPnode(f, 0)
	log.Printf("%+v\n", a)
}

