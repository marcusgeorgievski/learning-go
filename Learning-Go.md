# Learning Go

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

## 2 - Predeclared Types and Declarations

**Predeclared Types:** Types that are built into the language; booleans, integers, floats, strings.

**Zero Value**

Predeclared Go types have _zero values_ to remove the source of many bugs that can be found in C and C++. These are assigned to variables that are declared but not assigned a value.

**Literals**

Literals are explicilty defined numbers, characters, or strings. They are considered _untyped_.

- _Integer literals_ can have specified bases with the following prefixes: `0b` binary, `0x` hexadecimal, `0o` octal. You may also include `_` between numbers for readability. They must not be at the start, end, or beside one another.
- _Floating-point literals_ have decimal points to denote the fractional portion. They can have an exponent with an `e`, `6.43e2 = 6.43 * 10^2`.
- _Rune literals_ represent a single character in single quotes. Can be written as 8-bit octal numbers, 16-bit hexadecimals, 32-bit unicode, and backslash-escape characters.
- _String literals_ are surrounded by double quotes or backticks. Double quote strings are called interpreted because they interpret rune literals into a single character. Backticks strings are raw string literals. You can include characters such as `"` and `\` without escaping them.

**Choosing Which Integer to use**

- When working with binary files or network protocols, use corresponding type
- When writing a library function that should work with any integer, use a generic type parameter
- Else, use `int`

**Types**

Go does not allow implicit type conversions; _automatic type promotion_. Types must be explicitly casted to use with other types. Additionality, there are no truthy or falsey values in Go. You must check for empty or zero values with the respective operation (eg. `str == ""`, `int == 0`)

**`var` vs `:=`**

```go
var x int
var x int = 0
x := 0
```

`var` can be used to declare and initialize a value. If you declare a value with `var`, it is assigned its type's zero value.

`:=` can be used to assign a new variable without the type or var keywords. Legal as long as at least one of the variables on the left hand side are new. Not legal outside of function scope (illegal in package level). Avoid using `:=` in the following cases

1. When initializing a zero value - use `var`
2. When the default type for constant is not what you want for a variable - `const x = 10` will return an int64 to a variable assignment with `:=`
3. When you are unsure if you are creating or assigning new variables

Both methods can be used for multiple assignments, although multiple assignments with `:=` should only be used for function return values and the ok idiom.

**`const`**

Constants are a way to give names to literals - variables cannot be constant in Go. They can be typed or untyped. Untyped consts act like literals. Typed constants can be assigned to variables of the same type.

**More**

- Should only declare vars in package block when they are immutable; difficult to track how they are used throughout program
- Variables can be any unicode
- `rune` is an alias for `int32` and should be used for clarity when working with characters

### Summary

Chapter 2 reviews Go's built-in types, variables, and constants. Each type has its own characteristics. Idiomatic practices, such as using `rune` over `int32`, are explained.

## 3 - Composite Types

### Arrays

Arrays in Go limited in their usage. Unlike other languages, their size is a part of their type. Consequently;

- Variables cannot be used to specify the size of an array
- Arrays of different sizes cannot be compared
- Arrays of different sizes cannot be type converted to each other.

Those with the same size can be compared with `==` and `!=`, which checks for same length and values. Arrays should not be used unless the exact length is known ahead of time.

```go
var a = [3]int{1,2,3}
var b = [...]int{1,2,3}
var c = [5]int{1, 2, 3: 55} // sparse array [1,2,0,55,0,0]

a == b // true
```

### Slices

Slices are more flexible than arrays. Their length is not a part of their type, which allows for much more dynamic handling and operations such as appending, resizing, and slicing. Their underlying data structure is an array.

They have a zero value of `nil`; an identifier that represents the lack of a value for some types. A nil slice returns a length of 0.

Slices are not comparable; use `slices.Equal` to compare length and elements, and `slices.EqualFunc` to compare slices using a custom equality function. Avoid `reflect.DeepEqual` as it is legacy code mostly used for testing and is slower and less safe.

Common operations include `append`, `clear`, and `make`. `append` will append one or more values to the slice. `clear` will set all the values to the zero value. `make` will create the underlying backing array by specifying its type, size, and capacity.

Slices have two main properties, **length** and **capacity**.

- Length is the number of consecutive memory locations that have been assigned a value.
- Capacity is the number of consecutive memory locations that have been reserved

Appending past the capacity causes the go runtime to allocate a new backing array with a larger capacity. Slices will grow by `2x` if the length is under 256, or `(currCap + 768) / 4`, which converges to 25% as the length increases. It is far more efficient to size them once. Note that appending to a slice that has been made with `make([]int, 5)` will cause a resize, as appends occur at the end of the current length.

> Runtime note: memory allocation, garbage collection, networking, concurrency, etc, are built into the binary. While this simplifies the distribution and compatibility of Go programs, it causes even the smallest programs to be around 2mb

```go
var x []int // nil

s := make([]int, 5, 10)
s = append(s, 1, 2, 3)  // len now 8
clear(s) // [0 0 0 0 0 0 0 0]

slice.Equal(x, s) // false
```

It is important to understand when you should specify the length and capacity

- Slices used for buffers should have a non-zero length
- Slices with the exact size known should be indexed to get and set values
- Otherwise, use a zero length and non-zero capacity

Slices be sliced with `[start:stop]`. When you assign a subslice to a slice from another slice, the two slices are sharing the backing array's memory. The subslice will retain the same capacity as the original slice, or the original capacity minus the starting offset.

You should never append with a subslice. Use a _full slice expression_ `[start:stop:cap]` to make clear how much memory is shared. The `cap` in this expression denotes the capacity for the new slice. Note that each slice separately manages its length and capacity.

The following code creates a slice `x` and assigns a subslice to `y`. Since they share the underlying backing array, we should see the appended 5 in `x`, but we do not. The length of `x` did not change, so the 5 only exists the `x`'s reserved capacity. It's length prohibits it from accessing that memory; `x[4]` would cause a panic. If we were to now append something to `x`, the 5 in `y` would be overwritten. If either slice were to be appended to a length past their capacities, that slice would be resized and no longer share the underlying memory. It will retain the same values.

```go
x := make([]int, 0, 5)
x = append(x, 1, 2, 3, 4)

y := x[2:4:5]

y = append(y, 5)

fmt.Println(x, len(x), cap(x)) // [1 2 3 4] 4 5
fmt.Println(y, len(y), cap(y)) // [3 4 5] 3 4
```

Use `copy` to create slice independant of the original `copy(dest, src) int` where int is the number of elements copied over. `src` and `dest` can also be a subslices.

```go
x := []int{1,2,3,4}
copy(x[:3], x[1:]) // [2 3 4 4] 3
```

Arrays can be used as slices by slicing the entire array with `[:]`. This is useful for passing arrays to functions that take a slice. The same memory sharing properties are applied. Arrays can also be created with sliced using `arr := [3]int(slc)` which copies the values into new memory. Array sizes bigger than the slice will panic, `...` cannot be used to set the size. A slice can be converted into pointer to array with `arr = (*[4]int)(someSlice)`; the memory will be shared.

### Strings and Bytes and Runes

Under the hood, strings are a sequence of bytes - not runes. `byte` must be used to index a string `var ch byte = str[4]`. Like slices, strings can be sliced to other strings with the slice expression `[start:stop]`.

Be careful when handling strings due to their encoding. Code point in UTF-8 can be between 1 and 4 bytes long. This means that characters that are more than one byte long, such as emojis, cannot be accessed with a single index position.

> UTF-8 is the most commonly used encoding for Unicode. Unicode uses four bytes (32 bits) to represent each _code point_; the technical name for each character.

```go
var s string = "Hello 😀"
var s2 string = s[6:] // "😀"
var s3 string = s[4:7] // NOT "o 😀" but instead "o ?"
len(s) // 10, not 7 ("😀" is 4 bytes)
```

Strings can be converted to slices of `byte` and `rune`. For the string `"Hello 😀"`:

- `var bs []byte = []byte(s)` would have a len of 10
- `var rs []rune = []rune(s)` would have a len of 7, since an emoji can be contained in 1 rune

It is recommended to extract substrings and code points from strings using functions in the `strings` and `unicode/utf8` packages from the standard library.

It is possible to convert an integer to a string (eg. 65 would be "A") but `go vet` blocks type conversion to string from any integer type as of Go 1.15.

### Maps

`map[keyType]valueType`

Maps are key-value pair data structure implemented with an underlying hash map. They are not comparable, and will automatically grow as kv pairs are added. A map's keys must be any comparable type. Go 1.21 added `maps.Equal` and `maps.EqualFunc` to compare maps.

They have a zero value of `nil` with a length of 0. Reading a `nil` map will return the zero value, but writing to a `nil` map will cause a panic. Therefore, it is suggested to create a map with:

```go
x := map[string]int{}
x := make(map[string]int, size) // use size if it is known
```

Key-value pairs can be added by assigning a value to an index of the map where the index is the key `map[key] = value`. Pairs can be deleted with the `delete(map, key)` function. Deleting a key that does not exist will do nothing. `clear(map)` will empty a map and give it a length of zero.

When reading from a key that was never set, the value's zero value is returned. A value with the zero value can still have meaning, so it is best to use the _comma ok idiom_ to differentiate between keys that are associated with a zero value, and keys that are not in the map - both will return the same value.

```go
x := map[string]int{
  "M": 22,
}

x["E"] = 21 // add
v, ok := x["E"] // Gets "E" v with comma ok idiom
delete(x, "M") // delete "M" kv pair
```

While Go is limited in the built-in data structures it provides, a structure such as a set can be simulated with a map. Given a slice, we can create a map where the key is the value, and the value is a `bool` that represents if it exists in the set.

```go
intSet := map[int]bool{}
vals := []int{1,2,3,4}

for _, v := range vals {
  intSet[v] = true
}

// bool zero value is `false`
if intSet[100] {
  fmt.Println("100 is in the set")
}
```

When simulating sets, some people prefer an empty `struct` as value since it takes 0 bytes whereas a `bool` takes 1. However, it is clumsier and you have to use comma ok idiom.

### Struct

Structs are used for grouping related data together. They are comparable with `==` and `!=` as long as the field name, field types, field order, and struct type are same.

A struct's zero value will set all its fields to their zero value. Unlike maps, this is ok.

```go
// The following are the same
p := Person{}
var p Person
```

Struct literals can be written as comma separated list of values for the fields, or a map literal style

```go
x := person{
  "M",
  22
}
// or
x := person{ // dont need to provide val for every field with this style
  name: "M",
  age: 22
}
```

Anonymous structs not common. They are useful for translating external data into a struct and vice versa (json and protocol buffers), and testing.

A struct's type can converted if both type share the same field name, field types, and field order. They cannot be converted until converted to the same type; create a compare function. A typed struct and an anonymous struct can be compared and assigned to the same type as long as the same requirements above are met.

## 4 - Blocks, Shadows, and Control Structures

### Blocks and Shadowing

Go has 4 main block scopes; _block, package, file_, and _universe_.

**Shadowing**

Shadowing is when a variable with the same identifier within an inner block _shadows_ that same identifier from the outer block. It generally works like most other languages, but can get tricky when using `:=`.

```go
x := 10
if x > 5 {
  x, y := 5, 20
fmt.Println(x, y)
}
fmt.Println(x)
```

```
5 20
10
```

In the example above, we may seem like 5 would be printed at the end since `y` is the new variable on the lefthand side of `:=`, so `x` would be assigned with `=` under the hood. However, `:=` only reuses variables that are declared in the current block. When using `:=`, ensure you don't have variables from an outer scope on lefthand side unless you intend to shadow them.

The _universe block_ is interesting. It is a block that contains all other blocks. This means that

> Go only has 25 keywords. Those such as `true/false`, `int`, `string`, `make`, and a few others are not included. They are _predeclared identifiers_. This means that they can be shadowed. Never redefine identifiers in universe block.

### Control Structures

**`if`**

Go allows you to declare variables scoped to conditional blocks.

```go
if err := foo(); err != nil {
  // err legal here
} else {
  // and here
}
// not here
```

**`for`**

There are four ways to write a loop in Go.

1. C-style - `for i :=0; i < 10; i++`
2. condition-only - `for i < 10`
3. infinite - `for`
4. for-range - `for i, v := range list`

In `for-range` loops, the values returned are index and value for loops, and key and value for maps. You can ignore the first returned value with `_`, or just leave out the second one. `for-range` copies elements to the value variable.

`break` and `continue` are both included in the language.

**Iterating Over Maps**

The order in which key-value pairs are iterated is random. The map's hashing algorithm uses a random number every time a map is created. This is actually a security feature the prevents _Hash DoS_ attacks

> If maps always hash to the exact same values, and you know a server is storing user data in a map, you can slow down a server with a _Hash Dos_ attack by sending crafted data with keys all hash to the same bucket.

However, printing maps will sort them in ascending order.

**Iterating Over Strings**

- `for-range` loops iterate over _runes_, not the _bytes_.
- Converts multibyte rune UTF-8 representation into single 32-bit number and assigns it to value. Offset is incremented by number of bytes in the rune

**Labeling Your for Statements**

Labels allow `break` or `continue` statements to be applied to an outer loop from an inner loop. The label is indented at the same level as the braces for the block.

```go
func main() {
outer: // label
  for i := range 100 {
    for j := range 10 {
      if j % 2 == 0 {
        continue outer
      }
    }
  }
}
```

**Switch Statement**

The switch statements in Go is similar to most other languages. Go allows you to declare variables scoped to the `switch` block, just like `if`. Multiple cases that share the same functionality cannot be separated by commas. _Empty cases_ with no corresponding code is allowed, nothing will happen. A _blank switch_ which does not declare a variable we are comparing against allows the cases to check for a true equality instead.

Favour blank `switch` statements over `if/else` chains when you have multiple related cases.

```go
switch size := len(word); size {
case 1, 2: // multiple cases
  fmt.Println("Short word")
case 5:
  fmt.Println("Medium word")
case 6: // empty case
default:
  fmt.Println("Who knows")
}

// blank switch
switch{
case a == 1: // equality
  // ...
}
```

**goto**

> Goto is not recommended

`goto` allows you to jump to different parts of a program that are labeled. There are some rules, such as that you cannot jump over a variable declaration, or into an inner or parallel block.

```go
func main() {
  for condition {
    if cond {
      goto done
    }
  }
  fmt.Println("not found") // skip over if done
done:
  fmt.Println("found and done") // skip over if done
}
```

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

## 6 - Pointers

A pointer is a variable that contains the address of where the data is stored. All pointers are the same size at 4 bytes. Primitive types can have different size; and `int` is 4 bytes whereas a `boolean` is 1 byte. The smallest memory that can be independently addressed is a byte. Pointers have a zero value of `nil`.

- `&` _Address_ operator - precedes a value type and returns address of where value is stored. Cannot be used before primitive types because they don;t have memory addresses.
- `*` _Indirection_ operator - precedes a pointer type and returns the pointed-to value, _dereferencing_. Dereferencing a non-nil pointer will panic.

```go
x := 10
xp := &x

xp // memory address of x
*xp // value of memory address, 10
```

One problem that can occur is setting a struct pointer field a string literal for example - the literal has no address. Creating a helper function can help with assigning values to pointers.

```go
type Person struct {
  name *string
}

p := Person{
  name: "Em" // not valid
}

p := Person{
  name: makePointer("Em") // valid
}

func makePointer[T any](t T) *T {
  return &t
}
```

The following rules are true:

- If you pass a pointer to a function and modify a field, that change is refelcted in the original variable.
- If you reassign a pointer parameter, the change is not reflected in original variable. This creates a separate value and does not affect original variable.
- If you pass a `nil` pointer, you cannot modify the pointer to a non-nil value. Can only reassign if value exists

It is recommended to use values over pointers when possible. This makes it easier to understand how and when data is modified. If a pointer is being modified, it should be documented.

### Mutability

- Go's lack of immutabilty may seem like problematic, but being able to choose between value and pointer parameters solves this issue.
- Pointers indicate that a parameter is mutable, whereas a value parameters guarantee the passed in variables immutability
- Other languages such as python, java, and js do not have this feature

One implication is that you pass a `nil` pointer, you cannot make it non-nil. Since Go is call by value, you are essentially passing in a nil variable.

Another implication of copying a pointer is that if you want the value to still be there when you exit the function, you _must_ deference the pointer and set the value. If you change the pointer, you change the copy, not the original. Dereferencing puts the new value in the memory location pointed to by both the original and the copy.

```go
func failedUpdate(px *int) {
  a := 20
  px = &a
}

func update(px *int) {
  *px = 20
}

func main() {
  x := 10
  failedUpdate(&x)
  update(&x)
}
```

- Rather than populating a struct by passing a pointer to it into a function, have the function instantiate and return the struct
- Only time you should use pointer to struct is when function expects and interface

### Performance

Since pointers are always 4 bytes, the time to pass a pointer to a function is constant for all data sized, ~1 ns. Passing a value parameter takes longer as the data gets larger, ~0.7 ms for a 10mb struct.

The behaviour for returning a pointer vs. value from function is more interesting. For data structures smaller than 10mb, it is actually slower to return a pointer. It is ok to use pointers if you are passing megabytes of data, even if it is meant to be immutable.

> If the distinction between the zero value and no value is important, use a `nil` pointer to represent unassigned variables or struct fields

### Maps vs Slices

Maps are implemented as a pointer to a struct. This is why changes made to a map parameter affects the original variable. Avoid using maps when you can go with a struct instead. Maps are ideal if the keys for the data are not known at compile time.

Slice is not a pointer. A slice is implemented as a struct with three fields; two `int` and a pointer to a block of memory:

```go
type slice struct {
  array unsafe.Pointer
  len   int
  cap   int
}
```

Since there is a _pointer_ to a block of memory, its contents can be modifed when passed as a parameter, but its field cannot be modified since it is passed by value. The reason a slice of any size can be passed into a function is becuse they are all the same size - the size of that struct. Arrays are passed in entirely, so their size must be typed due to their variable size.

### Slices as Buffers

When reading from an external data source, it is common to iterate over the chunks of data being read. If we use a single variable to do this, the memory is constantly allocated and destroyed each iteration. Idomatic Go avoids unnecessary memory alloctions, which can be accomplished with buffers.

It would be more efficient create a slice of buffers (eg. 100 bytes) once and read the data into it. Each time the data is read in the loop the block of bytes (up to 100) is copied into the slice. It is allocated once and reused each loop.

```go
// error handling excluded
f := os.Open("File")
data := make([]byte, 100)
for {
  count := f.Read(data)
  process(data[:count])
}
```

### Reducing Garbage Collector's Workload

Garbage is data that has no more pointers pointing to it. It is the garbage collector's job to automatically detect this unused memory an recover it so it can be reused.

The **stack** is consecutive block of memory. Allocating memory to it is fast and simple. To store something on the stack, must know how big it is at compile time.

The **heap** is the memory managed by the garbage collector. In order for Go to allocate data a pointer points to page 137:

1. The data must be a local variable whose data size is known at compile time.
2. Pointer cannot be returned from the function.

#### What is so bad about storing data on the heap?

**1. Garbage collection takes time**

Many garbage-collection algorithms have been written. Some of them focus on **higher throughput**, finding the most garbage in a single scan. Some focus on **lower latency**, finish the garbage scan as quickly as possible. The paper _The Tail at Scale_ by Jeffrey Dean argues why systems should be optimized for latency.

**2. The fastest way to read data is sequentially**

A slice of struct would be stored sequentially. However, a slice of pointers to structs would be scattered across RAM which will cause _slower reading and processing_.

> _Mechanical Sympathy_: Approach of writing software that's aware of the hardware it's running on

Java and other languages that store classes as pointers have their memory scattered across RAM. Only the pointer is allocated to stack; the data within the object is allocated on the heap. The Java garbage collector has to do a great deal of more work bouncing thorugh memory which is inefficient. Although it may appear sequential, even the `List` interface is a pointer to a list of pointers. The JVM is optimized for both throughput and sometimes latency. Python, and Javascript garbage collectors are not optimized, and their performance suffers due to this.

This is why Go encourges to use pointers sparingly; reduce the garbage collector workload by making sure as much data is stored on the stack as possible. The key is to make less garbage in the first place.

### Tuning the Garbage Collector

The Go garbage collector can be tunes with environment variables.

`GOGC` used the formula uses formula `CURRENT_HEAP_SIZE + CURRENT_HEAP_SIZE*GOGC/100` to calculate the heap size that needs to be reached to trigger the next garbage-collection cycle. The default is set to `100`.

`GOMEMLIMIT` is the limit on the total memory the Go program can use. It is disabled by default. It is a _soft_ limit that can be exceeded in certain circumstances, such as when _trashing_ occurs.

**Thrashing**: Garbage-collection cycles are being rapidly triggered becuse program is repeatedly hitting the limit

### Excercise Notes

- **Inlining**: The compiler inlines small functions to reduce the overhead of function calls, which can improve performance.
- **Escape Analysis**: Determines whether variables can be safely allocated on the stack or must be moved to the heap. Variables that "escape" need to be heap-allocated to ensure they outlive the function call.

## 7 - Types, Methods, and Interfaces

- bulit-in and user-defined types
- methods and interfaces (an abstraction)
- can use primitive or compound type literal to define a concrete type
  - can convert between appropriate primitive and custom concrete types
  - these types serve as documentation (eg. Percentage vs int as a return value of a func)
  - When same underlying data has different operations to perform, use two different types

```go
type Score int
type TeamScores map[string]Score
```

### Methods

- methods. receivers - use short abbreviation of its type
- methods can only be defined at package level
- keep type definition and associated methods in same file
- pointer and value receivers
  - If method modifies receiver, must use pointer
  - If method needs to handle nil instances, must use pointer
  - If method does not modify receiver, must use value
- Able to call pointer receiver method with a value type, converts `c.Increment()` to `(&c).Increment`
- If you call value receiver on pointer variable, Go automatically dereferences it to `(*c).Increment()`
- Calling a value reciver on a pointer variable will compile but panic at runtime
- Go considers both value and pointer methods to be in same _method set_ for a pointer instance. For a value, only the value receivers are in the set.
- Reserve methods for business logic (not getters/setters)
- Properly handle `nil` instances in pointer receivers (eg. errors, or maybe you want to let it panic)
- mathods are like functions
  - can assign method to variable or pass it into a parameter of same function type - without receiver portion
  - can also create function from type itself where the first parameter is the receiver

```go
myAdder := Adder{start: 5}
f1 := myAdder.Add // can now call something like f1(5)

// Creating method from type
f2 := Adder.Add
f2(myAdder, 5) // first param is receiver
```

> Anytime some logic depends on values that are configured at startup or changed while your program is running, those values should be stored in a struct, and that logic shoyld be implemented as a method. If logic depends only on the input parameters, it should be a function.

### iota

- lets you assign an increasing value to a set of constants

1. define type based on int
2. use a `const` block to define a set of values for your type

```go
type Category int

const (
  Uncategorized Category = iota
  Spam
  Personal
  Silly
)
```

First constant has the type, and its value is set to `iota`. Every subsequent line has neither the type or the value. The Go compiler will see this and repeat the type and assignment to all subsequent constants in the block. The value of `iota` increases for each constant in `cosnt` block, starting at `0`. A new `const` block will reset `iota`.

```go
const (
  A = 0
  B = 1 + iota
  C = 20
  d = iota
)
fmt.Println(A, B, C, D) // 0 2 20 4
```

> Do not use `iota` for defining constants where it values are explicitly defined. Use `iota` for internal purposes only. This way you can insert new constants without the risk breaking everything.

### Use Embedding for Composition

- Favour object composition over inheritance
- All fields or methods declared on an embedded field are _promoted_ to the containing struct and can be invoked direct on it.

```go
type Employee struct {
  Name string
  ID   string
}

type Manager struct {
  Employee
  Reports []Employee
}

m := Manager{
  Employee: Employee{
    // name, id
  },
  Reports: []Employee{},
}

// promoted fields
m.Name
m.ID
m.Employee.ID // also valid
```

If acontaining struct has fields or methods with the same name as the embedded field, you need to use the embedded field's type to refer to obscured field or method.

- dynamic dispatch
- methods on an embedded field count towar the **method set** of the containing struct

### Quick Lesson on Interfaces

- Only abstract type in Go
- lists the methods that must be implemented by the concrete type to meet the interface
- usually named with "er" endings
- implemented _implicitly_, concrete class does not need to declare that it implements an interface; type safety and decoupling

> Program to an interface, not an implementation. Doing so allows you to code based on behaviour, allowing you to swap implementations as needed. This allows your coe to evolve over time, as requirement inevitably change

```go
// Actual logic
type LogicProvider struct{}

func (lp LogicProvider) process(data string) string {
  // business logic
}

// Interface
type Logic interface {
  process(data string) string
}

// Client code
type Client struct {
  L Logic
}

func (c Client) program() {
  // get some data from somewhere
  c.L.process(data)
}

// main
func main() {
  c := Client{
    L: LogicProvider{}
  }
  c.Program()
}
```

- can also embed an interface in an interface

```go
type Reader interface {
  read(p []byte) (int, error)
}

type Closer interface {
  Close() error
}

type ReadCloser interface {
  Reader,
  Closer,
}
```

**Accept Interfaces, Return Struct**

- Business lofic should be invoked via interfaces, but output should be a concrete type. If you add a nre method to an interface, you must update all existing implementations.
- Interface parameters invoke a heap allocation

**Interfaces and Nil**

- `nil` can represent the zero value for an interface instance, but there are some nuances.
  Interfaces in Go are implemented as a struct with two pointer fields: one for the type of the value and one for the value itself. For an interface to be considered `nil`, both the type and value fields must be `nil`. A `nil` interface indicates whether or not you can invoke methods on it.
- A `nil` interface will panic if its methods are invoked. However, a non-nil interface can have its methods invoked even if the value field of the concrete type is `nil`; the behavior in this case depends on how the value type handles `nil`.
- Must use Reflection to check if interface's value is `nil`

```go
var reader io.Reader
reader == nil // true
```

**Comparable**

- Interfaces are comparable with `==` to check if value and type are equal - comparing values that are not comparable will compile, but panic at runtime. Be careful with comparing interfaces - use Reflection
- `interface{}` / `any` satisfies all types, try to avoid it
- Common as palceholder for data of uncertain schema, such as JSON

**Type Assertions and Type Switches**

- Provides way to see if interface type has a specific concrete type
- _Type assertion_ names the concrete type that implemented the interface, or names another interface that is also implemented by concrete type of interface value. Incorrect assertion panics
- Type assertion must match type of value stored in the interface, even if they share the same underlying type. Use comma ok idiom avoid panics.

```go
var i any
var mine MyInt = 20
i = mine

i2 := i.(MyInt) // ok
i3 := i.(int)   // panic

i4, ok := i.(int)   // ok

if !ok {
  // ...handle failed conv
}
```

> **Type assertion** reveals the type of the value stored in an interface and is checked at runtime. **Type conversion** changes a value to a new type and is checked at compiile time.

- When an interface could be one of multiple possible types, use a type switch
- The type of the new variable depends on which case matches. If more than one type si listed on a case, the variable is `any`, otherwise it is the case's single type.

```go
function process(i any) {
  switch j := i.(type) {
    case nil:
      // j = nil
    case int:
      // j = int type
    case io.Reader:
      // j = io.Reader type
    case bool, rune:
      // j = any type (because more than one type is listed in this case)
    default:
      //
  }
}
```

**Functions Types are a Bridge to Interfaces**

- Go allows methods on _any_ user-defined type, including user-defined function types
- Most common usage is for HTTP handlers

```go
type Handler interface {
  ServeHTTP(http.ResponseWriter, *http.Request)
}

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(http.ResponseWriter, *http.Request) {
  f(w, r)
}
```

- Using a type conversion to `http.HandlerFunc` allows any function with signature `func(http.ResponseWriter, *http.Request)` to be used as an `http.Handler`

**Dependency Injection** is the concept that your code should explicitly specify the functionality it needs to perform its task.
