package assert_test

// This test executes the examples and compare the go test output.

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
	"testing"
)

func TestExamples(t *testing.T) {
	fixtures := []struct {
		Test, Expected string
		Comp           func(t *testing.T, a, b string)
	}{
		{Test: "examples/assert_test.go", Expected: "examples/assert_test.out", Comp: compareAssert},
		{Test: "examples/fatal_test.go", Expected: "examples/fatal_test.out", Comp: compareFatal},
		{Test: "examples/stack_test.go", Expected: "examples/stack_test.out", Comp: compareStack},
	}
	for _, fix := range fixtures {
		test(t, fix.Test, fix.Expected, fix.Comp)
	}
}

func test(t *testing.T, test, expected string, comp func(t *testing.T, a, b string)) {
	var outBuf, errBuf bytes.Buffer
	c := exec.Command("go", "test", "-v", test)
	c.Stdout = &outBuf
	c.Stderr = &errBuf

	if err := c.Run(); err != nil {
		exiterr, ok := err.(*exec.ExitError)
		if !ok {
			t.Fatal(err)
		}
		status, ok := exiterr.Sys().(syscall.WaitStatus)
		if !ok {
			t.Fatal(err)
		}
		if status.ExitStatus() != 1 {
			t.Fatalf("Error: %s\n, Status: %d\nStdout: %s\n, Stderr: %s\n", err,
				status.ExitStatus(), outBuf.String(), errBuf.String())
		}
	}

	fi, err := ioutil.ReadFile(expected)
	if err != nil {
		t.Fatal(err)
	}

	exp := string(fi)
	got := outBuf.String()
	comp(t, exp, got)
}

// compare the expected go test output and take into account the variadic lines.
func compareAssert(t *testing.T, a, b string) {
	sa := strings.Split(a, "\n")
	sb := strings.Split(b, "\n")
	msg := "Result mismatch at the line %d\nExp: %s\nGot: %s"
	for i, ra := range sa {
		// Warning with the temporary directory name, this name changes for each execution.
		if i == 84 {
			rule := `^WARNING: Temporary directory deletion canceled:`
			m, _ := regexp.MatchString(rule, ra)
			if !m {
				t.Fatalf(msg, i+1, ra, sb[i])
			}
			m, _ = regexp.MatchString(rule, sb[i])
			if !m {
				t.Fatalf(msg, i+1, ra, sb[i])
			}
			continue
		}
		if i == 107 {
			// The duration time can varied between the execution.
			rule := `^FAIL.*$`
			m, _ := regexp.MatchString(rule, ra)
			if !m {
				t.Fatalf(msg, i+1, ra, sb[i])
			}
			m, _ = regexp.MatchString(rule, sb[i])
			if !m {
				t.Fatalf(msg, i+1, ra, sb[i])
			}
			continue
		}

		if i > len(sb)-1 {
			t.Fatalf(msg, i+1, ra, "none")
		}
		if ra != sb[i] {
			t.Fatalf(msg, i+1, ra, sb[i])
		}
	}
	if len(sa) != len(sb) {
		t.Fatalf("Result mismatch, expected number of lines: %d, got: %d", len(sa), len(sb))
	}
}

func compareFatal(t *testing.T, a, b string) {
	sa := strings.Split(a, "\n")
	sb := strings.Split(b, "\n")
	msg := "Result mismatch at the line %d\nExp: %s\nGot: %s"
	for i, ra := range sa {
		if i == 11 {
			rule := `^FAIL	command-line-arguments	0\.00[0-9]s$`
			m, _ := regexp.MatchString(rule, ra)
			if !m {
				t.Fatalf(msg, i+1, ra, sb[i])
			}
			m, _ = regexp.MatchString(rule, sb[i])
			if !m {

				t.Fatalf(msg, i+1, ra, sb[i])
			}
			continue
		}
		if ra != sb[i] {
			t.Fatalf(msg, i+1, ra, sb[i])
		}
	}
}

func compareStack(t *testing.T, a, b string) {
	m, _ := regexp.MatchString(a, b)
	if !m {
		t.Fatalf("stack mismatch: %s", a)
	}
}
