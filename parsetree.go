package main
import "fmt"

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
//need a set of constants
type NodeTag int

const (

	//SQL clause types
	STATEMENT NodeTag = iota
	QUERY
	SELECT_LIST
	SELECT_LIST_ITEM
	//Other Clause metadata
	//

)

//no need, just use Stringer. To be removed. 
var sql_clause = [...]string {
	"statement",
	"query",
	"select_list",
	"select_list_item",
}

//need a thing to put in the parse tree nodes
//the thing can hold another slice of tree nodes,
//or an s-expression,

//need a parse_tree_node thing

type Pnode struct {
	tag	NodeTag
	val	interface{}
}

//need an interface for it 

//need a factory method



func main() {
/*	var slice1 []int = make([]int, 10)

	for i := 0; i < len(slice1); i++ {
	slice1[i] = 5 * i
	}

	for i := 0; i < len(slice1); i++ {
	fmt.Printf("Slice at %d is %d\n", i, slice1[i])
	}

	fmt.Printf("Length of slice1 is is %d\n", len(slice1))
	fmt.Printf("Capacity of slice1 is is %d\n", cap(slice1))
*/

	fmt.Printf("hello!\n")

}
