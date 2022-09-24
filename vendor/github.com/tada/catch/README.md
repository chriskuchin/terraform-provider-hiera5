Catch provides controlled recovery of panics using a special error implementation that has a cause (another error).

[![](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![](https://goreportcard.com/badge/github.com/tada/catch)](https://goreportcard.com/report/github.com/tada/catch)
[![](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/tada/catch)
[![](https://github.com/tada/catch/workflows/Catch%20Test/badge.svg)](https://github.com/tada/catch/actions)
[![](https://coveralls.io/repos/github/tada/catch/badge.svg?service=github)](https://coveralls.io/github/tada/catch)

### How to get:
```sh
go get github.com/tada/catch
```
### Sample usage

Consider the following code where more than half of the lines are devoted to error handling:
```go
func foo() (error, int) {
  a, err := x()
  if err != nil {
     return err
  }
  b, err := y()
  if err != nil {
     return err
  }
  c, err := z()
  if err != nil {
     return err
  }
  return a + b * c
}
```
using `catch` on all involved functions, the code can instead very distinct:
```go
func foo() int {
  return x() + y() * z()
}
```
Leaf functions such as x, y, and z, where errors are produced may look something like this without `catch`:
```go
func x() (error, int) {
  err, v := computeSomeValue()
  if err != nil {
  	return 0, err
  }
  return int(v)
}
```
and like this with `catch`:
```go
func x() int {
  err, v := computeSomeValue()
  if err != nil {
    panic(catch.Error(err))
  }
  return int(v)
}
```
At the very top, errors produced somewhere in the executed code can be recovered using the
`catch.Do` function which will recover only `catch.Error` and return its cause:
```go
func furtherUp() error {
  return catch.Do(func() {
    x := foo()
    ...
  })
}
```
