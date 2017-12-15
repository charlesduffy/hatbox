package main

//go:generate goyacc -o parser/y.go -p expr parser/gram.y
//go:generate build/mkstructures.sh

import (	"log"
		"bufio"
		"os"
		"io"
		"github.com/davecgh/go-spew/spew"
		"github.com/hatbox/parser"
)

var Parsetree ptree
func main(){
        in := bufio.NewReader(os.Stdin)
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

                exprParse(&exprLex{line: string(line)})

		log.Printf("Parse tree is: %+v", Parsetree)
		log.Printf("ok, calling walkParseTree")
		Parsetree.tree[0].walkParseTree()
		log.Printf("==========================")
		log.Printf("visualise with Spew:")
		spew.Dump(Parsetree)
		log.Printf("ok, calling get_rangetable")
		Parsetree.tree[0].getRangeTable()
		dg := Parsetree.tree[0].mkdot()
		log.Printf("====================================================================")
		dg.drawdot()
        }
}
