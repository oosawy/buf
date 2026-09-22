//go:build unix

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// bufBaseDir returns the buf directory under XDG_RUNTIME_DIR when set
// (already per-user and cleared on logout), or a per-uid directory under
// the system temp dir otherwise.
func bufBaseDir() string {
	if runtimeDir := os.Getenv("XDG_RUNTIME_DIR"); runtimeDir != "" {
		return filepath.Join(runtimeDir, "buf")
	}
	return filepath.Join(os.TempDir(), fmt.Sprintf("buf-%d", os.Getuid()))
}

// shellKey identifies the calling shell so each shell gets its own buffer.
// It prefers the controlling terminal's device number, since that stays
// stable across the shell's child processes regardless of which of their
// own fds are piped. It falls back to the parent process's PID when there
// is no controlling terminal (e.g. no tty at all).
func shellKey() string {
	if tty, err := os.OpenFile("/dev/tty", os.O_RDONLY, 0); err == nil {
		defer tty.Close()
		var st syscall.Stat_t
		if err := syscall.Fstat(int(tty.Fd()), &st); err == nil {
			return fmt.Sprintf("tty-%d", st.Rdev)
		}
	}
	return fmt.Sprintf("ppid-%d", os.Getppid())
}
