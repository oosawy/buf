package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mattn/go-isatty"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	path, err := bufPath()
	if err != nil {
		return fail(err)
	}

	if len(args) > 0 && args[0] == "--" {
		if len(args) < 2 {
			fmt.Fprintln(os.Stderr, "buf: usage: buf -- <command> [args...]")
			return 1
		}
		return runPipe(path, args[1:])
	}

	if isatty.IsTerminal(os.Stdin.Fd()) {
		return readMode(path)
	}
	if err := teeToBuffer(path, os.Stdin); err != nil {
		return fail(err)
	}
	return 0
}

// teeToBuffer copies src to stdout while also saving it as the buffer at
// path, replacing the previous contents atomically.
func teeToBuffer(path string, src io.Reader) error {
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	_, copyErr := io.Copy(io.MultiWriter(f, os.Stdout), src)
	f.Close()
	if copyErr != nil {
		os.Remove(tmp)
		return copyErr
	}
	return os.Rename(tmp, path)
}

// readMode writes the saved buffer to stdout.
func readMode(path string) int {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		return fail(err)
	}
	defer f.Close()
	if _, err := io.Copy(os.Stdout, f); err != nil {
		return fail(err)
	}
	return 0
}

// fail prints err to stderr and returns the process exit code for it.
func fail(err error) int {
	fmt.Fprintln(os.Stderr, "buf:", err)
	return 1
}

// runPipe feeds the saved buffer to cmd's stdin, tees cmd's stdout to
// os.Stdout, and saves that output as the new buffer.
func runPipe(path string, cmdArgs []string) int {
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stderr = os.Stderr

	if inFile, err := os.Open(path); err == nil {
		defer inFile.Close()
		cmd.Stdin = inFile
	} else if !os.IsNotExist(err) {
		return fail(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fail(err)
	}
	if err := cmd.Start(); err != nil {
		return fail(err)
	}

	teeErr := teeToBuffer(path, stdout)
	waitErr := cmd.Wait()

	if teeErr != nil {
		return fail(teeErr)
	}
	if exitErr, ok := waitErr.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	if waitErr != nil {
		return fail(waitErr)
	}
	return 0
}

// bufPath returns the path of the buffer file for the calling shell,
// creating its directory if needed. bufBaseDir is platform-specific (see
// platform_*.go).
func bufPath() (string, error) {
	dir := bufBaseDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "buf-"+shellKey()), nil
}
