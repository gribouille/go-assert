package assert

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCp(t *testing.T) {
	dir, err := ioutil.TempDir("", "go-assert-")
	if err != nil {
		panic(err)
	}

	dest := filepath.Join(dir, "a")
	err = Cp("examples/testdata/a", dest)
	if err != nil {
		t.Error(err)
	}

	st, err := os.Stat(dest)
	if err != nil {
		t.Error(err)
	}
	if !st.IsDir() {
		t.Errorf("%s is not a directory", dest)
	}

	_, err = os.Stat(filepath.Join(dest, "b", "c.txt"))
	if err != nil {
		t.Error(err)
	}
}
