package parser

import "log"

var P ParseTree
// Parse tree "wrapper" struct
type ParseTree struct {
	tree  []pnode
	query string
}

// Parse an SQL statement

func (pt *ParseTree) Parse(s string) {
	log.Printf("parsetree: string is %s", s)
	exprParse(&exprLex{line: string(s)})
	pt.query = "hello i am spiderman"
}
