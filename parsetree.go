package main
//import "fmt"

//Parse tree node types enumeration
//TODO: figure out a way to encapsualte this in 
//a more idiomatic fashion

type nodetype int

// Node types are autogenerated from gram.y in Makefile

// Parse tree "wrapper" struct
type ptree struct {
	tree	[]pnode
	query	string
}

// For now, we just let the parser
// add node values to the tree from 
// the bottom-up. We may wish to add fancier 
// methods that permit constructing the tree 
// in a more elegant way later

// Datum interface
// Need to add some methods here 
// but leave it blank for the time being

// Parse tree node
type pnode struct {
	tag	nodetype
	subtree []pnode
	value	datum
}

// 

/*

	
struct s_expression {
  tuple *value;
  llist list;
  s_expr *left;
  s_expr *right;
};

struct ord_pair {
    char *tag;
    ttype type;
    union {
        unsigned int v_long;
        int v_int;
        float v_float;
        char * v_text;
        tuple * v_tuple;
        s_expr * v_sexpr;
    };
    llist list;
};


node data payload ( within s-expr )
	
	sqltyp
	value (go datatype "union")


type TList struct {
	
}

type Tree1 struct {
	
}

*/
