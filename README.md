# gonv - Go Type Conversion Utilities

Package `gonv` provides comprehensive type conversion utilities for Go applications. It offers safe and flexible casting between different data types with generic support, making it easy to convert values between various types without worrying about runtime panics.

## Features

- Generic type conversion functions for all basic Go types
- Safe conversions with error handling options
- Support for slices, maps, and complex types
- Protobuf and database driver.Valuer support
- Zero-allocation string/bytes conversion (Go 1.20 and below)

## Installation

```bash
go get github.com/go-leo/gonv
```

## Basic Usage

```go
import "github.com/go-leo/gonv"

// Convert string to int
result := gonv.Int[int]("42") // returns 42

// Convert with error handling
result, err := gonv.IntE[int]("invalid") // returns 0, error

// Convert slice of strings to slice of ints
nums := gonv.IntS[[]int]([]string{"1", "2", "3"}) // returns []int{1, 2, 3}
```

## Supported Types

- Basic types: string, bool, int, uint, float, time.Duration, time.Time
- Slice types for all basic types
- Map types with string keys and various value types
- Protobuf wrapper types
- Database driver.Valuer implementations

## Error Handling

Functions ending with 'E' (e.g., [IntE](file:///Users/soyacen/Workspace/github.com/go-leo/gonv/int.go#L23-L25), [StringE](file:///Users/soyacen/Workspace/github.com/go-leo/gonv/string.go#L24-L26)) return both the converted value and an error, allowing for explicit error handling. Functions without 'E' (e.g., [Int](file:///Users/soyacen/Workspace/github.com/go-leo/gonv/int.go#L17-L20), [String](file:///Users/soyacen/Workspace/github.com/go-leo/gonv/string.go#L18-L21)) ignore errors and return the zero value of the target type when conversion fails.

## Examples

### String conversions

```go
str := gonv.String[string](42)           // "42"
str := gonv.String[string](true)         // "true"
str := gonv.String[string](3.14)         // "3.14"
```

### Boolean conversions

```go
b := gonv.Bool[bool]("true")             // true
b := gonv.Bool[bool]("false")            // false
b := gonv.Bool[bool](1)                  // true
b := gonv.Bool[bool](0)                  // false
```

### Numeric conversions

```go
i := gonv.Int[int64]("42")               // 42
f := gonv.Float[float64]("3.14")         // 3.14
u := gonv.Uint[uint64]("-1")             // error (negative not allowed)
```

### Time conversions

```go
t := gonv.Time("2023-01-01T12:00:00Z")   // time.Time
d := gonv.Duration("1h30m")              // 1 hour 30 minutes
```

### Slice conversions

```go
strs := gonv.StringS[[]string]([]int{1, 2, 3}) // []string{"1", "2", "3"}
bools := gonv.BoolS[[]bool]([]string{"true", "false"}) // []bool{true, false}
```

### Map conversions

```go
m := gonv.StringIntMap[string, int](map[string]string{"key": "42"}) // map[string]int{"key": 42}
```

## Safety

All conversions are safe and will not panic. When a conversion is not possible, functions either return the zero value of the target type or an error, depending on whether you use the error-handling variant.

## Performance

The package uses optimized conversion paths for common type combinations and falls back to reflection only when necessary. String/bytes conversions use unsafe operations for zero-allocation performance on Go 1.20 and earlier versions.

## License

MIT