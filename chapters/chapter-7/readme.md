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
