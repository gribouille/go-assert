package examples

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	T "github.com/gribouille/go-assert"
)

var (
	e1 = []string{"a", "b"}
	e2 = []string{"a", "b"}
	e3 = []string{"a", "c"}
	e4 = []string{"a", "b", "c"}
	e5 = []int{1, 2, 3}
	e6 = []int{2, 1, 3}
	e7 = []interface{}{"aaa", "bbb", "ccc", "ddd"}
	e8 = []interface{}{"aaa", "bbc", "ccc", "ddd"}
)

func TestSubTest(t *testing.T) {
	T.New(t).
		It("sub test 1", func(a *T.Assert) {
			// ...
		}).
		It("sub test 2", func(a *T.Assert) {
			// ...
		})
}

func TestEqual(t *testing.T) {
	a := T.New(t)
	a.Equal(3, 3, "Equal success").Equal("not", "equal", "Equal failed")
	a.Equal("a", "b") // message is optional
	a.True(true, "True success").True(2 == 3, "True failed")
	a.False(2 == 3, "False success").False(4 == 4, "False failed")
	a.NotEqual(2, 3, "NotEqual success").NotEqual("a", "a", "NotEqual failed")
	a.EqualDeep(e1, e2, "EqualDeep success").EqualDeep(e1, e3, "EqualDeep failed")
	a.NotEqualDeep(e1, e3).NotEqualDeep(e1, e2, "NotEqualDeep failed")
}

func TestSlice(t *testing.T) {
	T.New(t).
		EqualSlice(e7, e8).
		EqualStringSlice(e3, e4).
		EqualStringSlice(e3, e2)
}

func TestError(t *testing.T) {
	a := T.New(t)
	a.Nil(nil).Nil(3, "Nil failed")
	a.NotNil("ok").NotNil(nil, "NotNil failed")
	file := "/file/not/exist"
	_, err := ioutil.ReadFile(file)
	a.Error(err, "open %s: no such file or directory", file)
}

func TestMatch(t *testing.T) {
	a := T.New(t)
	re := `^[a-z]+\[[0-9]+\]$`
	a.Match(re, "adam[23]").
		Match(re, "Job[48]", "Match failed")
}

func TestEqualFile(t *testing.T) {
	got := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam sed tortor eget erat venenatis accumsan. Nunc interdum orci urna, quis aliquam leo facilisis interdum.`
	a := T.New(t)
	a.EqualFile(got, "testdata/lorem.txt").
		EqualFile("Nulla facilisy.", "testdata/ipsum.txt", "EqualFile failed")
}

func TestEqualTemplate(t *testing.T) {
	got := `Version: 2.0
Authors:
  - Bob Leponge
  - Mr. Krabs	
`
	data := struct {
		Version string
		Authors []string
	}{
		Version: "2.0",
		Authors: []string{"Bob Leponge", "Mr. Krabs"},
	}
	a := T.New(t)
	a.EqualTemplate(got, "testdata/version.tpl", data)
	data.Version = "3.0"
	a.EqualTemplate(got, "testdata/version.tpl", data, "EqualTemplate failed")
}

func TestIsFile(t *testing.T) {
	T.New(t).SetOS("linux").
		IsFile("/etc/passwd").
		IsFile("/etc", "NotExists IsFile: it is a folder").
		IsFile("/file/not/exist", "NotExists IsFile: does not exist")
}

func TestIsDir(t *testing.T) {
	T.New(t).SetOS("linux").
		IsDir("/usr/bin/").
		IsDir("/etc/passwd", "NotExists IsDir: it is a file").
		IsDir("/dir/not/exist", "NotExists IsDir: does not exist")
}

func TestNotExists(t *testing.T) {
	T.New(t).SetOS("linux").
		NotExists("/file/not/exist").
		NotExists("/dir/not/exists/").
		NotExists("/etc/passwd", "NotExists failed: the file exists").
		NotExists("/usr/bin/", "NotExists failed: the folder exists")
}

func TestItTmp(t *testing.T) {
	tmp := ""
	a := T.New(t)
	a.
		ItTmp("Sub test 1 with temp dir", func(a *T.Assert, dir string) {
			a.IsDir(dir)
			tmp = dir
		}).
		NotExists(tmp, "directory has been cleaned")

	if err := os.Setenv("GO_ASSERT_TMP_DISABLE", "1"); err != nil {
		panic(err)
	}

	a.
		ItTmp("Sub test 2 with temp dir", func(a *T.Assert, dir string) {
			a.IsDir(dir)
			tmp = dir
		}).
		IsDir(tmp, "clean up disable")

	if err := os.RemoveAll(tmp); err != nil {
		panic(err)
	}

	a.NotExists(tmp)

	if err := os.Setenv("GO_ASSERT_TMP_DISABLE", ""); err != nil {
		panic(err)
	}
}

func TestItEnv(t *testing.T) {
	tmp := ""
	T.New(t).
		ItEnv("Sub test 1 with temp dir",
			T.Copy{"testdata/a", "a"},
			T.Copy{"testdata/version.tpl", "version.tpl"},
			T.Copy{"testdata/ipsum.txt", "a/ipsum.txt"},
		)(func(a *T.Assert, dir string) {
		a.IsDir(dir).
			IsFile(filepath.Join(dir, "a/b/c.txt")).
			IsFile(filepath.Join(dir, "version.tpl")).
			IsFile(filepath.Join(dir, "a/ipsum.txt"))
		tmp = dir
	}).
		NotExists(tmp, "directory has been cleaned").
		NotExists(filepath.Join(tmp, "a/b/c.txt")).
		NotExists(filepath.Join(tmp, "version.tpl")).
		NotExists(filepath.Join(tmp, "a/ipsum.txt"))
}

func TestCapture(t *testing.T) {
	T.New(t).Capture("Sub test capture", func() {
		fmt.Fprintf(os.Stdout, "Hello")
		fmt.Fprintf(os.Stderr, "World")
	}, func(a *T.Assert, stdout, stderr string) {
		a.Equal("Hello", stdout).
			Equal("World", stderr)
	})
}

func TestCrash(t *testing.T) {
	T.New(t).Crash("Sub test crash", func() {
		fmt.Fprintf(os.Stdout, "Ping")
		fmt.Fprintf(os.Stderr, "Pong")
		os.Exit(33)
	}, func(a *T.Assert, code int, stdout, stderr string) {
		a.Equal("Hello", stdout, "No it is Ping").
			Equal("World", stderr, "No it is Pong").
			Equal("Ping", stdout).
			Equal("Pong", stderr).
			Equal(code, 33)
	})
}
