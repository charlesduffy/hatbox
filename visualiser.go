package main
import "log"
import "fmt"

//visualise a ptree using DOT

type dotNode struct {
	nodeId int
	parentId int
	label  string
	tag int
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

	var attstr string

	fmt.Printf("graph \"parsetree\" {\n")
	fmt.Printf(" node [ fontsize=11 ]; graph [ fontsize=10 ];\n")
	fmt.Printf(" subgraph parsetree_1 { color=\"blue\" \n")


	if (d.dn != nil) {
		for _,p := range d.dn {
			if (p.tag == scalar_expr) {
				attstr = "color=deepskyblue shape=egg style=filled"
			} else {
				attstr = "color=thistle2 shape=box style=filled"
			}
			fmt.Printf("ptree_%0.0d [ label = \"%v\" %s ];\n", p.nodeId, p.label, attstr);
		}
	}

	if (d.dl != nil) {
		for _,q := range d.dl {
			fmt.Printf("ptree_%0.0d -- ptree_%0.0d [ id = %d ]\n", q.from , q.to , q.linkId);
		}
	}

	fmt.Printf("\t}\n}")
}


//Trverse an Expr structure accumulating a
//dotGraph representing a scalar expression, 
//and return it
//for the time being we accept args for the 
//initial nodeId and LinkID values

//this could be alleviated with some dot language
//trickery or a wrapper function


func (e Pnode) mkdot() (dotGraph) {

	var dn  []dotNode
	var dt  []int
	var dl  []dotLink
	var nodeid int = 1
	var linkid int = 1
	var depth int
	var pid int

	var f = func(e Pnode, d int)(bool, Pnode) {

		switch {
			case d == 0:
				dt = append(dt, nodeid)
				depth = 0
				pid = 0
			case d == depth+1:
				pid = dt[len(dt)-1]

				if e.tree != nil {
					dt = append(dt, nodeid)
				}

				depth = d
			case d < depth:
				dd := depth - d
				dt = dt[:len(dt)-dd]
				pid = dt[len(dt)-1]

				if e.tree != nil {
					dt = append(dt, nodeid)
				}

				depth = d
		}

		dn = append(dn, dotNode{
					nodeId: nodeid,
					parentId: pid,
					label: fmt.Sprintf("tag: %d %s %+v", e.tag, typName(e.tag), e.dat),
					tag: e.tag})

		dl = append(dl, dotLink{
					linkId: linkid,
					from: pid,
					to: nodeid})

		nodeid += 1
		linkid += 1

		return false, Pnode{}
	}

	e.walkPnode(f,0)

	return dotGraph{
			dn,
			dl}

}
