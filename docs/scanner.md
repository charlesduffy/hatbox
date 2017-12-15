Lexical Analyser
================

Lexer token types:
------------------

IDENTIFIER

STRING

NUMERIC

OPERATOR


Lexer flowchart - functional version
-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

Order of token class lex attempts:

acceptKeyword() - tries to match against SQL keyword
	
	* acceptRun() on aA-zZ, at end then try to match case

* functional pseudocode



	lexText()


	  //try to lex keyword





Lexer flowchart - stateful version
----------------------------------

a. set quoted state = false

b. set ws_state = false (whitespace state)

c. set eat_quote = false

1. read one character from the input buffer (containing SQL query text)

2. if eat_quote is true, set eat_quote to false and go to 3. Else, push character on to the stack. 

3. peek ahead one character

4. if we are in a quoted state and the peek-ahead character is not a matching end-quote go to 1. 

5. if we are in a quoted state and the peek-ahead character is a matching end-quote go to z. Note we
   do not push the quote to the stack. 

6. if we are in an unquoted state and the peek-ahead character is a quote, enter the relevant quoted
   state (identifier or string). Set eat_quote to true. Go to z. 

7. if we are in unquoted state and whitespace_class of the peek-ahead character does not equal the 
   current ws_state, set ws_state to the whitespace_class of the peek-ahead character and go to z.

8. go to 1. 

z. Pop the contents of the stack into a variable tok. 

y. If quoted != false, set class to STRING
   Attempt to match the contents of tok to a token class.
   Return the token class and yylval (or error if the token cannot be matched). Go to 1.






Other notes
-----------

* Token class matching to be done by regex. We hold an array of compiled regexes and match the 
  stack against	each one in order in a loop. 

* Keywords are held in a map. The token is checked against these first before being classed a string, identifier or whatever.
