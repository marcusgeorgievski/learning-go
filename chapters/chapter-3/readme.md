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
