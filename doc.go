/*Package assert provides a simple and no intrusive Go test library.


Installation

Use go get tool:

	go get -u github.com/gribouille/go-assert

Or dep tool:

	dep ensure -add github.com/gribouille/go-assert


Behavior

All errors that are not assertions are panic, for example:

	a.EqualTemplate(got, filename, data)

will cause a panic if the filename does not exist or if the data does not match
with the template.

You could customize the behavior of the library with the environment variables:

	- GO_ASSERT_STACK: show the stacktrace if the test fails
	- GO_ASSERT_FATAL: uses fatal errors
	- GO_ASSERT_TMP_DISABLE: disable the deletion of temporary directory with ItTmp and ItEnv
	- GO_ASSERT_OS: test for a specific OS (the possible values are similar to runtime.GOOS)

or with the NewCustom constructor.

	a := T.NewCustom(t, true, false)


Examples

See the examples/*.go for more examples.

	func TestXXX(t *testing.T) {
		T.New(t).
			It("sub test 1", func(a *T.Assert) {
				a.Equal(...)
				// ...
			}).
			It("sub test 2", func(a *T.Assert) {
				// ...
			})
	}

*/
package assert
