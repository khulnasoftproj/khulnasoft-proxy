//go:build !windows
// +build !windows

package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func (runner *Runner) RunXSysExec(args ...string) error {
	cmdName := filepath.Base(args[0])
	if cmdName == "khulnasoft" {
		return errKhulnasoftCantBeExecuted
	}

	khulnasoftPath, err := absoluteKhulnasoftPath()
	if err != nil {
		return fmt.Errorf("get khulnasoft's absolute path: %w", err)
	}
	if err := unix.Exec(khulnasoftPath, append([]string{"khulnasoft", "exec", "--", cmdName}, args[1:]...), os.Environ()); err != nil {
		return fmt.Errorf("execute khulnasoft: %w", err)
	}
	return nil
}
