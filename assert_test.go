package assert

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func isAssert(t *testing.T, a interface{}) {
	to := reflect.TypeOf(a)
	if to.String() != "*assert.Assert" {
		t.Errorf("is not an Assert: %#v => %s", a, to.String())
	}
}

func TestAssert(t *testing.T) {
	a := New(t)
	if a == nil {
		t.Errorf("New failed")
	}
	isAssert(t, a)
}

func TestAssertNewCustom(t *testing.T) {
	b := NewCustom(t, true, true)
	if b == nil {
		t.Errorf("NewCustom failed")
	}
	isAssert(t, b)
}

func TestAssertSetStack(t *testing.T) {
	isAssert(t, New(t).SetStack(true))
}
func TestAssertSetFatal(t *testing.T) {
	isAssert(t, New(t).SetFatal(true))
}

func TestAssertSetOS(t *testing.T) {
	isAssert(t, New(t).SetOS("linux"))
}

func TestAssertSkip(t *testing.T) {
	var a *Assert
	t.Run("skip", func(t *testing.T) {
		a = New(t).Skip("skipped")
	})
	isAssert(t, a)
}

func TestAssertIt(t *testing.T) {
	var a *Assert
	isAssert(t, New(t).It("sub test", func(x *Assert) { a = x }))
	isAssert(t, a)
}

func TestAssertEqual(t *testing.T) {
	isAssert(t, New(t).Equal("a", "a", "format str"))
}

func TestAssertFAbs(t *testing.T) {
	isAssert(t, New(t).EqualFAbs(1.1, 1.1, 0.2, "format str"))
}

func TestAssertFRel(t *testing.T) {
	isAssert(t, New(t).EqualFRel(100, 102, 5, "format str"))
}

func TestAssertTrue(t *testing.T) {
	isAssert(t, New(t).True(true))
}

func TestAssertFalse(t *testing.T) {
	isAssert(t, New(t).False(false))
}

func TestAssertNotEqual(t *testing.T) {
	isAssert(t, New(t).NotEqual("a", "b"))
}

func TestAssertEqualDeep(t *testing.T) {
	isAssert(t, New(t).EqualDeep(struct{ X, Y int }{1, 2}, struct{ X, Y int }{1, 2}))
}

func TestAssertNotEqualDeep(t *testing.T) {
	isAssert(t, New(t).NotEqualDeep(struct{ X, Y int }{1, 2}, struct{ X, Y int }{1, 3}))
}

func TestAssertEqualSlice(t *testing.T) {
	isAssert(t, New(t).EqualSlice([]interface{}{1, 2, 3}, []interface{}{1, 2, 3}))
}

func TestAssertEqualStringSlice(t *testing.T) {
	isAssert(t, New(t).EqualStringSlice([]string{"a", "b", "c"}, []string{"a", "b", "c"}))
}

func TestAssertNil(t *testing.T) {
	isAssert(t, New(t).Nil(nil))
}

func TestAssertNotNil(t *testing.T) {
	isAssert(t, New(t).NotNil(fmt.Errorf("error")))
}

func TestAssertError(t *testing.T) {
	isAssert(t, New(t).Error(fmt.Errorf("my message 33"), "my message %d", 33))
}

func TestAssertMatch(t *testing.T) {
	isAssert(t, New(t).Match(`^[a-z]+\[[0-9]+\]$`, "adam[23]"))
}

func TestAssertEqualFile(t *testing.T) {
	isAssert(t, New(t).EqualFile("Nulla facilisi.", "examples/testdata/ipsum.txt"))
}

func TestAssertEqualTemplate(t *testing.T) {
	isAssert(t, New(t).EqualTemplate("aa - bb", "examples/testdata/temp.tpl", struct{ A, B string }{"aa", "bb"}))
}

func TestAssertIsFile(t *testing.T) {
	isAssert(t, New(t).IsFile("examples/testdata/temp.tpl"))
}

func TestAssertIsDir(t *testing.T) {
	isAssert(t, New(t).IsDir("examples/testdata"))
}

func TestAssertNotExists(t *testing.T) {
	isAssert(t, New(t).NotExists("examples/testdata/blabla"))
}

func TestAssertItTmp(t *testing.T) {
	isAssert(t, New(t).ItTmp("tmp", func(a *Assert, dir string) {
		isAssert(t, a)
		stat, err := os.Stat(dir)
		if err != nil {
			t.Error(err)
		}
		if !stat.IsDir() {
			t.Errorf("%s is not a directory", dir)
		}
	}))
}

func TestAssertItEnv(t *testing.T) {
	isAssert(t, New(t).ItEnv("env", Copy{"examples/testdata", "mydir"})(func(a *Assert, dir string) {
		isAssert(t, a)
		stat, err := os.Stat(dir)
		if err != nil {
			t.Error(err)
		}
		if !stat.IsDir() {
			t.Errorf("%s is not a directory", dir)
		}
		if _, err := os.Stat(filepath.Join(dir, "mydir", "lorem.txt")); err != nil {
			t.Error(err)
		}
		if _, err := os.Stat(filepath.Join(dir, "mydir", "ipsum.txt")); err != nil {
			t.Error(err)
		}
	}))
}

func TestAssertCapture(t *testing.T) {
	isAssert(t, New(t).Capture("caputre", func() {
		fmt.Println("Hello")
		fmt.Fprintf(os.Stderr, "World")
	}, func(a *Assert, stdout, stderr string) {
		isAssert(t, a)
		if stdout != "Hello\n" {
			t.Errorf("got: %s, exp: Hello\\n", stdout)
		}
		if stderr != "World" {
			t.Errorf("got: %s, exp: World", stderr)
		}
	}))
}

func TestAssertCrash(t *testing.T) {
	isAssert(t, New(t).Crash("crash", func() {
		fmt.Println("Hello")
		fmt.Fprintf(os.Stderr, "World")
		os.Exit(329)
	}, func(a *Assert, status int, stdout, stderr string) {
		isAssert(t, a)
		exp := 329 & 0377
		if stdout != "Hello\n" {
			t.Errorf("got: %s, exp: Hello\\n", stdout)
		}
		if stderr != "World" {
			t.Errorf("got: %s, exp: World", stderr)
		}
		if status != exp {
			t.Errorf("got: %d, exp: %d", status, exp)
		}
	}))
}
