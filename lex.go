package main

import ( 	"log"
		"bufio"
		"os"
		"io"
		"unicode"
		"unicode/utf8"
		"strings"
)

const eof = 0

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
        line		string                  //Buffer containing SQL input text to be lexed
	tok		int			//start of current token
	pos		int			//current character position in input
	typ		int			//current token type
	w		int			//width of rune
}

func (x *exprLex) peek() rune {
        if len(x.line) < 1 {
                return 0
        }
	r, _ := utf8.DecodeRuneInString(x.line[x.pos:])
	return r
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
log.Printf("Entering lexnumber()")
	x.typ = NUMERIC

	for {
	  n := x.next()
	  if unicode.IsDigit(n) != true && n != '.' {
log.Printf("Not a digit or dot")
	  //here we have encountered a rune that is not
	  //accepted within a numeric token. 

log.Printf("Peek character is: >>%c<<", n)

	    if  isOperator(n) || n == ';'|| isWhitespace(n) {
	     p:=x.emit() //just for debug
log.Printf("Token is: >>%s<<\n", p)
	     return
log.Printf("after return???")
	    } else {
	       //raise error. 
	    x.Error(x.line)
	    return
	    }
	  }
log.Printf("Another numeric rune...")
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
	x.typ = IDENTIFIER

	for {
	  n := x.next()
	  if ! isAlphaNumeric(n) {
	  //here we have encountered the end of the token. 
	  //determine if it is a keyword.
	    o := x.emit()

log.Printf("Token is: >>%s<<\n", o)
	    if l, ok := SQLkeys[strings.ToLower(o)]; ok {
	      x.typ = l
	    }
	    return
	  }
	}
}

// Lex a point character. In this context it 
// will be a delimiter between table / column 
// identifiers
func (x *exprLex) lexpoint() {

}

// Lex an operator. An operator is one or more
// characters from the list of:
// + - * / % ! =
func (x *exprLex) lexoper() {

}

func (x *exprLex) lexterm() {
	x.typ = SEMICOLON
	return
}

// Return the next rune for the lexer.
func (x *exprLex) next() rune {

	r, w := utf8.DecodeRuneInString(x.line[x.pos:])
	x.pos = x.pos + w
	return r
}

// Return the current token
func (x *exprLex) emit() string {

	if (x.pos - 1 > len(x.line)) {
		return ""
	}
	r := x.line[x.tok:x.pos-1]
	return r
}

//move up the tok pointer
func (x *exprLex) shift() {

	x.tok = x.pos
}

// The parser calls this method to get each new token.
func (x *exprLex) Lex(yylval *exprSymType) int {
log.Printf("Entering Lex function")
	//This is called either at the very beginning of the 
	//string to be parsed or at the start of a new 
	//token
	n := x.next()

log.Printf("Next rune is: %c",n)
	switch {
	  case n == eof:
log.Printf("Found EOF")
		return eof
	  case n == ' ':
		x.ignore()
	  case n == '"':
		x.lexdquote()
	  case n == '\'':
		x.lexsquote()
	  case n >= '0' && n <= '9':
		x.lexnumber()
	  case n == '_' || unicode.IsLetter(n):
		//Here we could match an identifier or a 
		//keyword
		x.lextext()
	  case n == '.':
		x.lexpoint()
	  case n == ';':
		x.lexterm()
	  default:
		x.lexoper()
	}
	log.Printf("Lexer: %+v\n", x)

	yylval.tokval = x.emit()
	x.shift()
log.Printf("Lexed token! yylval: %s toktyp: %d", yylval.tokval, x.typ)
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
        }
}
