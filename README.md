# go-lsp
Implementing an entire Language Server in Go

## What is a Language Server?
A Language Server is a program that provides language-specific features like auto complete, go to definition, find all references etc. to an editor or an IDE. The editor or IDE communicates with the language server over a protocol called Language Server Protocol (LSP). The LSP is a JSON-RPC based protocol that is designed to be stateless and "language agnostic".

ref) https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/


## Build
```go
go build -ldflags="-s -w" -o bin/golsp main.go
```