define keyword_preamble
package main\nvar SQLkeys = map[string]int{
endef

define nodetypes_preamble
package main\nconst (\n
endef

GOYACC=$$GOPATH/src/golang.org/x/tools/cmd/goyacc/goyacc
PARSER=lex.go y.go keywords.go nodetypes.go

parser: $(PARSER)
	go build -o parser $(PARSER)
#TODO use proper Makefile macros here for target / deps

y.go:	gram.y
	$(GOYACC) -p "expr" gram.y

#Below we have some extremely crude metaprogramming
#which absolutely has to be replaced

keywords.go: gram.y
	#produce the keywords map 
	#========================
	#We need a way to minimise the number of places we define 
	#all the SQL keywords. This seems the least painful for the 
	#time being. We define them in the Yacc grammar then parse that here,
	#then we generate a golang source file containing the keywords and corresponding
	#constant values stuffed into a map. 
	echo "$(keyword_preamble)" > keywords.go
	awk '/^\%token <keyword>/ {for (i=3;i<=NF;i++) { printf "\t\"%s\": %s,\n",tolower($$i),toupper($$i)}}' gram.y >> keywords.go
	echo "\n}" >> keywords.go

nodetypes.go: gram.y
	#produce the node types const
	#========================
	echo "$(nodetypes_preamble)" > nodetypes.go
	awk '/^\%type <Tuple>/ {for (i=3;i<=NF;i++) { printf "\t%s\n",$$i} }' gram.y >> nodetypes.go
	echo "\n)" >> nodetypes.go
	
.PHONY: clean

clean:
	rm -f keywords.go y.go nodetypes.go parser y.output
