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
