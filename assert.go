package assert

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"
	"text/template"
)

// Assert wraps the standard testing.T structure.
type Assert struct {
	t     *testing.T
	stack bool
	os    string
	as    func(format string, args ...interface{})
}

func (a *Assert) clone(t *testing.T) *Assert {
	return &Assert{t, a.stack, a.os, a.as}
}

// New creates a new Assert object.
//
// Uses the environments variables to customize the behavior:
// 	- GO_ASSERT_STACK: show the stacktrace if the test fails
// 	- GO_ASSERT_FATAL: uses fatal errors
// 	- GO_ASSERT_OS: test for a specific OS (the possible values are similar to runtime.GOOS)
func New(t *testing.T) *Assert {
	s := os.Getenv("GO_ASSERT_STACK")
	stack := false
	if s == "true" || s == "t" || s == "1" {
		stack = true
	}
	fn := t.Errorf
	f := os.Getenv("GO_ASSERT_FATAL")
	if f == "true" || f == "t" || f == "1" {
		fn = t.Fatalf
	}
	o := os.Getenv("GO_ASSERT_OS")
	if o != "" {

	}
	return &Assert{t, stack, "all", fn}
}

// NewCustom is similar to New but not uses the environment variables.
//
// If stack is true then the stacktrace is showed; if fatal is true then uses
// the fatal errors.
func NewCustom(t *testing.T, fatal, stack bool) *Assert {
	fn := t.Errorf
	if fatal {
		fn = t.Fatalf
	}
	return &Assert{t, stack, "all", fn}
}

// assert wraps the other methods. It should not used directly.
func (a *Assert) assert(fn func()) *Assert {
	if a.os != "all" {
		if runtime.GOOS != a.os {
			return a
		}
	}
	fn()
	return a
}

// errorMessage is an helper to show homogeneous messages.
func (a *Assert) errorMessage(f string, args ...interface{}) func(...interface{}) *Assert {
	return func(msg ...interface{}) *Assert {
		if a.stack {
			debug.PrintStack()
		}
		m := errorMessage(fmt.Sprintf(f, args...), msg...)
		a.as(m)
		return a
	}
}

// SetStack sets to true if the stacktrace must be showed after an error.
func (a *Assert) SetStack(v bool) *Assert {
	a.stack = v
	return a
}

// SetFatal sets to true if the assert must use the Fatal method.
func (a *Assert) SetFatal(v bool) *Assert {
	fn := a.t.Errorf
	if v {
		fn = a.t.Fatalf
	}
	a.as = fn
	return a
}

// SetOS specifies if the test is operating system dependant.
//
// The possible values are similar to the values of runtime.GOOS.
func (a *Assert) SetOS(v string) *Assert {
	a.os = v
	return a
}

// Skip the test.
func (a *Assert) Skip(args ...interface{}) *Assert {
	a.t.Skip(args...)
	return a
}

// It defines a new subtest.
func (a *Assert) It(msg string, fn func(*Assert)) *Assert {
	a.t.Run(msg, func(t *testing.T) {
		fn(a.clone(t))
	})
	return a
}

// Equal assertion.
//
// The assertion function can add an optional custom error message:
// 	a.Equal("a", "b"). // no message
// 	a.Equal("a", "b", "a is different of b").
// 	a.Equal("a", "b", "%s is different of %s", "a", "b")
//
// The assertion function can be chained:
// 	a.Nil(...).Equal(...).True(...)
func (a *Assert) Equal(exp, got interface{}, msg ...interface{}) *Assert {
	return a.assert(func() {
		if exp != got {
			a.errorMessage("Exp: %+v\nGot: %+v\n", exp, got)(msg...)
		}
	})
}

// EqualFAbs compares 2 floats numbers with an absolute tolerance.
func (a *Assert) EqualFAbs(exp, got, epsilon float64, msg ...interface{}) *Assert {
	return a.assert(func() {
		if !compareAbs(exp, got, epsilon) {
			a.errorMessage("Exp: %f\nGot: %f with an absolute tolerance: %f\n", exp, got, epsilon)(msg...)
		}
	})
}

// EqualFRel compares 2 floats numbers with an relative tolerance.
func (a *Assert) EqualFRel(exp, got, epsilon float64, msg ...interface{}) *Assert {
	return a.assert(func() {
		if !compareRel(exp, got, epsilon) {
			a.errorMessage("Exp: %f\nGot: %f with an relative tolerance: %f\n", exp, got, epsilon)(msg...)
		}
	})
}

// True assertion.
func (a *Assert) True(v bool, msg ...interface{}) *Assert {
	return a.assert(func() {
		if !v {
			a.errorMessage("Not true")(msg...)
		}
	})
}

// False assertion.
func (a *Assert) False(v bool, msg ...interface{}) *Assert {
	return a.assert(func() {
		if v {
			a.errorMessage("Not false")(msg...)
		}
	})
}

// NotEqual assertion.
func (a *Assert) NotEqual(exp, got interface{}, msg ...interface{}) *Assert {
	return a.assert(func() {
		if exp == got {
			a.errorMessage("Equal: %+v\n", exp)(msg...)
		}
	})
}

// EqualDeep is similar to Equal but uses reflect.DeepEqual to test the equality.
func (a *Assert) EqualDeep(exp, got interface{}, msg ...interface{}) *Assert {
	return a.assert(func() {
		if !reflect.DeepEqual(exp, got) {
			a.errorMessage("Exp: %#v\nGot: %#v\n", exp, got)(msg...)
		}
	})
}

// NotEqualDeep is the inverse of EqualDeep.
func (a *Assert) NotEqualDeep(exp, got interface{}, msg ...interface{}) *Assert {
	return a.assert(func() {
		if reflect.DeepEqual(exp, got) {
			a.errorMessage("Equal: %#v\n", exp)(msg...)
		}
	})
}

// EqualSlice compares two slices of interface{}.
func (a *Assert) EqualSlice(exp, got []interface{}, msg ...interface{}) *Assert {
	return a.assert(func() {
		if len(exp) != len(got) {
			a.errorMessage("Expected size: %d, got size: %d", len(exp), len(got))(msg...)
			return
		}
		dash := 4
		for i := range exp {
			s := fmt.Sprintf("%+v", exp[i])
			dash += len(s)
			if exp[i] != got[i] {
				d := strings.Repeat("━", dash)
				a.errorMessage("Exp: %+v\nGot: %+v\n ┗%s┛\n", exp, got, d)(msg...)
				return
			}
		}
	})
}

// EqualStringSlice compares tow slices of strings.
//
// If the assertion fails then a message shows the first differences:
// 	Error:
// 	  Exp: [aaa bbb ccc ddd]
// 	  Got: [aaa bbc ccc ddd]
// 	   ┗━━━━━━━━━┛
func (a *Assert) EqualStringSlice(exp, got []string, msg ...interface{}) *Assert {
	return a.assert(func() {
		if len(exp) != len(got) {
			a.errorMessage("Expected size: %d, got size: %d", len(exp), len(got))(msg...)
			return
		}
		dash := 4
		for i := range exp {
			s := fmt.Sprintf("%+v", exp[i])
			dash += len(s)
			if exp[i] != got[i] {
				d := strings.Repeat("━", dash)
				a.errorMessage("Exp: %s\nGot: %s\n ┗%s┛\n", exp, got, d)(msg...)
				return
			}
		}
	})
}

// Nil assertion.
func (a *Assert) Nil(v interface{}, msg ...interface{}) *Assert {
	return a.assert(func() {
		if v == nil {
			return
		}
		defer func() {
			if err := recover(); err != nil {
				a.errorMessage("%#v is not a nullable value", v)()
			}
		}()
		if !reflect.ValueOf(v).IsNil() {
			a.errorMessage("Nil expected: %#v", v)(msg...)
		}
	})
}

// NotNil assertion.
func (a *Assert) NotNil(v interface{}, msg ...interface{}) *Assert {
	return a.assert(func() {
		if v != nil {
			defer func() {
				if err := recover(); err != nil {
					a.errorMessage("%#v is not a nullable value", v)()
				}
			}()
			if reflect.ValueOf(v).IsNil() {
				a.errorMessage("Not nil expected")(msg...)
			}
		}
	})
}

// Error message assertion.
//
// Example:
// 	file := "/file/not/exist"
//	_, err := ioutil.ReadFile(file)
// 	a.Error(err, "open %s: no such file or directory", file)
func (a *Assert) Error(err error, errFormat string, errA ...interface{}) *Assert {
	return a.assert(func() {
		exp := fmt.Sprintf(errFormat, errA...)
		if err == nil {
			a.errorMessage("Expected error with message: %s", exp)()
			return
		}
		if err.Error() != exp {
			a.errorMessage("Error message mismatch\nExp: %s\nGot: %s\n", exp, err.Error())()
		}
	})
}

// Match a string to a regex expression.
//
// Example:
// 	a.Match(`^[a-z]+\[[0-9]+\]$`, "adam[23]")
func (a *Assert) Match(pattern, s string, msg ...interface{}) *Assert {
	return a.assert(func() {
		m, err := regexp.MatchString(pattern, s)
		if err != nil {
			panic(err)
		}
		if !m {
			a.errorMessage("Regex (%s) mismatch: %s", pattern, s)(msg...)
		}
	})
}

// EqualFile tests the equality between the content of file and the got string.
func (a *Assert) EqualFile(got string, filename string, msg ...interface{}) *Assert {
	return a.assert(func() {
		fi, err := ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		a.Equal(string(fi), got, msg...)
	})
}

// EqualTemplate is similar to EqualFile but the file can be a template. The spaces are trimed.
// Example:
// 	data := struct {
// 		Version string
// 		Authors []string
// 	}{
// 		Version: "2.0",
// 		Authors: []string{"Bob Leponge", "Mr. Krabs"},
// 	}
// 	a.EqualTemplate(got, "testdata/version.tpl", data)
func (a *Assert) EqualTemplate(got string, filename string, data interface{}, msg ...interface{}) *Assert {
	return a.assert(func() {
		tpl, err := template.ParseFiles(filename)
		if err != nil {
			fmt.Println("ok")
			panic(err)
		}
		var buf strings.Builder
		if err := tpl.Execute(&buf, data); err != nil {
			panic(err)
		}
		a.Equal(strings.TrimSpace(buf.String()), strings.TrimSpace(got), msg...)
	})
}

// IsFile tests if the file exists.
func (a *Assert) IsFile(pth string, msg ...interface{}) *Assert {
	return a.assert(func() {
		if !isFile(pth) {
			a.errorMessage("Not a file: %s", pth)(msg...)
		}
	})
}

// IsDir tests if the directory exists.
func (a *Assert) IsDir(pth string, msg ...interface{}) *Assert {
	return a.assert(func() {
		if !isDir(pth) {
			a.errorMessage("Not a directory: %s", pth)(msg...)
		}
	})
}

// NotExists tests if the path does not exist.
func (a *Assert) NotExists(pth string, msg ...interface{}) *Assert {
	return a.assert(func() {
		if exists(pth) {
			a.errorMessage("Expected not exists: %s", pth)(msg...)
		}
	})
}

// ItTmp creates a sub test with a temporary directory. This directory is clean
// up after the test.
//
// For debugging, the deletion of temporary folder can be disable with the
// environment variable GO_ASSERT_TMP_DISABLE=1.
//
// Example:
// 	a.ItTmp("Sub test 1 with temp dir", func(a *T.Assert, dir string) {
// 		// dir is the temporary directory
// 	})
func (a *Assert) ItTmp(msg string, fn func(*Assert, string)) *Assert {
	a.t.Run(msg, func(t *testing.T) {
		tmpDir(func(dir string) {
			fn(a.clone(t), dir)
		})
	})
	return a
}

// Copy is a file or a directory copied in the temporary directory.
type Copy struct {
	Source, Dest string
}

// ItEnv is similar to ItTmp but it copies the files or the folders in the
// temporary directory.
//
// Example:
// 	T.New(t).ItEnv("Sub test 1 with temp dir",
// 		T.Copy{"testdata/a", "a"},
// 		T.Copy{"testdata/version.tpl", "version.tpl"},
// 		T.Copy{"testdata/ipsum.txt", "a/ipsum.txt"},
// 	)(func(a *T.Assert, dir string) {
// 	// ...
// 	})
func (a *Assert) ItEnv(msg string, copies ...Copy) func(func(*Assert, string)) *Assert {
	return func(fn func(*Assert, string)) *Assert {
		a.t.Run(msg, func(t *testing.T) {
			tmpDir(func(dir string) {
				for _, c := range copies {
					if err := Cp(c.Source, filepath.Join(dir, c.Dest)); err != nil {
						panic(err)
					}
				}
				fn(a.clone(t), dir)
			})
		})
		return a
	}
}

// Capture the messages sends to the standard output and the standard error of the function.
//
// Example:
// 	T.New(t).Capture("Sub test capture", func() {
// 		fmt.Fprintf(os.Stdout, "Hello")
// 		fmt.Fprintf(os.Stderr, "World")
// 	}, func(a *T.Assert, stdout, stderr string) {
// 		a.Equal("Hello", stdout).Equal("World", stderr)
// 	})
func (a *Assert) Capture(msg string, act func(), fn func(*Assert, string, string)) *Assert {
	a.t.Run(msg, func(t *testing.T) {
		stdOut, stdErr := captureOutput(act)
		fn(a.clone(t), stdOut, stdErr)
	})
	return a
}

// Crash is similar to Capture but the function should exit the program too.
// The assertion captures the return code too.
func (a *Assert) Crash(msg string, act func(), fn func(*Assert, int, string, string)) *Assert {
	a.t.Run(msg, func(t *testing.T) {
		rc, stdOut, stdErr := crashTest(t, act)
		fn(a.clone(t), rc, stdOut, stdErr)
	})
	return a
}
