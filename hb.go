package main

//go:generate rm -f parser/y.go
//go:generate goyacc -o parser/y.go -p expr parser/gram.y
//go:generate build/mkstructures.sh

import (	"log"
		"bufio"
		"os"
		"io"
		"github.com/davecgh/go-spew/spew"
		"github.com/hatbox/parser"
		//"github.com/hatbox/planner"
)


func main(){
        in := bufio.NewReader(os.Stdin)
	var ptree parser.ParseTree
        for {
                if _, err := os.Stdout.WriteString("> "); err != nil {
                        log.Fatalf("WriteString: %s", err)
                }
                line, err := in.ReadBytes('\n')
                if err == io.EOF {
                        return
                }
                if err != nil {
                        log.Fatalf("ReadBytes: %s", err)
                }

                //exprParse(&exprLex{line: string(line)})
		ptree.Parse(string(line))

		log.Printf("Parse tree is: %+v", ptree)
		log.Printf("ok, calling walkParseTree")
	//	p.tree[0].walkParseTree()
	//	log.Printf("==========================")
		log.Printf("visualise with Spew:")
		spew.Dump(ptree)
		log.Printf("ok, dumping P")
		spew.Dump(parser.P)
//		log.Printf("ok, calling get_rangetable")
	//	p.tree[0].getRangeTable()
	//	dg := p.tree[0].mkdot()
		log.Printf("BEEEEF====================================================================")
//	spew.Dump(parser.P)
		//P.planner.PlanQuery()
	//	dg.drawdot()
        }
}
