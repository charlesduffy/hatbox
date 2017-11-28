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
	from	string
	to	string
}

func (t Pnode) mkdot() ([]dotNode, []dotLink) {

	var dn  []dotNode
	var dt  []int
	var dl  []dotLink
	var nodeid int
	var depth int
	var parentid int

	nodeid = 1

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


		nodeid += 1
		return false,Pnode{}
	}

	t.walkPnode(f,0)

	log.Printf("node list is: %+v  link list is: %+v ", dn, dl)
return dn,dl
}
