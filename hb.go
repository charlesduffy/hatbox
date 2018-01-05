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

var p parser.ParseTree

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

                //exprParse(&exprLex{line: string(line)})
		p.Parse(string(line))

		log.Printf("Parse tree is: %+v", p)
		log.Printf("ok, calling walkParseTree")
	//	p.tree[0].walkParseTree()
	//	log.Printf("==========================")
		log.Printf("visualise with Spew:")
		spew.Dump(p)
		log.Printf("ok, calling get_rangetable")
	//	p.tree[0].getRangeTable()
	//	dg := p.tree[0].mkdot()
		log.Printf("====================================================================")
	//	dg.drawdot()
        }
}
