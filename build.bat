@echo off

REM Go get commands
go get github.com/99designs/gqlgen/internal/imports
go get github.com/99designs/gqlgen/internal/code
go get github.com/99designs/gqlgen/codegen/config@v0.17.34
go get github.com/99designs/gqlgen@v0.17.34

REM Generate GraphQL
go run github.com/99designs/gqlgen generate