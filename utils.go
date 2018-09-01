package assert

import (
	"fmt"
	"os"
)

// errorMessage homogenizes the error message.
func errorMessage(s string, msg ...interface{}) string {
	if len(msg) == 0 {
		return fmt.Sprintf("Error:\n%s", s)
	}
	f, ok := msg[0].(string)
	if !ok {
		panic(fmt.Sprintf("first message should be a string: %#v", msg[0]))
	}
	if len(msg) == 1 {
		return fmt.Sprintf("%s\n%s", f, s)
	}
	if len(msg) > 1 {
		return fmt.Sprintf("%s\n%s", fmt.Sprintf(f, msg[1:]...), s)
	}
	return s
}

// isDir returns true if the path is an existing directory.
func isDir(path string) bool {
	if path == "" {
		return false
	}
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}

// isFile returns true if the path is an existing regular file.
func isFile(path string) bool {
	if path == "" {
		return false
	}
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.Mode().IsRegular()
}

// exists returns true if the directory or the file exists.
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
