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
