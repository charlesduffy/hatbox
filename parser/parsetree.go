package parser

//import "log"

var P pnode
// Parse tree interface
type ParseTree interface {
	Parse(string)
	GetRangeTable() RangeTable
	appendNode(pnode )
}


// Parse an SQL statement

func (pt *pnode) Parse(s string) {
	//log.Printf("parsetree: string is %s", s)
	P = pnode{}
	exprParse(&exprLex{line: string(s)})
}
