package parser

import "log"

var P ParseTree
// Parse tree "wrapper" struct
type ParseTree interface {
	Parse(string)
	GetRangeTable() RangeTable
	appendNode(pnode )
}


// Parse an SQL statement

func (pt *pnode) Parse(s string) {
	log.Printf("parsetree: string is %s", s)
	exprParse(&exprLex{line: string(s)})
}
