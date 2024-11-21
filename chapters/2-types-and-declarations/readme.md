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
