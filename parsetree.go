package main
//import "fmt"
import "log"

//Parse tree node types enumeration
//TODO: figure out a way to encapsualte this in 
//a more idiomatic fashion

type nodetype int
type datumtype int

// Node types are autogenerated from gram.y in Makefile

// Pnode 'attribute' map keys
const (

	att_alias = iota	// Node alias (IDENTIFIER)
	att_order		// Order in a list (int)
	att_subq		// Is a subquery (bool)
)

// Parse tree "wrapper" struct
type ptree struct {
	tree	[]Pnode
	query	string
}

// Parse tree node
type Pnode struct {
	tag	int
	tree	[]Pnode
	attr	map[int]interface{}
	dat	Datum
}

// Datum interface
// Need to add some methods here 
// but leave it blank for the time being

type Datum struct {
	value	datumval
	dtype	datumtype
}

func (p *Pnode) appendNode(n Pnode) {
	p.tree = append(p.tree, n)
}

func makeNode(n int ) Pnode {
	return Pnode{ tag: n }
}

func (p *Pnode) addAttr(att int , val interface{}) {

	if (p.attr == nil) {
		p.attr = make(map[int]interface{})
	}

	p.attr[att] = val
}

func (p *Pnode) addDatum(d datumval, t datumtype) {
	p.dat = Datum{
			value: d,
			dtype: t}
}

func (p *Pnode) addDatum0(d Datum) {
	p.dat = d
}

/* --flag for removal
func makeIdentifier(i string) Datum {
	return Datum{
				value: i,
				dtype: IDENTIFIER}
}
*/

func makeScalarExpr(d Datum, l Pnode, r Pnode) Pnode {

	n := Pnode{tag:scalar_expr}
	n.dat = d
	n.appendNode(l)
	n.appendNode(r)
	return n
}

type PUserFunc func(Pnode,int) (bool, Pnode)

func (t Pnode) walkPnode(fn PUserFunc, depth int) (bool, Pnode) {

	//traverse 'tree' slice left-depth first

	var p Pnode

	ret , q := fn(t,depth)
	if ( ret == true ) {
		 return true, q
	}

	if (t.tree != nil) {
		for _ , p = range t.tree {

			ret , q = p.walkPnode(fn, depth+1)

			if (ret == true) {
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
func (t Pnode) getIdent() string {

	if t.dat.dtype == IDENTIFIER {
		return t.dat.value.(string)
	} else {
		return ""
	}
}

// For a Pnode that has an alias, return it as a string

func (t Pnode) getIdentAlias() string {

	if alias,ok := t.attr[att_alias].(string); ok {
		return alias
	} else {
		return ""
	}
}

//Gets the list of tables that we need to scan from. 
//Produces a table with relation catalogue name , schema name , 
//relation name , alias , projection list

func (t Pnode) getRangeTable() RangeTable {

	var rt RangeTable
	var planid int

	planid = 0

	var f = func(l Pnode, _ int)(bool,Pnode) {

		if (l.tag == table_ref) {
			rt = append(rt, TRange{
						catId: 0,
						planId: planid,
						physName: l.getIdent(),
						relName:  l.getIdent(),
						schemaName: "public",
						aliasName: l.getIdentAlias()})
			planid += 1
		}
		return false,Pnode{}
	}

	t.walkPnode(f,0)

	log.Printf("range table is: %+v", rt)
return rt
}

func (t Pnode) getSelection() SelectionTable {

return nil
}

func (t Pnode) getProjection() ProjectionTable {

return nil
}

// Walk parse tree for debugging purposes
func (t Pnode) walkParseTree() {

	var f = func(l Pnode, _ int)(bool,Pnode){
		log.Print("calling func\n")
		//log.Printf("fn: current pnode: %s %d %+v ",typName(l.tag), l.tag, l.val)
		return false,Pnode{}
	}
	_, a := t.walkPnode(f,0)
	log.Printf("%+v\n", a)
}
