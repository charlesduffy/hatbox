package main
import "log"
import "fmt"

//visualise a ptree using DOT

type dotNode struct {
	nodeId int
	parentId int
	label  string
}

type dotLink struct {
	linkId	int
	from	int
	to	int
}

type dotGraph struct {
	dn	[]dotNode
	dl	[]dotLink
}

func (d dotGraph) drawdot() {

	fmt.Printf("graph \"parsetree\" { node [ fontsize=12 ]; graph [ fontsize=10 ]; label = \"query text goes here\" subgraph parsetree_1 { color=\"blue\" ")


	if (d.dn != nil) {
		for _,p := range d.dn {
			fmt.Printf("ptree_%0.0d [ label = \"%v\" ];", p.nodeId, p.label);
		}
	}

	if (d.dl != nil) {
		for _,q := range d.dl {
			fmt.Printf("ptree_%0.0d -- ptree_%0.0d [ id = %d ]", q.from , q.to , q.linkId);
		}
	}

	fmt.Printf("} }")
}


//Trverse an Expr structure accumulating a
//dotGraph representing a scalar expression, 
//and return it
//for the time being we accept args for the 
//initial nodeId and LinkID values

//this could be alleviated with some dot language
//trickery or a wrapper function

/* -- temp hide
func (e Expr) mkdot(n int, linkId int) (dotGraph) {

	var dn  []dotNode
	var dt  []int
	var dl  []dotLink
	var nodeid int
	var linkid int
	var depth int
	var pid int

	nodeid = n
	linkid = linkId

	var f = func(e Expr, d int)(bool, Expr) {

			log.Printf("in expr processing")

		switch d {
			case 0:
				dt = append(dt, nodeid)
				log.Printf("Expr depth initial: %d %+v", d, e.data )
				depth = 0
				pid = n
			case depth+1:
				pid = n + dt[len(dt)-1]
				if e.left != nil && e.right != nil {
					dt = append(dt, n)
				}
				depth = d
			case depth-1:
				dt = dt[:len(dt)-1]
				pid = dt[len(dt)-1]
				if e.left != nil && e.right != nil {
					dt = append(dt, n)
				}
				depth = d
		}

		dn = append(dn, dotNode{
					nodeId: nodeid,
					parentId: pid,
					label: fmt.Sprintf("EX %+v", e.data)})

		dl = append(dl, dotLink{
					linkId: linkid,
					from: pid,
					to: nodeid})

		nodeid += 1
		linkid += 1

		return false, Expr{}
	}

	e.walkExpr(f,0)


return dotGraph{
		dn,
		dl}

	return dotGraph{ dn, dl}


}
*/
func (t Pnode) mkdot() (dotGraph) {

	var dn  []dotNode
	var dt  []int
	var dl  []dotLink
	var nodeid int
	var linkid int
	var depth int
	var parentid int

	nodeid = 1
	linkid = 1

	var f = func(l Pnode, d int)(bool,Pnode) {


		if (d == 0) {
			dt = append(dt, nodeid)

		log.Printf("depth 0 %+v ",typName(l.tag))
		depth = 0
		parentid = 0
		} else if (d == (depth+1)) {
			parentid = dt[len(dt)-1]
			if (l.tree != nil) {
				dt = append(dt, nodeid)
			}
			depth = d
		log.Printf("d %d depth %d %+v ",d,depth, typName(l.tag))

		} else if (d == (depth-1)) {
			dt = dt[:len(dt)-1]
			parentid = dt[len(dt)-1]
			if (l.tree != nil) {
				dt = append(dt, nodeid)
			}
			depth = d
		log.Printf("d %d %+v ",d,typName(l.tag))
		}


		dn = append(dn, dotNode{
					nodeId: nodeid,
					parentId: parentid,
					label: fmt.Sprintf("%s", typName(l.tag))})

		dl = append(dl, dotLink{
					linkId: linkid,
					from: parentid,
					to: nodeid})

		nodeid += 1
		linkid += 1

		//check if we are an Expr node, if so, draw the expression graph
		return false,Pnode{}
	}

	t.walkPnode(f,0)

	log.Printf("node list is: %+v  link list is: %+v ", dn, dl)
return dotGraph{
		dn,
		dl}
}
