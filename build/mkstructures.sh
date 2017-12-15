#!/bin/bash

G=parser/gram.y
K=parser/keywords.go
N=parser/nodetypes.go
PKG=parser

KP="package ${PKG}\nvar SQLkeys = map[string]int{"
NP="package ${PKG}\ntype ParseNode int\n const (\n"
NN="var NodeYNames = []string{\n"

echo -e "$KP" > $K
awk '/^\%token <keyword>/ {for (i=3;i<=NF;i++) { printf "\t\"%s\": %s,\n",tolower($i),toupper($i)}}' $G >> $K
echo -e "\n}" >> $K

#produce the node types const
echo -e "$NP" > $N
awk 'BEGIN { x = "\tParseNode = iota"; }; /^\%type <node>/ {for (i=3;i<=NF;i++) { printf "\t%s%s\n",$i,x; x="";} }' $G >> $N
echo -e "\n)" >> $N

echo -e "\n$NN" >> $N
awk '/^\%type <node>/ {for (i=3;i<=NF;i++) { printf "\t\"%s\",\n",$i} }' $G >> $N
echo -e "\n}" >> $N
