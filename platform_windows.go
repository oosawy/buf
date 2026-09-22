//go:build windows

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

var (
	kernel32             = syscall.NewLazyDLL("kernel32.dll")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")
)

// bufBaseDir returns the buf directory under the system temp dir. Windows'
// temp dir is already per-user (e.g. %LOCALAPPDATA%\Temp), so no extra
// disambiguation is needed.
func bufBaseDir() string {
	return filepath.Join(os.TempDir(), "buf")
}

// shellKey identifies the calling shell so each shell gets its own buffer.
// It prefers the console window handle attached to the process, since that
// stays stable across the shell's child processes regardless of which of
// their own fds are piped. It falls back to the parent process's PID when
// there is no console attached (e.g. running detached).
func shellKey() string {
	hwnd, _, _ := procGetConsoleWindow.Call()
	if hwnd != 0 {
		return fmt.Sprintf("con-%x", hwnd)
	}
	return fmt.Sprintf("ppid-%d", os.Getppid())
}
