%{

package parser

import "log"
//import "github.com/davecgh/go-spew/spew"



%}

%union 
{
	tokval		datumval
	node		pnode
	sexpr		pnode
	datum		Datum
}

//%error sql error : "error select_statement"

/* note there is no type 'keyword' in the yylval struct

   the yylval of 'keyword' tokens is never used
*/

/* SQL Datatypes */

%token <keyword> INTEGER BIGINT SMALLINT INT2 INT4 INT8 NUMERIC REAL DOUBLE 
%token <keyword> BIT DATE TIME TIMESTAMP ZONE INTERVAL PRECISION FLOAT TEXT CHAR VARCHAR

/* SQL keywords */
%token <keyword> SELECT INSERT UPDATE DELETE WHERE FROM VALUES CREATE DROP SUM 
%token <keyword> COUNT SET INTO TABLE WITH ORDER BY HAVING GROUP CASE WHEN THEN END
%token <keyword> ELSE DESC ASC FIRST LAST NULLS _NULL TRUE FALSE IS UNKNOWN

/* Literal values */
%token <tokval> INT_LIT
%token <tokval> NUM_LIT
%token <tokval> STRING_LIT

/* punctuation */
%token <keyword> QUOTE NEWLINE 

/* operators */
%token <operator> ADD SUB MUL DIV MOD EQ NE LT GT LE GE AND OR NOT IN BETWEEN

/* operator precedence */

%left           OR
%left           AND
%left		NE
%left 		IN
%right		NOT
%right		EQ
%nonassoc	LT GT
%nonassoc	LE GE
%nonassoc	BETWEEN
%left           ADD SUB
%left           MUL DIV MOD
%left           EXP
/* Unary Operators */
%right          UMINUS
%left		LPAREN RPAREN
%left		SEMICOLON COMMA
//%left         TYPECAST
%left           POINT
%left 		AS

%start sql
/* The integer value of the tok names start from zero (with query_statement)  */
%type <node>	query_statement select_statement select_list u_select_list_item select_list_item table_ref
%type <node>	table_ref_list table_expr
%type <node>	function case_expr case_expr_when_list case_expr_when from_clause 
%type <node>	order_by_list order_by_list_item order_by_clause
%type <node>	column_definition column_definition_list data_type insert_statement insert_value_list column_list
%type <node>	ddl_table_ref create_table_stmt drop_table_stmt in_predicate

/* These emit expr structures. Review this */
%type <tokval>	group_by_clause having_clause where_clause 
%type <sexpr> 	scalar_expr 

%type <datum>	value_expr colref

%type <keyword>	order_by_direction order_by_nulls boolean sqlval

%token  <tokval>  IDENTIFIER 



%%
/* SQL
 * ----------------------------------------------------------------------------
 * Multi-statement query string, delimited by semicolon
 */


sql:
    query_statement SEMICOLON
    {
	// Assign query_statement to the first parse node in ParseTree
	P.tree = append(P.tree,$1)	
//	spew.Dump(p)
    }
    |
    sql query_statement SEMICOLON
    {
	P.tree = append(P.tree,$2)	
	log.Printf("PARSER: sql query_statement SEMICOLON")
    }
;

query_statement:
    select_statement 
    { 
	$$ = makeNode(query_statement)
	$$.appendNode($1)
    } 
    |
    insert_statement
    { 
    } 
    |
    create_table_stmt
    {
    }
    |
    drop_table_stmt
    {
    }
;

/* INSERT
 * ----------------------------------------------------------------------------
 * Insert statement 
 */

insert_statement:
    INSERT INTO ddl_table_ref LPAREN column_list RPAREN VALUES LPAREN insert_value_list RPAREN 
    { 
    }
    |
    INSERT INTO ddl_table_ref VALUES LPAREN insert_value_list RPAREN 
    {
    }
    |
    INSERT INTO ddl_table_ref select_statement 
    {
    }
;

column_list:
    IDENTIFIER
    {
    } 
    |
    column_list COMMA IDENTIFIER
    {
    }
;	

insert_value_list:
	scalar_expr { 
	}	
	|		
	insert_value_list COMMA scalar_expr {	
	}
;

/* SELECT
 * ----------------------------------------------------------------------------
 * Select statement 
 */

select_list:
    select_list_item
    {
	$$ = makeNode(select_list)
	$$.appendNode($1)
    }
    |
    select_list COMMA select_list_item
    { 
	$$.appendNode($3)
    } 
;

select_list_item:
    u_select_list_item
    {
	$$=$1;
    }
    |
    u_select_list_item AS IDENTIFIER
    {
	//TODO: incorporate Alias code
	$$=$1;
    }
;

u_select_list_item:
    scalar_expr
    {
	$$ = makeNode(select_list_item)
	$$.appendNode($1)
    }
    |
    MUL
    {
//	new_tuple($$, v_text, "value", "wildcard");
    }	 
    |
    LPAREN select_statement RPAREN
    {
//	new_tuple($$, v_tuple, "subquery", $2);
    }
;

select_statement:
    SELECT select_list table_expr
    {
	$$ = makeNode(select_statement)
	$$.appendNode($2)
	$$.appendNode($3)
    }
;

/* UPDATE 
 * ----------------------------------------------------------------------------
 * Update statement 
 */

/* DELETE
 * ----------------------------------------------------------------------------
 * Delete statement 
 */



table_ref:
    IDENTIFIER
    {
	$$ = makeNode(table_ref)
	$$.addDatum($1, IDENTIFIER)
    }
    |
    IDENTIFIER IDENTIFIER
    {
	$$ = makeNode(table_ref)
	$$.addDatum($1,IDENTIFIER)	
	$$.addAttr(att_alias,$2)
    }
    |
    IDENTIFIER AS IDENTIFIER
    {
	$$ = makeNode(table_ref)
	$$.addDatum($1,IDENTIFIER)	
	$$.addAttr(att_alias,$3)
    }
    |
    LPAREN select_statement RPAREN IDENTIFIER
    {
	// Subquery as source table
    }
    |
    LPAREN select_statement RPAREN AS IDENTIFIER
    {
	// Subquery as source table
    }
;

table_ref_list:
    table_ref
    {
	$$ = makeNode(table_ref_list)
	$$.appendNode($1)
    }
    |
    table_ref_list COMMA table_ref
    {
	$$.appendNode($3)
    }
;

from_clause:
    FROM table_ref_list
    {
	$$ = makeNode(from_clause)
	$$.appendNode($2)
    }
;

where_clause:
    empty
    {
//	$$=NULL;
    }
    |
    WHERE scalar_expr
    {
//	$$=$2;
    }
;

having_clause:
    empty
    {
//	$$=NULL;
    }
    |
    HAVING scalar_expr
    {
//	$$=$2;
    }
;

order_by_clause:
    empty
    {
//	$$=NULL;
    }
    |
    ORDER BY order_by_list
    {
//	$$=$3;
    }
;

order_by_list:
    order_by_list_item
    {
//	new_tuple($$, v_tuple, "order_by_expression", $1);
    }
    |
    order_by_list COMMA order_by_list_item
    {
//	tuple_append($$, v_tuple, "order_by_expression", $3);
    }
;

order_by_list_item:
    scalar_expr order_by_direction order_by_nulls  
    {
//	new_tuple($$, v_sexpr, "value", $1);	
//	if ($2 != NULL) tuple_append($$, v_text, "direction", $2); 
//	if ($3 != NULL) tuple_append($$, v_text, "nulls", $3); 
    }
;

order_by_direction:
    empty
    {
//	$$=NULL;
    }
    |
    ASC
    {
//	$$="asc";
    }
    |
    DESC
    {
//	$$="desc";
    }
;

order_by_nulls:
    empty
    {
//	$$=NULL;
    }
    |
    NULLS FIRST
    {
//	$$="first";
    }
    |
    NULLS LAST
    {
//	$$="last";
    }
; 

group_by_clause:
    empty
    {
//	$$=NULL;
    }
    |
    GROUP BY scalar_expr
    {
//	$$=$3;
    }
;

table_expr:
    from_clause where_clause group_by_clause having_clause order_by_clause
    {
	$$ = makeNode(table_expr)
	$$.appendNode($1)
    }
;


/*

EXPRESSIONS

*/


scalar_expr:
    value_expr
    { 
	log.Printf("Parser: found value_expr %+v",$1)
	$$ = makeNode(scalar_expr)
	$$.addDatum0($1)
    }
    |
    LPAREN scalar_expr RPAREN
    { 
	$$ = $2;
    }
    |
    scalar_expr ADD scalar_expr 
    {
	log.Printf("parser: found Exp ADD Exp")
	$$ = makeOperScalarExpr(ADD,$1,$3)
    }
    |
    scalar_expr MUL scalar_expr 		
    {
	$$ = makeOperScalarExpr(MUL,$1,$3)
    }
    |
    scalar_expr DIV scalar_expr 		
    {
	$$ = makeOperScalarExpr(DIV,$1,$3)
    }
    |
    scalar_expr MOD scalar_expr 		
    {
	$$ = makeOperScalarExpr(MOD,$1,$3)
    }
    |
    scalar_expr AND scalar_expr 		
    {
	$$ = makeOperScalarExpr(AND,$1,$3)
    }
    |
    scalar_expr OR scalar_expr 		
    {
	$$ = makeOperScalarExpr(OR,$1,$3)
    }
    |
    scalar_expr EQ scalar_expr 		
    {
	$$ = makeOperScalarExpr(EQ,$1,$3)
    }
    |
    scalar_expr NE scalar_expr 		 
    {
	$$ = makeOperScalarExpr(NE,$1,$3)
    }
    |
    scalar_expr GT scalar_expr 		
    {
	$$ = makeOperScalarExpr(GT,$1,$3)
    }
    |
    scalar_expr LT scalar_expr 		
    {
	$$ = makeOperScalarExpr(LT,$1,$3)
    }
    |
    scalar_expr GE scalar_expr 		
    {
	$$ = makeOperScalarExpr(GE,$1,$3)
    }
    |
    scalar_expr LE scalar_expr 	
    {
	$$ = makeOperScalarExpr(LE,$1,$3)
    }
    |
    scalar_expr SUB scalar_expr 	
    {
	$$ = makeOperScalarExpr(ADD,$1,$3)
    };

/*
    |
    scalar_expr IN LPAREN in_predicate RPAREN
    {
    }
    |	
    scalar_expr NOT IN LPAREN in_predicate RPAREN
    {
    }
    |	
    scalar_expr BETWEEN scalar_expr
    {
    }
    |
    scalar_expr NOT BETWEEN scalar_expr
    {
    }
    |
    scalar_expr IS scalar_expr
    {
//	mk_s_expr_oper($$, "IS", $1, $3);
    }
    |
    scalar_expr IS NOT scalar_expr
    {
//	mk_s_expr_oper($$, "ISNOT", $1, $4);
    }
;
*/

value_expr:
	colref
	{ 
	log.Printf("PARSER: value_expr->Found colref %+v", $1)
	    $$=$1;
	log.Printf("PARSER: value_expr->Found colref %+v", $$)
	}
	|
	boolean
	{
//	    mk_tuplist_lit($$, v_text, "BOOL", $1);
	}
	|
	sqlval
	{
//	    mk_tuplist_lit($$, v_text, "SQLV", $1);
	}
	|
	INT_LIT
	{
//	    mk_tuplist_lit($$, v_int, "INT", $1);
	}
	|	
	NUM_LIT 
	{
		$$ = Datum{
				value: $1,
				dtype: NUM_LIT}
	}
	|	
	STRING_LIT
	{
//	    mk_tuplist_lit($$, v_text, "TEXT", $1);
	}
	|
	function
	{
//	    $$=$1;
	}
;

boolean:
    TRUE    
    {
//	$$="true";	
    }
    |
    FALSE
    {
//	$$="false";	
    }
;

sqlval:
    _NULL
    {
//	$$="sqlnull";	
    }
    |
    UNKNOWN
    {
//	$$="unknown";	
    }
;

function:
    case_expr 	
    {
//	$$=$1;
    }
;

case_expr:
    CASE case_expr_when_list ELSE scalar_expr END 
    {
//	new_tuple($$, v_tuple, "when_list", $2);
//	tuple_append($$, v_sexpr, "else", $4);
    }
    |
    CASE case_expr_when_list END
    {
//	new_tuple($$, v_tuple, "when_list", $2);
    }
;

case_expr_when_list:
    case_expr_when
    {
//	new_tuple($$, v_tuple, "when", $1);
    }
    |
    case_expr_when_list case_expr_when
    {
//	tuple_append($$, v_tuple, "when", $2);
    }
;

case_expr_when:
    WHEN scalar_expr THEN scalar_expr
    {
//	new_tuple($$, v_sexpr, "condition", $2);
//	tuple_append($$, v_sexpr, "result", $4);	
    }
; 

colref:
	IDENTIFIER 
	{ 
	log.Printf("Parser: I found an Identifier!!")
	log.Printf("IDENTIFIER: node value %s", $1) 
	log.Printf("------")
	$$ = Datum{
		value: $1,
		dtype: IDENTIFIER}

	log.Printf("IDENTIFIER: Datum is: %+v", $$)
	}

	|
	IDENTIFIER POINT IDENTIFIER  
	{
//		mk_tuplist_ident($$, $1, $3);
	}
;

in_predicate:
	scalar_expr {
		} |
	in_predicate COMMA scalar_expr {
	}
;


/* Data definition language commands */

/* Drop Table */

drop_table_stmt:
	DROP TABLE ddl_table_ref
	{
	}
;

/* Create Table */

ddl_table_ref:
	IDENTIFIER 
	{
	}
	|
	IDENTIFIER POINT IDENTIFIER  
	{
	}
;

data_type:
	INTEGER	
		{
		}
		|
	NUMERIC 
		{
		}
		|
	CHAR	{
		}
;

create_table_stmt:
	CREATE TABLE ddl_table_ref LPAREN column_definition_list RPAREN
	{
	}
;

column_definition_list:
	column_definition
	{
	} 
	|
	column_definition_list COMMA column_definition
	{
	}
;

column_definition: 
	IDENTIFIER data_type
	{
	}

;

empty: ;
%%
