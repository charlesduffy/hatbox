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
import "regexp"


const (
	SELECT int = iota
)

var SQLkeys = map[string]int{
        "select": SELECT,
}

// Token type contains code val and regexp pattern, not used for keywords
type token struct {
        code int
        pat regexp.Regexp
}
// Precompiled regular expressions for matching tokens
var Regint = regexp.MustCompile(regintpat)                      //integer token
var Regnum = regexp.MustCompile(regnumpat)                      //numeric token
var Regnym = regexp.MustCompile(regnympat)                      //identifier token

// Constant values for regular expression patterns
const (
        regintpat string = "^[0-9]+$"
        regnumpat string = "^[0-9]+\\.[0-9]+$"
        regnympat string = "^[a-zA-Z_]+\\w*"
)

// The parser expects the lexer to return 0 on EOF.  Give it a name
// for clarity.
const eof = 0

// Quoting state enum
const (
        unquoted int = iota
        string_quoted
        ident_quoted
)

// The parser uses the type <prefix>Lex as a lexer.  It must provide
// the methods Lex(*<prefix>SymType) int and Error(string).
type exprLex struct {
        line string                     //Buffer containing SQL input text to be lexed
        stack string                    //Token character stack
        c rune				//"current" character
        qs int                          //quoting state
        ws bool                         //whitespace state
        eat bool                        //"eat quote" flag
}

func (x *exprLex) popbuf() bool {
        if len(x.line) < 1 {
                return false
        }
        x.c = x.line[:1]
        x.line = x.line[1:]
        return true
}

func (x *exprLex) peek() rune {
        if len(x.line) < 1 {
                return nil
        }
        return x.line[:1]
}

func (x *exprLex) pushstack() bool {

        if x.eat == true {
                x.eat = false
                return false
        }

        //push the char on to stack...
        append(x.stack, x.c)

}

func (x *exprLex) matchTok(yylval *exprSymType) int {
  //Pop the contents of the stack into a variable tok.
  //y. If quoted != false, set class to STRING
  //Attempt to match the contents of tok to a token class.
  //Return the token class and yylval (or error if the token cannot be matched). Go to 1.

  // Token match detail
  // If quoted - we are a string or a identifier (hopefully)

  tokval := x.stack
  x.stack = nil
  yylval = nil

  switch x.qs {
        case string:
                yylval.text = tokval
                return STRING
        case identifier:
                yylval.text = tokval
                return IDENTIFIER
        }

  // Attempt to match a keyword
  n := x.matchKeyword(tokval)
  if (n != nil) {
        return n
  }

  // Attempt to match an operator
  switch tokval {
   case '+' : return SUM
   case '-' : return SUB
   case '*' : return MUL
   case '/' : return DIV
   case '=' : return EQ
//   case '\!=': return NE
   case ')' : return RPAREN
   case '(' : return LPAREN
   case '%' : return MOD
   case ',' : return COMMA
  }

  // Attempt to match an integer


  // Attempt to match a numeric

  // Attempt to match an unquoted identifier

}

// The parser calls this method to get each new token.
func (x *exprLex) Lex(yylval *exprSymType) int {

        var rune p

        //start point
        start:

        //pop char from  input buffer
        if (x.popbuf() == false) {
                return eof
        }

        //if eat_quote is true, do not push the char to the stack
        x.pushstack()


        p = x.peek()

        // if we are in a quoted state and the peek-ahead character is not a matching end-quote
        // keep pushing characters on to the stack.
        // consider introducing a token-length parameter here
        for (x.qs != unquoted && x.matchQuote(p) == false) {
                x.popbuf()
                x.pushstack()
                p = x.peek()
        }

        // if we are in a quoted state and the peek-ahead character is a matching end-quote go to z.
        // Note we do not push the quote to the stack

        if (x.qs != unquoted && x.matchQuote(p) == true) {
                tc = x.matchTok(yylval)
                return tc
        }

        // if we are in an unquoted state and the peek-ahead character is a quote, enter the relevant quoted state (identifier or string).
        // Set eat_quote to true. Go to z.
        //fix the bug in the logic around "eat" here
        if (x.qs == unquoted && x.matchQuote(p) != notquote) {
                x.qs = x.matchQuote(p)
                x.eat = true
                tc = x.matchTok(yylval)
                return tc
        }

        //if we are in unquoted state and whitespace_class of the peek-ahead character does not
        //equal the current ws_state, set ws_state to the whitespace_class of the peek-ahead character and go to z.
        if (x.qs == unquoted && x.matchWs(p) != x.sw) {
                x.ws = x.matchWs(p)
                tc = x.matchTok(yylval)
                return tc
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
