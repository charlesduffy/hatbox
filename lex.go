package main

import (	"log"
		"bufio"
		"os"
		"io"
		"unicode"
		"unicode/utf8"
		"strings"
)

const eof = 0

// Global parse tree. There's got to be 
// a better way. 
var Parsetree ptree

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
        line		string                  //Buffer containing SQL input text to be lexed
	tok		int			//start of current token
	pos		int			//current character position in input
	typ		int			//current token type
	width		int			//width of rune
}

// This interface is for the token value 
// inside yylval. 

type datum interface {

}

// Eventually we will define methods on the Datum
// interface to transform string token values
// to internal types.

func (x *exprLex) peek() rune {
        if len(x.line) < 1 {
                return 0
        }
	r, _ := utf8.DecodeRuneInString(x.line[x.pos:])
	return r
}

// Ignore this character
func (x *exprLex) consume() {
	x.pos += x.width
	x.shift()
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
	  //TODO consider inserting a check for max identifier length
	}
}

// We're inside single quotes. 
// Lex a string literal
func (x *exprLex) lexsquote() {

	x.typ = STRING

	for {
	  n := x.next()
	  if n == '\'' {
	    if x.peek() != '\'' {
	      return
	    }
	  }
	}
}

// We've hit what looks like the 
// start of a number. 
// Try to lex a numeric literal 
func (x *exprLex) lexnumber() {
//log.Printf("Entering lexnumber()")
	x.typ = NUMERIC

	for {
	  n := x.next()
	  if unicode.IsDigit(n) != true && n != '.' {
//log.Printf("Not a digit or dot")
	  //here we have encountered a rune that is not
	  //accepted within a numeric token. 

//log.Printf("Peek character is: >>%c<<", n)

	    if  isOperator(n) || n == ';'|| isWhitespace(n) {
	     //p:=x.emit() //just for debug
//log.Printf("Token is: >>%s<<\n", p)
	     return
//log.Printf("after return???")
	    } else {
	       //raise error. 
	    x.Error(x.line)
	    return
	    }
	  }
//log.Printf("Another numeric rune...")
	}

}

// This can either be an unquoted identifier,
// or a keyword
func (x *exprLex) lextext() {

// consume characters until
// we encounter one not in the accepted
// set for a keyword or an identifier.
// The accepted set is:
// a-zA-Z0-9_
//log.Printf("Entering lextext")

	for {
	  n := x.next()
	  if ! isAlphaNumeric(n) {
	  //here we have encountered the end of the token. 
	  //determine if it is a keyword.
	    o := x.emit()

//log.Printf("Token is: >>%s<<\n", o)
	    if l, ok := SQLkeys[strings.ToLower(o)]; ok {
//log.Printf("Matched token >>%s<< to keyword value %d", o, l)
	      x.typ = l
	    } else {
//log.Printf("Token >>%s<< is an identifier")
		x.typ = IDENTIFIER
	   }
	    return
	  }
	}
}

// Lex a point character. In this context it 
// will be a delimiter between table / column 
// identifiers
func (x *exprLex) lexpoint() {
//log.Printf("Entering Lexpoint")
	x.typ = POINT
	x.consume()
	return
}

// Lex a comma. Delimiter between 
// select list items or table expression
// items.
func (x *exprLex) lexcomma() {
//log.Printf("Entering Lexcomma")
	x.typ =  COMMA
	x.consume()
	return
}

// Lex an operator. An operator is one or more
// characters from the list of:
// + - * / % ! =
func (x *exprLex) lexoper() {
}

func (x *exprLex) lexterm() {
//log.Printf("Entering Lexterm")
	x.typ = SEMICOLON
	x.consume()
	return
}

//func get tok
//get the 'current' token
func (x *exprLex) gettok() (rune, int) {
//log.Printf("Entering gettok() %s", x.line[x.pos:])
	//x.ptok()

	if (x.tok == x.pos && x.pos == len(x.line) - 1) {
		//end of the line
		return eof, 0
	}

	r, w := utf8.DecodeRuneInString(x.line[x.pos:])
	if w == 0 {
		//log.Printf("NO MORE RUNES!!!")
		return eof, 0
	} else {
	return r,w
	}

}

//func curr
//wrapper around get current
func (x *exprLex) curr() (rune) {
//log.Printf("Entering curr()")
	r,w := x.gettok()
	x.width = w
	return r
}

//func mext
//get next



// Return the next rune for the lexer.
func (x *exprLex) next() (rune) {
//log.Printf("Entering next()")
	x.pos = x.pos + x.width
	r,w := x.gettok()
	x.width = w
	return r
}

func (x *exprLex) ptok() {
//	log.Printf("          0  1   2   3   4   5   6   7   8   9   10 11 12 13 14")
	var i int = 0
	var tokspace string = ""
	var posspace string = ""
	for i < x.tok {
		tokspace += " "
		i++
	}
	i = 0
	for i < x.pos {
		posspace += " "
		i++
	}
	log.Printf("line          |%s|", x.line)
	log.Printf("tok[%2.2d]       |%s^",x.tok, tokspace)
	log.Printf("pos[%2.2d]       |%s^",x.pos, posspace)

}

// Return the current token
func (x *exprLex) emit() string {
//log.Printf("Entering emit()")
//log.Printf("len: %d tok: %d pos: %d\n", len(x.line), x.tok, x.pos)

	if (x.pos >= len(x.line)) {
//log.Printf("Emit() says, token zero length, pos at EOL")
		return ""
	}
	r := x.line[x.tok:x.pos]
	return r
}

//move up the tok pointer
func (x *exprLex) shift() {

	x.tok = x.pos
}

// The parser calls this method to get each new token.
func (x *exprLex) Lex(yylval *exprSymType) int {
//log.Printf("=====================")
//log.Printf("Entering Lex function")
	//This is called either at the very beginning of the 
	//string to be parsed or at the start of a new 
	//token
	L:
	n := x.curr()

//log.Printf("Next rune is: %c",n)
	switch {
	  case n == '\n':
//log.Printf("short circuit eof")
		return eof
	  case n == eof:
//log.Printf("Found EOF")
		return eof
	  case n == ' ' || n == '\n':
//log.Printf("Found space");
		x.consume()
		goto L
	  case n == '"':
//log.Printf("Found double quote")
		x.lexdquote()
	  case n == '\'':
//log.Printf("Found single quote")
		x.lexsquote()
	  case n >= '0' && n <= '9':
//log.Printf("Found digit")
		x.lexnumber()
	  case n == '.':
//log.Printf("Found point")
		x.lexpoint()
	  case n == ';':
//log.Printf("Found semicolon")
		x.lexterm()
	  case n == ',':
//log.Printf("Found comma")
		x.lexcomma()
	  case n == '_' || unicode.IsLetter(n):
//log.Printf("Found text")
		//Here we could match an identifier or a 
		//keyword
		//could be sensitive to order in switch stmt
		x.lextext()
	  default:
//log.Printf("Found default oper")
		x.lexoper()
	}
	//log.Printf("Lexer: %+v\n", x)

	yylval.tokval = x.emit()
	x.shift()
//log.Printf("Lexed token! yylval: %s toktyp: %d", yylval.tokval, x.typ)
	return x.typ
}

func isAlphaNumeric(r rune) bool {
   return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isWhitespace(r rune) bool {
   wschars := " " //TODO add tabs, newlines

   if strings.IndexRune(wschars, r) == -1 {
     return false
   } else {
     return true
   }
}

func isOperator(r rune) bool {
   operchars := "+-/*%!=^~"

   if strings.IndexRune(operchars, r) == -1 {
     return false
   } else {
     return true
   }
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

                exprParse(&exprLex{line: string(line)})

		log.Printf("Parse tree is: %+v", Parsetree)
        }
}
