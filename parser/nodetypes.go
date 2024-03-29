package parser
const (

	query_statement	= iota
	select_statement
	select_list
	u_select_list_item
	select_list_item
	table_ref
	table_ref_list
	table_expr
	function
	case_expr
	case_expr_when_list
	case_expr_when
	from_clause
	order_by_list
	order_by_list_item
	order_by_clause
	column_definition
	column_definition_list
	data_type
	insert_statement
	insert_value_list
	column_list
	ddl_table_ref
	create_table_stmt
	drop_table_stmt
	in_predicate
	scalar_expr

)

var NodeYNames = []string{

	"query_statement",
	"select_statement",
	"select_list",
	"u_select_list_item",
	"select_list_item",
	"table_ref",
	"table_ref_list",
	"table_expr",
	"function",
	"case_expr",
	"case_expr_when_list",
	"case_expr_when",
	"from_clause",
	"order_by_list",
	"order_by_list_item",
	"order_by_clause",
	"column_definition",
	"column_definition_list",
	"data_type",
	"insert_statement",
	"insert_value_list",
	"column_list",
	"ddl_table_ref",
	"create_table_stmt",
	"drop_table_stmt",
	"in_predicate",

}
