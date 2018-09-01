package assert

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
	"testing"
)

// tmpDir creates a temporary directory.
func tmpDir(fn func(dir string)) {
	dir, err := ioutil.TempDir("", "go-testing-")
	if err != nil {
		panic(err)
	}
	if os.Getenv("GO_ASSERT_TMP_DISABLE") == "1" {
		defer fmt.Printf("WARNING: Temporary directory deletion canceled: %s\n", dir)
	} else {
		defer os.RemoveAll(dir) // clean up
	}
	fn(dir)
}

// captureOutput captures the standard output.
func captureOutput(fn func()) (string, string) {
	oldOut := os.Stdout
	oldErr := os.Stderr

	ro, wo, _ := os.Pipe()
	re, we, _ := os.Pipe()
	os.Stdout = wo
	os.Stderr = we
	fn()
	outC := make(chan string)
	errC := make(chan string)

	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var outB, errB bytes.Buffer
		io.Copy(&outB, ro)
		io.Copy(&errB, re)
		outC <- outB.String()
		errC <- errB.String()
	}()

	// back to normal state
	wo.Close()
	we.Close()
	os.Stdout = oldOut // restoring the real stdout
	os.Stderr = oldErr // restoring the real stdout
	out := <-outC
	err := <-errC
	return out, err
}

// crashTest captures the return code of functions that uses os.Exit.
func crashTest(t *testing.T, fn func()) (int, string, string) {
	// The forked process executes only fn and exit.
	if os.Getenv("GO_TESTING_CRASH_TEST") == "1" {
		fn()
		panic("the function must exit to use crashTest")
	}

	// Fork the process and run the function in this process.
	var outBuf, errBuf bytes.Buffer
	cmd := exec.Command(os.Args[0], fmt.Sprintf("-test.run=^%s$", t.Name()))
	cmd.Env = append(os.Environ(), "GO_TESTING_CRASH_TEST=1")
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if err := cmd.Start(); err != nil {
		t.Fatalf("unexpected error in crash test: %s", err.Error())
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), outBuf.String(), errBuf.String()
			}
		} else {
			t.Fatalf("unexpected error in crash test: %s", err.Error())
		}
	}
	return 0, outBuf.String(), errBuf.String()
}
