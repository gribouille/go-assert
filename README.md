# go-assert

Simple and no intrusive Go test library ([documentation](https://godoc.org/github.com/gribouille/go-assert)).

## Getting started

Install the library:

```txt
go get -u github.com/gribouille/go-assert
```

```txt
dep ensure -add github.com/gribouille/go-assert
```

Create a test file:

```go
package module_test

import (
  "testing"

  T "github.com/gribouille/go-assert"
)

func TestFunc(t *testing.T) {
  a := T.New(t)
  // ready...
}
```

Chain the assertions:

```go
a.
  Equal(3, 3, "Equal success").
  Equal("not", "equal", "Equal failed").
  Equal("a", "b").
  True(true, "True success").
  True(2 == 3, "True failed").
  False(2 == 3, "False success").
  False(4 == 4, "False failed").
  NotEqual(2, 3, "NotEqual success").
  NotEqual("a", "a", "NotEqual failed")
```

Use the BDD style:

```go
T.New(t).
  It("sub test 1", func(a *T.Assert) {
    // ...
  }).
  It("sub test 2", func(a *T.Assert) {
    // ...
  })
```

Add custom messages if the assertion failed:

```go
a.
  Equal("a", "b"). // no message
  Equal("a", "b", "a is different of b").
  Equal("a", "b", "%s is different of %s", "a", "b")
```

Many utils comparison functions:

```go
// Classic assertion
a.Equal(...)
a.True(...).
a.False(...)
a.NotEqual(...)
a.EqualDeep(...)
a.NotEqualDeep(...)
a.EqualSlice(...)
a.EqualStringSlice(...)
a.Nil(...)
a.NotNil(...)
a.Error(err, "open %s: no such file or directory", file)

// Advanced matching
a.Match(`^[a-z]+\[[0-9]+\]$`, "adam[23]")
a.EqualFile(got, "testdata/lorem.txt")

data := struct {
  Version string
  Authors []string
}{
  Version: "2.0",
  Authors: []string{"Bob Leponge", "Mr. Krabs"},
}
a.EqualTemplate(got, "testdata/version.tpl", data)

// Filesystem
a.IsFile("/etc/passwd")
a.IsDir("/usr/bin/")
a.NotExists("/file/not/exist")
```

Create temporary testing environments to execute your tests:

```go
T.New(t).
  ItEnv("Sub test 1 with temp dir",
    T.Copy{"testdata/a", "a"},
    T.Copy{"testdata/version.tpl", "version.tpl"},
    T.Copy{"testdata/ipsum.txt", "a/ipsum.txt"},
  )(func(a *T.Assert, dir string) {
  // here you have a temporary directory with the test data copied inside
  // this directory will be deleted at the end of this function
})
```

Test crashable function:

```go
T.New(t).Crash("Sub test crash", func() {
  // function that uses os.Exit(...)
}, func(a *T.Assert, code int, stdout, stderr string) {
  // here your test
})
```

Capture the standard output and error:

```go
T.New(t).Capture("Sub test capture", func() {
  // function that prints messages on stdout and stderr
}, func(a *T.Assert, stdout, stderr string) {
  // here your test
})
```

Customize the behavior with environment variables or with custom constructor:

```go
// fatal = true: use t.Fatal for the error
// stack = true: show the stack trace
a := T.NewCustom(t, true, true)
```

See more details in the [examples](./examples) and in the [documentation](https://godoc.org/github.com/gribouille/go-assert).

## Tests

To execute the tests:

```txt
go test -v .
```

## TODO

- [x] Improve the documentation
- [ ] Add net functions

## References

- [testing](https://golang.org/pkg/testing/)

## License

This project is licensed under [Mozilla Public License Version 2.0](./LICENSE).
