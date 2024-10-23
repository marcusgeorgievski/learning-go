## 5 - Functions

Functions operate as they do in other languages. A function declarations includes; `func` keyword, function name, input parameters, and return type.

Go requires all paramters to be provided. It is possible to simulate **named and optional parameters** by defining a struct with fields that match the desired parameters, then pass struct into function. The unspecified fields will contains the zero value.

**Variadic parameters** must be the last or only parameter in the input parameter list. They are indicated with `...` before the type. A slice variable of the specified type is created.

**Named returns** allow you to give return types in the declaration an identifier. This will predeclare them as variables and initialize them to their type's zero value. If you only want to give some return values a name, use `_` for the others. Even one named return is required to be in parenthesis. Other values can be returned in their place. Naming them just gives a way to declare the intent to use the variables, but does not require those specific variables to be returned.

> Avoid blank/naked returns

```go
func foo(options UserOptions) // struct param
func foo(id int, vals ...string) // variadic
func foo() (int, int, error) // multiple returns

// named return values
func foo() (result int, err error) {
  if cond {
    err = errors.New("some error")
  } else {
    int = 10
  }
  return result, err
}
```

**Functions Are Values**. The `func` keyword, parameter list, and return values make up the function _signature_. Since functions are values, we can declare them as variables. They have a zero value of `nil`.

```go
var foo func(string) int
foo = someOtherFunc

// example:

type opFunc func(int int) int
var calculatorMap = map[string]opFunc
// or just
var calculatorMap = map[string]func(int int) int {
  "+": add // these are funtions
  "-": sub
  // ...
}
```

Anonymous function exclude their name. They are useful for `defer` statements and goroutines.

### Closures

Closures are functions declared inside function that are able to access and modify variables declared in the outer function. They are useful for functions that will only be called from one other function. This hides it from the package scope and reduces number of declarations at package scope, which can make it easier to find an unused identifier.

Functions can be passed into parameters, and can be returned from other functions. Function that accept or return function are called _Higher-Order Functions_.

Here is an example of an anonymous function and a closure in one.

```go
users := []User{
  // ...
}

  // Func as param: sort by age
slices.Sort(users, func(i, j int) bool {
  return users[i].age > users[j].age
})
```

### Defer

`defer` is a keyword that will execute a function at the end of its block, after any return statement. It is useful for cleaning up temporary resources, such as handling files and database connection. Any function, method, or closer can be used with `defer` Multiple defers will run in LIFO order.

A common pattern in Go is for a function that allocates a resource to also return a closure that cleans up the resource. This allows that closure to be defered in the block the function that returns it is called in.

```go
defer foo()
defer func() {

}() // dont' forget `()` at the end
```

### Call by Value

Go is call by value. When a variable is supplied to a function parameter, Go _always_ makes a copy. There is different behaviour for slice and map types:

A **slice** passed into a function can have its values modified, but its length cannot be increased with `append`. Appending values will place them in the following memory locations and increase the length of the _copy_, but they cannot be accessed by the _original_ due to the original slice length not being increased. Any changes made to a **map** parameter are reflected in the variable passed into the function.

Call by value is one reason Go's limited support for constants is only a minor handicap. We can be sure that calling a function does not modify the variables passed in (except for maps and slices).
