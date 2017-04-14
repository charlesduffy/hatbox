define keyword_preamble
package main\nvar SQLkeys = map[string]int{
endef

GOYACC=$$GOPATH/src/golang.org/x/tools/cmd/goyacc/goyacc

parser: lex.go y.go keywords.go
	go build -o parser lex.go y.go keywords.go

y.go:	gram.y
	$(GOYACC) -p "expr" gram.y

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

.PHONY: clean

clean:
	rm -f keywords.go y.go
