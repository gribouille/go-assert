package assert

import "testing"

func TestErrorMessage(t *testing.T) {
	a := errorMessage("one")
	if a != "Error:\none" {
		t.Errorf("Got: %s, exp: Error:\\none", a)
	}

	defer func() {
		if r := recover(); r != nil {
			if r != "first message should be a string: 2" {
				t.Error(r)
			}
		}
	}()
	errorMessage("one", 2)

	b := errorMessage("one", "two")
	if b != "one\ntwo" {
		t.Errorf("Got: %s, exp: one\\ntwo", b)
	}

	c := errorMessage("one", "%d %d %d", 1, 2, 3)
	if c != "one\ntwo" {
		t.Errorf("Got: %s, exp: one\\n1 2 3", c)
	}
}

func TestIsDir(t *testing.T) {
	if !isDir("/etc") {
		t.Errorf("isDir failed 1")
	}

	if isDir("testdata/a/b/c.txt") {
		t.Errorf("isDir failed 2")
	}

	if isDir("testdata/a/b/notexist/") {
		t.Errorf("isDir failed 3")
	}
}

func TestIsFile(t *testing.T) {
	if isFile("/etc") {
		t.Errorf("isFile failed 1")
	}

	if !isFile("/etc/passwd") {
		t.Errorf("isFile failed 2")
	}

	if isFile("testdata/a/b/filenotexist") {
		t.Errorf("isFile failed 3")
	}
}
