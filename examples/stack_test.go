package examples

import (
	"os"
	"testing"

	T "github.com/gribouille/go-assert"
)

func TestStack(t *testing.T) {
	a := T.NewCustom(t, false, true)
	a.
		Equal(3, 3, "Equal success").
		Equal("not", "equal", "Equal failed and show stack") // <= stack
	a.
		Equal("a", "b", "this message is showed")
}

func TestStackEnv(t *testing.T) {
	os.Setenv("GO_ASSERT_STACK", "1")
	a := T.New(t)
	a.
		Equal(3, 3, "Equal success").
		Equal("not", "equal", "Equal failed and show stack") // <= stack
	a.
		Equal("a", "b", "this message is showed")
}
