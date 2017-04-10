// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is an example of a goyacc program.
// To build it:
// go tool yacc -p "expr" expr.y (produces y.go)
// go build -o expr y.go
// expr
// > <type an expression>

package main
//import "fmt"

var SQLkeys = map[string]int{
        "select": SELECT,
}

const eof = 0

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
        line string                     //Buffer containing SQL input text to be lexed
	tokstart int			//start of current token
	tokend int			//end of current token 
}

func (x *exprLex) peek() rune {
        if len(x.line) < 1 {
                return nil
        }
        return x.line[:x.tokstart]
}

// Ignore this space
func (x *exprLex) ignore() {

}

// We're inside double quotes. 
// Lex an identifier
func (x *exprLex) lexdquote() {

/*
	* consume input
	* permit any character including spaces
	* peek when quote encountered, to see if
	* escaped quote or end of token
*/
	x.typ = IDENT

	for n := x.next()  {
	  if n == "\"" {
		if x.peek() != "\"" {
			return
		}
	  }
	}
}

// We're inside single quotes. 
// Lex a string literal
func (x *exprLex) lexsquote() {

	x.typ = STRING

	for n := x.next()  {
	  if n == "'" {
	    if x.peek() != "'" {
	      return
	    }
	  }
	}
}

// We've hit what looks like the 
// start of a number. 
// Try to lex a numeric literal 


// The parser calls this method to get each new token.
func (x *exprLex) Lex(yylval *exprSymType) int {

	//This is called either at the very beginning of the 
	//string to be parsed or at the start of a new 
	//token

        start:

	switch n := x.next() {
	  case n == eof:
		//do something at end-of-input
	  case isSpace(n):
		x.ignore()
	  case r == "\"":
		x.lexdquote()
	  case r == "'":
		x.lexsquote()
	  case r >= '0' && r <= '9':
		x.lexnumber()
	  case isAlphaNumeric(r):
		//Here we could match an identifier or a 
		//keyword
		x.lextext()
	  case r == ".":
		x.lexpoint()
	  case r == ";":
		x.lexterm()
	  default:
		x.lexoper()
	}

        goto start
}


// Return the next rune for the lexer.
func (x *exprLex) next() rune {
        if x.peek != eof {
                r := x.peek
                x.peek = eof
                return r
        }
        if len(x.line) == 0 {
                return eof
        }
        c, size := utf8.DecodeRune(x.line)
        x.line = x.line[size:]
        if c == utf8.RuneError && size == 1 {
                log.Print("invalid utf8")
                return x.next()
        }
        return c
}

// The parser calls this method on a parse error.
func (x *exprLex) Error(s string) {
        log.Printf("parse error: %s", s)
}

func main() {
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

                exprParse(&exprLex{line: line})
        }
}
