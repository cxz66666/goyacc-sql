# goyacc-sql
A simple sql interpreter created by goyacc, use goyacc to interpret sql statement, you can check spanner.go.y to find which type of sql statement is support. Also you can find [another repo](https://github.com/cxz66666/MiniSQL) to find a complete implementation of minisql.

## How to use it
- go run main.go
- input some sql statement.
- check the type of statement and try to cast it to it's struct from duck type
- print the whole struct and compare with what you want



## What is lex and yacc
[This](https://pingcap.com/blog-cn/tidb-source-code-reading-5/) is a very good sample for me to learn more about, even Tidb also use goyacc as their Parser, but their codes is too difficult to read, and I write this for beginner as me.


## Why we use channel

You can write the another version without channel quickly, but it will not have a better performance when you ready to interpret more than 10w statements.
So channel and pipeline is a greater idea.

##What kind of statement it can parse
Almost every sql statement you can, you can find [this repo](https://github.com/cxz66666/MiniSQL) to learn more about, also you can give it more support to other statements.