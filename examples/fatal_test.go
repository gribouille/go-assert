package examples

import (
	"os"
	"testing"

	T "github.com/gribouille/go-assert"
)

func TestFatal(t *testing.T) {
	a := T.NewCustom(t, true, false)
	a.
		Equal(3, 3, "Equal success").
		Equal("not", "equal", "Equal failed with fatal") // <= fatal
	a.
		Equal("a", "b", "this message will not showed")
}

func TestFatalEnv(t *testing.T) {
	os.Setenv("GO_ASSERT_FATAL", "1")
	a := T.New(t)
	a.
		Equal(3, 3, "Equal success").
		Equal("not", "equal", "Equal failed with fatal") // <= fatal
	a.
		Equal("a", "b", "this message will not showed")
}
