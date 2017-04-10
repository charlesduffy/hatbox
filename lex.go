package main

import "log"
import "bufio"
import "os"
import "io"

const eof = 0

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
        line		string                  //Buffer containing SQL input text to be lexed
	tokstart	int			//start of current token
	tokend		int			//end of current token 
	pos		int			//current position in input
	typ		int			//current token type
}

func (x *exprLex) peek() rune {
        if len(x.line) < 1 {
                return 0
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
	x.typ = IDENTIFIER

	for {
	  n := x.next()
	  if n == '"' {
	    if x.peek() != '"' {
	      return
	    }
	  }
	  //TODO insert a check for max identifier length
	}
}

// We're inside single quotes. 
// Lex a string literal
func (x *exprLex) lexsquote() {

	x.typ = STRING

	for {
	  n := x.next()
	  if n == `'` {
	    if x.peek() != `'` {
	      return
	    }
	  }
	}
}

func (x *exprLex) lexnumber() {

}
// We've hit what looks like the 
// start of a number. 
// Try to lex a numeric literal 

func (x *exprLex) lextext() {

}


func (x *exprLex) lexpoint() {

}


func (x *exprLex) lexoper() {

}

func (x *exprLex) lexterm() {

}

// Return the next rune for the lexer.
func (x *exprLex) next() rune {

	r, w := utf8.DecodeRuneInString(x.line[x.pos:])
	return r
}

// The parser calls this method to get each new token.
func (x *exprLex) Lex(yylval *exprSymType) int {

	//This is called either at the very beginning of the 
	//string to be parsed or at the start of a new 
	//token
	n := x.next()

	switch {
		//case n == eof:
		//do something at end-of-input
	  case n == " ":
		x.ignore()
	  case n == "\"":
		x.lexdquote()
	  case n == "'":
		x.lexsquote()
	  case n >= "0" && n <= "9":
		x.lexnumber()
	  case isAlphaNumeric(n):
		//Here we could match an identifier or a 
		//keyword
		x.lextext()
	  case n == ".":
		x.lexpoint()
	  case n == ";":
		x.lexterm()
	  default:
		x.lexoper()
	}

	return 0
}

func isAlphaNumeric(r string) bool {
   return r == "_" || unicode.IsLetter(r) || unicode.IsDigit(r)
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
