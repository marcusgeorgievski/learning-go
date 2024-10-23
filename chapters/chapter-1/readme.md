## 1 - Environment and Tooling

**`go fmt`**

Formatting tool that enforces a standard format for Go. Simplifies the compiler and allows for the creation of clever tools for generating code. Created for better tooling. Can be run throughout a project with `go fmt ./...` to apply to all files and subdirectories.

**`go vet`**

Tool to detect errors that are syntactically valid, but quite likely incorrect. For example, `fmt.Printf("Name: %s")` is valid Go, but the template is missing a value.

**`go build`**

Creates an executable (single native binary) with the same name as the module

**Semicolon Insertion Rule**

Semicolons are required at the end of every statement. The Go compiler does this itself and the developer should never do so. They are inserted after the following tokens: identifier, basic literal, and more such as `break continue return ++ -- ) }`. This can lead to errors if the opening bracket `{` is on a new line after a function signature, as a semicolon will be placed before it

```go
func main(); <- compiler will insert this
{
    ...
}
```

**Makefiles**

Tool to specify a set of operations

**Go Compatibility Promise**

The Go team has made a commitment to ensuring that all their 1.x versions are backwards compatible. Go releases are incremental rather than expansive for this reason.

### Summary

This chapter discusses basic Go tools for formatting, vetting, and building Go code, as well as how `make` can boost your productivity through automation.
