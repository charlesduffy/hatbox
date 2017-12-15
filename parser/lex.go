package parser

import (	"log"
		"bufio"
		"os"
		"io"
		"unicode"
		"unicode/utf8"
		"strings"
		"github.com/davecgh/go-spew/spew"
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

type datumval interface {

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

	x.typ = STRING_LIT

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
	x.typ = NUM_LIT

	for {
	  n := x.next()
	  if unicode.IsDigit(n) != true && n != '.' {
	  //here we have encountered a rune that is not
	  //accepted within a numeric token. 
	    if  isOperChar(n) || n == ';'|| isWhitespace(n) {
	     return
	    } else {
	    x.Error(x.line)
	    return
	    }
	  }
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

// TODO: Modify to handle the case where identifiers
//       may only begin with a particular set of 
//       characters

	for {
	  n := x.next()
	  if ! isAlphaNumeric(n) {
	  //here we have encountered the end of the token. 
	  //determine if it is a keyword.
	    o := x.emit()

	    if l, ok := SQLkeys[strings.ToLower(o)]; ok {
	      x.typ = l
	    } else {
		x.typ = IDENTIFIER
		log.Printf("lexer: IDENTIFIER: %s", o)
	   }
	    return
	  }
	}
}

// Lex an operator. An operator is one or more
// characters from the list of:
// + - * / % ! = > <
//
// multi-char opers
// 
// >= <= != <> 

func (x *exprLex) lexoper() {

	for {
		n := x.next()
		if ! isOperChar(n) {

			o := x.emit()

			switch(o) {
				case "+":
					x.typ = ADD
				case "-":
					x.typ = SUB
				case "*":
					x.typ = MUL
				case "/":
					x.typ = DIV
				case "%":
					x.typ = MOD
				case "=":
					x.typ = EQ
				case ">":
					x.typ = GT
				case "<":
					x.typ = LT
				case ">=":
					x.typ = GE
				case "<=":
					x.typ = LE
				case "!=":
					x.typ = NE
				case "<>":
					x.typ = NE
			}
			log.Printf("lexer: OPER lval: %s toktyp: %d", o, x.typ)
			return
		}
	}
}


// Lex a point character. In this context it 
// will be a delimiter between table / column 
// identifiers
func (x *exprLex) lexpoint() {
	x.typ = POINT
	x.consume()
	return
}

// Lex a comma. Delimiter between 
// select list items or table expression
// items.
func (x *exprLex) lexcomma() {
	x.typ =  COMMA
	x.consume()
	return
}


func (x *exprLex) lexterm() {
	x.typ = SEMICOLON
	x.consume()
	return
}

//func get tok
//get the 'current' token
func (x *exprLex) gettok() (rune, int) {

	if (x.tok == x.pos && x.pos == len(x.line) - 1) {
		//end of the line
		return eof, 0
	}

	r, w := utf8.DecodeRuneInString(x.line[x.pos:])
	if w == 0 {
		return eof, 0
	} else {
	return r,w
	}

}

//func curr
//wrapper around get current
func (x *exprLex) curr() (rune) {
	r,w := x.gettok()
	x.width = w
	return r
}

// Return the next rune for the lexer.
func (x *exprLex) next() (rune) {
	x.pos = x.pos + x.width
	r,w := x.gettok()
	x.width = w
	return r
}

func (x *exprLex) ptok() {
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

	if (x.pos >= len(x.line)) {
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
	//This is called either at the very beginning of the 
	//string to be parsed or at the start of a new 
	//token
	L:
	n := x.curr()

	switch {
	  case n == '\n':
		return eof
	  case n == eof:
		return eof
	  case n == ' ' || n == '\n':
		x.consume()
		goto L
	  case n == '"':
		x.lexdquote()
	  case n == '\'':
		x.lexsquote()
	  case n >= '0' && n <= '9':
		x.lexnumber()
	  case n == '.':
		x.lexpoint()
	  case n == ';':
		x.lexterm()
	  case n == ',':
		x.lexcomma()
	  case isOperChar(n):
		//Lex a symbolic operator character
		x.lexoper()
	  case n == '_' || unicode.IsLetter(n):
		//Here we could match an identifier or a 
		//keyword
		//could be sensitive to order in switch stmt
		x.lextext()
	}

	yylval.tokval = x.emit()
	x.shift()

	log.Printf("Lexer: Token text is: %v token type is: %d", yylval.tokval, x.typ)

	return x.typ
}

func isAlphaNumeric(r rune) bool {
   return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isWhitespace(r rune) bool {
   wschars := " \t\n" //TODO add tabs, newlines

   if strings.IndexRune(wschars, r) == -1 {
     return false
   } else {
     return true
   }
}

func isOperChar(r rune) bool {
   operchars := "+-/*%!=^~><"

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

